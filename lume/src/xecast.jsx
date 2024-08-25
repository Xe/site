export const title = "Xecast";
export const layout = "base.njk";

export default ({ search }, { date }) => {
  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>

      <p className="mb-4">Welcome to Xecast! Here you can find all the episodes of the podcast. Xecast is an experimental podcast where I talk about computing, AI, technology, and other topics that I find interesting.</p>

      <p className="mb-4"><a href="/xecast.rss">Subscribe via RSS 📡</a></p>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("is_xecast=true", "order date=desc")
          .filter(post => post.index)
          .map((post) => (
            <li>
              <time datetime={date(post.date)}>{date(post.date, "DATE_US")}</time> -{" "}
              <a href={post.url}>{post.title}</a>
            </li>
          ))}
      </ul>
    </>
  );
};
