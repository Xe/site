use maud::{html, Markup, PreEscaped, Render};
use serde::{
    de::{self, Visitor},
    Deserialize, Deserializer, Serialize,
};
use std::fmt;

struct StringVisitor;

impl<'de> Visitor<'de> for StringVisitor {
    type Value = MarkdownString;

    fn visit_borrowed_str<E>(self, value: &'de str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        Ok(MarkdownString(xesite_markdown::render(value).map_err(
            |why| de::Error::invalid_value(de::Unexpected::Other(&format!("{why}")), &self),
        )?))
    }

    fn visit_str<E>(self, value: &str) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        Ok(MarkdownString(xesite_markdown::render(value).map_err(
            |why| de::Error::invalid_value(de::Unexpected::Other(&format!("{why}")), &self),
        )?))
    }

    fn visit_string<E>(self, value: String) -> Result<Self::Value, E>
    where
        E: de::Error,
    {
        Ok(MarkdownString(xesite_markdown::render(&value).map_err(
            |why| de::Error::invalid_value(de::Unexpected::Other(&format!("{why}")), &self),
        )?))
    }

    fn expecting(&self, formatter: &mut fmt::Formatter) -> fmt::Result {
        formatter.write_str("a string with xesite-flavored markdown")
    }
}

#[derive(Serialize, Clone, Default)]
pub struct MarkdownString(String);

impl<'de> Deserialize<'de> for MarkdownString {
    fn deserialize<D>(deserializer: D) -> Result<MarkdownString, D::Error>
    where
        D: Deserializer<'de>,
    {
        deserializer.deserialize_string(StringVisitor)
    }
}

impl Render for MarkdownString {
    fn render(&self) -> Markup {
        html! {
            (PreEscaped(&self.0))
        }
    }
}
