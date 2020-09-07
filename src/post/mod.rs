use anyhow::{anyhow, Result};
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
    pub date: DateTime<FixedOffset>,
}

impl Into<jsonfeed::Item> for Post {
    fn into(self) -> jsonfeed::Item {
        let mut result = jsonfeed::Item::builder()
            .title(self.front_matter.title)
            .content_html(self.body_html)
            .content_text(self.body)
            .id(format!("https://christine.website/{}", self.link))
            .url(format!("https://christine.website/{}", self.link))
            .date_published(self.date.to_rfc3339())
            .author(
                jsonfeed::Author::new()
                    .name("Christine Dodrill")
                    .url("https://christine.website")
                    .avatar("https://christine.website/static/img/avatar.png"),
            );

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

impl Post {
    pub fn detri(&self) -> String {
        self.date.format("M%m %d %Y").to_string()
    }
}

pub fn load(dir: &str) -> Result<Vec<Post>> {
    let mut result: Vec<Post> = vec![];

    for path in glob(&format!("{}/*.markdown", dir))?.filter_map(Result::ok) {
        log::debug!("loading {:?}", path);
        let body = fs::read_to_string(path.clone()).expect("things to work");
        let (fm, content_offset) = frontmatter::Data::parse(body.clone().as_str()).expect("stuff to work");
        let markup = &body[content_offset..];
        let date = NaiveDate::parse_from_str(&fm.clone().date, "%Y-%m-%d")?;

        result.push(Post {
            front_matter: fm,
            link: format!("{}/{}", dir, path.file_stem().unwrap().to_str().unwrap()),
            body: markup.to_string(),
            body_html: crate::app::markdown(&markup),
            date: {
                DateTime::<Utc>::from_utc(
                    NaiveDateTime::new(date, NaiveTime::from_hms(0, 0, 0)),
                    Utc,
                )
                .with_timezone(&Utc)
                .into()
            },
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
    fn blog() {
        let _ = pretty_env_logger::try_init();
        load("blog").expect("posts to load");
    }

    #[test]
    fn gallery() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("gallery")?;
        Ok(())
    }

    #[test]
    fn talks() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("talks")?;
        Ok(())
    }
}
