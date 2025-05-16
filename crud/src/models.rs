use chrono::{DateTime, Utc};
use serde::{Serialize, Deserialize};
use diesel::{Queryable, Insertable, Selectable};

use super::schema::messages;

#[derive(Queryable, Serialize, Selectable)]
#[diesel(table_name = messages)]
pub struct Message {
    pub id: i64,
    pub content: String,
    pub timestamp: DateTime<Utc>
}

#[derive(Insertable, Deserialize)]
#[diesel(table_name = messages)]
pub struct NewMessage {
    pub content: String,
}
