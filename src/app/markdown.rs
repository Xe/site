use crate::templates::Html;
use color_eyre::eyre::{Result, WrapErr};
use comrak::nodes::{Ast, AstNode, NodeValue};
use comrak::{format_html, markdown_to_html, parse_document, Arena, ComrakOptions};
use std::cell::RefCell;
use url::Url;

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

    iter_nodes(root, &|node| {
        let mut data = node.data.borrow_mut();
        match &mut data.value {
            &mut NodeValue::Link(ref mut link) => {
                let base = Url::parse("https://christine.website/")?;
                let u = base.join(std::str::from_utf8(&link.url.clone())?)?;
                if u.scheme() != "conversation" {
                    return Ok(());
                }
                let parent = node.parent().unwrap();
                node.detach();
                let mut message = vec![];
                for child in node.children() {
                    format_html(child, &options, &mut message)?;
                }
                let message = std::str::from_utf8(&message)?;
                let mut message = markdown_to_html(message, &options);
                crop_letters(&mut message, 3);
                message.drain((message.len() - 5)..);
                let mood = without_first(u.path());
                let name = u.host_str().unwrap_or("Mara");

                let mut html = vec![];
                crate::templates::mara(&mut html, mood, name, Html(message.trim().into()))?;

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
    format_html(root, &options, &mut html).unwrap();

    String::from_utf8(html).wrap_err("post is somehow invalid UTF-8")
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
