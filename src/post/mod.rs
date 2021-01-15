use chrono::prelude::*;
use color_eyre::eyre::{eyre, Result, WrapErr};
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
    pub mentions: Vec<mi::WebMention>,
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

pub async fn load(dir: &str, mi: Option<&mi::Client>) -> Result<Vec<Post>> {
    let mut result: Vec<Post> = vec![];

    for path in glob(&format!("{}/*.markdown", dir))?.filter_map(Result::ok) {
        log::debug!("loading {:?}", path);
        let body =
            fs::read_to_string(path.clone()).wrap_err_with(|| format!("can't read {:?}", path))?;
        let (fm, content_offset) = frontmatter::Data::parse(body.clone().as_str())
            .wrap_err_with(|| format!("can't parse frontmatter of {:?}", path))?;
        let markup = &body[content_offset..];
        let date = NaiveDate::parse_from_str(&fm.clone().date, "%Y-%m-%d")
            .map_err(|why| eyre!("error parsing date in {:?}: {}", path, why))?;
        let link = format!("{}/{}", dir, path.file_stem().unwrap().to_str().unwrap());
        let mentions: Vec<mi::WebMention> = match mi {
            None => vec![],
            Some(mi) => mi
                .mentioners(format!("https://christine.website/{}", link))
                .await
                .map_err(|why| tracing::error!("error: can't load mentions for {}: {}", link, why))
                .unwrap_or(vec![]),
        };

        result.push(Post {
            front_matter: fm,
            link: link,
            body: markup.to_string(),
            body_html: crate::app::markdown::render(&markup)
                .wrap_err_with(|| format!("can't parse markdown for {:?}", path))?,
            date: {
                DateTime::<Utc>::from_utc(
                    NaiveDateTime::new(date, NaiveTime::from_hms(0, 0, 0)),
                    Utc,
                )
                .with_timezone(&Utc)
                .into()
            },
            mentions: mentions,
        })
    }

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
        load("blog", None).await.expect("posts to load");
    }

    #[tokio::test]
    async fn gallery() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("gallery", None).await?;
        Ok(())
    }

    #[tokio::test]
    async fn talks() -> Result<()> {
        let _ = pretty_env_logger::try_init();
        load("talks", None).await?;
        Ok(())
    }
}
