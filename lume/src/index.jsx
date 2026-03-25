export const layout = "base.njk";
export const date = "2012-12-21";

const GitHubIconSmall = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
    <path d="M9 19c-4.3 1.4 -4.3 -2.5 -6 -3m12 5v-3.5c0 -1 .1 -1.4 -.5 -2c2.8 -.3 5.5 -1.4 5.5 -6a4.6 4.6 0 0 0 -1.3 -3.2a4.2 4.2 0 0 0 -.1 -3.2s-1.1 -.3 -3.5 1.3a12.3 12.3 0 0 0 -6.2 0c-2.4 -1.6 -3.5 -1.3 -3.5 -1.3a4.2 4.2 0 0 0 -.1 3.2a4.6 4.6 0 0 0 -1.3 3.2c0 4.6 2.7 5.7 5.5 6c-.6 .6 -.6 1.2 -.5 2v3.5" />
  </svg>
);

const PatreonIconSmall = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
    <path d="M20 8.408c-.003 -2.299 -1.746 -4.182 -3.79 -4.862c-2.54 -.844 -5.888 -.722 -8.312 .453c-2.939 1.425 -3.862 4.545 -3.896 7.656c-.028 2.559 .22 9.297 3.92 9.345c2.75 .036 3.159 -3.603 4.43 -5.356c.906 -1.247 2.071 -1.599 3.506 -1.963c2.465 -.627 4.146 -2.626 4.142 -5.273z" />
  </svg>
);

const StarIconSmall = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
    <path stroke="none" d="M0 0h24v24H0z" fill="none" />
    <path d="M12 17.75l-6.172 3.245l1.179 -6.873l-5 -4.867l6.9 -1l3.086 -6.253l3.086 6.253l6.9 1l-5 4.867l1.179 6.873z" />
  </svg>
);

const SponsorCard = () => (
  <div className="relative my-6 mx-auto max-w-xl overflow-hidden rounded-md border border-solid border-fg-4 shadow-sm dark:border-fgDark-4 bg-bg-2 dark:bg-bgDark-2">
    <div className="h-[2px] w-full dark:hidden" style={{ background: "linear-gradient(90deg, #d65d0e, #b16286)" }} />
    <div className="hidden h-[2px] w-full dark:block" style={{ background: "linear-gradient(90deg, #fe8019, #d3869b)" }} />
    <div className="px-6 py-5 text-center">
      <h3 className="mb-2 text-xl font-semibold font-serif text-fg-0 dark:text-fgDark-0">
        <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className="inline-block align-text-bottom mr-1">
          <path stroke="none" d="M0 0h24v24H0z" fill="none" />
          <path d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>
        Support My Work
      </h3>
      <p className="mb-4 text-sm text-fg-3 dark:text-fgDark-1">
        Help me keep creating technical content and open source software.
      </p>
      <div className="flex flex-wrap items-center justify-center gap-3">
        <a
          href="https://patreon.com/cadey"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-xl border-0 bg-orange-light px-4 py-2 text-sm font-medium text-white no-underline shadow-sm transition-all duration-200 hover:-translate-y-px hover:bg-orange-dark hover:text-white hover:no-underline hover:shadow-md visited:text-white visited:hover:bg-orange-dark visited:hover:text-white dark:bg-orangeDark-light dark:hover:bg-orangeDark-dark"
        >
          <PatreonIconSmall /> Patreon
        </a>
        <a
          href="https://github.com/sponsors/Xe"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-xl border-0 bg-fg-0 px-4 py-2 text-sm font-medium text-white no-underline shadow-sm transition-all duration-200 hover:-translate-y-px hover:bg-fg-1 hover:text-white hover:no-underline hover:shadow-md visited:text-white visited:hover:bg-fg-1 visited:hover:text-white dark:bg-purpleDark-light dark:hover:bg-purpleDark-dark"
        >
          <GitHubIconSmall /> GitHub Sponsors
        </a>
        <a
          href="https://sponsors.xeiaso.net"
          target="_blank"
          rel="noopener noreferrer"
          className="inline-flex items-center gap-2 rounded-xl border-0 bg-purple-light px-4 py-2 text-sm font-medium text-white no-underline shadow-sm transition-all duration-200 hover:-translate-y-px hover:bg-purple-dark hover:text-white hover:no-underline hover:shadow-md visited:text-white visited:hover:bg-purple-dark visited:hover:text-white dark:bg-blueDark-light dark:hover:bg-blueDark-dark"
        >
          <StarIconSmall /> Sponsor Panel
        </a>
      </div>
    </div>
  </div>
);

