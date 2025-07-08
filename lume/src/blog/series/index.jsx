export const title = "Post Series";
export const layout = "base.njk";
export const date = "2012-01-01";

export default ({ seriesDescriptions }) => {
  const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>
      <p className="mb-4">
        Some posts of mine are intended to be read in order. This is a list of all the series I have written along with a little description of what it's about.
      </p>

      <ul class="list-disc ml-4 mb-4">
        {seriesDescriptions.map((v) => (
          <li>
            <a href={`/blog/series/${v.name}`}>{v.name}</a>: {v.details}
          </li>
        ))}
      </ul>
    </>
  );
};
