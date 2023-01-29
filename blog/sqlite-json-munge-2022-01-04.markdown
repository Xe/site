---
title: Bashing JSON into Shape with SQLite
date: 2022-01-04
series: howto
tags: 
 - sqlite
 - json
---

It is clear that most of the world has decided that they want to use JSON for
their public-facing API endpoints. However, most of the time you will need to
deal with storage engines that don't deal with JSON very well. This can be
confusing to deal with because you need to fit a square peg into a round hole.

However, [SQLite](https://www.sqlite.org) added [JSON
functions](https://www.sqlite.org/json1.html) to allow you to munge and modify
JSON data in whatever creative ways you want. You can use these and SQLite
[triggers](https://www.sqlite.org/lang_createtrigger.html) in order to
automatically massage JSON into whatever kind of tables you want. Throw in
upserts and you'll be able to make things even more automated. This support
was added in SQLite 3.9.0 (released in 2015), so assuming Debian didn't disable
it for no good reason, you should be able to use it today.

For this example, we're going to be querying publicly available
[JSONFeed](https://www.jsonfeed.org/) endpoints and turning that into SQL
tables. Let's start with a table schema that looks like this:

```sql
CREATE TABLE IF NOT EXISTS jsonfeed_raw
  ( feed_url     TEXT  PRIMARY KEY
  , scrape_date  TEXT  NOT NULL  DEFAULT (DATE('now'))
  , raw          TEXT  NOT NULL
  );
```

[The scrape date is essentially the date that the JSONFeed row was inserted into
the database. This can be useful when writing other parts of the stack to
automatically query feeds for changes. This is left as an exercise to the
reader.](conversation://Mara/hacker)

You can then insert things into the SQLite database using Python's `sqlite3`
module:

```python
#!/usr/bin/env nix-shell
#! nix-shell -p python39 --run python

import sqlite3
import urllib.request

con = sqlite3.connect("data.db")

def get_feed(feed_url):
    req = urllib.request.Request(feed_url, headers={"User-Agent": "Xe/feedfetch"})
    with urllib.request.urlopen(req) as response:
        cur = con.cursor()
        body = response.read()
        cur.execute("""
           INSERT INTO jsonfeed_raw
             (feed_url, raw)
           VALUES
             (?, json(?))
        """, (feed_url, body))
        con.commit()
        print("got feed %s" % (feed_url))

get_feed("https://xeiaso.net/blog.json")
```

So now let's play with the data! Let's load the database schema in with the
`sqlite3` command:

```console
$ sqlite3 data.db < schema.sql
```

[The less-than symbol there is a redirect, it loads the data from `schema.sql`
as standard input to the `sqlite` command. See <a
href="https://xeiaso.net/blog/fun-with-redirection-2021-09-22">here</a>
for more information on redirections.](conversation://Mara/hacker)

Then run that python script to populate the database:

```console
$ python ./jsonfeedfetch.py
got feed https://xeiaso.net/blog.json
```

Then open up the SQLite command line:

```console
$ sqlite3 data.db
SQLite version 3.36.0 2021-06-18 18:36:39
Enter ".help" for usage hints.
sqlite>
```

And now we can play with a few of the JSON functions. First let's show off
[`json_extract`](https://www.sqlite.org/json1.html#the_json_extract_function).
This lets you pull a value out of a JSON object. For example, let's get the feed
title out of my website's JSONFeed:

```console
sqlite> select json_extract(raw, '$.title') from jsonfeed_raw;
Xe's Blog
```

We can use this function to help us create a table that stores the metadata we
care about from a JSONFeed, such as this:

```sql
CREATE TABLE IF NOT EXISTS jsonfeed_metadata
  ( feed_url      TEXT  PRIMARY KEY
  , title         TEXT  NOT NULL
  , description   TEXT
  , home_page_url TEXT
  , updated_at    TEXT  NOT NULL  DEFAULT (DATE('now'))
  );
```

[If you ask my coworkers, they can confirm that I actually do real life
unironcally write SQL like that.](conversation://Cadey/coffee)

Then we can populate that table with a query like this:

```sql
INSERT INTO jsonfeed_metadata
            ( feed_url
            , title
            , description
            , home_page_url
            , updated_at
            )
SELECT jsonfeed_raw.feed_url AS feed_url
     , json_extract(jsonfeed_raw.raw, '$.title') AS title
     , json_extract(jsonfeed_raw.raw, '$.description') AS description
     , json_extract(jsonfeed_raw.raw, '$.home_page_url') AS home_page_url
     , DATE('now') AS updated_at
FROM jsonfeed_raw;
```

[The `AS` keyword lets you bind values in a `SELECT` statement to names for use
elsewhere in the query. I don't know if it's _strictly_ needed, however it makes
the names line up and SQLite doesn't complain about it, so it's probably
fine.](conversation://Mara/hacker)

Now this is workable, however you know what's easier than writing statements in
the SQLite console like that? Not having to! SQLite triggers allow us to run
database statements automatically when certain conditions happen. The main
condition we want to care about right now is when we insert new data. We can
turn that statement into an after-insert trigger like this:

```sql
CREATE TRIGGER IF NOT EXISTS jsonfeed_raw_ins
  AFTER INSERT ON jsonfeed_raw
  BEGIN
    INSERT INTO jsonfeed_metadata
                ( feed_url
                , title
                , description
                , home_page_url
                )
    VALUES ( NEW.feed_url
           , json_extract(NEW.raw, '$.title')
           , json_extract(NEW.raw, '$.description')
           , json_extract(NEW.raw, '$.home_page_url')
           );
  END;
```

Then we can run a few commands to nuke all the database state:

```console
sqlite3> DELETE FROM jsonfeed_metadata;
sqlite3> DELETE FROM jsonfeed_raw;
```

And run that python script again, then the data should automatically show up:

```
sqlite3> SELECT * FROM jsonfeed_metadata;
https://xeiaso.net/blog.json|Xe's Blog|My blog posts and rants about various technology things.|https://xeiaso.net|2022-01-04
```

It's like magic!

However, if you run that python script again without deleting the rows, you will
get a primary key violation. We can fix this by turning the insert into an
[upsert](https://www.sqlite.org/lang_UPSERT.html) with something like this:

```python
cur.execute("""
    INSERT INTO jsonfeed_raw
      (feed_url, raw)
    VALUES
      (?, json(?))
    ON CONFLICT DO
      UPDATE SET raw = json(?)
""", (feed_url, body, body))
```

And also make a complementary update trigger for the `jsonfeed_raw` table:

```sql
CREATE TRIGGER IF NOT EXISTS jsonfeed_raw_upd
  AFTER UPDATE ON jsonfeed_raw
  BEGIN
    INSERT INTO jsonfeed_metadata
                ( feed_url
                , title
                , description
                , home_page_url
                )
    VALUES ( NEW.feed_url
            , json_extract(NEW.raw, '$.title')
            , json_extract(NEW.raw, '$.description')
            , json_extract(NEW.raw, '$.home_page_url')
            )
    ON CONFLICT DO
       UPDATE SET
             title         = json_extract(NEW.raw, '$.title')
           , description   = json_extract(NEW.raw, '$.description')
           , home_page_url = json_extract(NEW.raw, '$.home_page_url')
           ;
```

[You should probably update the original trigger to be an upsert too. You can
follow this trigger as a guide. Be sure to `DROP TRIGGER jsonfeed_raw_upd;`
first though!](conversation://Mara/hacker)

We can also scrape the feed items out too with `json_each`. `json_each` lets you
iterate a JSON array and returns SQLite rows for every value in that array.
Let's take this for example:

```console
sqlite> select * from json_each('["foo", "bar"]');
0|foo|text|foo|1||$[0]|$
1|bar|text|bar|2||$[1]|$
```

The schema for the temporary table that `json_each` (and the related
`json_tree`) uses can be found [here](https://www.sqlite.org/json1.html#jeach).
You can also grab things out of a list in an object with the second argument to
`json_each`, so you can do things like this:

```console
sqlite> select * from json_each('{"spam": ["foo", "bar"]}', '$.spam');
0|foo|text|foo|3||$.spam[0]|$.spam
1|bar|text|bar|4||$.spam[1]|$.spam
```

Using this, we can make a table for each of the feed items that looks something
like this:

```sql
CREATE TABLE IF NOT EXISTS jsonfeed_posts
  ( url             TEXT  PRIMARY KEY
  , feed_url        TEXT  NOT NULL
  , title           TEXT  NOT NULL
  , date_published  TEXT  NOT NULL
  );
```

And then munge everything out of the data in the database with a query like
this:

```sql
INSERT INTO jsonfeed_posts
            ( url
            , feed_url
            , title
            , date_published
            )
SELECT
  json_extract(json_each.value, '$.url') AS url
, jsonfeed_raw.feed_url AS feed_url
, json_extract(json_each.value, '$.title') AS title
, json_extract(json_each.value, '$.date_published') AS date_published
FROM
  jsonfeed_raw
, json_each(jsonfeed_raw.raw, '$.items');
```

This will fetch all of the values of the `items` field in every JSONFeed and
then automatically populate them into the `jsonfeed_posts` table. However
turning this into a trigger with the naiive approach will not instantly work.

Let's say we have the trigger form that looks like this:

```sql
CREATE TRIGGER IF NOT EXISTS jsonfeed_raw_upd_posts
  AFTER INSERT ON jsonfeed_raw
  BEGIN
    INSERT INTO jsonfeed_posts
                ( url
                , feed_url
                , title
                , date_published
                )
    SELECT
        json_extract(json_each.value, '$.url') AS url
      , NEW.feed_url AS feed_url
      , json_extract(json_each.value, '$.title') AS title
      , json_extract(json_each.value, '$.date_published') AS date_published
    FROM json_each(NEW.raw, '$.items')
    ON CONFLICT DO
      UPDATE SET title = excluded.title
               , date_published = excluded.date_published
               ;
  END;
```

If you paste this into your SQLite console, you'll get this error:

```
Error: near "DO": syntax error
```

This is actually due to a [parsing ambiguity in
SQLite](https://www.sqlite.org/lang_UPSERT.html). In order to fix this you will
need to add `WHERE TRUE` between the `FROM` and `ON CONFLICT` clauses of the
trigger:

```sql
-- ...
FROM json_each(NEW.raw, '$.items')
WHERE TRUE
ON CONFLICT DO
-- ...
```

[And thus the day is saved by the wheretrue, the hidden apex predator of the
SQLite realm, a fated value that is only non-falsy at night. Weep in terror lest
it add you to its table of victims!](conversation://Numa/delet)

[The correlating insert trigger change is also an exercise for the
reader.](conversation://Mara/hacker)

Now you can add JSONFeeds how you want and all of the data will automatically be
updated. This can probably be vastly simplified further with the use of
[generated columns](https://dgl.cx/2020/06/sqlite-json-support), however this
should work admirably for most needs.

SQLite is able to be a NOSQL database. It's good enough for your needs. If you
want to play with the code I wrote while writing this article, check it out
[here](https://git.io/JSDVR). This post was written live on
[twitch.tv](https://www.twitch.tv/princessxen). Please follow or subscribe to be
kept up to date on when I go live!

The VOD for this post is [here](https://www.twitch.tv/videos/1253083566). The
corresponding YouTube upload is [here](https://youtu.be/zkM_lY65Lcw). It won't
be available immediately after this post goes live, but it will go up in time.

Here is my favorite message from the chat while I was researching this post:

> jbpratt: if you were married to sqlite, i'd be reporting domestic abuse. This
> is awesome
