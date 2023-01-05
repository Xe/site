use axum::extract::Json;
use serde_json::{json, Value};

pub async fn pronoun_set() -> Json<Value> {
    Json(json!({
        "@type": "PronounSet",
        "nominative": "string",
        "accusative": "string",
        "possessiveDeterminer": "string",
        "possessive": "string",
        "reflexive": "string",
        "singular": "boolean",
    }))
}
