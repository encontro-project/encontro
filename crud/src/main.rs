use std::sync::{Arc, Mutex};
use std::env;

use diesel::prelude::*;
use diesel::pg::PgConnection;
use diesel::r2d2::{self, ConnectionManager};
use dotenvy::dotenv;

use actix::{Actor, Handler, Recipient, StreamHandler, AsyncContext};
use actix_web::{web, App, HttpServer, HttpRequest, HttpResponse, Error};
use actix_web_actors::ws;

mod models;
mod schema;

use models::{Message, NewMessage};
use schema::messages::dsl as dsl;

type DbPool = r2d2::Pool<ConnectionManager<PgConnection>>;

#[derive(actix::Message)]
#[rtype(result = "()")]
struct Broadcast(String);

type SharedUsers = Arc<Mutex<Vec<Recipient<Broadcast>>>>;
type SharedMessages = Arc<Mutex<Vec<String>>>;

#[derive(Clone)]
struct AppState {
    users: SharedUsers,
    messages: SharedMessages,
    pool: DbPool
}

struct WsUser {
    users: SharedUsers,
    app_data: web::Data<AppState>,
}

impl Actor for WsUser {
    type Context = ws::WebsocketContext<Self>;

    fn started(&mut self, ctx: &mut Self::Context) {
        let addr = ctx.address().recipient();
        self.users.lock().unwrap().push(addr);
    }

    fn stopped(&mut self, _ctx: &mut Self::Context) {
        // Clean up users that are no longer connected
        let mut users = self.users.lock().unwrap();
        users.retain(|r| !r.connected());
    }
}

impl StreamHandler<Result<ws::Message, ws::ProtocolError>> for WsUser {
    fn handle(&mut self, msg: Result<ws::Message, ws::ProtocolError>, _ctx: &mut Self::Context) {
        if let Ok(ws::Message::Text(text)) = msg {
            let recipients = self.users.lock().unwrap().clone();

            self.app_data.messages.lock().unwrap().push(text.to_string());

            for user in recipients {
                let _ = user.do_send(Broadcast(text.to_string()));
            }
        }
    }
}

impl Handler<Broadcast> for WsUser {
    type Result = ();

    fn handle(&mut self, msg: Broadcast, ctx: &mut Self::Context) {
        ctx.text(msg.0);
    }
}

// REST: Add a message
async fn add_message(data: web::Data<AppState>, msg: web::Json<NewMessage>) -> HttpResponse {
    let mut conn = data.pool.get().expect("Couldn't get DB connection");

    let new_msg = diesel::insert_into(dsl::messages)
        .values(&msg.into_inner())
        .get_result::<Message>(&mut conn);

    match new_msg {
        Ok(_) => HttpResponse::Created().body("Message added"),
        Err(_) => HttpResponse::InternalServerError().body("Failed to insert"),
    }
}

// REST: Delete message by index
async fn delete_message(data: web::Data<AppState>, id: web::Path<usize>) -> HttpResponse {
    let mut conn = data.pool.get().expect("Couldn't get DB connection");

    match diesel::delete(dsl::messages.filter(dsl::id.eq(*id as i32))).execute(&mut conn) {
        Ok(1) => HttpResponse::Ok().body("Deleted"),
        Ok(_) => HttpResponse::NotFound().body("Message not found"),
        Err(_) => HttpResponse::InternalServerError().body("DB error"),
    }
}

// REST: Get all messages
async fn get_messages(data: web::Data<AppState>) -> HttpResponse {
    let mut conn = data.pool.get().expect("Couldn't get DB connection");
    match dsl::messages.load::<Message>(&mut conn) {
        Ok(list) => HttpResponse::Ok().json(list),
        Err(_) => HttpResponse::InternalServerError().body("DB error"),
    }
}

// REST: Broadcast a message
async fn broadcast_message(data: web::Data<AppState>, msg: web::Json<String>) -> HttpResponse {
    let users = data.users.lock().unwrap();
    for user in users.iter() {
        let _ = user.do_send(Broadcast(msg.clone()));
    }
    HttpResponse::Ok().body("Broadcasted")
}

// WebSocket route
async fn ws_route(
    req: HttpRequest,
    stream: web::Payload,
    data: web::Data<AppState>,
) -> Result<HttpResponse, Error> {
    let actor = WsUser {
        users: Arc::clone(&data.users),
        app_data: data.clone(),
    };
    ws::start(actor, &req, stream)
}

// Main entry point
#[actix_web::main]
async fn main() -> std::io::Result<()> {
    dotenv().ok();

    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL must be set");
    let mgr = ConnectionManager::<PgConnection>::new(db_url);
    let pool = r2d2::Pool::builder().build(mgr).expect("Failed to create pool");

    let mut conn = pool.get().expect("Couldn't get DB connection");

    let messages = match dsl::messages.load::<Message>(&mut conn) {
        Ok(list) => list.iter().map(|m| m.content.clone()).collect(),
        Err(err) => panic!("DB error: {err:?}")
    };

    let state = AppState {
        messages: Arc::new(Mutex::new(messages)),
        users: Arc::new(Mutex::new(Vec::new())),
        pool: pool.clone(),
    };

    println!("Running server at http://0.0.0.0:8080");

    HttpServer::new(move || {
        App::new()
            .app_data(web::Data::new(state.clone()))
            .route("/ws", web::get().to(ws_route))
            .route("/messages", web::get().to(get_messages))
            .route("/messages", web::post().to(add_message))
            .route("/messages/{id}", web::delete().to(delete_message))
            .route("/broadcast", web::post().to(broadcast_message))
    })
    .bind("0.0.0.0:8080")?
    .run()
    .await
}
