export const title = "Blog Articles";
export const layout = "base.njk";

export default ({ search }) => {
  const dateOptions = { year: "numeric", month: "numeric", day: "numeric" };

  return (
    <>
      <h1 className="text-3xl mb-4">Blog Articles</h1>
      <p className="mb-4">
        If you have a compatible reader, be sure to check out my{" "}
        <a href="/blog.rss">RSS feed</a>{" "}
        for automatic updates. Also check out the{" "}
        <a href="/blog.json">JSONFeed</a>.
      </p>

      <div class="bg-white-800 rounded-xl m-2 px-2 py-1 shadow-md max-w-xl">
        <div className=" my-4" id="search"></div>
      </div>

      <p className="mb-4">
        For a breakdown by post series, see <a href="/blog/series">here</a>.
      </p>
      <ul class="list-disc ml-4 mb-4">
        {search.pages("type=blog", "order date=desc").map((post) => (
          <li>
            {post.data.date.toLocaleDateString("en-US", dateOptions)} -{" "}
            <a href={post.data.url}>{post.data.title}</a>
          </li>
        ))}
      </ul>
    </>
  );
};
