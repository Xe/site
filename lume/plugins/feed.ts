import { getExtension } from "lume/core/utils/path.ts";
import { isPlainObject, merge } from "lume/core/utils/object.ts";
import { getGenerator } from "lume/core/utils/lume_version.ts";
import { getDataValue, getPlainDataValue } from "lume/core/utils/data_values.ts";
import { cdata, stringify } from "lume/deps/xml.ts";
import { Page } from "lume/core/file.ts";
import { log } from "lume/core/utils/log.ts";
import { parseDate } from "lume/core/utils/date.ts";

import type Site from "lume/core/site.ts";
import type { Data } from "lume/types.ts";
import type { stringifyable } from "lume/deps/xml.ts";

export interface Options {
    /** The output filenames */
    output?: string | string[];

    /** The query to search the pages */
    query?: string;

    /** The sort order */
    sort?: string;

    /** The maximum number of items */
    limit?: number;

    /** The feed info */
    info?: FeedInfoOptions;

    /** The feed items configuration */
    items?: FeedItemOptions;
}

export interface FeedInfoOptions {
    /** The feed title */
    title?: string;

    /** The feed subtitle */
    subtitle?: string;

    /**
     * The feed published date
     * @default `new Date()`
     */
    published?: Date;

    /** The feed description */
    description?: string;

    /** The feed language */
    lang?: string;

    /** The feed generator. Set `true` to generate automatically */
    generator?: string | boolean;
}

export interface FeedItemOptions {
    /** The item title */
    title?: string | ((data: Data) => string | undefined);

    /** The item description */
    description?: string | ((data: Data) => string | undefined);

    /** The item published date */
    published?: string | ((data: Data) => Date | undefined);

    /** The item updated date */
    updated?: string | ((data: Data) => Date | undefined);

    /** The item content */
    content?: string | ((data: Data) => string | undefined);

    /** The item categories */
    categories?: string | ((data: Data) => string[] | undefined);

    /** The item language */
    lang?: string | ((data: Data) => string | undefined);
}

export const defaults = {
    /** The output filenames */
    output: "/feed.rss",

    /** The query to search the pages */
    query: "",

    /** The sort order */
    sort: "date=desc",

    /** The maximum number of items */
    limit: 10,

    /** The feed info */
    info: {
        title: "My RSS Feed",
        published: new Date(),
        description: "",
        lang: "en",
        generator: true,
    },
    items: {
        title: "=title",
        description: "=description",
        published: "=date",
        content: "=children",
        categories: "=tags",
        lang: "=lang",
    },
} satisfies Options;

export interface FeedData {
    title: string;
    url: string;
    description: string;
    published: Date;
    lang: string;
    generator?: string;
    items: FeedItem[];
}

export interface FeedItem {
    title: string;
    url: string;
    description: string;
    published: Date;
    updated?: Date;
    content: string;
    lang: string;
    categories?: string[];
}

const defaultGenerator = getGenerator();

export default function (userOptions?: Options) {
    const options = merge(defaults, userOptions);

    return (site: Site) => {
        site.process(function processFeed() {
            const output = Array.isArray(options.output)
                ? options.output
                : [options.output];

            const pages = site.search.pages(
                options.query,
                options.sort,
                options.limit,
            );

            const { info, items } = options;
            const rootData = site.source.data.get("/") || {};

            const feed: FeedData = {
                title: getPlainDataValue(rootData, info.title),
                description: getPlainDataValue(rootData, info.description),
                published: toDate(getDataValue(rootData, info.published)) ??
                    new Date(),
                lang: getDataValue(rootData, info.lang),
                url: site.url("", true),
                generator: info.generator === true
                    ? defaultGenerator
                    : info.generator || undefined,
                items: pages.map((data): FeedItem => {
                    const content = getDataValue(data, items.content)?.toString();
                    const pageUrl = site.url(data.url, true);
                    const fixedContent = fixUrls(new URL(pageUrl), content || "");

                    // Linkposts point at the external article instead of the
                    // page on this site, and get an extra "external" tag.
                    const redirectTo = getDataValue(data, "=redirect_to");
                    const isLinkpost = redirectTo !== undefined;
                    const link = isLinkpost ? redirectTo : pageUrl;

                    const existingTags = toStringArray(
                        getDataValue(data, items.categories),
                    );
                    const categories = isLinkpost
                        ? [...existingTags, "external"]
                        : existingTags;

                    return {
                        title: getPlainDataValue(data, items.title),
                        url: link,
                        description: getPlainDataValue(data, items.description),
                        published: toDate(getDataValue(data, items.published)) ??
                            new Date(),
                        updated: toDate(getDataValue(data, items.updated)),
                        content: fixedContent,
                        lang: getDataValue(data, items.lang),
                        categories,
                    };
                }),
            };

            for (const filename of output) {
                const file = site.url(filename, true);

                switch (getMimeType(filename)) {
                    case "application/rss+xml":
                        site.pages.push(
                            Page.create({ url: filename, content: generateRss(feed, file) }),
                        );
                        break;

                    case "application/feed+json":
                        site.pages.push(
                            Page.create({ url: filename, content: generateJson(feed, file) }),
                        );
                        break;

                    default:
                        log.error(`[feed plugin] Invalid output format: ${filename}`);
                }
            }
        });
    };
}

