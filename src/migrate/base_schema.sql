CREATE TABLE IF NOT EXISTS notes
  ( id           INTEGER PRIMARY KEY
  , content      TEXT    NOT NULL
  , content_html TEXT    NOT NULL
  , created_at   TEXT    NOT NULL     -- Unix epoch timestamp
  , updated_at   TEXT                 -- Unix epoch timestamp
  , deleted_at   TEXT                 -- Unix epoch timestamp
  , reply_to     TEXT
  );

CREATE INDEX IF NOT EXISTS notes_reply_to ON notes(reply_to);
