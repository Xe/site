export const title = "Conference Talks";
export const layout = "base.njk";

export default ({ search }, { date }) => {
  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("layout=talk.njk", "order date=desc")
          .filter(post => post.index)
          .map((post) => {
            const url = post.redirect_to ? post.redirect_to : post.url;
            return (
              <li>
                <time datetime={date(post.date)} className="font-mono">{post.date.toISOString().split('T')[0]}</time> -{" "}
                <a href={url}>{post.title}</a>
              </li>
            );
          })}
      </ul>
    </>
  );
};
