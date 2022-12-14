use crate::post::Post;
use lazy_static::lazy_static;
use maud::{html, Markup};
use regex::Regex;

lazy_static! {
    static ref LOBSTERS: Regex = Regex::new(r#"^https?://lobste.rs"#).unwrap();
}

#[cfg(debug_assertions)]
pub fn referer(_: Option<String>) -> Markup {
    html! {
        .warning {
            "This is a development instance of xesite. Things here are probably unfinished or in drafting. Don't take anything here super seriously. If you want to share this to an online aggregator, please don't. Drafts are not finalized yet for a reason."
        }
        br;
    }
}

#[cfg(not(debug_assertions))]
pub fn referer(referer: Option<String>) -> Markup {
    use xesite_templates::conv as xeblog_conv;
    if referer.is_none() {
        return xesite_templates::advertiser_nag();
    }

    let referer = referer.unwrap();

    if LOBSTERS.is_match(&referer) {
        return html! {
            (xeblog_conv("Mara".into(), "happy".into(), html!{
                "Hey, thanks for reading Lobsters! We've disabled the ads to thank you for choosing to use a more ethical aggregator."
            }))
        };
    }

    if DEV_SERVER.is_match(&referer) {}

    xesite_templates::advertiser_nag()
}

#[cfg(debug_assertions)]
pub fn prerelease(_: &Post) -> Markup {
    html! {}
}

#[cfg(not(debug_assertions))]
pub fn prerelease(post: &Post) -> Markup {
    use chrono::prelude::*;
    if Utc::now().date_naive().num_days_from_ce() < post.date.num_days_from_ce() {
        html! {
            .warning {
                (xeblog_conv("Mara".into(), "hacker".into(), html!{
                    "Hey, this post is set to go live on "
                    (format!("{}", post.detri()))
                    " UTC. Right now you are reading a pre-publication version of this post. Please do not share this on social media. This post will automatically go live for everyone on the intended publication date. If you want access to these posts, please join the "
                    a href="https://patreon.com/cadey" { "Patreon" }
                    ". It helps me afford the copyeditor that I contract for the technical content I write."
                    br;
                }))
            }
        }
    } else {
        html! {}
    }
}
