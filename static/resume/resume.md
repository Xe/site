# Christine Dodrill

#### Full-stack Engineer

##### Montreal, QC &emsp; [christine.website][homepage]

`Docker`, `Git`, `Go`, `Rust`, `C`, `Stenography`, `DevOps`, `Heroku`, `Continuous
Integration/Delivery`, `WebAssembly`, `Lua`, `Mindfulness`, `HTTP/2`, `Alpine
Linux`, `Ubuntu`, `Linux`, `GraphViz`, `Progressive Web Apps`, `yaml`, `SQL`,
`Postgres`, `MySQL`, `SQLite`, `Ordained Minister`, `Dudeism`, `Tech Writing`,
`Kubernetes`, `Command Line Apps`

## Experience

### Tailscale - Software Designer &emsp; <small>*2020 - present*</small>

> [Tailscale][tailscale] is a zero config VPN for building secure networks.
> Install on any device in minutes. Remote access from any network or physical
> location.

#### Highlights

- Go programming
- Nix and NixOS

### Lightspeed - Expert principal en fiabilité du site &emsp; <small>*2019 - 2020*</small>

(Senior Site Reliability Expert)

> [Lightspeed][lightspeedhq] is a provider of retail, ecommerce and
> point-of-sale solutions for small and medium scale businesses. 

#### Highlights

- Migration from cloud to cloud
- Work on the cloud platform initiative
- Crafting reliable infrastructure for clients of customers
- Creation of an internally consistent and extensible command line interface for
  internal tooling

### Heroku - Senior Software Engineer &emsp; <small>*2017 - 2019*</small>

> [Heroku][heroku] is a cloud Platform-as-a-Service (PaaS) that created the term
> "platform as a service". Heroku currently supports several programming
> languages that are commonly used on the web. Heroku, one of the first cloud
> platforms, has been in development since June 2007, when it supported only the
> Ruby programming language, but now supports Java, Node.js, Scala, Clojure,
> Python, PHP, and Go. 

#### Highlights