export default ({ search, resume, notableProjects, contactLinks }, { date }) => {
  const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

  return (
    <>
      <h1 class="text-3xl mb-4">{resume.name}</h1>
      <p class="text-xl mb-4">{resume.tagline} - {resume.location.city}, {resume.location.country}</p>
      <div className="my-4 flex space-x-4 rounded-md border border-solid border-fg-4 p-3 dark:border-fgDark-4 bg-bg-2 dark:bg-bgDark-2 max-w-full min-h-fit">
        <div className="flex max-h-16 shrink-0 items-center justify-center self-center">
          <img
            style="max-height:6rem"
            alt="A pink haired orca character"
            loading="lazy"
            src="/static/img/avatar.png"
          />
        </div>
        <div className="convsnippet min-w-0 self-center">
          <p className="">
            I'm Xe Iaso, I'm a technical educator, <a href="/talks/">conference speaker</a>, <a href="https://twitch.tv/princessxen">twitch streamer</a>, vtuber, and philosopher that focuses on ways to help make technology easier to understand and do cursed things in the process. I live in {resume.location.city} with my husband and I do developer relations professionally. I am an avid writer for my <a href="/blog">blog</a>, where I have over 400 articles. I regularly experiment with new technologies and find ways to mash them up with old technologies for my own amusement.
          </p>
        </div>
      </div>

      <h2 class="text-2xl mb-4">Recent Articles</h2>
      <ul class="list-disc ml-4 mb-4">
        {search.pages("index=true", "order date=desc", 8).map((post) => {
          const url = post.redirect_to ? post.redirect_to : post.url;
          return (
            <li>
              <time datetime={date(post.date)} className="font-mono">{post.date.toLocaleDateString("en-US", dateOptions)}</time> -{" "}
              <a href={url}>{post.title}</a>
            </li>
          );
        })}
      </ul>

      <SponsorCard />

      <h2 class="text-2xl mb-4">Notable Publications</h2>
      <ul class="list-disc ml-4 mb-4">
        {resume.notablePublications.map((publication) => (
          <li>
            <a href={publication.url}>{publication.title}</a><br />{publication.description}
          </li>
        ))}
      </ul>

      <h2 class="text-2xl mb-4">Highlighted Projects</h2>
      <ul class="list-disc ml-4 mb-4">
        {notableProjects.map((project) => (
          <li>
            <a href={project.url}>{project.title}</a> - {project.description}
          </li>
        ))}
      </ul>

      <h2 class="text-2xl mb-4">Quick Links</h2>
      <ul class="list-disc ml-4 mb-4">
        {contactLinks.map((link) => (
          <li>
            <a rel="me" target="_blank" href={link.url}>{link.title}</a>
          </li>
        ))}
      </ul>

      <p class="mb-4">Looking for someone for your team? Check <a href="/signalboost">here</a>.</p>

      <div class="flex flex-wrap items-start justify-center p-5">
        {resume.buzzwords.map((buzzword) => (
          <span class="m-2 p-2 text-sm bg-bg-1 dark:bg-bgDark-1 text-fg-1 dark:text-fgDark-1 rounded-lg">{buzzword}</span>
        ))}
      </div>
    </>
  );
}
