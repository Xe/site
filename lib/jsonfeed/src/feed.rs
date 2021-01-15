use std::default::Default;

use builder::Builder;
use item::Item;

const VERSION_1: &'static str = "https://jsonfeed.org/version/1";

/// Represents a single feed
///
/// # Examples
///
/// ```rust
/// // Serialize a feed object to a JSON string
///
/// # extern crate jsonfeed;
/// # use std::default::Default;
/// # use jsonfeed::Feed;
/// # fn main() {
/// let feed: Feed = Feed::default();
/// assert_eq!(
///     jsonfeed::to_string(&feed).unwrap(),
///     "{\"version\":\"https://jsonfeed.org/version/1\",\"title\":\"\",\"items\":[]}"
/// );
/// # }
/// ```
///
/// ```rust
/// // Deserialize a feed objects from a JSON String
///
/// # extern crate jsonfeed;
/// # use jsonfeed::Feed;
/// # fn main() {
/// let json = "{\"version\":\"https://jsonfeed.org/version/1\",\"title\":\"\",\"items\":[]}";
/// let feed: Feed = jsonfeed::from_str(&json).unwrap();
/// assert_eq!(
///     feed,
///     Feed::default()
/// );
/// # }
/// ```
#[derive(Debug, Clone, PartialEq, Deserialize, Serialize)]
pub struct Feed {
    pub version: String,
    pub title: String,
    pub items: Vec<Item>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub home_page_url: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub feed_url: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub description: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub user_comment: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub next_url: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub icon: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub favicon: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub author: Option<Author>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub expired: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub hubs: Option<Vec<Hub>>,
}

impl Feed {
    /// Used to construct a Feed object
    pub fn builder() -> Builder {
        Builder::new()
    }
}

impl Default for Feed {
    fn default() -> Feed {
        Feed {
            version: VERSION_1.to_string(),
            title: "".to_string(),
            items: vec![],
            home_page_url: None,
            feed_url: None,
            description: None,
            user_comment: None,
            next_url: None,
            icon: None,
            favicon: None,
            author: None,
            expired: None,
            hubs: None,
        }
    }
}

/// Represents an `attachment` for an item
#[derive(Debug, Clone, PartialEq, Deserialize, Serialize)]
pub struct Attachment {
    url: String,
    mime_type: String,
    title: Option<String>,
    size_in_bytes: Option<u64>,
    duration_in_seconds: Option<u64>,
}

/// Represents an `author` in both a feed and a feed item
#[derive(Debug, Clone, PartialEq, Deserialize, Serialize)]
pub struct Author {
    name: Option<String>,
    url: Option<String>,
    avatar: Option<String>,
}

impl Author {
    pub fn new() -> Author {
        Author {
            name: None,
            url: None,
            avatar: None,
        }
    }

    pub fn name<I: Into<String>>(mut self, name: I) -> Self {
        self.name = Some(name.into());
        self
    }

    pub fn url<I: Into<String>>(mut self, url: I) -> Self {
        self.url = Some(url.into());
        self
    }

    pub fn avatar<I: Into<String>>(mut self, avatar: I) -> Self {
        self.avatar = Some(avatar.into());
        self
    }
}

/// Represents a `hub` for a feed
#[derive(Debug, Clone, PartialEq, Deserialize, Serialize)]
pub struct Hub {
    #[serde(rename = "type")]
    type_: String,
    url: String,
}

#[cfg(test)]
mod tests {
    use super::*;
    use serde_json;
    use std::default::Default;

    #[test]
    fn serialize_feed() {
        let feed = Feed {
            version: "https://jsonfeed.org/version/1".to_string(),
            title: "some title".to_string(),
            items: vec![],
            home_page_url: None,
            description: None,
            expired: Some(true),
            ..Default::default()
        };
        assert_eq!(
            serde_json::to_string(&feed).unwrap(),
            r#"{"version":"https://jsonfeed.org/version/1","title":"some title","items":[],"expired":true}"#
        );
    }

