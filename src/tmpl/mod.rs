use crate::{app::*, post::Post, signalboost::Person};
use chrono::prelude::*;
use lazy_static::lazy_static;
use maud::{html, Markup, PreEscaped, Render, DOCTYPE};
use patreon::Users;

pub mod blog;
pub mod nag;

lazy_static! {
    static ref CACHEBUSTER: String = uuid::Uuid::new_v4().to_string().replace("-", "");
}

pub fn base(title: Option<&str>, styles: Option<&str>, content: Markup) -> Markup {
    let now = Utc::now();
    html! {
        (DOCTYPE)
        (PreEscaped(include_str!("./asciiart.txt")))
        html lang="en" {
            head {
                title {
                    @if let Some(title) = title {
                        (title)
                        " - Xe Iaso"
                    } @else {
                        "Xe Iaso"
                    }
                }
                meta name="viewport" content="width=device-width, initial-scale=1.0";
                link rel="stylesheet" href={"/static/css/hack.css?bustCache=" (*CACHEBUSTER)};
                link rel="stylesheet" href={"/static/css/gruvbox-dark.css?bustCache=" (*CACHEBUSTER)};
                link rel="stylesheet" href={"/static/css/shim.css?bustCache=" (*CACHEBUSTER)};
                @match now.month() {
                    12|1|2 => {
                        link rel="stylesheet" href={"/static/css/snow.css?bustCache=" (*CACHEBUSTER)};
                    }
                    _ => {},
                }
                link rel="manifest" href="/static/manifest.json";
                link rel="alternate" title="Xe's Blog" type="application/rss+xml" href="https://xeiaso.net/blog.rss";
                link rel="alternate" title="Xe's Blog" type="application/json" href="https://xeiaso.net/blog.json";
                link rel="apple-touch-icon" sizes="57x57" href="/static/favicon/apple-icon-57x57.png";
                link rel="apple-touch-icon" sizes="60x60" href="/static/favicon/apple-icon-60x60.png";
                link rel="apple-touch-icon" sizes="72x72" href="/static/favicon/apple-icon-72x72.png";
                link rel="apple-touch-icon" sizes="76x76" href="/static/favicon/apple-icon-76x76.png";
                link rel="apple-touch-icon" sizes="114x114" href="/static/favicon/apple-icon-114x114.png";
                link rel="apple-touch-icon" sizes="120x120" href="/static/favicon/apple-icon-120x120.png";
                link rel="apple-touch-icon" sizes="144x144" href="/static/favicon/apple-icon-144x144.png";
                link rel="apple-touch-icon" sizes="152x152" href="/static/favicon/apple-icon-152x152.png";
                link rel="apple-touch-icon" sizes="180x180" href="/static/favicon/apple-icon-180x180.png";
                link rel="icon" type="image/png" sizes="192x192" href="/static/favicon/android-icon-192x192.png";
                link rel="icon" type="image/png" sizes="32x32" href="/static/favicon/favicon-32x32.png";
                link rel="icon" type="image/png" sizes="32x32" href="/static/favicon/favicon-32x32.png";
                link rel="icon" type="image/png" sizes="96x96" href="/static/favicon/favicon-96x96.png";
                link rel="icon" type="image/png" sizes="16x16" href="/static/favicon/favicon-16x16.png";
                meta name="msapplication-TileColor" content="#ffffff";
                meta name="msapplication-TileImage" content="/static/favicon/ms-icon-144x144.png";
                meta name="theme-color" content="#ffffff";
                link href="https://mi.within.website/api/webmention/accept" rel="webmention";
                @if let Some(styles) = styles {
                    style {
                        (PreEscaped(styles))
                    }
                }
            }
            body.snow.hack.gruvbox-dark {
                .container {
                    header {
                        span.logo {}
                        nav {
                            a href="/" { "Xe" }
                            " - "
                            a href="/blog" { "Blog" }
                            " - "
                            a href="/contact" { "Contact" }
                            " - "
                            a href="/resume" { "Resume" }
                            " - "
                            a href="/talks" { "Talks" }
                            " - "
                            a href="/signalboost" { "Signal Boost" }
                            " - "
                            a href="/vods" { "VODs" }
                            " | "
                            a target="_blank" rel="noopener noreferrer" href="https://graphviz.christine.website" { "Graphviz" }
                            " - "
                            a target="_blank" rel="noopener noreferrer" href="https://when-then-zen.christine.website/" { "When Then Zen" }
                        }
                    }

                    br;
                    br;

                    .snowframe {
                        (content)
                    }
                    hr;
                    footer {
                        blockquote {
                            "Copyright 2012-"
                            (now.year())
                            " Xe Iaso (Christine Dodrill). Any and all opinions listed here are my own and not representative of my employers; future, past and present."
                        }
                        p {
                            "Like what you see? Donate on "
                            a href="https://www.patreon.com/cadey" { "Patreon" }
                            " like "
                            a href="/patrons" { "these awesome people" }
                            "!"
                        }
                        p {
                            "Looking for someone for your team? Take a look "
                            a href="/signalboost" { "here" }
                            "."
                        }
                        p {
                            "See my salary transparency data "
                            a href="/salary-transparency" {"here"}
                            "."
                        }
                        p {
                            "Served by "
                            (env!("out"))
                            "/bin/xesite, see "
                            a href="https://github.com/Xe/site" { "source code here" }
                            "."
                        }
                    }
                    script src="/static/js/installsw.js" defer {}
                }
            }
        }
    }
}

