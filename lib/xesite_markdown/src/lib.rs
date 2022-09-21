use color_eyre::eyre::{Result, WrapErr};
use comrak::nodes::{Ast, AstNode, NodeValue};
use comrak::plugins::syntect::SyntectAdapter;
use comrak::{
    format_html_with_plugins, markdown_to_html_with_plugins, parse_document, Arena, ComrakOptions,
    ComrakPlugins,
};
use lazy_static::lazy_static;
use lol_html::{element, html_content::ContentType, rewrite_str, RewriteStrSettings};
use maud::PreEscaped;
use std::{cell::RefCell, io::Write};
use url::Url;

lazy_static! {
    static ref SYNTECT_ADAPTER: SyntectAdapter<'static> = SyntectAdapter::new("base16-mocha.dark");
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
                let u = base.join(std::str::from_utf8(&link.url.clone())?)?;
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

                let mut html = vec![];
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

    let html = rewrite_str(&html, RewriteStrSettings{
        element_content_handlers: vec![
            element!("xeblog-conv", |el| {
                let name = el.get_attribute("name").expect("wanted xeblog-conv to contain name");
                let name_lower = name.clone().to_lowercase();
                let mood = el.get_attribute("mood").expect("wanted xeblog-conv to contain mood");

                el.before(&format!(r#"
<div class="conversation">
    <div class="conversation-picture conversation-smol">
        <picture>
            <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/{name_lower}/{mood}.avif" type="image/avif">
            <source srcset="https://cdn.xeiaso.net/file/christine-static/stickers/{name_lower}/{mood}.webp" type="image/webp">
            <img src="https://cdn.xeiaso.net/file/christine-static/stickers/{name_lower}/{mood}.png" alt="{name} is {mood}">
        </picture>
    </div>
    <div class="conversation-chat">&lt;<b>{name}</b>&gt; "#), ContentType::Html);
                el.after("</div></div>", ContentType::Html);

                el.remove_and_keep_content();
                Ok(())
            }),
            element!("xeblog-picture", |el| {
                let path = el.get_attribute("path").expect("wanted xeblog-picture to contain path");
                el.replace(&xesite_templates::picture(path).0, ContentType::Html);
                Ok(())
            }),
            element!("xeblog-hero", |el| {
                let file = el.get_attribute("file").expect("wanted xeblog-hero to contain file");
                el.replace(&xesite_templates::hero(file, el.get_attribute("prompt"), el.get_attribute("ai")).0, ContentType::Html);
                Ok(())
            }),
            element!("xeblog-sticker", |el| {
                let name = el.get_attribute("name").expect("wanted xeblog-sticker to contain name");
                let mood = el.get_attribute("mood").expect("wanted xeblog-sticker to contain mood");
                el.replace(&xesite_templates::sticker(name, mood).0, ContentType::Html);

                Ok(())
            }),
            element!("xeblog-slide", |el| {
                let name = el.get_attribute("name").expect("wanted xeblog-slide to contain name");
                let essential = el.get_attribute("essential").is_some();
                el.replace(&xesite_templates::slide(name, essential).0, ContentType::Html);

                Ok(())
            }),
            element!("xeblog-talk-warning", |el| {
                el.replace(&xesite_templates::talk_warning().0, ContentType::Html);
                Ok(())
            }),
        ],
        ..RewriteStrSettings::default()
    }).unwrap();

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
