use serde::{Deserialize, Serialize};
use chrono::prelude::*;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Pledge {
    pub id: String,
    pub amount_cents: u32,
    pub created_at: DateTime<Utc>,
    pub declined_since: DateTime<Utc>,
    pub pledge_cap_cents: u32,
    pub patron_pays_fees: bool,
    pub total_historical_amount_cents: Option<u32>,
    pub is_paused: Option<bool>,
    pub has_shipping_address: Option<bool>,
    pub outstanding_payment_amount_cents: Option<u32>,
}
