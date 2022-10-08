use chrono::prelude::*;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct User {
    #[serde(rename = "id")]
    id: String,

    #[serde(rename = "type")]
    user_type: String,

    #[serde(rename = "following")]
    following: String,

    #[serde(rename = "followers")]
    followers: String,

    #[serde(rename = "inbox")]
    inbox: String,

    #[serde(rename = "outbox")]
    outbox: String,

    #[serde(rename = "featured")]
    featured: String,

    #[serde(rename = "featuredTags")]
    featured_tags: String,

    #[serde(rename = "preferredUsername")]
    preferred_username: String,

    #[serde(rename = "name")]
    name: String,

    #[serde(rename = "summary")]
    summary: String,

    #[serde(rename = "url")]
    url: String,

    #[serde(rename = "manuallyApprovesFollowers")]
    manually_approves_followers: bool,

    #[serde(rename = "discoverable")]
    discoverable: bool,

    #[serde(rename = "published")]
    published: String,

    #[serde(rename = "devices")]
    devices: String,

    #[serde(rename = "icon")]
    icon: Icon,

    #[serde(rename = "image")]
    image: Icon,
}

#[derive(Serialize, Deserialize)]
pub struct Icon {
    #[serde(rename = "type")]
    icon_type: String,

    #[serde(rename = "mediaType")]
    media_type: String,

    #[serde(rename = "url")]
    url: String,
}

#[derive(Serialize, Deserialize)]
pub struct Toot {
    #[serde(rename = "id")]
    id: String,

    #[serde(rename = "type")]
    toot_type: String,

    #[serde(rename = "inReplyTo")]
    in_reply_to: Option<String>,

    #[serde(rename = "published")]
    published: DateTime<Utc>,

    #[serde(rename = "url")]
    url: String,

    #[serde(rename = "attributedTo")]
    attributed_to: String,

    #[serde(rename = "to")]
    to: Vec<String>,

    #[serde(rename = "cc")]
    cc: Vec<String>,

    #[serde(rename = "sensitive")]
    sensitive: bool,

    #[serde(rename = "atomUri")]
    atom_uri: String,

    #[serde(rename = "inReplyToAtomUri")]
    in_reply_to_atom_uri: Option<String>,

    #[serde(rename = "conversation")]
    conversation: String,

    #[serde(rename = "content")]
    content: String,

    #[serde(rename = "contentMap")]
    content_map: ContentMap,

    #[serde(rename = "attachment")]
    attachment: Vec<Attachment>,

    #[serde(rename = "tag")]
    tag: Vec<Tag>,

    #[serde(rename = "replies")]
    replies: Replies,
}

#[derive(Serialize, Deserialize)]
pub struct Tag {
    #[serde(rename = "type")]
    tag_type: String,

    #[serde(rename = "href")]
    href: String,

    #[serde(rename = "name")]
    name: String,
}

#[derive(Serialize, Deserialize)]
pub struct Attachment {
    #[serde(rename = "type")]
    attachment_type: String,

    #[serde(rename = "mediaType")]
    media_type: String,

    #[serde(rename = "url")]
    url: String,

    #[serde(rename = "name")]
    name: Option<serde_json::Value>,

    #[serde(rename = "blurhash")]
    blurhash: String,

    #[serde(rename = "width")]
    width: i64,

    #[serde(rename = "height")]
    height: i64,
}

#[derive(Serialize, Deserialize)]
pub struct ContentMap {
    #[serde(rename = "en")]
    en: String,
}

#[derive(Serialize, Deserialize)]
pub struct Replies {
    #[serde(rename = "id")]
    id: String,

    #[serde(rename = "type")]
    replies_type: String,

    #[serde(rename = "first")]
    first: Page,
}

#[derive(Serialize, Deserialize)]
pub struct Page {
    #[serde(rename = "type")]
    first_type: String,

    #[serde(rename = "next")]
    next: String,

    #[serde(rename = "partOf")]
    part_of: String,

    #[serde(rename = "items")]
    items: Vec<String>,
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
