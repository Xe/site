use serde::Deserialize;

#[derive(Clone, Debug, Deserialize)]
pub struct Person {
    pub name: String,
    pub tags: Vec<String>,

    #[serde(rename = "gitLink")]
    pub git_link: String,

    pub twitter: String,
}
