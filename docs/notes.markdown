# Notes

## Goals

- Have somewhere other than Twitter or Mastodon to host short-form content or
  list articles I "like".
- Authenticate over Tailscale
- Simple API to automate posting with iOS Shortcuts
- Have a JSONFeed for people to subscribe
- Send WebMentions when I reply to things
- Store things in SQLite

## Schema

```sql
CREATE TABLE IF NOT EXISTS notes
  ( id           INTEGER PRIMARY KEY
  , content      TEXT    NOT NULL
  , content_html TEXT    NOT NULL
  , created_at   INTEGER NOT NULL -- Unix epoch timestamp
  , updated_at   INTEGER          -- Unix epoch timestamp
  , deleted_at   INTEGER          -- Unix epoch timestamp
  , reply_to     TEXT
  );
```

```rust
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
```
