use std::fmt;
use std::default::Default;

use feed::{Author, Attachment};
use builder::ItemBuilder;

use serde::ser::{Serialize, Serializer, SerializeStruct};
use serde::de::{self, Deserialize, Deserializer, Visitor, MapAccess};

/// Represents the `content_html` and `content_text` attributes of an item
#[derive(Debug, Clone, PartialEq, Deserialize, Serialize)]
pub enum Content {
    Html(String),
    Text(String),
    Both(String, String),
}

/// Represents an item in a feed
#[derive(Debug, Clone, PartialEq)]
pub struct Item {
    pub id: String,
    pub url: Option<String>,
    pub external_url: Option<String>,
    pub title: Option<String>,
    pub content: Content,
    pub summary: Option<String>,
    pub image: Option<String>,
    pub banner_image: Option<String>,
    pub date_published: Option<String>, // todo DateTime objects?
    pub date_modified: Option<String>,
    pub author: Option<Author>,
    pub tags: Option<Vec<String>>,
    pub attachments: Option<Vec<Attachment>>,
}

impl Item {
    pub fn builder() -> ItemBuilder {
        ItemBuilder::new()
    }
}

impl Default for Item {
    fn default() -> Item {
        Item {
            id: "".to_string(),
            url: None,
            external_url: None,
            title: None,
            content: Content::Text("".into()),
            summary: None,
            image: None,
            banner_image: None,
            date_published: None,
            date_modified: None,
            author: None,
            tags: None,
            attachments: None,
        }
    }
}

impl Serialize for Item {
    fn serialize<S>(&self, serializer: S) -> Result<S::Ok, S::Error>
        where S: Serializer
    {
        let mut state = serializer.serialize_struct("Item", 14)?;
        state.serialize_field("id", &self.id)?;
        if self.url.is_some() {
            state.serialize_field("url", &self.url)?;
        }
        if self.external_url.is_some() {
            state.serialize_field("external_url", &self.external_url)?;
        }
        if self.title.is_some() {
            state.serialize_field("title", &self.title)?;
        }
        match self.content {
            Content::Html(ref s) => {
                state.serialize_field("content_html", s)?;
                state.serialize_field("content_text", &None::<Option<&str>>)?;
            },
            Content::Text(ref s) => {
                state.serialize_field("content_html", &None::<Option<&str>>)?;
                state.serialize_field("content_text", s)?;
            },
            Content::Both(ref s, ref t) => {
                state.serialize_field("content_html", s)?;
                state.serialize_field("content_text", t)?;
            },
        };
        if self.summary.is_some() {
            state.serialize_field("summary", &self.summary)?;
        }
        if self.image.is_some() {
            state.serialize_field("image", &self.image)?;
        }
        if self.banner_image.is_some() {
            state.serialize_field("banner_image", &self.banner_image)?;
        }
        if self.date_published.is_some() {
            state.serialize_field("date_published", &self.date_published)?;
        }
        if self.date_modified.is_some() {
            state.serialize_field("date_modified", &self.date_modified)?;
        }
        if self.author.is_some() {
            state.serialize_field("author", &self.author)?;
        }
        if self.tags.is_some() {
            state.serialize_field("tags", &self.tags)?;
        }
        if self.attachments.is_some() {
            state.serialize_field("attachments", &self.attachments)?;
        }
        state.end()
    }
}

