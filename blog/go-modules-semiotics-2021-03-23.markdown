---
title: |
  The Semiotics of "Simple": Go Modules
date: 2021-03-23
tags:
  - philosophy
  - golang
  - simple
---

<big>DISCLAIMER</big>

In this post there are opinions. These opinions are my opinions and are purely targeted against ideas, not people. I have a lot of respect for the things I am talking about. These arguments are coming from a place of passion and wanting to make software understandable. Please keep this in mind before reacting to any opinions in this article.

# The Semiotics of "Simple": Go Modules

This sentence is a series of words that conveys a meaning. Each word in this sentence has its own meanings, but individually they don't really mean much. Pastry. Pasta. However, these words have meanings that can be as simple as a name for a color (ex: pink) and as complicated as an entire way of thinking (ex: fatalism). Semiotics is a branch of philosophy that focuses on the meanings and interpretations of symbols, such as words, logos, icons and similar things. In this article, we're going to examine the word "simple" and how different interpretations of that word can result in vastly different outcomes, specifically in the context of [Go Modules](https://blog.golang.org/using-go-modules).

There is no central body that dictates what is and is not correct in English. This lack of a central body makes analysis of semiotics in English very hard. However, [Merriam-Webster](https://www.merriam-webster.com/dictionary/simple) is generally considered to be an authoritative dictionary in Anglophone communities, so let's start with the definitions it lists:

> Definition of _simple_
>
> 1. free from guile
> 2. free from vanity
> 3. free from ostentation or display 
> 4. of humble origin or modest position 
> 5. lacking in knowledge or expertise
> 6. lacking in intelligence
> 7. not socially or culturally sophisticated
> 8. free of secondary complications
> 9. having only one main clause and no subordinate clauses 
> 10. _of a subject or predicate_ : having no modifiers, complements, or objects
> 11. constituting a basic element
> 12. not made up of many like units
> 13. free from elaboration or figuration
> 14. not subdivided into branches or leaflets
> 15. consisting of a single carpel
> 16. developing from a single ovary
> 17. controlled by a single gene
> 18. not limited or restricted : unconditional
> 19. readily understood or performed 
> 20. _of a statistical hypothesis_ : specifying exact values for one or more statistical parameters

"simple" is one of the more "core" parts of English vocabulary. As such it ends up having a lot of meanings and senses of the word, which makes it very easy for people to get confused.

Let's start this analysis with [Minimal Version Selection](https://research.swtch.com/vgo-mvs), which is described as a "simple" alternative to other solutions of the [version selection problem](https://en.wikipedia.org/wiki/Boolean_satisfiability_problem). To put a complicated problem in very simple terms, this is solving the problem of picking versions of libraries needed to build a program. The Minimal Version Selection paper does a really great job of exploring the problem space. If you are at all unfamiliar with this topic, I really suggest reading it.

Something interesting about this problem is that it's [NP-complete](https://en.wikipedia.org/wiki/NP-completeness), meaning that it's one of the hardest problems you can solve with a computer. So, looking at it again, let's see which sense of the word "simple" is likely being used here. I think Russ Cox was getting at "readily understood or performed" with the word "simple". 

Minimal Version Selection is very readily understood and performed. It's mainly described in a paper on Russ' website and the paper is very concise and easy to understand. It does its job. It works, does its thing and then gets out of your way so you can go back to writing the programs you want to write.

However, as a side effect of this easy to understand (and likely very easy to program) implementation of the version selection problem, there are some practical issues that can crop up like landmines, especially so when you consider "simple" to mean "free of secondary complications".

Let's take this image from Russ' article as an example:

<center>

![A dependency tree showing an old version of package D selecting an old version of package E](https://cdn.christine.website/file/christine-static/blog/Screenshot+2021-03-19+214555.png)

</center>

Say that you want to depend on package B, which depends on package D, which depends on package E. With older versions of the Go toolchain, installing package B into your project grabs the newest possible version of that package and uses that. In Go modules, it will grab version 1.3 of package D, which depends on version 1.2 of package E. For people that have used Go for years, this behavior violates the expectation that you will always have the newest possible version of the program. 

Even worse, this makes it very easy to accidentally pull buggy packages into your program too. Let's say that version 1.3 of package D has a known security flaw that will result in code execution. Your solution will be to upgrade package D to version 1.4 to fix this, but unless the maintainers of package B set an [exclude](https://golang.org/ref/mod#go-mod-file-exclude) directive on that module forbidding the Go toolchain from using that version.

On one hand, yes this is simple as in it's easy to understand, but it is not simple in the sense that it has few hidden costs. They have very good reasons to avoid making the version selection code too complicated; however any way you look at this problem it is still a very difficult problem, and the complication has to go somewhere. In this case, doing the naiive thing that Go does can and will actively harm people in the worst ways.

But the Go team themselves may not end up actually hearing about that harm. The algorithm itself [has greater consequences beyond the code they write](https://youtu.be/ajGX7odA87k). Maybe it would be better to describe Minimal Version Selection as "easily understandable" rather than "simple".

[-Wpedantic: It may also be a good idea to try to avoid too much clever thinking into critical parts of infrastructure like this. Not to say there's no room for research into this branch of computer science, because there _definitely_ is; maybe it's something that shouldn't be put into production across hundreds of thousands of companies production environments.<br /><br />Operators of software are often distracted and not always knowledgeable about the minutae of how something they are using is implemented (and arguably they shouldn't have to be to prevent the <a href="https://www.techrepublic.com/article/why-pgp-is-fundamentally-flawed-and-needs-to-be-fixed/">PGP problem</a>).](conversation://Mara/hacker)

