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
    pub date: NaiveDate,
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
