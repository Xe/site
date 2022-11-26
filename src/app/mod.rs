use crate::{post::Post, signalboost::Person};
use chrono::prelude::*;
use color_eyre::eyre::Result;
use std::{path::PathBuf, sync::Arc};
use tracing::{error, instrument};

pub mod config;
pub mod poke;

pub use config::*;

#[instrument]
async fn patrons() -> Result<Option<patreon::Users>> {
    let mut p = dirs::home_dir().unwrap_or(".".into());
    p.push(".patreon.json");
    if !p.exists() {
        info!("{:?} does not exist", p);
        return Ok(None);
    }

    let mut cli = patreon::Client::new()?;

    if let Err(why) = cli.refresh_token().await {
        error!("error getting refresh token: {}", why);
    }

    match cli.campaign().await {
        Ok(camp) => {
            let id = camp.data[0].id.clone();

            match cli.pledges(id).await {
                Ok(users) => Ok(Some(users)),
                Err(why) => {
                    error!("error getting pledges: {}", why);
                    Ok(None)
                }
            }
        }
        Err(why) => {
            error!("error getting patreon campaign: {}", why);
            Ok(None)
        }
    }
}

pub const ICON: &'static str = "https://xeiaso.net/static/img/avatar.png";

pub struct State {
    pub cfg: Arc<Config>,
    pub signalboost: Vec<Person>,
    pub blog: Vec<Post>,
    pub gallery: Vec<Post>,
    pub talks: Vec<Post>,
    pub everything: Vec<Post>,
    pub jf: xe_jsonfeed::Feed,
    pub sitemap: Vec<u8>,
    pub patrons: Option<patreon::Users>,
    pub mi: mi::Client,
}

pub async fn init(cfg: PathBuf) -> Result<State> {
    let cfg: Arc<Config> = Arc::new(serde_dhall::from_file(cfg).parse()?);
    let sb = cfg.signalboost.clone();
    let mi = mi::Client::new(
        cfg.clone().mi_token.clone(),
        crate::APPLICATION_NAME.to_string(),
    )?;
    let blog = crate::post::load("blog").await?;
    let gallery = crate::post::load("gallery").await?;
    let talks = crate::post::load("talks").await?;
    let mut everything: Vec<Post> = vec![];

    {
        let blog = blog.clone();
        let gallery = gallery.clone();
        let talks = talks.clone();
        everything.extend(blog.iter().cloned());
        everything.extend(gallery.iter().cloned());
        everything.extend(talks.iter().cloned());
    };

    everything.sort();
    everything.reverse();

    let today = Utc::now().date_naive();
    let everything: Vec<Post> = everything
        .into_iter()
        .filter(|p| today.num_days_from_ce() >= p.date.num_days_from_ce())
        .take(5)
        .collect();

    let mut jfb = xe_jsonfeed::Feed::builder()
        .title("Xe's Blog")
        .description("My blog posts and rants about various technology things.")
        .author(
            xe_jsonfeed::Author::new()
                .name("Xe")
                .url("https://xeiaso.net")
                .avatar(ICON),
        )
        .feed_url("https://xeiaso.net/blog.json")
        .user_comment("This is a JSON feed of my blogposts. For more information read: https://jsonfeed.org/version/1")
        .home_page_url("https://xeiaso.net")
        .icon(ICON)
        .favicon(ICON);

    for post in &everything {
        let post = post.clone();
        jfb = jfb.item(post.clone().into());
    }

    let mut sm: Vec<u8> = vec![];
    let smw = sitemap::writer::SiteMapWriter::new(&mut sm);
    let mut urlwriter = smw.start_urlset()?;
    for url in &[
        "https://xeiaso.net/resume",
        "https://xeiaso.net/contact",
        "https://xeiaso.net/",
        "https://xeiaso.net/blog",
        "https://xeiaso.net/signalboost",
    ] {
        urlwriter.url(*url)?;
    }

    for post in &blog {
        urlwriter.url(format!("https://xeiaso.net/{}", post.link))?;
    }
    for post in &gallery {
        urlwriter.url(format!("https://xeiaso.net/{}", post.link))?;
    }
    for post in &talks {
        urlwriter.url(format!("https://xeiaso.net/{}", post.link))?;
    }

    urlwriter.end()?;

    Ok(State {
        mi,
        cfg,
        signalboost: sb,
        blog,
        gallery,
        talks,
        everything,
        jf: jfb.build(),
        sitemap: sm,
        patrons: patrons().await?,
    })
}

#[cfg(test)]
mod tests {
    use color_eyre::eyre::Result;
    #[tokio::test]
    async fn init() -> Result<()> {
        super::init("./config.dhall".into()).await?;
        Ok(())
    }
}
