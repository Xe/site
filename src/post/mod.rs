use chrono::prelude::*;
use color_eyre::eyre::{eyre, Result, WrapErr};
use glob::glob;
use serde::Serialize;
use std::{cmp::Ordering, path::PathBuf};
use tokio::fs;

pub mod frontmatter;

#[derive(Eq, PartialEq, Debug, Clone)]
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
#[derive(Eq, PartialEq, Debug, Clone, Serialize)]
pub struct NewPost {
    pub title: String,
    pub summary: String,
    pub link: String,
}

impl Into<jsonfeed::Item> for Post {
    fn into(self) -> jsonfeed::Item {
        let mut result = jsonfeed::Item::builder()
            .title(self.front_matter.title)
            .content_html(self.body_html)
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

async fn read_post(dir: &str, fname: PathBuf) -> Result<Post> {
    let body = fs::read_to_string(fname.clone())
        .await
        .wrap_err_with(|| format!("can't read {:?}", fname))?;
    let (front_matter, content_offset) = frontmatter::Data::parse(body.clone().as_str())
        .wrap_err_with(|| format!("can't parse frontmatter of {:?}", fname))?;
    let body = &body[content_offset..];
    let date = NaiveDate::parse_from_str(&front_matter.clone().date, "%Y-%m-%d")
        .map_err(|why| eyre!("error parsing date in {:?}: {}", fname, why))?;
    let link = format!("{}/{}", dir, fname.file_stem().unwrap().to_str().unwrap());
    let body_html = crate::app::markdown::render(&body)
        .wrap_err_with(|| format!("can't parse markdown for {:?}", fname))?;
    let date: DateTime<FixedOffset> =
        DateTime::<Utc>::from_utc(NaiveDateTime::new(date, NaiveTime::from_hms(0, 0, 0)), Utc)
            .with_timezone(&Utc)
            .into();

    let mentions: Vec<mi::WebMention> = match std::env::var("MI_TOKEN") {
        Ok(token) => mi::Client::new(token.to_string(), crate::APPLICATION_NAME.to_string())?
            .mentioners(format!("https://christine.website/{}", link))
            .await
            .map_err(|why| tracing::error!("error: can't load mentions for {}: {}", link, why))
            .unwrap_or(vec![])
            .into_iter()
            .filter(|wm| {
                wm.title.as_ref().unwrap_or(&"".to_string()) != &"Bridgy Response".to_string()
            })
            .collect(),
        Err(_) => vec![],
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
        link: format!("https://christine.website/{}", link),
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
    let futs = glob(&format!("{}/*.markdown", dir))?
        .filter_map(Result::ok)
        .map(|fname| read_post(dir, fname));

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
