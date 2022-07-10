use crate::templates;
use axum::{extract::Path, response::Html, Json};
use chrono::prelude::*;
use maud::{html, Markup, PreEscaped};
use rusqlite::params;
use serde::{Deserialize, Serialize};

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct Note {
    pub id: u64,
    pub content: String,
    pub content_html: String,
    pub created_at: DateTime<Utc>,
    pub updated_at: Option<DateTime<Utc>>,
    pub deleted_at: Option<DateTime<Utc>>,
    pub reply_to: Option<String>,
}

impl Note {
    pub fn to_html(&self) -> Markup {
        html! {
            article."h-entry" {
                time."dt-published" datetime=(self.created_at) {
                    {(self.detrytemci())}
                }
                a href={"/notes/" (self.id)} {
                    "🔗"
                }

                @if let Some(reply_to) = &self.reply_to {
                    p {
                        "In reply to "
                        a href=(reply_to) {(reply_to)}
                        "."
                    }
                }

                div."e-content" {
                    (PreEscaped(self.content_html.clone()))
                }
            }
        }
    }

    pub fn detrytemci(&self) -> String {
        self.created_at.format("M%m %d %Y %M:%H").to_string()
    }
}

impl Into<xe_jsonfeed::Item> for Note {
    fn into(self) -> xe_jsonfeed::Item {
        let url = format!("https://xeiaso.net/note/{}", self.id);
        let mut builder = xe_jsonfeed::Item::builder()
            .content_html(self.content_html)
            .id(url.clone())
            .url(url.clone())
            .date_published(self.created_at.to_rfc3339())
            .author(
                xe_jsonfeed::Author::new()
                    .name("Xe Iaso")
                    .url("https://xeiaso.net")
                    .avatar("https://xeiaso.net/static/img/avatar.png"),
            );

        if let Some(updated_at) = self.updated_at {
            builder = builder.date_modified(updated_at.to_rfc3339());
        }

        builder.build().unwrap()
    }
}

#[instrument(err)]
pub async fn index() -> super::Result {
    let conn = crate::establish_connection()?;

    let mut stmt = conn.prepare("SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes ORDER BY id DESC LIMIT 25")?;
    let notes = stmt
        .query_map(params![], |row| {
            Ok(Note {
                id: row.get(0)?,
                content: row.get(1)?,
                content_html: row.get(2)?,
                created_at: row.get(3)?,
                updated_at: row.get(4)?,
                deleted_at: row.get(5)?,
                reply_to: row.get(6)?,
            })
        })?
        .filter(Result::is_ok)
        .map(Result::unwrap)
        .collect::<Vec<Note>>();

    let mut result: Vec<u8> = vec![];
    templates::notesindex_html(&mut result, notes)?;
    Ok(Html(result))
}

#[instrument(err)]
pub async fn feed() -> super::Result<Json<xe_jsonfeed::Feed>> {
    let conn = crate::establish_connection()?;

    let mut stmt = conn.prepare("SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes ORDER BY id DESC LIMIT 25")?;
    let notes = stmt
        .query_map(params![], |row| {
            Ok(Note {
                id: row.get(0)?,
                content: row.get(1)?,
                content_html: row.get(2)?,
                created_at: row.get(3)?,
                updated_at: row.get(4)?,
                deleted_at: row.get(5)?,
                reply_to: row.get(6)?,
            })
        })?
        .filter(Result::is_ok)
        .map(Result::unwrap)
        .collect::<Vec<Note>>();

    let mut feed = xe_jsonfeed::Feed::builder()
        .author(
            xe_jsonfeed::Author::new()
                .name("Xe Iaso")
                .url("https://xeiaso.net")
                .avatar("https://xeiaso.net/static/img/avatar.png"),
        )
        .description("Short posts that aren't to the same quality level as mainline blogposts")
        .feed_url("https://xeiaso.net/notes.json")
        .title("Xe's Notes");

    for note in notes {
        feed = feed.item(note.into());
    }

    Ok(Json(feed.build()))
}

#[instrument(err)]
pub async fn view(Path(id): Path<u64>) -> super::Result {
    let conn = crate::establish_connection()?;

    let mut stmt = conn.prepare(
        "SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes WHERE id = ?1"
    )?;

    let note = stmt.query_row(params![id], |row| {
        Ok(Note {
            id: row.get(0)?,
            content: row.get(1)?,
            content_html: row.get(2)?,
            created_at: row.get(3)?,
            updated_at: row.get(4)?,
            deleted_at: row.get(5)?,
            reply_to: row.get(6)?,
        })
    })?;

    let mut result: Vec<u8> = vec![];
    templates::notepost_html(&mut result, note)?;
    Ok(Html(result))
}
