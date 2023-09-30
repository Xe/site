export const title = "Conference Talks";
export const layout = "base.njk";

export default ({ search }) => {
  const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>

      <div class="bg-bg-1 dark:bg-bg-1 rounded-xl m-2 px-2 py-1 shadow-md max-w-xl">
        <div className=" my-4" id="search"></div>
      </div>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("layout=talk.njk", "order date=desc").map((post) => (
          <li>
            {post.data.date.toLocaleDateString("en-US", dateOptions)} -{" "}
            <a href={post.data.url}>{post.data.title}</a>
          </li>
        ))}
      </ul>
    </>
  );
};
