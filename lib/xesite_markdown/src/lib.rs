use color_eyre::eyre::{Result, WrapErr};
use comrak::nodes::{Ast, AstNode, LineColumn, NodeValue};
use comrak::plugins::syntect::SyntectAdapter;
use comrak::{
    format_html_with_plugins, markdown_to_html_with_plugins, parse_document, Arena, ComrakOptions,
    ComrakPlugins,
};
use lazy_static::lazy_static;
use lol_html::{element, html_content::ContentType, rewrite_str, RewriteStrSettings};
use maud::PreEscaped;
use sha2::{Digest, Sha256};
use std::{cell::RefCell, fmt::Write};
use url::Url;
use xesite_types::mastodon::{Toot, User};

pub fn hash_string(inp: String) -> String {
    let mut h = Sha256::new();
    h.update(&inp.as_bytes());
    hex::encode(h.finalize())
}

lazy_static! {
    static ref SYNTECT_ADAPTER: SyntectAdapter = SyntectAdapter::new("base16-mocha.dark");
}

#[derive(thiserror::Error, Debug, Clone)]
pub enum Error {
    #[error("missing element attribute {0}")]
    MissingElementAttribute(String),
}

pub fn render(inp: &str) -> Result<String> {
    let mut options = ComrakOptions::default();

    options.extension.autolink = true;
    options.extension.table = true;
    options.extension.description_lists = true;
    options.extension.superscript = true;
    options.extension.strikethrough = true;
    options.extension.footnotes = true;

    options.render.unsafe_ = true;

    let arena = Arena::new();
    let root = parse_document(&arena, inp, &options);

    let mut plugins = ComrakPlugins::default();
    plugins.render.codefence_syntax_highlighter = Some(&*SYNTECT_ADAPTER);

    iter_nodes(root, &|node| {
        let mut data = node.data.borrow_mut();
        match &mut data.value {
            &mut NodeValue::Link(ref mut link) => {
                let base = Url::parse("https://xeiaso.net/")?;
                let u = base.join(&link.url.clone())?;
                if u.scheme() != "conversation" {
                    return Ok(());
                }
                let parent = node.parent().unwrap();
                node.detach();
                let mut message = vec![];
                for child in node.children() {
                    format_html_with_plugins(child, &options, &mut message, &plugins)?;
                }
                let message = std::str::from_utf8(&message)?;
                let mut message = markdown_to_html_with_plugins(message, &options, &plugins);
                crop_letters(&mut message, 3);
                message.drain((message.len() - 5)..);
                let mood = without_first(u.path());
                let name = u.host_str().unwrap_or("Mara");

                let mut html = String::new();
                write!(
                    html,
                    "{}",
                    xesite_templates::conv(
                        name.to_string(),
                        mood.to_string(),
                        PreEscaped(message.trim().into())
                    )
                    .0
                )?;

                let new_node = arena.alloc(AstNode::new(RefCell::new(Ast::new(
                    NodeValue::HtmlInline(html),
                    LineColumn { line: 0, column: 0 },
                ))));
                parent.append(new_node);

                Ok(())
            }
            _ => Ok(()),
        }
    })?;

    let mut html = vec![];
    format_html_with_plugins(root, &options, &mut html, &plugins).unwrap();

    let html = String::from_utf8(html).wrap_err("post is somehow invalid UTF-8")?;

    let html = rewrite_str(
        &html,
        RewriteStrSettings {
            element_content_handlers: vec![
                element!("xeblog-conv", |el| {
                    let name = el
                        .get_attribute("name")
                        .ok_or(Error::MissingElementAttribute("name".to_string()))?;
                    let name_lower = name.clone().to_lowercase();
                    let mood = el
                        .get_attribute("mood")
                        .ok_or(Error::MissingElementAttribute("mood".to_string()))?;
                    let name = name.replace("_", " ");

                    let (size, class) = el
                        .get_attribute("standalone")
                        .map_or((64, "conversation-smol"), |_| {
                            (128, "conversation-standalone")
                        });

                    el.before(
                        &format!(
                            r#"
<div class="conversation">
    <div class="{class}">
        <img src="https://cdn.xeiaso.net/sticker/{name_lower}/{mood}/{size}" alt="{name} is {mood}">
    </div>
    <div class="conversation-chat">&lt;<a href="/characters#{name_lower}"><b>{name}</b></a>&gt; "#
                        ),
                        ContentType::Html,
                    );
                    el.after("</div></div>", ContentType::Html);

                    el.remove_and_keep_content();
                    Ok(())
                }),
                element!("xeblog-picture", |el| {
                    let path = el
                        .get_attribute("path")
                        .expect("wanted xeblog-picture to contain path");
                    el.replace(&xesite_templates::picture(path).0, ContentType::Html);
                    Ok(())
                }),
                element!("xeblog-hero", |el| {
                    let file = el
                        .get_attribute("file")
                        .ok_or(Error::MissingElementAttribute("file".to_string()))?;
                    el.replace(
                        &xesite_templates::hero(
                            file,
                            el.get_attribute("prompt"),
                            el.get_attribute("ai"),
                        )
                        .0,
                        ContentType::Html,
                    );
                    Ok(())
                }),
                element!("xeblog-sticker", |el| {
                    let name = el
                        .get_attribute("name")
                        .ok_or(Error::MissingElementAttribute("name".to_string()))?;
                    let mood = el
                        .get_attribute("mood")
                        .ok_or(Error::MissingElementAttribute("mood".to_string()))?;
                    el.replace(&xesite_templates::sticker(name, mood).0, ContentType::Html);

                    Ok(())
                }),
                element!("xeblog-slide", |el| {
                    let name = el
                        .get_attribute("name")
                        .ok_or(Error::MissingElementAttribute("name".to_string()))?;
                    let essential = el.get_attribute("essential").is_some();
                    el.replace(
                        &xesite_templates::slide(name, essential).0,
                        ContentType::Html,
                    );

                    Ok(())
                }),
                element!("xeblog-talk-warning", |el| {
                    el.replace(&xesite_templates::talk_warning().0, ContentType::Html);
                    Ok(())
                }),
                element!("xeblog-video", |el| {
                    let path = el
                        .get_attribute("path")
                        .ok_or(Error::MissingElementAttribute("path".to_string()))?;

                    el.replace(&xesite_templates::video(path).0, ContentType::Html);
                    Ok(())
                }),
                #[cfg(not(target_arch = "wasm32"))]
                element!("xeblog-toot", |el| {
                    use serde_json::from_reader;
                    use std::fs;

                    let mut toot_url = el
                        .get_attribute("url")
                        .ok_or(Error::MissingElementAttribute("url".to_string()))?;

                    if !toot_url.ends_with(".json") {
                        toot_url = format!("{toot_url}.json");
                    }

                    let toot_fname = format!("./data/toots/{}.json", hash_string(toot_url.clone()));
                    tracing::debug!("opening {toot_fname}");
                    let mut fin = fs::File::open(&toot_fname).context(toot_url)?;
                    let t: Toot = from_reader(&mut fin)?;

                    let user_fname = format!(
                        "./data/users/{}.json",
                        hash_string(format!("{}.json", t.attributed_to.clone()))
                    );
                    tracing::debug!("opening {user_fname}");
                    let mut fin = fs::File::open(&user_fname).context(t.attributed_to.clone())?;

                    let u: User = from_reader(&mut fin)?;

                    el.replace(&xesite_templates::toot_embed(u, t).0, ContentType::Html);
                    Ok(())
                }),
            ],
            ..RewriteStrSettings::default()
        },
    )?;

    Ok(html)
}

fn iter_nodes<'a, F>(node: &'a AstNode<'a>, f: &F) -> Result<()>
where
    F: Fn(&'a AstNode<'a>) -> Result<()>,
{
    f(node)?;
    for c in node.children() {
        iter_nodes(c, f)?;
    }
    Ok(())
}

fn without_first(string: &str) -> &str {
    string
        .char_indices()
        .nth(1)
        .and_then(|(i, _)| string.get(i..))
        .unwrap_or("")
}

fn crop_letters(s: &mut String, pos: usize) {
    match s.char_indices().nth(pos) {
        Some((pos, _)) => {
            s.drain(..pos);
        }
        None => {
            s.clear();
        }
    }
}
