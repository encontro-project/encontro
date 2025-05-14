use diesel::{Queryable, Insertable};
use serde::{Serialize, Deserialize};

use super::schema::messages;

#[derive(Queryable, Serialize)]
pub struct Message {
    pub id: i32,
    pub content: String,
}

#[derive(Insertable, Deserialize)]
#[diesel(table_name = messages)]
pub struct NewMessage {
    pub content: String,
}
