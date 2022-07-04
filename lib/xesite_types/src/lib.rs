use serde::{Deserialize, Serialize};

#[derive(Eq, PartialEq, Deserialize, Default, Debug, Serialize, Clone)]
pub struct Frontmatter {
    #[serde(default = "frontmatter_about")]
    pub about: String,
    #[serde(skip_serializing)]
    pub title: String,
    #[serde(skip_serializing)]
    pub date: String,
    #[serde(skip_serializing)]
    pub author: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub series: Option<String>,
    #[serde(skip_serializing)]
    pub tags: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub slides_link: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub image: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub thumb: Option<String>,
    #[serde(skip_serializing)]
    pub redirect_to: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub vod: Option<Vod>,
}

fn frontmatter_about() -> String {
    "https://xeiaso.net/blog/api-jsonfeed-extensions#_xesite_frontmatter".to_string()
}

#[derive(Eq, PartialEq, Deserialize, Default, Debug, Serialize, Clone)]
pub struct Vod {
    pub twitch: String,
    pub youtube: String,
}
