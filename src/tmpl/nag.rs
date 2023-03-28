use crate::post::Post;
use lazy_static::lazy_static;
use maud::{html, Markup};
use regex::Regex;

// XXX(S): needed because release builds need this
#[allow(unused)]
use chrono::prelude::*;

lazy_static! {
    static ref LOBSTERS: Regex = Regex::new(r#"^https?://lobste.rs"#).unwrap();
    static ref HACKER_NEWS: Regex = Regex::new(r#"^https?://news.ycombinator.com"#).unwrap();
}

#[cfg(debug_assertions)]
pub fn referer(_: &Post, _: Option<String>) -> Markup {
    html! {
        .warning {
            "This is a development instance of xesite. Things here are probably unfinished or in drafting. Don't take anything here super seriously. If you want to share this to an online aggregator, please don't. Drafts are not finalized yet for a reason. Please don't be the reason I need to implement more advanced security than just obscurity."
        }
        br;
    }
}

#[cfg(not(debug_assertions))]
pub fn referer(post: &Post, referer: Option<String>) -> Markup {
    use xesite_templates::conv as xeblog_conv;

    if referer.is_none() {
        return xesite_templates::advertiser_nag(None);
    }

    let referer = referer.unwrap();

    if Utc::now().date_naive().num_days_from_ce() < post.date.num_days_from_ce() {
        return html! {};
    }

    if HACKER_NEWS.is_match(&referer) {
        return xesite_templates::advertiser_nag(Some(xeblog_conv(
            "Mara".into(),
            "hacker".into(),
            html! {
                "Hello. Before commenting about the author, please read "
                a href="/pronouns" {"this page"}
                " that explains the pronouns that you should be using. tl;dr: the author of this website is NOT male. Please do not use \"he\" or \"him\" when referring to the author."
            },
        )));
    }

    if LOBSTERS.is_match(&referer) {
        return xeblog_conv(
            "Mara".into(),
            "happy".into(),
            html! {
                "Hey, thanks for reading Lobsters! We've disabled the ads to thank you for choosing to use a more ethical aggregator."
            },
        );
    }

    xesite_templates::advertiser_nag(None)
}

#[cfg(debug_assertions)]
pub fn prerelease(_: &Post) -> Markup {
    html! {}
}

#[cfg(not(debug_assertions))]
pub fn prerelease(post: &Post) -> Markup {
    use chrono::prelude::*;
    use xesite_templates::conv as xeblog_conv;

    if Utc::now().date_naive().num_days_from_ce() < post.date.num_days_from_ce() {
        html! {
            .warning {
                (xeblog_conv("Mara".into(), "hacker".into(), html!{
                    "Hey, this post is set to go live on "
                    (format!("{}", post.detri()))
                    " UTC. Right now you are reading a pre-publication version of this post. Please do not share this on social media. This post will automatically go live for everyone on the intended publication date. If you want access to these posts, please join the "
                    a href="https://www.patreon.com/cadey" { "Patreon" }
                    ". It helps me afford the copyeditor that I contract for the technical content I write."
                    br;
                }))
            }
        }
    } else {
        html! {}
    }
}
