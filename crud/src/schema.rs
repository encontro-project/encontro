// @generated automatically by Diesel CLI.

diesel::table! {
    messages (id) {
        id -> Int4,
        content -> Text,
    }
}
