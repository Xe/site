export const title = "Stream VODs";
export const layout = "base.njk";

export default ({ search }, { date }) => {
  return (
    <>
      <h1 className="text-3xl mb-4">{title}</h1>

      <p class="my-4">
        I'm a VTuber and I stream every other weekend on{" "}
        <a href="https://twitch.tv/princessxen">Twitch</a>{" "}
        about technology, the weird art of programming, and sometimes video
        games. This page will contain copies of my stream recordings/VODs so
        that you can watch your favorite stream again. All VOD pages support
        picture-in-picture mode so that you can have the recordings open in the
        background while you do something else.
      </p>

      <p class="my-4">
        Please note that to save on filesize, all videos are rendered at 720p
        and optimized for viewing at that resolution or on most mobile phone
        screens. If you run into video quality issues, please contact me as I am
        still trying to find the correct balance between video quality and
        filesize. These videos have been tested and known to work on most of the
        browser and OS combinations that visit this site.
      </p>

      <ul class="list-disc ml-4 mb-4">
        {search.pages("layout=vod.njk", "order date=desc").map((post) => (
          <li>
            <time datetime={date(post.date)}>{date(post.date, "DATE_US")}</time> -{" "}
            <a href={post.url}>
              {post.title}
            </a>
          </li>
        ))}
      </ul>
    </>
  );
};
