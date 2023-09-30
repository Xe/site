---
title: "[ANN] majc 0.2.1"
date: 2020-07-28
series: flightJournal
---

A parsing bug has been found in majc 0.2.0. Specifically when parsing pages that include a comment in a preformatted text block, majc's parser could accidentally eat the entire document and only render the contents of that preformatted text block. For an example of this, see the source of gemlog.blue.

A fix has been made and this behavior should never happen again. If you find other cases where this kind of unexpected eating the document thing happens, please contact me and I will fix it as soon as I am able to.

You can download majc 0.2.1 from greedo.

Be well.
