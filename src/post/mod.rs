use chrono::prelude::*;
use color_eyre::eyre::{eyre, Result, WrapErr};
use glob::glob;
use serde::{Deserialize, Serialize};
use std::{borrow::Borrow, cmp::Ordering, path::PathBuf};
use tokio::fs;

pub mod frontmatter;
pub mod schemaorg;

#[derive(Eq, PartialEq, Debug, Clone, Serialize, Deserialize)]
pub struct Post {
    pub front_matter: frontmatter::Data,
    pub link: String,
    pub body_html: String,
    pub date: DateTime<FixedOffset>,
    pub mentions: Vec<mi::WebMention>,
    pub new_post: NewPost,
    pub read_time_estimate_minutes: u64,
}

/// Used with the Android app to show information in a widget.
#[derive(Eq, PartialEq, Debug, Clone, Serialize, Deserialize)]
pub struct NewPost {
    pub title: String,
    pub summary: String,
    pub link: String,
}

impl Into<schemaorg::Article> for &Post {
    fn into(self) -> schemaorg::Article {
        schemaorg::Article {
            context: "https://schema.org".to_string(),
            r#type: "Article".to_string(),
            headline: self.front_matter.title.clone(),
            image: "https://xeiaso.net/static/img/avatar.png".to_string(),
            url: format!("https://xeiaso.net/{}", self.link),
            date_published: self.date.format("%Y-%m-%d").to_string(),
        }
    }
}

impl Into<xe_jsonfeed::Item> for Post {
    fn into(self) -> xe_jsonfeed::Item {
        let mut result = xe_jsonfeed::Item::builder()
            .title(self.front_matter.title.clone())
            .content_html(self.body_html)
            .id(format!("https://xeiaso.net/{}", self.link))
            .url(if let Some(url) = self.front_matter.redirect_to.as_ref() {
                url.clone()
            } else {
                format!("https://xeiaso.net/{}", self.link)
            })
            .date_published(self.date.to_rfc3339())
            .author(
                xe_jsonfeed::Author::new()
                    .name("Xe Iaso")
                    .url("https://xeiaso.net")
                    .avatar("https://xeiaso.net/static/img/avatar.png"),
            )
            .xesite_frontmatter(self.front_matter.clone());

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

async fn read_post(dir: &str, fname: PathBuf, cli: &Option<mi::Client>) -> Result<Post> {
    debug!(
        "loading {}",
        fname.clone().into_os_string().into_string().unwrap()
    );

    let body = fs::read_to_string(fname.clone())
        .await
        .wrap_err_with(|| format!("can't read {:?}", fname))?;
    let (front_matter, content_offset) = frontmatter::parse(body.clone().as_str())
        .wrap_err_with(|| format!("can't parse frontmatter of {:?}", fname))?;
    let body = &body[content_offset..];
    let date = NaiveDate::parse_from_str(&front_matter.clone().date, "%Y-%m-%d")
        .map_err(|why| eyre!("error parsing date in {:?}: {}", fname, why))?;
    let link = format!("{}/{}", dir, fname.file_stem().unwrap().to_str().unwrap());
    let body_html = xesite_markdown::render(&body)
        .wrap_err_with(|| format!("can't parse markdown for {:?}", fname))?;
    let date: DateTime<FixedOffset> = DateTime::<Utc>::from_utc(
        NaiveDateTime::new(date, NaiveTime::from_hms_opt(0, 0, 0).unwrap()),
        Utc,
    )
    .with_timezone(&Utc)
    .into();

    let mentions: Vec<mi::WebMention> = match cli {
        Some(cli) => cli
            .mentioners(format!("https://xeiaso.net/{}", link))
            .await
            .map_err(|why| tracing::error!("error: can't load mentions for {}: {}", link, why))
            .unwrap_or(vec![])
            .into_iter()
            .filter(|wm| {
                wm.title.as_ref().unwrap_or(&"".to_string()) != &"Bridgy Response".to_string()
            })
            .map(|wm| {
                let mut wm = wm.clone();
                wm.title = Some(
                    mastodon2text::convert(
                        wm.title.as_ref().unwrap_or(&"".to_string()).to_string(),
                    )
                    .unwrap(),
                );
                wm
            })
            .collect(),
        None => vec![],
    };

    let time_taken = estimated_read_time::text(
        &body,
        &estimated_read_time::Options::new()
            .technical_document(true)
            .technical_difficulty(1)
            .build()
            .unwrap_or_default(),
    );
    let read_time_estimate_minutes = time_taken.seconds() / 60;

    let new_post = NewPost {
        title: front_matter.title.clone(),
        summary: format!("{} minute read", read_time_estimate_minutes),
        link: format!("https://xeiaso.net/{}", link),
    };

    Ok(Post {
        front_matter,
        link,
        body_html,
        date,
        mentions,
        new_post,
        read_time_estimate_minutes,
    })
}

pub async fn load(dir: &str) -> Result<Vec<Post>> {
    let cli = match std::env::var("MI_TOKEN") {
        Ok(token) => mi::Client::new(token.to_string(), crate::APPLICATION_NAME.to_string()).ok(),
        Err(_) => None,
    };

    let futs = glob(&format!("{}/*.markdown", dir))?
        .filter_map(Result::ok)
        .map(|fname| read_post(dir, fname, cli.borrow()));

    let mut result: Vec<Post> = futures::future::join_all(futs)
        .await
        .into_iter()
        .map(Result::unwrap)
        .collect();

    if result.len() == 0 {
        Err(eyre!("no posts loaded"))
    } else {
        result.sort();
        result.reverse();
        Ok(result)
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use color_eyre::eyre::Result;

    #[tokio::test]
    async fn blog() {
        let _ = pretty_env_logger::try_init();
        load("blog").await.expect("posts to load");
    }

    #[tokio::test]
    async fn gallery() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("gallery").await?;
        Ok(())
    }

    #[tokio::test]
    async fn talks() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("talks").await?;
        Ok(())
    }
}
