export const title = "Conference Talks";
export const layout = "base.njk";

export default ({ search }, { date }) => {
  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("layout=talk.njk", "order date=desc").map((post) => (
          <li>
            <time datetime={date(post.data.date)}>{date(post.data.date, "DATE_US")}</time> -{" "}
            <a href={post.data.url}>{post.data.title}</a>
          </li>
        ))}
      </ul>
    </>
  );
};
