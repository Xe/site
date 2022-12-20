use chrono::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct User {
    #[serde(rename = "id")]
    pub id: String,

    #[serde(rename = "type")]
    pub user_type: String,

    #[serde(rename = "following")]
    pub following: String,

    #[serde(rename = "followers")]
    pub followers: String,

    #[serde(rename = "inbox")]
    pub inbox: String,

    #[serde(rename = "outbox")]
    pub outbox: String,

    #[serde(rename = "featured")]
    pub featured: String,

    #[serde(rename = "featuredTags")]
    pub featured_tags: Option<String>,

    #[serde(rename = "preferredUsername")]
    pub preferred_username: String,

    #[serde(rename = "name")]
    pub name: String,

    #[serde(rename = "summary")]
    pub summary: String,

    #[serde(rename = "url")]
    pub url: String,

    #[serde(rename = "manuallyApprovesFollowers")]
    pub manually_approves_followers: bool,

    #[serde(rename = "discoverable")]
    pub discoverable: bool,

    #[serde(rename = "published")]
    pub published: Option<String>,

    #[serde(rename = "devices")]
    pub devices: Option<String>,

    #[serde(rename = "icon")]
    pub icon: Icon,

    #[serde(rename = "image")]
    pub image: Option<Icon>,
}

#[derive(Serialize, Deserialize)]
pub struct Icon {
    #[serde(rename = "type")]
    pub icon_type: String,

    #[serde(rename = "mediaType")]
    pub media_type: Option<String>,

    #[serde(rename = "url")]
    pub url: String,
}

#[derive(Serialize, Deserialize)]
pub struct Toot {
    #[serde(rename = "id")]
    pub id: String,

    #[serde(rename = "type")]
    pub toot_type: String,

    #[serde(rename = "inReplyTo")]
    pub in_reply_to: Option<String>,

    #[serde(rename = "published")]
    pub published: DateTime<Utc>,

    #[serde(rename = "url")]
    pub url: Option<String>,

    #[serde(rename = "attributedTo")]
    pub attributed_to: String,

    #[serde(rename = "to")]
    pub to: Vec<String>,

    #[serde(rename = "cc")]
    pub cc: Vec<String>,

    #[serde(rename = "sensitive")]
    pub sensitive: Option<bool>,

    #[serde(rename = "conversation")]
    pub conversation: String,

    #[serde(rename = "summary")]
    pub summary: Option<String>,

    #[serde(rename = "content")]
    pub content: String,

    #[serde(rename = "contentMap")]
    pub content_map: Option<ContentMap>,

    #[serde(rename = "attachment")]
    pub attachment: Vec<Attachment>,

    #[serde(rename = "tag")]
    pub tag: Vec<Tag>,

    #[serde(rename = "replies")]
    pub replies: Option<Replies>,
}

impl Toot {
    pub fn content_text(&self) -> String {
        html2text::from_read(std::io::Cursor::new(&self.content), 80)
    }
}

#[derive(Serialize, Deserialize)]
pub struct Source {
    #[serde(rename = "content")]
    pub content: String,
    #[serde(rename = "mediaType")]
    pub content_type: String,
}

#[derive(Serialize, Deserialize)]
pub struct Tag {
    #[serde(rename = "type")]
    pub tag_type: String,

    #[serde(rename = "href")]
    pub href: String,

    #[serde(rename = "name")]
    pub name: String,
}

#[derive(Serialize, Deserialize)]
pub struct Attachment {
    #[serde(rename = "type")]
    pub attachment_type: String,

    #[serde(rename = "mediaType")]
    pub media_type: String,

    #[serde(rename = "url")]
    pub url: String,

    #[serde(rename = "name")]
    pub name: Option<String>,

    #[serde(rename = "blurhash")]
    pub blurhash: String,

    #[serde(rename = "width")]
    pub width: i64,

    #[serde(rename = "height")]
    pub height: i64,
}

#[derive(Serialize, Deserialize)]
pub struct ContentMap {
    #[serde(rename = "en")]
    pub en: String,
}

#[derive(Serialize, Deserialize)]
pub struct Replies {
    #[serde(rename = "id")]
    pub id: String,

    #[serde(rename = "type")]
    pub replies_type: String,

    #[serde(rename = "first")]
    pub first: Page,
}

#[derive(Serialize, Deserialize)]
pub struct Page {
    #[serde(rename = "type")]
    pub first_type: String,

    #[serde(rename = "next")]
    pub next: String,

    #[serde(rename = "partOf")]
    pub part_of: String,

    #[serde(rename = "items")]
    pub items: Vec<String>,
}

#[cfg(test)]
pub mod test {
    use super::*;
    use serde_json::from_str;

    #[test]
    fn user() {
        let _: User = from_str(include_str!("./testdata/robocadey.json")).unwrap();
    }

    #[test]
    fn toot() {
        let _attachment: Toot = from_str(include_str!("./testdata/post_attachment.json")).unwrap();
        let _hashtags: Toot = from_str(include_str!("./testdata/post_hashtags.json")).unwrap();
        let _mention: Toot = from_str(include_str!("./testdata/post_mention.json")).unwrap();
    }
}