impl<'de> Deserialize<'de> for Item {
    fn deserialize<D>(deserializer: D) -> Result<Self, D::Error> 
        where D: Deserializer<'de>
    {
        enum Field {
            Id,
            Url,
            ExternalUrl,
            Title,
            ContentHtml,
            ContentText,
            Summary,
            Image,
            BannerImage,
            DatePublished,
            DateModified,
            Author,
            Tags,
            Attachments,
        };

        impl<'de> Deserialize<'de> for Field {
            fn deserialize<D>(deserializer: D) -> Result<Self, D::Error>
                where D: Deserializer<'de>
            {
                struct FieldVisitor;

                impl<'de> Visitor<'de> for FieldVisitor {
                    type Value = Field;

                    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
                        formatter.write_str("non-expected field")
                    }

                    fn visit_str<E>(self, value: &str) -> Result<Field, E>
                        where E: de::Error
                    {
                        match value {
                            "id" => Ok(Field::Id),
                            "url" => Ok(Field::Url),
                            "external_url" => Ok(Field::ExternalUrl),
                            "title" => Ok(Field::Title),
                            "content_html" => Ok(Field::ContentHtml),
                            "content_text" => Ok(Field::ContentText),
                            "summary" => Ok(Field::Summary),
                            "image" => Ok(Field::Image),
                            "banner_image" => Ok(Field::BannerImage),
                            "date_published" => Ok(Field::DatePublished),
                            "date_modified" => Ok(Field::DateModified),
                            "author" => Ok(Field::Author),
                            "tags" => Ok(Field::Tags),
                            "attachments" => Ok(Field::Attachments),
                            _ => Err(de::Error::unknown_field(value, FIELDS)),
                        }
                    }
                }

                deserializer.deserialize_identifier(FieldVisitor)
            }
        }

        struct ItemVisitor;

        impl<'de> Visitor<'de> for ItemVisitor {
            type Value = Item;
            fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
                formatter.write_str("non-expected thing")
            }

