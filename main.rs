use std::sync::Arc;
use std::collections::HashMap;
use std::sync::atomic::{Ordering, AtomicUsize};

use tokio::sync::Mutex;
use actix_files::Files;
use rcgen::CertifiedKey;
use actix_web_actors::ws;
use rustls::ServerConfig;
use actix::{Actor, Handler, Message, Recipient, ActorContext, AsyncContext, StreamHandler};
use actix_web::{get, App, Error, HttpRequest, HttpResponse, HttpServer, web::{Data, Path, Payload}};

static NEXT_ID: AtomicUsize = AtomicUsize::new(0);

struct WsConn {
    id: usize,
    #[allow(unused)] room: Box::<str>,
    addr: Recipient::<Broadcast>
}

struct WsActor {
    id: usize,
    room: Box::<str>,
    rooms: Data::<AtomicRooms>
}

#[derive(Message)]
#[rtype(result = "()")]
struct Broadcast(String);

type Rooms = HashMap::<Box::<str>, HashMap::<usize, WsConn>>;
type AtomicRooms = Arc::<Mutex::<Rooms>>;

impl Actor for WsActor {
    type Context = ws::WebsocketContext::<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        let addr = ctx.address().recipient();
        let conn = WsConn {
            id: self.id,
            addr,
            room: Box::clone(&self.room)
        };

        actix::spawn({
            let room = Box::clone(&self.room);
            let rooms = Data::clone(&self.rooms);
            async move {
                let mut rooms = rooms.lock().await;
                rooms.entry(room)
                    .or_default()
                    .insert(conn.id, conn);
            }
        });
    }

    fn stopped(&mut self, _: &mut Self::Context) {
        actix::spawn({
            let id = self.id;
            let room = Box::clone(&self.room);
            let rooms = Data::clone(&self.rooms);
            async move {
                let mut rooms = rooms.lock().await;
                if let Some(room_conns) = rooms.get_mut(&room) {
                    room_conns.remove(&id);
                    if room_conns.is_empty() {
                        rooms.remove(&room);
                    }
                }
            }
        });
    }
}

impl StreamHandler::<Result<ws::Message, ws::ProtocolError>> for WsActor {
    fn handle(&mut self, msg: Result::<ws::Message, ws::ProtocolError>, ctx: &mut Self::Context) {
        match msg {
            Ok(ws::Message::Text(text)) => {
                let src_id = self.id;
                let room = Box::clone(&self.room);
                let rooms = Data::clone(&self.rooms);

                // broadcast message
                actix::spawn(async move {
                    let rooms = rooms.lock().await;
                    if let Some(conns) = rooms.get(&room) {
                        conns.iter().filter(|(id, ..)| **id != src_id).for_each(|(.., conn)| {
                            _ = conn.addr.do_send(Broadcast(text.to_string()))
                        })
                    }
                });
            }
            Ok(ws::Message::Ping(msg)) => ctx.pong(&msg),
            Ok(ws::Message::Close(reason)) => {
                ctx.close(reason);
                ctx.stop()
            }
            _ => {}
        }
    }
}

impl Handler<Broadcast> for WsActor {
    type Result = ();

    #[inline]
    fn handle(&mut self, msg: Broadcast, ctx: &mut Self::Context) {
        ctx.text(msg.0)
    }
}

#[get("/ws/{room}")]
async fn ws_route(rq: HttpRequest, stream: Payload, room: Path::<String>, rooms: Data::<AtomicRooms>) -> Result::<HttpResponse, Error> {
    let id = NEXT_ID.fetch_add(1, Ordering::Relaxed);
    let actor = WsActor {
        id,
        room: room.into_inner().into_boxed_str(),
        rooms: rooms.clone(),
    };

    ws::start(actor, &rq, stream)
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let CertifiedKey { cert, key_pair } = rcgen::generate_simple_self_signed(vec![
        "127.0.0.1".to_owned(), "localhost".to_owned()
    ]).unwrap();

    let key = rustls::PrivateKey(key_pair.serialize_der());
    let cert_chain = rustls::Certificate(cert.der().to_vec());

    let cfg = ServerConfig::builder()
        .with_safe_defaults()
        .with_no_client_auth()
        .with_single_cert(vec![cert_chain], key)
        .expect("failed to create server config");

    let rooms = Arc::new(Mutex::new(Rooms::new()));

    println!("starting server at <https://localhost:8443>");
    println!("please note:");
    println!("  * Note the HTTPS in the URL; there is no HTTP -> HTTPS redirect.");
    println!("  * You'll need to accept the invalid TLS certificate as it is self-signed.");
    println!("  * Some browsers or OSs may not allow the webcam to be used by multiple pages at once.");

    HttpServer::new(move || {
        App::new()
            .service(ws_route)
            .app_data(Data::new(Arc::clone(&rooms)))
            .service(Files::new("/", ".").index_file("index.html"))
    }).bind_rustls_021("0.0.0.0:8443", cfg)?.run().await
}
