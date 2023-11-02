export const title = "Blog Articles";
export const layout = "base.njk";

export default ({ search }) => {
  const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

  return (
    <>
      <h1 className="text-3xl mb-4">Blog Articles</h1>
      <p className="mb-4">
        If you have a compatible reader, be sure to check out my{" "}
        <a href="/blog.rss">RSS feed</a>{" "}
        for automatic updates. Also check out the{" "}
        <a href="/blog.json">JSONFeed</a>.
      </p>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("type=blog", "order date=desc").map((post) => {
          const url = post.data.redirect_to ? post.data.redirect_to : post.data.url;
          return (
          <li>
            <span className="font-mono">{post.data.date.toLocaleDateString("en-US", dateOptions)}</span> -{" "}
            <a href={url}>{post.data.title}</a>
          </li>
        );
        })}
      </ul>
    </>
  );
};
