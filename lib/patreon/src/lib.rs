use std::{fs, io};

use chrono::prelude::*;
use serde::{Deserialize, Serialize};
use thiserror::Error;
use tracing::{debug, error, instrument};
use url::Url;

pub type Campaigns = Vec<Object<Campaign>>;
pub type Pledges = Vec<Object<Pledge>>;
pub type Users = Vec<Object<User>>;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Campaign {
    pub summary: String,
    pub creation_name: String,
    pub display_patron_goals: bool,
    pub pay_per_name: String,
    pub one_liner: Option<String>,
    pub main_video_embed: Option<String>,
    pub main_video_url: Option<String>,
    pub image_small_url: String,
    pub image_url: String,
    pub thanks_video_url: Option<String>,
    pub thanks_embed: Option<String>,
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
    pub creation_count: u32,
    pub outstanding_payment_amount_cents: u64,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Pledge {
    pub amount_cents: u32,
    pub created_at: String,
    pub declined_since: Option<String>,
    pub pledge_cap_cents: u32,
    pub patron_pays_fees: bool,
    pub total_historical_amount_cents: Option<u32>,
    pub is_paused: Option<bool>,
    pub has_shipping_address: Option<bool>,
    pub outstanding_payment_amount_cents: Option<u32>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct User {
    pub first_name: String,
    pub last_name: String,
    pub full_name: String,
    pub vanity: Option<String>,
    pub about: Option<String>,
    pub gender: i32,
    pub image_url: String,
    pub thumb_url: String,
    pub created: DateTime<Utc>,
    pub url: String,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct RefreshGrant {
    pub access_token: String,
    pub refresh_token: String,
    pub expires_in: serde_json::Value,
    pub scope: serde_json::Value,
    pub token_type: String,
}

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Error, Debug)]
pub enum Error {
    #[error("json error: {0}")]
    Json(#[from] serde_json::Error),

    #[error("request error: {0}")]
    Request(#[from] reqwest::Error),

    #[error("{0}")]
    IO(#[from] io::Error),

    #[error("url parse error: {0}")]
    URLParse(#[from] url::ParseError),
}

#[derive(Debug, Serialize, Deserialize, Clone, Default, Eq, PartialEq)]
pub struct Credentials {
    pub client_id: String,
    pub client_secret: String,
    pub access_token: String,
    pub refresh_token: String,
}

pub struct Client {
    cli: reqwest::Client,
    base_url: String,
    creds: Credentials,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Data<T, U> {
    pub data: T,
    pub included: Option<Vec<U>>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Object<T> {
    pub id: String,
    pub attributes: T,
    pub r#type: String,
    pub links: Option<Links>,
}

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct Links {
    related: String,
}

impl Client {
    pub fn new() -> Result<Self> {
        let mut p = dirs::home_dir().unwrap_or(".".into());
        p.push(".patreon.json");
        let config = fs::read_to_string(p)?;
        let creds = serde_json::from_str(&config)?;

        Ok(Self {
            cli: reqwest::Client::new(),
            base_url: "https://api.patreon.com".into(),
            creds,
        })
    }

    #[instrument(skip(self))]
    pub async fn campaign(&self) -> Result<Data<Vec<Object<Campaign>>, ()>> {
        let data = self
            .cli
            .get(&format!(
                "{}/oauth2/api/current_user/campaigns",
                self.base_url
            ))
            .query(&[("include", "patron.null"), ("includes", "")])
            .header(
                "Authorization",
                format!("Bearer {}", self.creds.access_token),
            )
            .send()
            .await?
            .error_for_status()?
            .text()
            .await?;
        debug!("campaign response: {}", data);
        Ok(serde_json::from_str(&data)?)
    }

    #[instrument(skip(self))]
    pub async fn pledges(&self, camp_id: String) -> Result<Vec<Object<User>>> {
        let data = self
            .cli
            .get(&format!(
                "{}/oauth2/api/campaigns/{}/pledges",
                self.base_url, camp_id
            ))
            .query(&[("include", "patron.null"), ("page[count]", "100")])
            .header(
                "Authorization",
                format!("Bearer {}", self.creds.access_token),
            )
            .send()
            .await?
            .error_for_status()?
            .text()
            .await?;
        debug!("pledges for {}: {}", camp_id, data);
        let data: Data<Vec<Object<Pledge>>, Object<User>> = serde_json::from_str(&data)?;
        Ok(data.included.unwrap())
    }

    /*
    POST www.patreon.com/api/oauth2/token
    ?grant_type=refresh_token
    &refresh_token=<the user‘s refresh_token>
    &client_id=<your client id>
    &client_secret=<your client secret>

    1. grab new creds
    2. serialize new creds to disk
    3. reload current creds in ram
    4. ???
    5. profit!
    */
    #[instrument(skip(self))]
    pub async fn refresh_token(&mut self) -> Result<()> {
        let mut u = Url::parse(&self.base_url)?;
        u.set_path("/api/oauth2/token");
        u.query_pairs_mut()
            .append_pair("grant_type", "refresh_token")
            .append_pair("refresh_token", &self.creds.refresh_token)
            .append_pair("client_id", &self.creds.client_id)
            .append_pair("client_secret", &self.creds.client_secret);

        let rg: RefreshGrant = self
            .cli
            .post(&u.to_string())
            .header(
                "Authorization",
                format!("Bearer {}", self.creds.access_token),
            )
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?;

        let mut creds = self.creds.clone();

        creds.access_token = rg.access_token;
        creds.refresh_token = rg.refresh_token;

        let mut p = dirs::home_dir().unwrap_or(".".into());
        p.push(".patreon.json");
        if p.exists() {
            fs::remove_file(&p)?;
        }
        let mut fout = fs::File::create(p)?;
        serde_json::to_writer(&mut fout, &creds)?;

        self.creds = creds;

        Ok(())
    }
}
