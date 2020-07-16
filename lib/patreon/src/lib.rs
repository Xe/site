#[macro_use]
extern crate jsonapi;

use jsonapi::api::*;
use jsonapi::model::*;
use serde::{Deserialize, Serialize};
use thiserror::Error;

mod campaign;
mod pledge;
mod user;
pub use campaign::Campaign;
pub use pledge::Pledge;
pub use user::User;

jsonapi_model!(Campaign; "campaign");
jsonapi_model!(Pledge; "pledge");
jsonapi_model!(User; "user");

pub type Result<T> = std::result::Result<T, Error>;

#[derive(Error, Debug)]
pub enum Error {
    #[error("json error: {0:?}")]
    Json(#[from] serde_json::Error),
    #[error("request error: {0:?}")]
    Request(#[from] reqwest::Error),
}

#[derive(Debug, Serialize, Deserialize, Clone)]
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

impl Client {
    pub fn new(creds: Credentials) -> Self {
        Self {
            cli: reqwest::Client::new(),
            base_url: "https://api.patreon.com".into(),
            creds: creds,
        }
    }

    pub async fn campaign(&self) -> Result<Campaign> {
        Ok(self
            .cli
            .get(&format!(
                "{}/oauth2/api/current_user/campaigns",
                self.base_url
            ))
            .header(
                "Authorization",
                format!("Bearer {}", self.creds.access_token),
            )
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    pub async fn pledges(&self, camp_id: u32) -> Result<Vec<Pledge>> {
        Ok(self
            .cli
            .get(&format!(
                "{}/oauth2/api/campaigns/{}/pledges",
                self.base_url, camp_id
            ))
            .query(&[("include", "patron.null")])
            .header(
                "Authorization",
                format!("Bearer {}", self.creds.access_token),
            )
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }
}
