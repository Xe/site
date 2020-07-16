use serde::{Deserialize, Serialize};
use chrono::prelude::*;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Campaign {
    pub id: String,
    pub summary: String,
    pub creation_name: String,
    pub display_patron_goals: bool,
    pub pay_per_name: String,
    pub one_liner: String,
    pub main_video_embed: String,
    pub main_video_url: String,
    pub image_small_url: String,
    pub image_url: String,
    pub thanks_video_url: String,
    pub thanks_embed: String,
    pub thanks_msg: String,
    pub is_charged_immediately: bool,
    pub is_monthly: bool,
    pub is_nsfw: bool,
    pub is_plural: bool,
    pub created_at: DateTime<Utc>,
    pub published_at: DateTime<Utc>,
    pub pledge_url: String,
    pub pledge_sum: i32,
    pub patron_count: u32,
    pub creation_coult: u32,
    pub outstanding_payment_amount_cents: u64,
}
