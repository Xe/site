use crate::{post::Post, signalboost::Person};
use anyhow::Result;
use atom_syndication as atom;
use comrak::{markdown_to_html, ComrakOptions};
use serde::Deserialize;
use std::{fs, path::PathBuf};

#[derive(Clone, Deserialize)]
pub struct Config {
    #[serde(rename = "clackSet")]
    clack_set: Vec<String>,
    signalboost: Vec<Person>,
    port: u16,
    #[serde(rename = "resumeFname")]
    resume_fname: PathBuf,
}

pub fn markdown(inp: &str) -> String {
    let mut options = ComrakOptions::default();

    options.extension.autolink = true;
    options.extension.table = true;
    options.extension.description_lists = true;
    options.extension.superscript = true;
    options.extension.strikethrough = true;
    options.extension.footnotes = true;

    options.render.unsafe_ = true;

    markdown_to_html(inp, &options)
}

async fn patrons() -> Result<Option<patreon::Users>> {
    use patreon::*;
    let creds: Credentials = envy::prefixed("PATREON_").from_env().unwrap();
    let cli = Client::new(creds);

    if let Ok(camp) = cli.campaign().await {
        let id = camp.data[0].id.clone();
        if let Ok(users) = cli.pledges(id).await {
            Ok(Some(users))
        } else {
            Ok(None)
        }
    } else {
        Ok(None)
    }
}

pub const ICON: &'static str = "https://christine.website/static/img/avatar.png";

pub struct State {
    pub cfg: Config,
    pub signalboost: Vec<Person>,
    pub resume: String,
    pub blog: Vec<Post>,
    pub gallery: Vec<Post>,
    pub talks: Vec<Post>,
    pub everything: Vec<Post>,
    pub jf: jsonfeed::Feed,
    pub rf: rss::Channel,
    pub af: atom::Feed,
    pub sitemap: Vec<u8>,
    pub patrons: Option<patreon::Users>,
}

pub async fn init(cfg: PathBuf) -> Result<State> {
    let cfg: Config = serde_dhall::from_file(cfg).parse()?;
    let sb = cfg.signalboost.clone();
    let resume = fs::read_to_string(cfg.resume_fname.clone())?;
    let resume: String = markdown(&resume);
    let blog = crate::post::load("blog")?;
    let gallery = crate::post::load("gallery")?;
    let talks = crate::post::load("talks")?;
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

    let mut ri: Vec<rss::Item> = vec![];
    let mut ai: Vec<atom::Entry> = vec![];

    let mut jfb = jsonfeed::Feed::builder()
        .title("Christine Dodrill's Blog")
        .description("My blog posts and rants about various technology things.")
        .author(
            jsonfeed::Author::new()
                .name("Christine Dodrill")
                .url("https://christine.website")
                .avatar(ICON),
        )
        .feed_url("https://christine.website/blog.json")
        .user_comment("This is a JSON feed of my blogposts. For more information read: https://jsonfeed.org/version/1")
        .home_page_url("https://christine.website")
        .icon(ICON)
        .favicon(ICON);

    for post in &everything {
        let post = post.clone();
        jfb = jfb.item(post.clone().into());
        ri.push(post.clone().into());
        ai.push(post.clone().into());
    }

    let af = {
        let mut af = atom::FeedBuilder::default();
        af.title("Christine Dodrill's Blog");
        af.id("https://christine.website/blog");
        af.generator({
            let mut generator = atom::Generator::default();
            generator.set_value(env!("CARGO_PKG_NAME"));
            generator.set_version(env!("CARGO_PKG_VERSION").to_string());
            generator.set_uri("https://github.com/Xe/site".to_string());

            generator
        });
        af.entries(ai);

        af.build().unwrap()
    };

    let rf = {
        let mut rf = rss::ChannelBuilder::default();
        rf.title("Christine Dodrill's Blog");
        rf.link("https://christine.website/blog");
        rf.generator(crate::APPLICATION_NAME.to_string());
        rf.items(ri);

        rf.build().unwrap()
    };

    let mut sm: Vec<u8> = vec![];
    let smw = sitemap::writer::SiteMapWriter::new(&mut sm);
    let mut urlwriter = smw.start_urlset()?;
    for url in &[
        "https://christine.website/resume",
        "https://christine.website/contact",
        "https://christine.website/",
        "https://christine.website/blog",
        "https://christine.website/signalboost",
    ] {
        urlwriter.url(*url)?;
    }

    for post in &everything {
        urlwriter.url(format!("https://christine.website/{}", post.link))?;
    }

    urlwriter.end()?;

    Ok(State {
        cfg: cfg,
        signalboost: sb,
        resume: resume,
        blog: blog,
        gallery: gallery,
        talks: talks,
        everything: everything,
        jf: jfb.build(),
        af: af,
        rf: rf,
        sitemap: sm,
        patrons: patrons().await?,
    })
}

#[cfg(test)]
mod tests {
    use anyhow::Result;
    #[tokio::test]
    async fn init() -> Result<()> {
        super::init("./config.dhall".into()).await?;
        Ok(())
    }
}
