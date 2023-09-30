---
title: "Site Update: Revue backups are live"
date: 2023-01-05
series: site-update
hero:
  ai: "Waifu Diffusion v1.3 (float16)"
  file: "kirito"
  prompt: "Kirito Sword Art Online, intricate, elegant, highly detailed, anime, genshin impact, thick outlines, concept art, smooth, sharp focus, illustration, third eye, portrait"
---

As I mentioned [previously](https://xeiaso.net/blog/site-update-no-more-revue),
I used to run a newsletter with Revue. Revue was bought by Twitter a while back,
and I chose to go with Revue instead of writing my own platform because I wanted
something a bit more turnkey and also because transactional email sending is
hell.

Turns out that [Revue is getting shut
down](https://www.getrevue.co/app/offboard). This means that I won't be able to
use it anymore (and I stopped using it because it wasn't getting much traction
vs the amount of work I put into it).

I have backed up all of my posts on Revue to my blog save one that I want to
rewrite and extend. You can see a list of all the posts [in the revueBackup
series](https://xeiaso.net/blog/series/revueBackup).

These posts are more "traditional" reviews and essays. They probably should have
been put on my main blog in the first place, but I wanted to make them exclusive
to Revue in order to make my newsletter more enticing.

- [My Thoughts on Paper Mario and the Origami
  King](https://xeiaso.net/blog/paper-mario-origami-king-2021-01-30)
- [Plurality as Portrayed in Cyberpunk 2077 and Xenoblade Chronicles
  2](https://xeiaso.net/blog/plurality-cyberpunk-xenoblade-2021-02-14)
- [Animal Crossing New Horizons: An Island of Stability in an Unstable
  World](https://xeiaso.net/blog/animal-crossing-stability-2021-02-28)

These posts are all sketches of fiction that I rescued from my notes. Not much
here is totally complete, but nearly all of it helps convey a mood that I wanted
to convey with writing.

- [Creation](https://xeiaso.net/blog/creation)
- [Mara's Ransack of Castle
  Charon](https://xeiaso.net/blog/mara-ransack-castle-2021-03-28)
- [Immigration](https://xeiaso.net/blog/immigration-2021-04-11)
- [New You](https://xeiaso.net/blog/new-you)
- [Untitled Furry Cyberpunk
  Story](https://xeiaso.net/blog/untitled-furry-cyberpunk-story)
- [Alone](https://xeiaso.net/blog/alone)
- [Second Go Around](https://xeiaso.net/blog/second-go-around)
- [The New Gods](https://xeiaso.net/blog/the-new-gods)
- [I Forgive You](https://xeiaso.net/blog/i-forgive-you-2021-08-08)
- [Theseus](https://xeiaso.net/blog/theseus)

If I ever do a newsletter again, I don't think I'm going to be using a third
party platform to do it. I'm going to self-host it so a billionaire buying up
the parent company doesn't kill it.

## Pronouns

I have made a page on my website that lists the pronouns that you should use
when talking about me. [Anything on this list](https://xeiaso.net/pronouns) is
fair game. I even made an [API endpoint](https://xeiaso.net/api/pronouns) in
case you want to query them programmatically for some reason. Here is the schema
of the reply:

The pronouns route replies with a JSON-formatted list of pronoun objects. Each
pronoun object has the following fields that correlate to common usage in
English:

- `nominative` - The pronoun used when the person is the subject of the
  sentence. The "she" in the sentence "She went to the park."
- `accusative` - The pronoun used when the person is the object of the sentence.
  The "her" in "I went with her."
- `possessiveDeterminer` - The pronoun used when talking about an object's
  owner. Also known as the dependent possessive case. It is the "his" in "He
  brought his frisbee."
- `possessive` - The pronoun used when talking about the owner of something,
  except used as an adjective directly in place of the owner's name. It is the
  "hers" in the sentence "At least I think it was hers."
- `reflexive` - The pronoun used when someone is talking about themselves in the
  third person. It is the "themselves" in the sentence "They threw the frisbee
  to themselves".
- `singular` - If true, then the pronoun is to be treated as singular. If not,
  it is to be treated as plural. This gets somewhat pedantically weird when
  singular "they" is in play, but most of the time you can get a hang of it with
  a few minutes of writing sentences.

---

I hope this is interesting! I hope 2023 will continue to be a prolific year in
the wild world of the xe iaso dot net cinematic universe!
