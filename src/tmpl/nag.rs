use super::xeblog_conv;
use crate::post::Post;
use chrono::prelude::*;
use lazy_static::lazy_static;
use maud::{html, Markup};
use regex::Regex;

lazy_static! {
    static ref HN: Regex = Regex::new(r#"^https?://news.ycombinator.com"#).unwrap();
    static ref REDDIT: Regex = Regex::new(r#"^https?://((.+).)?reddit.com"#).unwrap();
}

pub fn referer(referer: Option<String>) -> Markup {
    if referer.is_none() {
        return html! {};
    }

    let referer = referer.unwrap();

    let nag = html! {
        script r#async src="https://media.ethicalads.io/media/client/ethicalads.min.js" { "" }
        div.adaptive data-ea-publisher="christinewebsite" data-ea-type="image" data-ea-style="stickybox" {
            .warning {
                (xeblog_conv(
                    "Cadey".into(),
                    "coffee".into(),
                    html! {
                        "Hello! Thank you for visiting my website. You seem to be visiting from a news aggregator and have ads disabled. These ads help pay for running the website and are done by "
                        a href="https://www.ethicalads.io/" { "Ethical Ads" }
                        ". I do not receive detailed analytics on the ads and from what I understand neither does Ethical Ads. If you don't want to disable your ad blocker, please consider donating on "
                        a href="https://patreon.com/cadey" { "Patreon" }
                        ". It helps fund the website's hosting bills and pay for the expensive technical editor that I use for my longer articles. Thanks and be well!"
                    },
                ))
            }
        }
    };

    if HN.is_match(&referer) {
        return nag;
    }

    if REDDIT.is_match(&referer) {
        return nag;
    }

    html! {}
}

pub fn prerelease(post: &Post) -> Markup {
    if Utc::today().num_days_from_ce() < post.date.num_days_from_ce() {
        html! {
            .warning {
                (xeblog_conv("Mara".into(), "hacker".into(), html!{
                    "Hey, this post is set to go live on "
                    (format!("{}", post.detri()))
                    " UTC. Right now you are reading a pre-publication version of this post. Please do not share this on social media. This post will automatically go live for everyone on the intended publication date. If you want access to these posts, please join the "
                    a href="https://patreon.com/cadey" { "Patreon" }
                    ". It helps me afford the copyeditor that I contract for the technical content I write."
                }))
            }
        }
    } else {
        html! {}
    }
}