pub fn post_index(posts: &Vec<Post>, title: &str, show_extra: bool) -> Markup {
    let today = Utc::now().date_naive();
    base(
        Some(title),
        None,
        html! {
            h1 { (title) }
            @if show_extra {
                p {
                    "If you have a compatible reader, be sure to check out my "
                    a href="/blog.rss" { "RSS feed" }
                    " for automatic updates. Also check out the "
                    a href="/blog.json" { "JSONFeed" }
                    "."
                }
                p {
                    "For a breakdown by post series, see "
                    a href="/blog/series" { "here" }
                    "."
                }
            }
            p {
                ul {
                    @for post in posts.iter().filter(|p| today.num_days_from_ce() >= p.date.num_days_from_ce()) {
                        li {
                            (post.detri())
                            " - "
                                a href={ @if post.front_matter.redirect_to.as_ref().is_some() {(post.front_matter.redirect_to.as_ref().unwrap())} @else {"/" (post.link)}} { (post.front_matter.title) }
                        }
                    }
                }
            }
        },
    )
}

pub fn gallery_index(posts: &Vec<Post>) -> Markup {
    base(
        Some("Gallery"),
        None,
        html! {
            h1 {"Gallery"}

            p {"Here are links to a lot of the art I have done in the last few years."}

            .grid {
                @for post in posts {
                    .card.cell."-4of12".blogpost-card {
                        header."card-header" {
                            (post.front_matter.title)
                        }
                        .card-content {
                            center {
                                p {
                                    "Posted on "
                                    (post.detri())
                                    br;
                                    a href={"/" (post.link)} {
                                        img src=(post.front_matter.thumb.as_ref().unwrap());
                                    }
                                }
                            }
                        }
                    }
                }
            }
        },
    )
}

pub fn contact(links: &Vec<Link>) -> Markup {
    base(
        Some("Contact Information"),
        None,
        html! {
            h1 {"Contact Information"}

            .grid {
                .cell."-6of12" {
                    h3 {"Email"}
                    p {"me@xeiaso.net"}

                    h3 {"Social Media"}
                    ul {
                        @for link in links {
                            li {(link)}
                        }
                    }
                }
                .cell."-6of12" {
                    h3 {"Other Information"}
                    h4 {"Discord"}
                    p {
                        code {"Cadey~#1337"}
                        " Please note that Discord will automatically reject friend requests if you are not in a mutual server with me. I don't have control over this behavior."
                    }
                }
            }
        },
    )
}

pub fn characters(characters: &Vec<Character>) -> Markup {
    base(
        Some("Characters"),
        None,
        html! {
            h1 {"Characters"}
            p{
                "When I am writing articles on this blog, sometimes I will use "
                a href="https://en.wikipedia.org/wiki/Socratic_method" {"the Socratic method"}
                " to help illustrate my point. These characters are written off of a set of tropes to help give them a place in the discussions. The characters are just that, characters. Their dialogues are fiction, unless otherwise indicated everything that happens in those dialogues are products of the author's imagination or are used in a fictitious manner. Any resemblance to actual persons (living or dead) is purely coincidental."
            }

            @for character in characters {
                (character)
            }

            h2 {"Other People"}

            p{
                "Some of the characters you see in posts aren't figments of my imagination, but instead the OCs of other people."
            }

            h3 #scoots {"Scoots"}
            p {"My husband. The picture he uses is a screenshot of his VRChat avatar."}
        },
    )
}

pub fn patrons(patrons: &Users) -> Markup {
    base(
        Some("Patrons"),
        None,
        html! {
            h1 {"Patrons"}

            p {
                "These awesome people donate to me on "
                a href="https://www.patreon.com/cadey" {"Patreon"}
                ". If you would like to show up in this list, please donate to me on Patreon. This is refreshed every time the site is deployed."
            }

            .grid {
                @for patron in patrons {
                    .cell."-3of12" {
                        center {
                            p {(patron.attributes.full_name)}
                            img src=(patron.attributes.thumb_url) loading="lazy";
                        }
                    }
                }
            }
        },
    )
}

