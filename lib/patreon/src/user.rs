use serde::{Deserialize, Serialize};
use chrono::prelude::*;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    pub id: String,
    pub first_name: String,
    pub last_name: String,
    pub full_name: String,
    pub vanity: String,
    pub email: String,
    pub about: String,
    pub facebook_id: String,
    pub gender: i32,
    pub has_password: bool,
    pub image_url: String,
    pub thumb_url: String,
    pub youtube: String,
    pub twitter: String,
    pub facebook: String,
    pub is_email_verified: bool,
    pub is_suspended: bool,
    pub is_deleted: bool,
    pub is_nuked: bool,
    pub created: DateTime<Utc>,
    pub url: String,
    pub discord_id: String,
}