            fn visit_map<V>(self, mut map: V) -> Result<Item, V::Error>
                where V: MapAccess<'de>
            {
                let mut id = None;
                let mut url = None;
                let mut external_url = None;
                let mut title = None;
                let mut content_html: Option<String> = None;
                let mut content_text: Option<String> = None;
                let mut summary = None;
                let mut image = None;
                let mut banner_image = None;
                let mut date_published = None;
                let mut date_modified = None;
                let mut author = None;
                let mut tags = None;
                let mut attachments = None;

                while let Some(key) = map.next_key()? {
                    match key {
                        Field::Id => {
                            if id.is_some() {
                                return Err(de::Error::duplicate_field("id"));
                            }
                            id = Some(map.next_value()?);
                        },
                        Field::Url => {
                            if url.is_some() {
                                return Err(de::Error::duplicate_field("url"));
                            }
                            url = map.next_value()?;
                        },
                        Field::ExternalUrl => {
                            if external_url.is_some() {
                                return Err(de::Error::duplicate_field("external_url"));
                            }
                            external_url = map.next_value()?;
                        },
                        Field::Title => {
                            if title.is_some() {
                                return Err(de::Error::duplicate_field("title"));
                            }
                            title = map.next_value()?;
                        },
                        Field::ContentHtml => {
                            if content_html.is_some() {
                                return Err(de::Error::duplicate_field("content_html"));
                            }
                            content_html = map.next_value()?;
                        },
                        Field::ContentText => {
                            if content_text.is_some() {
                                return Err(de::Error::duplicate_field("content_text"));
                            }
                            content_text = map.next_value()?;
                        },
                        Field::Summary => {
                            if summary.is_some() {
                                return Err(de::Error::duplicate_field("summary"));
                            }
                            summary = map.next_value()?;
                        },
                        Field::Image => {
                            if image.is_some() {
                                return Err(de::Error::duplicate_field("image"));
                            }
                            image = map.next_value()?;
                        },
                        Field::BannerImage => {
                            if banner_image.is_some() {
                                return Err(de::Error::duplicate_field("banner_image"));
                            }
                            banner_image = map.next_value()?;
                        },
                        Field::DatePublished => {
                            if date_published.is_some() {
                                return Err(de::Error::duplicate_field("date_published"));
                            }
                            date_published = map.next_value()?;
                        },
                        Field::DateModified => {
                            if date_modified.is_some() {
                                return Err(de::Error::duplicate_field("date_modified"));
                            }
                            date_modified = map.next_value()?;
                        },
                        Field::Author => {
                            if author.is_some() {
                                return Err(de::Error::duplicate_field("author"));
                            }
                            author = map.next_value()?;
                        },
                        Field::Tags => {
                            if tags.is_some() {
                                return Err(de::Error::duplicate_field("tags"));
                            }
                            tags = map.next_value()?;
                        },
                        Field::Attachments => {
                            if attachments.is_some() {
                                return Err(de::Error::duplicate_field("attachments"));
                            }
                            attachments = map.next_value()?;
                        },
                    }
                }

                let id = id.ok_or_else(|| de::Error::missing_field("id"))?;
                let content = match (content_html, content_text) {
                    (Some(s), Some(t)) => {
                        Content::Both(s.to_string(), t.to_string())
                    },
                    (Some(s), _) => {
                        Content::Html(s.to_string())
                    },
                    (_, Some(t)) => {
                        Content::Text(t.to_string())
                    },
                    _ => return Err(de::Error::missing_field("content_html or content_text")),
                };

                Ok(Item {
                    id,
                    url,
                    external_url,
                    title,
                    content,
                    summary,
                    image,
                    banner_image,
                    date_published,
                    date_modified,
                    author,
                    tags,
                    attachments,
                })
            }
        }

        const FIELDS: &'static [&'static str] = &[
            "id",
            "url",
            "external_url",
            "title",
            "content",
            "summary",
            "image",
            "banner_image",
            "date_published",
            "date_modified",
            "author",
            "tags",
            "attachments",
        ];
        deserializer.deserialize_struct("Item", FIELDS, ItemVisitor)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use feed::Author;
    use serde_json;

    #[test]
    #[allow(non_snake_case)]
    fn serialize_item__content_html() {
        let item = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Html("<p>content</p>".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(
            serde_json::to_string(&item).unwrap(),
            r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":"<p>content</p>","content_text":null,"summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#
        );
    }

    #[test]
    #[allow(non_snake_case)]
    fn serialize_item__content_text() {
        let item = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Text("content".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(
            serde_json::to_string(&item).unwrap(),
            r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":null,"content_text":"content","summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#
        );
    }

    #[test]
    #[allow(non_snake_case)]
    fn serialize_item__content_both() {
        let item = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Both("<p>content</p>".into(), "content".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(
            serde_json::to_string(&item).unwrap(),
            r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":"<p>content</p>","content_text":"content","summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#
        );
    }

    #[test]
    #[allow(non_snake_case)]
    fn deserialize_item__content_html() {
        let json = r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":"<p>content</p>","content_text":null,"summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#;
        let item: Item = serde_json::from_str(&json).unwrap();
        let expected = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Html("<p>content</p>".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(item, expected);
    }

    #[test]
    #[allow(non_snake_case)]
    fn deserialize_item__content_text() {
        let json = r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":null,"content_text":"content","summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#;
        let item: Item = serde_json::from_str(&json).unwrap();
        let expected = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Text("content".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(item, expected);
    }

    #[test]
    #[allow(non_snake_case)]
    fn deserialize_item__content_both() {
        let json = r#"{"id":"1","url":"http://example.com/feed.json","external_url":"http://example.com/feed.json","title":"feed title","content_html":"<p>content</p>","content_text":"content","summary":"feed summary","image":"http://img.com/blah","banner_image":"http://img.com/blah","date_published":"2017-01-01 10:00:00","date_modified":"2017-01-01 10:00:00","author":{"name":"bob jones","url":"http://example.com","avatar":"http://img.com/blah"},"tags":["json","feed"],"attachments":[]}"#;
        let item: Item = serde_json::from_str(&json).unwrap();
        let expected = Item {
            id: "1".into(),
            url: Some("http://example.com/feed.json".into()),
            external_url: Some("http://example.com/feed.json".into()),
            title: Some("feed title".into()),
            content: Content::Both("<p>content</p>".into(), "content".into()),
            summary: Some("feed summary".into()),
            image: Some("http://img.com/blah".into()),
            banner_image: Some("http://img.com/blah".into()),
            date_published: Some("2017-01-01 10:00:00".into()),
            date_modified: Some("2017-01-01 10:00:00".into()),
            author: Some(Author::new().name("bob jones").url("http://example.com").avatar("http://img.com/blah")),
            tags: Some(vec!["json".into(), "feed".into()]),
            attachments: Some(vec![]),
        };
        assert_eq!(item, expected);
    }
}

