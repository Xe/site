use serde::{Deserialize, Serialize};

#[derive(Eq, PartialEq, Debug, Clone, Serialize, Deserialize)]
pub struct Article {
    #[serde(rename = "@context")]
    pub context: String,
    #[serde(rename = "@type")]
    pub r#type: String,
    pub headline: String,
    pub image: String,
    pub url: String,
    #[serde(rename = "datePublished")]
    pub date_published: String,
}
