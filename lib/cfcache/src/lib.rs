use reqwest::header;
use tracing::instrument;

pub type Result<T = ()> = std::result::Result<T, Error>;

#[derive(thiserror::Error, Debug)]
pub enum Error {
    #[error("json error: {0}")]
    Json(#[from] serde_json::Error),

    #[error("request error: {0}")]
    Request(#[from] reqwest::Error),

    #[error("invalid header value: {0}")]
    InvalidHeaderValue(#[from] reqwest::header::InvalidHeaderValue),
}

pub struct Client {
    zone_id: String,
    cli: reqwest::Client,
}

static USER_AGENT: &str = concat!(
    "xesite ",
    env!("CARGO_PKG_NAME"),
    "/",
    env!("CARGO_PKG_VERSION")
);

impl Client {
    pub fn new(api_key: String, zone_id: String) -> Result<Self> {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            header::AUTHORIZATION,
            header::HeaderValue::from_str(&format!("Bearer {}", api_key))?,
        );

        let cli = reqwest::Client::builder()
            .user_agent(USER_AGENT)
            .default_headers(headers)
            .build()?;

        Ok(Self { zone_id, cli })
    }

    #[instrument(skip(self), err)]
    pub async fn purge(&self, urls: Vec<String>) -> Result {
        #[derive(serde::Serialize)]
        struct Files {
            files: Vec<String>,
        }

        self.cli
            .post(&format!(
                "https://api.cloudflare.com/client/v4/zones/{}/purge_cache",
                self.zone_id
            ))
            .json(&Files { files: urls })
            .send()
            .await?
            .error_for_status()?;
        Ok(())
    }
}