- [JVM Application Metrics](https://devcenter.heroku.com/changelog-items/1133)
- [Go Runtime Metrics
  Agent](https://github.com/heroku/x/tree/master/runtime-metrics)
- Other backend fixes and improvements on [Threshold
  Autoscaling](https://blog.heroku.com/heroku-autoscaling) and [Threshold
  Alerting](https://devcenter.heroku.com/articles/metrics#threshold-alerting)
- [How to Make a Progressive Web App From Your Existing
  Website](https://blog.heroku.com/how-to-make-progressive-web-app)

### Backplane.io - Software Engineer &emsp; <small>*2016 - 2016*</small>

> [Backplane](https://backplane.io) (now defunct) was an innovative reverse reverse proxy that
> helps administrators and startups simplify their web application routing.

#### Highlights

- Performance monitoring of production servers
- Continuous deployment and development in Go
- Learning a lot about HTTP/2 and load balancing

### Pure Storage - Member of Technical Staff &emsp; <small>*2016 - 2016*</small>

> Pure Storage is a Mountain View, California-based enterprise data flash storage
> company founded in 2009. It is traded on the NYSE (PSTG).

#### Highlights

- Code maintenance

### IMVU - Site Reliability Engineer &emsp; <small>*2015 - 2016*</small>

> IMVU, inc is a company whose mission is to help people find and communicate
> with eachother. Their main product is a 3D avatar-based chat client and its
> surrounding infrastructure allowing creators to make content for the avatars
> to wear.

#### Highlights

- Wrote up technical designs
- Implemented technical designs on an over 800 machine cluster
- Continuous learning of a lot of very powerful systems and improving upon them
  when it is needed 

### VTCSecure - Deis Consultant (contract) &emsp; <small>*2014 - 2015*</small>

> VTCSecure is a company dedicated to helping with custom and standard
> audio/video conferencing solutions. They specialize in helping the deaf and
> blind communicate over today's infrastructure without any trouble on their end.

#### Highlights

- Started groundwork for a dynamically scalable infrastructure on a project for
  helping the blind see things 
- Developed a prototype of a new website for VTCSecure
- Education on best practices using Docker and CoreOS
- Learning Freeswitch

### Crowdflower - Deis Consultant (Contract) &emsp; <small>*2014 - 2014*</small>

> Crowdflower is a company that uses crowdsourcing to have its customers submit
> tasks to be done, similar to Amazon's Mechanical Turk. CrowdFlower has over 50
> labor channel partners, and its network has more than 5 million contributors
> worldwide.

#### Highlights

- Research and development on scalable Linux deployments on AWS via CoreOS and
  Docker
- Development of in-house tools to speed instance creation
- Laid groundwork on the creation and use of better tools for managing large
  clusters of CoreOS and Fleet machines

### OpDemand - Software Engineering Intern &emsp; <small>*2014 - 2014*</small>

> OpDemand is the company behind the open source project Deis, a distributed
> platform-as-a-service (PaaS) designed from the ground up to emulate Heroku but
> on privately owned servers.

#### Highlights

- Built new base image for Deis components
- Research and development on a new builder component

## Portfolio Highlights

### [Olin](https://github.com/Xe/olin)

An embeddable userspace kernel for executing WebAssembly programs.
The main goal of this is to allow for an easier migration to another CPU
architecture (such as RISC-V, aarch64 or ppc64be) without having to recompile
existing code.

I have written multiple blogposts on this project:

- https://christine.website/blog/olin-1-why-09-1-2018
- https://christine.website/blog/olin-2-the-future-09-5-2018
- https://christine.website/blog/olin-progress-2019-12-14

As of March 21, 2019, Olin is able to run binaries compiled with [Go 1.12.x
WebAssembly support](https://github.com/golang/go/wiki/WebAssembly). Olin also
is known to work on big-endian systems with no changes needed to source code or
binaries.

It also supports security policies similar to a combination of the OSX sandbox
profiles and OpenBSD's pledge() system call. This allows users to limit the
scope of what resources an Olin program can access, including file URLs, the
amount of ram that can be used or the number of WebAssembly instructions that
can be executed.

### [Wasmcloud](https://tulpa.dev/within/wasmcloud)

Wasmcloud is a Heroku or AWS Lambda-like functions as a service backend and
platform for event-driven architecture built on top of WebAssembly. It wraps
[Olin](https://github.com/Xe/olin) and provides a lot of higher-level
conveniences for users. I have written [a
blogpost](https://christine.website/blog/wasmcloud-progress-2019-12-08) on my
progress and where I'm wanting to go with this project.

### [ilo Kesi](https://github.com/Xe/x/tree/master/discord/ilo-kesi)

A chatbot that parses its commands through the grammar of the constructed
language [Toki Pona](http://tokipona.org), then figures out what the user is
asking for using a lookup table and executes that request.

### When Then Zen

[When Then Zen](https://when-then-zen.christine.website) is meditation instructions translated into Gherkin, a-la:

```
Feature: Anapana (mindfulness via breathing) meditation
  Background:
    Given no assumption about meditation background
    And a willingness to learn
    And no significant problems with breathing through the body's nose
    And I am seated or laying down comfortably
    And no music is playing

  Scenario Outline: mindfulness of breathing
    As a meditator
    In order to be mindful of the body's breath
    When I <verb> through the body's nose
    Then I focus on the sensations of breath
    Then I focus on the feelings of breath through the nasal cavity
    Then I focus on the feelings of breath interacting with the nostrils
    Then I repeat until done

    Examples:
      | verb   |
      | inhale |
      | exhale |
```

This has been well-recieved by coworkers, friends and others. I have written more on the subject [here](https://christine.website/blog/when-then-zen-anapana-2018-08-15).

## Writing

> Articles listed below will be either personal or professional and do not reflect the views of any company or group I am affiliated with. The writing is my own, with the help of others to make things legible.

- [My Blog](https://christine.website/blog)
- [I Put Words on This Webpage so You Have to Listen to Me Now](https://christine.website/blog/experimental-rilkef-2018-11-30)

I have gotten to the front page of [Hacker News](https://news.ycombinator.com) several times. Here are a few of the comment threads:

- [I Put Words on This Webpage so You Have to Listen to Me now](https://news.ycombinator.com/item?id=18577758)
- [TempleOS: 1 - Installation](https://news.ycombinator.com/item?id=19961082)
- [WebAssembly on the Server: How System Calls Work](https://news.ycombinator.com/item?id=20066204)
- [Olin: Defining a New Primitive for Event-Driven Services](https://news.ycombinator.com/item?id=17896307)

## Ordination

I am an ordained minister with the [Church of the Latter-day Dude](https://dudeism.com). This allows me to officiate religious ceremonies in at least the United States. I would be honored if you were to choose me to officiate anything for any reason. Please [contact](/contact) me if you have any questions.

[homepage]: https://christine.website
[twitter]: https://twitter.com/theprincessxena
[twit]: http://cdn-careers.sstatic.net/careers/Img/icon-twitter.png?v=b1bd58ad2034
[heroku]: https://www.heroku.com
[lightspeedhq]: https://www.lightspeedhq.com
[tailscale]: https://tailscale.com/