pub fn signalboost(people: &Vec<Person>) -> Markup {
    base(
        Some("Signal Boosts"),
        None,
        html! {
            h1 {"Signal Boosts"}

            p {"These awesome people are currently looking for a job. If you are looking for anyone with these skills, please feel free to reach out to them."}

            p {
                "To add yourself to this list, fork "
                a href="https://github.com/Xe/site" {"this website's source code"}
                " and send a pull request with edits to "
                code {"/dhall/signalboost.dhall"}
                "."
            }

            p {"With COVID-19 raging across the world, these people are in need of a job now more than ever."}

            h2 {"People"}

            .grid.signalboost {
                @for person in people {
                    .cell."-4of12".content {
                        big {(person.name)}

                        p {
                            @for tag in &person.tags {(tag) " "}
                        }

                        p {
                            @for link in &person.links {(link) " "}
                        }
                    }
                }
            }
        },
    )
}

pub fn error(why: impl Render) -> Markup {
    base(
        Some("Error"),
        None,
        html! {
            h1 {"Error"}

            pre {
                (why)
            }

            p {
                "You could try to "
                a href="/" {"go home"}
                " or "
                a href="https://github.com/Xe/site/issues/new" {"report this issue"}
                " so it can be fixed."
            }
        },
    )
}

pub fn not_found(path: impl Render) -> Markup {
    base(
        Some("Not found"),
        None,
        html! {
            h1 {"Not found"}
            p {
                "The path at "
                code {(path)}
                " could not be found. If you expected this path to exist, please "
                a href="https://github.com/Xe/site/issues/new" {"report this issue"}
                " so it can be fixed."
            }
        },
    )
}

pub fn gitea(pkg_name: &str, git_repo: &str, branch: &str) -> Markup {
    html! {
        (DOCTYPE)
        html {
            head {
                meta http-equiv="Content-Type" content="text/html; charset=utf-8";
                meta name="go-import" content={(pkg_name)" git " (git_repo)};
                meta name="go-source" content={(format!("{pkg_name} {git_repo} {git_repo}/src/{branch}{{/dir}} {git_repo}/src/{branch}{{/dir}}/{{file}}#L{{line}}"))};
                meta http-equiv="refresh" content={(format!("0; url=https://pkg.go.dev/{pkg_name}"))};
            }
            body {
                p {
                    "Please see"
                    a href={"https://pkg.go.dev/" (pkg_name)} {"here"}
                    " for documentation on this package."
                }
            }
        }
    }
}

pub fn resume() -> Markup {
    base(
        Some("Resume"),
        None,
        html! {
            h1 {"Resume"}

            p {"This resume is automatically generated when the website gets deployed."}

            iframe src="/static/resume/resume.pdf" width="100%" height="900px" {}

            hr;

            a href="/static/resume/resume.pdf" { "PDF version" }
        },
    )
}

fn schema_person(a: &Author) -> Markup {
    let data = PreEscaped(serde_json::to_string(&a).unwrap());

    html! {
        script type="application/ld+json" { (data) }
    }
}

pub fn index(xe: &Author, projects: &Vec<Link>) -> Markup {
    base(
        None,
        None,
        html! {
            link rel="authorization_endpoint" href="https://idp.christine.website/auth";
            link rel="canonical" href="https://xeiaso.net/";
            meta name="google-site-verification" content="rzs9eBEquMYr9Phrg0Xm0mIwFjDBcbdgJ3jF6Disy-k";
            (schema_person(&xe))

            meta name="twitter:card" content="summary";
            meta name="twitter:site" content="@theprincessxena";
            meta name="twitter:title" content=(xe.name);
            meta name="twitter:description" content=(xe.job_title);
            meta property="og:type" content="website";
            meta property="og:title" content=(xe.name);
            meta property="og:site_name" content=(xe.job_title);
            meta name="description" content=(xe.job_title);
            meta name="author" content=(xe.name);

            .grid {
                .cell."-3of12".content {
                    img src="/static/img/avatar.png" alt="My Avatar";
                    br;
                    a href="/contact" class="justify-content-center" { "Contact me" }
                }
                .cell."-9of12".content {
                    h1 {(xe.name)}
                    h4 {(xe.job_title)}
                    h5 { "Skills" }
                    ul {
                        li { "Go, Lua, Haskell, C, Rust and other languages" }
                        li { "Docker (deployment, development & more)" }
                        li { "Mashups of data" }
                        li { "kastermakfa" }
                    }

                    h5 { "Highlighted Projects" }
                    ul {
                        @for project in projects {
                            li {(project)}
                        }
                    }

                    h5 { "Quick Links" }
                    ul {
                        li {a href="https://github.com/Xe" rel="me" {"GitHub"}}
                        li {a href="https://twitter.com/theprincessxena" rel="me" {"Twitter"}}
                        li {a href="https://pony.social/@cadey" rel="me" {"Fediverse"}}
                        li {a href="https://www.patreon.com/cadey" rel="me" {"Patreon"}}
                    }

                    p {
                        "Looking for someone for your team? Check "
                        a href="/signalboost" { "here" }
                        "."
                    }
                }
            }
        },
    )
}

