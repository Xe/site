use anyhow::{anyhow, Result};
use atom_syndication as atom;
use chrono::prelude::*;
use glob::glob;
use std::{cmp::Ordering, fs};

pub mod frontmatter;

#[derive(Eq, PartialEq, Debug, Clone)]
pub struct Post {
    pub front_matter: frontmatter::Data,
    pub link: String,
    pub body: String,
    pub body_html: String,
    pub date: NaiveDate,
}

impl Into<jsonfeed::Item> for Post {
    fn into(self) -> jsonfeed::Item {
        let mut result = jsonfeed::Item::builder()
            .title(self.front_matter.title)
            .content_html(self.body_html)
            .id(format!("https://christine.website/{}", self.link))
            .url(format!("https://christine.website/{}", self.link))
            .date_published(self.front_matter.date);

        let mut tags: Vec<String> = vec![];

        if let Some(series) = self.front_matter.series {
            tags.push(series);
        }

        if let Some(mut meta_tags) = self.front_matter.tags {
            tags.append(&mut meta_tags);
        }

        if tags.len() != 0 {
            result = result.tags(tags);
        }

        if let Some(image_url) = self.front_matter.image {
            result = result.image(image_url);
        }

        result.build().unwrap()
    }
}

impl Into<atom::Entry> for Post {
    fn into(self) -> atom::Entry {
        let mut content = atom::ContentBuilder::default();

        content.src(format!("https://christine.website/{}", self.link));
        content.content_type(Some("html".into()));
        content.value(Some(xml::escape::escape_str_pcdata(&self.body_html).into()));

        let content = content.build().unwrap();

        let mut result = atom::EntryBuilder::default();
        result.title(self.front_matter.title);
        let mut link = atom::Link::default();
        link.href = format!("https://christine.website/{}", self.link);
        result.links(vec![link]);
        result.content(content);
        result.published(Some(
            DateTime::<Utc>::from_utc(
                NaiveDateTime::new(self.date, NaiveTime::from_hms(0, 0, 0)),
                Utc,
            )
            .with_timezone(&Utc)
            .into(),
        ));

        result.build().unwrap()
    }
}

impl Into<rss::Item> for Post {
    fn into(self) -> rss::Item {
        let mut guid = rss::Guid::default();
        guid.set_value(format!("https://christine.website/{}", self.link));
        let mut result = rss::ItemBuilder::default();
        result.title(Some(self.front_matter.title));
        result.link(format!("https://christine.website/{}", self.link));
        result.guid(guid);
        result.author(Some("me@christine.website".to_string()));
        result.content(self.body_html);
        result.pub_date(self.front_matter.date);

        result.build().unwrap()
    }
}

impl Ord for Post {
    fn cmp(&self, other: &Self) -> Ordering {
        self.partial_cmp(&other).unwrap()
    }
}

impl PartialOrd for Post {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        Some(self.date.cmp(&other.date))
    }
}

pub fn load(dir: &str) -> Result<Vec<Post>> {
    let mut result: Vec<Post> = vec![];

    for path in glob(&format!("{}/*.markdown", dir))?.filter_map(Result::ok) {
        let body = fs::read_to_string(path.clone())?;
        let (fm, content_offset) = frontmatter::Data::parse(body.clone().as_str())?;
        let markup = &body[content_offset..];
        let date = NaiveDate::parse_from_str(&fm.clone().date, "%Y-%m-%d")?;

        result.push(Post {
            front_matter: fm,
            link: format!("{}/{}", dir, path.file_stem().unwrap().to_str().unwrap()),
            body: markup.to_string(),
            body_html: crate::app::markdown(&markup),
            date: date,
        })
    }

    if result.len() == 0 {
        Err(anyhow!("no posts loaded"))
    } else {
        result.sort();
        result.reverse();
        Ok(result)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use anyhow::Result;

    #[test]
    fn blog() -> Result<()> {
        load("blog")?;
        Ok(())
    }

    #[test]
    fn gallery() -> Result<()> {
        load("gallery")?;
        Ok(())
    }

    #[test]
    fn talks() -> Result<()> {
        load("talks")?;
        Ok(())
    }
}
