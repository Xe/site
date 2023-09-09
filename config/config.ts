import types from "./types.ts";

import authors from "./authors.ts";
import characters from "./characters.ts";
import contactLinks from "./contactLinks.ts";
import jobHistory from "./jobHistory.ts";
import notableProjects from "./notableProjects.ts";
import pronouns from "./pronouns.ts";
import seriesDescriptions from "./seriesDescriptions.ts";
import signalBoost from "./signalboost.ts";
import vods from "./vods.ts";

export type Config = {
    signalBoost: types.Person[];
    defaultAuthor: types.Author;
    authors: Record<string, types.Author>;
    clackSet: string[];
    webMentionEndpoint: string;
    jobHistory: types.Job[];
    seriesDescriptions: Record<string, string>;
    notableProjects: types.Link[];
    contactLinks: types.Link[];
    pronouns: types.PronounSet[];
    characters: types.Character[];
    vods : types.StreamVOD[];
    redirects?: Record<string, string>;
};

const config: Config = {
    signalBoost,
    defaultAuthor: authors["xe"],
    authors,
    clackSet: [
        "Ashlynn",
        "Terry Davis",
        "Dennis Ritchie",
        "Steven Hawking",
        "Kris Nova",
    ],
    webMentionEndpoint: "https://mi.within.website/api/webmention/accept",
    jobHistory,
    seriesDescriptions,
    notableProjects,
    contactLinks,
    pronouns,
    characters,
    vods,
    redirects: {
        "/blog/🥺": "/blog/xn--ts9h",
    },
};

export default config;