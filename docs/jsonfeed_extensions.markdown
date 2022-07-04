# JSON Feed Extensions

Here is the documentation of all of my JSON Feed extensions. I have created
these JSON Feed extensions in order to give users more metadata about my
articles and talks.

## `_xesite_frontmatter`

This extension is added to [JSON Feed
Items](https://www.jsonfeed.org/version/1.1/#items-a-name-items-a) and gives
readers a copy of the frontmatter data that I annotate my posts with. The
contents of this will vary by post, but will have any of the following fields:

* `about` (required, string) is a link to this documentation. It gives readers
  of the JSON Feed information about what this extension does. This is for
  informational purposes only and can safely be ignored by programs.
* `series` (optional, string) is the optional blogpost series name that this
  item belongs to. When I post multiple posts about the same topic, I will
  usually set the `series` to the same value so that it is more discoverable [on
  my series index page](https://xeiaso.net/blog/series).
* `slides_link` (optional, string) is a link to the PDF containing the slides
  for a given talk. This is always set on talks, but is technically optional
  because not everything I do is a talk.
* `vod` (optional, string) is an object that describes where you can watch the
  Video On Demand (vod) for the writing process of a post. This is an object
  that always contains the fields `twitch` and `youtube`. These will be URLs to
  the videos so that you can watch them on demand.
