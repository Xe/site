---
title: "Ten Thousand Laughs"
date: "2018-12-01"
series: conlangs
---

```
pemci zo'e la xades  
ni'o pano ki'o nu cmila  
.i cmila cei broda  
.i ke broda jo'u broda jo'u broda jo'u broda jo'u broda jo'u broda 
 jo'u broda jo'u broda jo'u broda jo'u broda ke'e cei brode  
.i ke brode jo'u brode jo'u brode jo'u brode jo'u brode jo'u brode
 jo'u brode jo'u brode jo'u brode jo'u brode ke'e cei brodi  
.i ke brodi jo'u brodi jo'u brodi jo'u brodi jo'u brodi jo'u brodi
 jo'u brodi jo'u brodi jo'u brodi jo'u brodi ke'e cei brodo  
.i ke brodo jo'u brodo jo'u brodo jo'u brodo jo'u brodo jo'u brodo
 jo'u brodo jo'u brodo jo'u brodo jo'u brodo ke'e cei brodu  
.i mi brodu
```

This is a synthesis of the [broda](https://lojban.org/publications/cll/cll_v1.1_xhtml-section-chunks/section-koha-broda-series.html) family of gismu in Lojban. In order to properly understand this lojban text, you must conceive laughter ten thousand times. This is a reference to the [Billion laughs attack](https://en.wikipedia.org/wiki/Billion_laughs_attack) that XML parsers can suffer from.

Translation:

```
Poem by Cadey
Ten Thousand Laughs

I laugh, and then I laugh, and then I laugh, and then I laugh (... 10,000 times in total).
```

This is roughly equivalent to the following XML document:

```xml
<?xml version="1.0"?>
<!DOCTYPE lolz [
 <!ENTITY lol "lol">
 <!ELEMENT lolz (#PCDATA)>
 <!ENTITY lol1 "&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;&lol;">
 <!ENTITY lol2 "&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;&lol1;">
 <!ENTITY lol3 "&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;&lol2;">
 <!ENTITY lol4 "&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;&lol3;">
]>
<lolz>&lol4;</lolz>
```