    #[test]
    fn deserialize_feed() {
        let json =
            r#"{"version":"https://jsonfeed.org/version/1","title":"some title","items":[]}"#;
        let feed: Feed = serde_json::from_str(&json).unwrap();
        let expected = Feed {
            version: "https://jsonfeed.org/version/1".to_string(),
            title: "some title".to_string(),
            items: vec![],
            ..Default::default()
        };
        assert_eq!(feed, expected);
    }

    #[test]
    fn serialize_attachment() {
        let attachment = Attachment {
            url: "http://example.com".to_string(),
            mime_type: "application/json".to_string(),
            title: Some("some title".to_string()),
            size_in_bytes: Some(1),
            duration_in_seconds: Some(1),
        };
        assert_eq!(
            serde_json::to_string(&attachment).unwrap(),
            r#"{"url":"http://example.com","mime_type":"application/json","title":"some title","size_in_bytes":1,"duration_in_seconds":1}"#
        );
    }

    #[test]
    fn deserialize_attachment() {
        let json = r#"{"url":"http://example.com","mime_type":"application/json","title":"some title","size_in_bytes":1,"duration_in_seconds":1}"#;
        let attachment: Attachment = serde_json::from_str(&json).unwrap();
        let expected = Attachment {
            url: "http://example.com".to_string(),
            mime_type: "application/json".to_string(),
            title: Some("some title".to_string()),
            size_in_bytes: Some(1),
            duration_in_seconds: Some(1),
        };
        assert_eq!(attachment, expected);
    }

    #[test]
    fn serialize_author() {
        let author = Author {
            name: Some("bob jones".to_string()),
            url: Some("http://example.com".to_string()),
            avatar: Some("http://img.com/blah".to_string()),
        };
        assert_eq!(
            serde_json::to_string(&author).unwrap(),
            r#"{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"}"#
        );
    }

    #[test]
    fn deserialize_author() {
        let json =
            r#"{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"}"#;
        let author: Author = serde_json::from_str(&json).unwrap();
        let expected = Author {
            name: Some("bob jones".to_string()),
            url: Some("http://example.com".to_string()),
            avatar: Some("http://img.com/blah".to_string()),
        };
        assert_eq!(author, expected);
    }

    #[test]
    fn serialize_hub() {
        let hub = Hub {
            type_: "some-type".to_string(),
            url: "http://example.com".to_string(),
        };
        assert_eq!(
            serde_json::to_string(&hub).unwrap(),
            r#"{"type":"some-type","url":"http://example.com"}"#
        )
    }

    #[test]
    fn deserialize_hub() {
        let json = r#"{"type":"some-type","url":"http://example.com"}"#;
        let hub: Hub = serde_json::from_str(&json).unwrap();
        let expected = Hub {
            type_: "some-type".to_string(),
            url: "http://example.com".to_string(),
        };
        assert_eq!(hub, expected);
    }

    #[test]
    fn deser_podcast() {
        let json = r#"{
  "version": "https://jsonfeed.org/version/1",
  "title": "Timetable",
  "home_page_url": "http://timetable.manton.org/",
  "items": [
    {
      "id": "http://timetable.manton.org/2017/04/episode-45-launch-week/",
      "url": "http://timetable.manton.org/2017/04/episode-45-launch-week/",
      "title": "Episode 45: Launch week",
      "content_html": "I’m rolling out early access to Micro.blog this week. I talk about how the first 2 days have gone, mistakes with TestFlight, and what to do next.",
      "date_published": "2017-04-26T01:09:45+00:00",
      "attachments": [
        {
          "url": "http://timetable.manton.org/podcast-download/139/episode-45-launch-week.mp3",
          "mime_type": "audio/mpeg",
          "size_in_bytes": 5236920
        }
      ]
    }
  ]
}"#;
        serde_json::from_str::<Feed>(&json).expect("Failed to deserialize podcast feed");
    }
}