pub fn blog_series(series: &Vec<SeriesDescription>) -> Markup {
    base(
        Some("Blogposts by series"),
        None,
        html! {
            h1 { "Blogposts by series" }
            p {
                "Some posts of mine are intended to be read in order. This is a list of all the series I have written along with a little description of what it's about."
            }
            p {
                ul {
                    @for set in series {
                        li {(set)}
                    }
                }
            }
        },
    )
}

pub fn series_view(name: &str, desc: &str, posts: &Vec<Post>) -> Markup {
    base(
        Some(&format!("{name} posts")),
        None,
        html! {
            h1 {"Series: " (name)}

            p {(desc)}

            ul {
                @for post in posts {
                    li {
                        (post.detri())
                        " - "
                        a href={"/" (post.link)} {(post.front_matter.title)}
                    }
                }
            }
        },
    )
}

pub fn feeds() -> Markup {
    base(
        Some("My Feeds"),
        None,
        html! {
            h1 { "My Feeds" }

            ul {
                li {
                    "Blog: "
                    a href="/blog.atom" { "Atom" }
                    " - "
                    a href="/blog.rss" { "RSS" }
                    " - "
                    a href="/blog.json" { "JSONFeed" }
                }
                li {
                    "Mastodon: "
                    a href="https://pony.social/users/cadey.rss" { "RSS" }
                }
            }
        },
    )
}

pub fn salary_transparency(jobs: &Vec<Job>) -> Markup {
    base(
        Some("Salary Transparency"),
        None,
        html! {
            h1 {"Salary Transparency"}

            p {
                "This page lists my salary for every job I've had in tech. I have had this data open to the public "
                a href="https://xeiaso.net/blog/my-career-in-dates-titles-salaries-2019-03-14" {"for years"}
                ", but I feel this should be more prominently displayed on my website. Other people have copied my approach of having a list of every salary they have ever been paid on their websites, and I would like to set the example by making it prominent on my website."
            }
            p {
                "As someone who has seen pay discrimination work in action first-hand, data is one of the ways that we can end this pointless hiding of information that leads to people being uninformed and hurt by their lack of knowledge. By laying my hand out in the open like this, I hope to ensure that people are better informed about how much money they "
                em {"can"}
                " make, so that they can be paid equally for equal work."
            }

            p {
                "Please keep in mind that this table doesn't tell the complete story. If you feel like judging me about any entry in this table, please do not do it around me."
            }

            h2 {"Salary Data"}

            p {
                "To get this data, I have scoured over past emails, contracts and everything so that I can be sure that this information is as accurate as possible. The data on this page intentionally omits employer names. Some information may also be omitted if relevant non-disclosure agreements or similar prohibit it."
            }

            (salary_history(jobs))

            p {
                "I typically update this page once any of the following things happens:"
            }

            ul {
                li {"I quit a job."}
                li {"I get a raise/title change at the same company."}
                li {"I get terminated from a job."}
                li {"I get converted from a contracter to a full-time employee."}
                li {"Other unspecified extranormal events happen."}
            }

            p {
                "Please consider publishing your salary data like this as well. By open, voluntary transparency we can help to end stigmas around discussing pay and help ensure that the next generations of people in tech are treated fairly. Stigmas thrive in darkness but die in the light of day. You can help end the stigma by playing your cards out in the open like this."
            }
        },
    )
}

fn salary_history(jobs: &Vec<Job>) -> Markup {
    html! {
        table.salary_history {
            tr {
                th { "Title" }
                th { "Start Date" }
                th { "End Date" }
                th { "Days Worked" }
                th { "Salary" }
                th { "How I Left" }
            }
            @for job in jobs {
                (job)
            }
        }
    }
}

pub fn pronoun_page(pronouns: &Vec<PronounSet>) -> Markup {
    base(
        Some("Pronouns"),
        None,
        html! {
            h1 {"Pronouns"}
            p {"This page lists the pronouns you should use for me. Please try to use one of these sets:"}
            .grid {
                @for ps in pronouns {
                    .card.cell."-4of12" {
                        (ps)
                    }
                }
            }

            (xesite_templates::conv("Mara".to_string(), "happy".to_string(), html!{
                "You can access this data with "
                a href="/api/pronouns" {"an API call"}
                " too!"
            }))
        },
    )
}
