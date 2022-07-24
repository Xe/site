CREATE TABLE IF NOT EXISTS notes
  ( id           INTEGER PRIMARY KEY
  , content      TEXT    NOT NULL
  , content_html TEXT    NOT NULL
  , created_at   TEXT    NOT NULL     -- RFC 3339 timestamp
  , updated_at   TEXT                 -- RFC 3339 timestamp
  , deleted_at   TEXT                 -- RFC 3339 timestamp
  , reply_to     TEXT
  );

CREATE INDEX IF NOT EXISTS notes_reply_to ON notes(reply_to);
