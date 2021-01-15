use std::default::Default;

use errors::*;
use feed::{Attachment, Author, Feed};
use item::{Content, Item};

/// Feed Builder
///
/// This is used to programmatically build up a Feed object,
/// which can be serialized later into a JSON string
pub struct Builder(Feed);

impl Builder {
    pub fn new() -> Builder {
        Builder(Feed::default())
    }

    pub fn title<I: Into<String>>(mut self, t: I) -> Builder {
        self.0.title = t.into();
        self
    }

    pub fn home_page_url<I: Into<String>>(mut self, url: I) -> Builder {
        self.0.home_page_url = Some(url.into());
        self
    }

    pub fn feed_url<I: Into<String>>(mut self, url: I) -> Builder {
        self.0.feed_url = Some(url.into());
        self
    }

    pub fn description<I: Into<String>>(mut self, desc: I) -> Builder {
        self.0.description = Some(desc.into());
        self
    }

    pub fn user_comment<I: Into<String>>(mut self, cmt: I) -> Builder {
        self.0.user_comment = Some(cmt.into());
        self
    }

    pub fn next_url<I: Into<String>>(mut self, url: I) -> Builder {
        self.0.next_url = Some(url.into());
        self
    }

    pub fn icon<I: Into<String>>(mut self, url: I) -> Builder {
        self.0.icon = Some(url.into());
        self
    }

    pub fn favicon<I: Into<String>>(mut self, url: I) -> Builder {
        self.0.favicon = Some(url.into());
        self
    }

    pub fn author(mut self, author: Author) -> Builder {
        self.0.author = Some(author);
        self
    }

    pub fn expired(mut self) -> Builder {
        self.0.expired = Some(true);
        self
    }

    pub fn item(mut self, item: Item) -> Builder {
        self.0.items.push(item);
        self
    }

    pub fn build(self) -> Feed {
        self.0
    }
}

/// Builder object for an item in a feed
pub struct ItemBuilder {
    pub id: Option<String>,
    pub url: Option<String>,
    pub external_url: Option<String>,
    pub title: Option<String>,
    pub content: Option<Content>,
    pub summary: Option<String>,
    pub image: Option<String>,
    pub banner_image: Option<String>,
    pub date_published: Option<String>,
    pub date_modified: Option<String>,
    pub author: Option<Author>,
    pub tags: Option<Vec<String>>,
    pub attachments: Option<Vec<Attachment>>,
}

impl ItemBuilder {
    pub fn new() -> ItemBuilder {
        ItemBuilder {
            id: None,
            url: None,
            external_url: None,
            title: None,
            content: None,
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

    pub fn title<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.title = Some(i.into());
        self
    }

    pub fn image<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.image = Some(i.into());
        self
    }

    pub fn id<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.id = Some(i.into());
        self
    }

    pub fn url<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.url = Some(i.into());
        self
    }

    pub fn external_url<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.external_url = Some(i.into());
        self
    }

    pub fn date_modified<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.date_modified = Some(i.into());
        self
    }

    pub fn date_published<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        self.date_published = Some(i.into());
        self
    }

    pub fn tags(mut self, tags: Vec<String>) -> ItemBuilder {
        self.tags = Some(tags);
        self
    }

    pub fn author(mut self, who: Author) -> ItemBuilder {
        self.author = Some(who);
        self
    }

    pub fn content_html<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        match self.content {
            Some(Content::Text(t)) => {
                self.content = Some(Content::Both(i.into(), t));
            }
            _ => {
                self.content = Some(Content::Html(i.into()));
            }
        }
        self
    }

    pub fn content_text<I: Into<String>>(mut self, i: I) -> ItemBuilder {
        match self.content {
            Some(Content::Html(s)) => {
                self.content = Some(Content::Both(s, i.into()));
            }
            _ => {
                self.content = Some(Content::Text(i.into()));
            }
        }
        self
    }

    pub fn build(self) -> Result<Item> {
        if self.id.is_none() || self.content.is_none() {
            return Err("missing field 'id' or 'content_*'".into());
        }
        Ok(Item {
            id: self.id.unwrap(),
            url: self.url,
            external_url: self.external_url,
            title: self.title,
            content: self.content.unwrap(),
            summary: self.summary,
            image: self.image,
            banner_image: self.banner_image,
            date_published: self.date_published,
            date_modified: self.date_modified,
            author: self.author,
            tags: self.tags,
            attachments: self.attachments,
        })
    }
}