function fixUrls(base: URL, html: string): string {
    return html.replaceAll(
        /\s(href|src)="([^"]+)"/g,
        (_match, attr, value) => ` ${attr}="${new URL(value, base).href}"`,
    );
}

function generateRss(data: FeedData, file: string): string {
    const feed: stringifyable = {
        "@version": "1.0",
        "@encoding": "UTF-8",
        rss: {
            "@xmlns:content": "http://purl.org/rss/1.0/modules/content/",
            "@xmlns:wfw": "http://wellformedweb.org/CommentAPI/",
            "@xmlns:dc": "http://purl.org/dc/elements/1.1/",
            "@xmlns:atom": "http://www.w3.org/2005/Atom",
            "@xmlns:sy": "http://purl.org/rss/1.0/modules/syndication/",
            "@xmlns:slash": "http://purl.org/rss/1.0/modules/slash/",
            "@version": "2.0",
            channel: {
                title: data.title,
                link: data.url,
                "atom:link": {
                    "@href": file,
                    "@rel": "self",
                    "@type": "application/rss+xml",
                },
                description: data.description,
                lastBuildDate: data.published.toUTCString(),
                language: data.lang,
                generator: data.generator,
                item: data.items.map((item) => ({
                    title: item.title,
                    link: item.url,
                    guid: {
                        "@isPermaLink": false,
                        "#text": item.url,
                    },
                    description: item.description,
                    "content:encoded": cdata(item.content),
                    category: item.categories,
                    pubDate: item.published.toUTCString(),
                    "atom:updated": item.updated?.toISOString(),
                })),
            },
        },
    };

    return stringify(clean(feed));
}

function generateJson(data: FeedData, file: string): string {
    const feed = {
        version: "https://jsonfeed.org/version/1.1",
        title: data.title,
        home_page_url: data.url,
        feed_url: file,
        description: data.description,
        language: data.lang,
        items: data.items.map((item) => ({
            id: item.url,
            url: item.url,
            title: item.title,
            language: item.lang,
            content_html: item.content,
            tags: item.categories,
            date_published: item.published.toISOString(),
            date_modified: item.updated?.toISOString(),
        })),
    };

    return JSON.stringify(clean(feed));
}

/** Remove undefined values of an object recursively */
function clean(obj: Record<string, unknown>): Record<string, unknown> {
    return Object.fromEntries(
        Object.entries(obj)
            .map(([key, value]): [string, unknown] => {
                if (isPlainObject(value)) {
                    const cleanValue = clean(value);
                    return [
                        key,
                        Object.keys(cleanValue).length > 0 ? cleanValue : undefined,
                    ];
                }
                if (Array.isArray(value)) {
                    const cleanValue = value
                        .map((v) => isPlainObject(v) ? clean(v) : v)
                        .filter((v) => v !== undefined);
                    return [
                        key,
                        cleanValue.length > 0 ? cleanValue : undefined,
                    ];
                }
                return [key, value];
            })
            .filter(([, value]) => value !== undefined),
    );
}

function toStringArray(value: unknown): string[] {
    if (!value) {
        return [];
    }
    const array = Array.isArray(value) ? value : [value];
    return array.map((v) => typeof v === "string" ? v : String(v));
}

function toDate(date?: string | number | Date): Date | undefined {
    if (date instanceof Date) {
        return date;
    }
    if (date === undefined) {
        return;
    }
    return parseDate(date);
}

function getMimeType(filename: string): string | undefined {
    const format = getExtension(filename).slice(1);

    switch (format) {
        case "rss":
        case "feed":
        case "xml":
            return "application/rss+xml";
        case "json":
            return "application/feed+json";
    }
}
