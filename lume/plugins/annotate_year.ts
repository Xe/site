import type Site from "lume/core/site.ts";

export default function () {
    return (site: Site) => {
        site.preprocess([".html"], (pages) => {
            for (const page of pages) {
                page.data.year = page.data.date.getFullYear();
            }
        });
    };
}