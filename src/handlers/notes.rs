use std::{net::SocketAddr, sync::Arc};

use crate::{app::State, templates};
use axum::{extract::Path, http::HeaderMap, response::Html, Extension, Json};
use chrono::prelude::*;
use maud::{html, Markup, PreEscaped};
use rusqlite::params;
use serde::{Deserialize, Serialize};

use super::Error;

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

#[derive(Clone, Debug, Serialize, Deserialize)]
pub struct NewNote {
    pub content: String,
    pub reply_to: Option<String>,
}

impl Default for Note {
    fn default() -> Self {
        Self {
            id: 0,
            content: "".into(),
            content_html: "".into(),
            created_at: Utc::now(),
            updated_at: None,
            deleted_at: None,
            reply_to: None,
        }
    }
}

impl Note {
    pub fn to_html(&self) -> Markup {
        html! {
            article."h-entry" {
                a href={"/notes/" (self.id)} {
                    "🔗"
                }
                " "
                time."dt-published" datetime=(self.created_at) {
                    {(self.detrytemci())}
                }
                " "
                @if let Some(_updated_at) = &self.updated_at {
                    "📝 "
                    (self.update_detrytemci().unwrap())
                }

                @if let Some(deleted_at) = &self.deleted_at {
                    p {
                        " ⚠️ This post was deleted at "
                        (deleted_at.to_rfc3339())
                        ". Please do not treat this note as a genuine expression of my views or opinions."
                    }
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
        self.created_at.format("M%m %d %Y %H:%M").to_string()
    }

    pub fn update_detrytemci(&self) -> Option<String> {
        if self.updated_at.is_none() {
            return None;
        }

        Some(
            self.updated_at
                .as_ref()
                .unwrap()
                .format("M%m %d %Y %H:%M")
                .to_string(),
        )
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

#[instrument(err, skip(state))]
pub async fn index(Extension(state): Extension<Arc<State>>) -> super::Result {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 25")?;
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

#[instrument(err, skip(state))]
pub async fn feed(
    Extension(state): Extension<Arc<State>>,
) -> super::Result<Json<xe_jsonfeed::Feed>> {
    let conn = state.pool.get().await?;

    let mut stmt = conn.prepare("SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes WHERE deleted_at IS NULL ORDER BY id DESC LIMIT 25")?;
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

#[instrument(err, skip(state))]
pub async fn view(Extension(state): Extension<Arc<State>>, Path(id): Path<u64>) -> super::Result {
    let conn = state.pool.get().await?;

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

#[instrument(err, skip(state, headers))]
pub async fn delete(
    Extension(state): Extension<Arc<State>>,
    Path(id): Path<u64>,
    headers: HeaderMap,
) -> super::Result<String> {
    let conn = state.pool.get().await?;

    let ip = headers.get("X-Real-Ip").clone();

    if ip.is_none() {
        return Err(Error::Catchall("Cannot get X-Real-Ip header".into()));
    }

    let ip: SocketAddr = (ip.unwrap().to_str()?.to_owned() + ":0").parse()?;
    let whois = ts_localapi::whois(ip).await?;

    if whois.user_profile.login_name != "Xe@github" {
        return Err(Error::Catchall(format!(
            "expected Tailscale user Xe@github, got: {}",
            whois.user_profile.login_name
        )));
    }

    conn.execute(
        "UPDATE notes SET deleted_at=?2 WHERE id=?1",
        params![id, Utc::now().to_rfc3339()],
    )?;

    Ok("deleted".into())
}

#[instrument(err, skip(state, headers))]
pub async fn update(
    Extension(state): Extension<Arc<State>>,
    Path(id): Path<u64>,
    headers: HeaderMap,
    data: Json<NewNote>,
) -> super::Result<Json<Note>> {
    let conn = state.pool.get().await?;

    let ip = headers.get("X-Real-Ip").clone();

    if ip.is_none() {
        return Err(Error::Catchall("Cannot get X-Real-Ip header".into()));
    }

    let ip: SocketAddr = (ip.unwrap().to_str()?.to_owned() + ":0").parse()?;
    let whois = ts_localapi::whois(ip).await?;

    if whois.user_profile.login_name != "Xe@github" {
        return Err(Error::Catchall(format!(
            "expected Tailscale user Xe@github, got: {}",
            whois.user_profile.login_name
        )));
    }

    info!(
        "authenticated as {} from machine {}",
        whois.user_profile.login_name, whois.node.hostinfo.hostname,
    );

    let content_html = crate::app::markdown::render(state.clone().cfg.clone(), &data.content)?;

    let mut stmt = conn.prepare(
        "SELECT id, content, content_html, created_at, updated_at, deleted_at, reply_to FROM notes WHERE id = ?1"
    )?;

    let old_note = stmt.query_row(params![id], |row| {
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

    let mut note = Note {
        content: data.content.clone(),
        content_html,
        created_at: old_note.created_at,
        updated_at: Some(Utc::now()),
        reply_to: old_note.reply_to,
        ..Default::default()
    };

    conn.execute(
        "UPDATE notes SET content=?, content_html=?, created_at=?, updated_at=?, deleted_at=?, reply_to=? where id=?",
        params![
            note.content,
            note.content_html,
            note.created_at,
            note.updated_at,
            note.deleted_at,
            note.reply_to,
            old_note.id,
        ],
    )?;

    note.id = conn.last_insert_rowid() as u64;

    Ok(Json(note))
}

#[instrument(err, skip(state, headers))]
pub async fn create(
    Extension(state): Extension<Arc<State>>,
    headers: HeaderMap,
    data: Json<NewNote>,
) -> super::Result<Json<Note>> {
    let conn = state.pool.get().await?;

    let ip = headers.get("X-Real-Ip").clone();

    if ip.is_none() {
        return Err(Error::Catchall("Cannot get X-Real-Ip header".into()));
    }

    let ip: SocketAddr = (ip.unwrap().to_str()?.to_owned() + ":0").parse()?;
    let whois = ts_localapi::whois(ip).await?;

    if whois.user_profile.login_name != "Xe@github" {
        return Err(Error::Catchall(format!(
            "expected Tailscale user Xe@github, got: {}",
            whois.user_profile.login_name
        )));
    }

    info!(
        "authenticated as {} from machine {}",
        whois.user_profile.login_name, whois.node.hostinfo.hostname,
    );

    let content_html = crate::app::markdown::render(state.clone().cfg.clone(), &data.content)?;

    let mut note = Note {
        content: data.content.clone(),
        content_html,
        reply_to: data.reply_to.clone(),
        ..Default::default()
    };

    conn.execute(
        "INSERT INTO notes(content, content_html, created_at, updated_at, deleted_at, reply_to) VALUES(?, ?, ?, ?, ?, ?)",
        params![
            note.content,
            note.content_html,
            note.created_at,
            note.updated_at,
            note.deleted_at,
            note.reply_to
        ],
    )?;

    note.id = conn.last_insert_rowid() as u64;

    Ok(Json(note))
}
