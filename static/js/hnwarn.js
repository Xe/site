import { g, x, r, t } from "./xeact.min.js";
import { div, ahref, br } from "./xeact-html.min.js";
import { mkConversation } from "./conversation.js";

// list of regexps for potentially problematic referrers to display the nag to
const FLAGGED_REFERRERS = [
    /^https?:\/\/((.+)\.)?reddit\.com/i,
    /^https?:\/\/news\.ycombinator\.com/i,
];

const addNag = () => {
    let root = g("refererNotice");
    x(root);
    root.appendChild(
        div(
            {style: "padding:1em"},
            mkConversation("Cadey", "coffee", [
                t("Thank you for reading this article. If you have any questions or thoughts about its contents, please comment civilly on it and remember the human on the other side of the screen. Due to facts and circumstances surrounding our fundamentally subjective reality, I may experience things differently than you do. If this is somehow unacceptable to you, please feel free to "),
                ahref("https://zombo.com", "go somewhere else"),
                t(". Have a good day and be well!")
            ], "warning"),
            br(),
            br()
        )
    );
};

r(() => {
    const ref = document.referrer;
    if (!ref) return;

    if (FLAGGED_REFERRERS.some(r => r.test(ref))) {
        addNag();
    }
});
