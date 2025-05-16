// @generated automatically by Diesel CLI.

diesel::table! {
    messages (id) {
        id -> Int8,
        content -> Text,
        timestamp -> Timestamptz,
    }
}
