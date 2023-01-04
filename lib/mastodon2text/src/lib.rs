use lol_html::{element, html_content::ContentType, HtmlRewriter, Settings};
use std::error::Error;

pub fn convert(input: String) -> Result<String, Box<dyn Error>> {
    let mut output = Vec::new();

    let mut rewriter = HtmlRewriter::new(
        Settings {
            element_content_handlers: vec![
                element!("span", |el| {
                    el.remove_and_keep_content();
                    Ok(())
                }),
                element!("p", |el| {
                    el.append(" ", ContentType::Html);
                    el.remove_and_keep_content();
                    Ok(())
                }),
                element!("br", |el| {
                    el.append(" ", ContentType::Html);
                    el.remove_and_keep_content();
                    Ok(())
                }),
                element!("a[href]", |el| {
                    el.remove_and_keep_content();

                    Ok(())
                }),
            ],
            ..Settings::default()
        },
        |c: &[u8]| output.extend_from_slice(c),
    );

    rewriter.write(input.as_bytes())?;
    rewriter.end()?;

    Ok(String::from_utf8_lossy(&output).to_string())
}
