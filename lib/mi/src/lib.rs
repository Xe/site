use color_eyre::eyre::Result;
use reqwest::header;
use serde::Deserialize;
use tracing::instrument;

const USER_AGENT_BASE: &str = concat!(
    "library/",
    env!("CARGO_PKG_NAME"),
    "/",
    env!("CARGO_PKG_VERSION")
);

pub struct Client {
    cli: reqwest::Client,
    base_url: String,
}

impl Client {
    pub fn new(token: String, user_agent: String) -> Result<Self> {
        let mut headers = header::HeaderMap::new();
        headers.insert(
            header::AUTHORIZATION,
            header::HeaderValue::from_str(&token.clone())?,
        );

        let cli = reqwest::Client::builder()
            .user_agent(&format!("{} {}", user_agent, USER_AGENT_BASE))
            .default_headers(headers)
            .build()?;

        Ok(Self {
            cli: cli,
            base_url: "https://mi.within.website".to_string(),
        })
    }

    #[instrument(skip(self), err)]
    pub async fn mentioners(&self, url: String) -> Result<Vec<WebMention>> {
        Ok(self
            .cli
            .get(&format!("{}/api/webmention/for", self.base_url))
            .query(&[("target", &url)])
            .send()
            .await?
            .error_for_status()?
            .json()
            .await?)
    }

    #[instrument(skip(self), err)]
    pub async fn refresh(&self) -> Result<()> {
        self.cli
            .post("https://mi.within.website/api/blog/refresh")
            .send()
            .await?
            .error_for_status()?;
        Ok(())
    }
}

#[derive(Debug, Deserialize, Eq, PartialEq, Clone)]
pub struct WebMention {
    pub source: String,
    pub title: Option<String>,
}

#[cfg(test)]
mod tests {
    #[test]
    fn it_works() {
        assert_eq!(2 + 2, 4);
    }
}
