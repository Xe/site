export const layout = "base.njk";
export const date = "2012-12-21";

export default ({ search, resume, notableProjects, contactLinks }, { date }) => {
    const dateOptions = { year: "numeric", month: "2-digit", day: "2-digit" };

    return (
        <>
            <h1 class="text-3xl mb-4">{resume.name}</h1>
            <p class="text-xl mb-4">{resume.tagline} - {resume.location.city}, {resume.location.country}</p>
            <div className="my-4 flex space-x-4 rounded-md border border-solid border-fg-4 bg-bg-2 p-3 dark:border-fgDark-4 dark:bg-bgDark-2 max-w-full min-h-fit">
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
