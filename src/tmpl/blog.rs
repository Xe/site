use super::{base, nag};
use crate::post::{schemaorg::Article, Post};
use maud::{html, Markup, PreEscaped};
use xesite_templates::xeact_component;

fn post_metadata(post: &Post) -> Markup {
    let art: Article = post.into();
    let json = PreEscaped(serde_json::to_string(&art).unwrap());

    html! {
        meta name="twitter:card" content="summary";
        meta name="twitter:site" content="@theprincessxena";
        meta name="twitter:title" content={(post.front_matter.title)};
        meta property="og:type" content="website";
        meta property="og:title" content={(post.front_matter.title)};
        meta property="og:site_name" content="Xe's Blog";
        meta name="description" content={(post.front_matter.title) " - Xe's Blog"};
        meta name="author" content="Xe Iaso";

        @if let Some(redirect_to) = &post.front_matter.redirect_to {
            link rel="canonical" href=(redirect_to);
            meta http-equiv="refresh" content=(format!("0;URL='{redirect_to}'"));
        } @else {
            link rel="canonical" href={"https://xeiaso.net/" (post.link)};
        }

        script type="application/ld+json" {(json)}
    }
}

fn share_button(post: &Post) -> Markup {
    return xeact_component("MastodonShareButton", serde_json::json!({
        "title": post.front_matter.title,
        "series": post.front_matter.series,
        "tags": post.front_matter.tags.as_ref().unwrap_or(&Vec::new())
    }));
}

fn twitch_vod(post: &Post) -> Markup {
    html! {
        @if let Some(vod) = &post.front_matter.vod {
            p {
                "This post was written live on "
                a href="https://www.twitch.tv/princessxen" {"Twitch"}
                ". You can check out the stream recording on "
                a href=(vod.twitch) {"Twitch"}
                " and on "
                a href=(vod.youtube) {"YouTube"}
                ". If you are reading this in the first day or so of this post being published, you will need to watch it on Twitch."
            }
        }
    }
}

pub fn blog(post: &Post, body: PreEscaped<&String>, referer: Option<String>) -> Markup {
    base(
        Some(&post.front_matter.title),
        None,
        html! {
            (post_metadata(post))
            @if !post.front_matter.skip_ads {
                (nag::referer(post, referer))
            }

            article {
                h1 {(post.front_matter.title)}

                (nag::prerelease(post))

                small {
                    "Read time in minutes: "
                    (post.read_time_estimate_minutes)
                }

                div {
                    (body)
                }
            }

            hr;

            (share_button(post))
            (twitch_vod(post))

            p {
                "This article was posted on "
                (post.detri())
                ". Facts and circumstances may have changed since publication. Please "
                a href="/contact" {"contact me"}
                " before jumping to conclusions if something seems wrong or unclear."
            }

            @if let Some(series) = &post.front_matter.series {
                p {
                    "Series: "
                    a href={"/blog/series/" (series)} {(series)}
                }
            }

            @if let Some(tags) = &post.front_matter.tags {
               p {
                   "Tags: "
                    @for tag in tags {
                        code {(tag)}
                        " "
                    }
               }
            }

            @if post.mentions.is_empty() {
                p {
                    "This post was not "
                    a href="https://www.w3.org/TR/webmention/" {"WebMention"}
                    "ed yet. You could be the first!"
                }
            } @else {
                ul {
                    @for mention in &post.mentions {
                        li {
                            a href=(mention.source) {(mention.title.as_ref().unwrap_or(&mention.source))}
                        }
                    }
                }
            }

            p {
                "The art for Mara was drawn by "
                a href="https://selic.re/" {"Selicre"}
                "."
            }

            p {
                "The art for Cadey was drawn by "
                a href="https://artzorastudios.weebly.com/" {"ArtZora Studios"}
                "."
            }

            p {
                "Some of the art for Aoi was drawn by "
                a href="https://twitter.com/Sandra_Thomas01" {"@Sandra_Thomas01"}
                "."
            }
        },
    )
}

pub fn gallery(post: &Post) -> Markup {
    base(
        Some(&post.front_matter.title),
        None,
        html! {
            (post_metadata(post))
             h1 {(post.front_matter.title)}

            (PreEscaped(&post.body_html))

                center {
                    img src=(post.front_matter.image.as_ref().unwrap());
                }

            hr;

            p {
                "This artwork was posted on "
                    (post.detri())
                    "."
            }

            @if let Some(tags) = &post.front_matter.tags {
                p {
                    "Tags: "
                        @for tag in tags {
                            code {(tag)}
                            " "
                        }
                }
            }

            (share_button(post))
        },
    )
}

pub fn talk(post: &Post, body: PreEscaped<&String>, referer: Option<String>) -> Markup {
    base(
        Some(&post.front_matter.title),
        None,
        html! {
            (post_metadata(post))

            @if !post.front_matter.skip_ads {
                (nag::referer(post, referer))
            }

            article {
                h1 {(post.front_matter.title)}

                (nag::prerelease(post))

                (body)
            }

            @if let Some(slides) = &post.front_matter.slides_link {
                a href=(slides) {"Link to the slides"}
            }

            hr;

            (share_button(post))

            p {
                "This talk was posted on "
                (post.detri())
                ". Facts and circumstances may have changed since publication Please "
                a href="/contact" {"contact me"}
                " before jumping to conclusions if something seems wrong or unclear."
            }

            p {
                "The art for Mara was drawn by "
                    a href="https://selic.re/" {"Selicre"}
                "."
            }

            p {
                "The art for Cadey was drawn by "
                    a href="https://artzorastudios.weebly.com/" {"ArtZora Studios"}
                "."
            }

            p {
                "Some of the art for Aoi was drawn by "
                a href="https://twitter.com/Sandra_Thomas01" {"@Sandra_Thomas01"}
                "."
            }
        },
    )
}
