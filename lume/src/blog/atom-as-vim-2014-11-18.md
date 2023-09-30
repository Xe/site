---
title: My Experience with Atom as A Vim User
date: 2014-11-18
series: medium-archive
tags:
 - atom
 - vim
---

Historically, I am a Vim user. People know me as a very very heavy vim
user. I have spent almost the last two years customizing [my .vimrc
file](https://github.com/Xe/dotfiles/blob/master/.vimrc) and I have parts 
of it mapped to the ways I think. Recently I have acquired both a Mac Pro 
and a Surface Pro 3, and my vim configuration didn't work on them. For a 
while I had used Docker and the image I made of my preferred dev 
environment to shim and hack around this.

Then I took a fresh look at [Atom](https://atom.io/){.markup--anchor
.markup--p-anchor}, Github's text editor that claims to be a replacment
for Sublime. Since then I have moved to using Atom as my main text
editor for programming in OSX and Windows, but still using my fine-tuned
vim setup in Linux. I like how I have Atom set up. It uses a lot of (but
not all sadly) the features I have come to love in my vim setup.

I also like that I can have the same setup on both my Mac and in
Windows. I have the same
[vim-mode](https://github.com/atom/vim-mode) bindings on both my machines 
(I only customize so far as to add :w and :q bindings), and easily jump 
from one to the other with Synergy and have little to no issues with 
editor differences. I typically end up taking my surface out with me to
a lot of places and will code some new ideas on the bus or in the food 
court of the mall.

Atom gets a lot of things right with the plugins I have. I have
Autocomplete+ and a plugin for it that uses GoCode for autocompletion as
I type like I have with vim-go and YouCompleteMe in Vim. Its native
pacakge support and extensibility is bar none the easiest way to be able
to add things to the editor I have ever seen.

But there are problems with Atom that are mostly based on my usage of
text editors and my understanding of programming with Javascript,
Coffeescript, HTML and CSS. Atom is a mostly Coffeescript editor, it
does mean that at runtime I can customize almost any aspect of the
editor, but I would have to learn one if not 5 more languages to be able
to describe the layouts or interfaces I would like to add to this
editor. It also being a hybrid between a web application and a normal
desktop application means that I am afraid to add things I normally
would such as raw socket support for being able to collaborate on a
single document, PiratePad style. Additionally, the Vim emulation mode
in Atom doesn't support ex-style :-commands nor \<Leader\>, meaning that
a fair bit of my editing is toned down and done more manually to make up
for this.

I wish I could just use vim natively with my preferred setup on Windows,
OSX and Linux, but for now Atom is the lesser of all the evils.

---

Update: I am now atom-free on my surface pro 3
