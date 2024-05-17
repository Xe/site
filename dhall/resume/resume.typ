#let primary_colour = rgb("#3E0C87") // vivid purple

#let coloredLine() = {
     line(length: 100%, stroke: 1pt + primary_colour)
}

#let sigil() = {
  align(center + horizon)[#box(
    height: 24pt,
    image("icons/xeiaso.svg")
  )]
}

#let icon(name, shift: 1.5pt) = {
  box(
    baseline: shift,
    height: 10pt,
    image("icons/" + name + ".svg")
  )
  h(3pt)
}

#let findMe(services) = {
  set text(8pt)
  let icon = icon.with(shift: 2.5pt)

  services.map(service => {
      icon(service.name)

      if "display" in service.keys() {
        link(service.link)[#{service.display}]
      } else {
        link(service.link)
      }
    }).join(h(10pt))
  [
    
  ]
}

#let term(period, location) = {
  text(9pt)[#icon("calendar") #period #h(1fr) #icon("location") #location]
}

#let alta(
  name: "",
  links: (),
  tagline: [],
  doc,
) = {
  set text(9.8pt, font: "Iosevka Aile Iaso", tracking: 0.0225pt)
  set par(justify: true)
  set page(
    paper: "us-letter",
    margin: (x: 54pt, y: 52pt),
  )

  show heading.where(
    level: 2
  ): it => text(
      fill: primary_colour,
      font: "Iosevka Etoile Iaso",
    [
      #{it.body}
      #v(-7pt)
      #line(length: 100%, stroke: 1pt + primary_colour)
    ]
  )

  show heading.where(
    level: 4
  ): it => text(
    fill: primary_colour,
    font: "Iosevka Etoile Iaso",
    it.body
  )
  
  [= #name]
  findMe(links)

  columns(2, doc)
}

#let resume = json("resume.json")

#let titleCompany(title, company) = [
     #text(10pt)[#title] #h(1fr) _#[#company]_
]

#let publication(details) = [
     #text(font: "Iosevka Etoile Iaso", size: 10.5pt, weight: 600)[#link(details.url)[#details.title]]\
     #details.description\
]

#let sigilAside(body) = {
  let cell = rect.with(
    width: 100%,
    radius: 0pt,
    stroke: none,
  )
  grid(
    columns: (auto, 4pt, 35pt),
    gutter: -1pt,
    cell(body),
    cell(box(
      width: 37pt,
      image("icons/xeiaso.svg")
    )),
  )
}

#show: doc => alta(
  name: resume.name,
  links: (
    (name: "email", link: "mailto:me@xeiaso.net"),
    (name: "website", link: "https://xeiaso.net/", display: "xeiaso.net"),
    (name: "github", link: "https://github.com/Xe", display: "@Xe"),
    (name: "linkedin", link: "https://www.linkedin.com/in/xe-iaso/", display: resume.name),
    (name: "mastodon", link: "https://pony.social/@cadey", display: "@cadey@pony.social"),
    (name: "twitch", link: "https://twitch.tv/princessxen", display: "@princessxen"),
  ),
  tagline: resume.tagline,
  doc,
)

#sigilAside[Hello, I'm Xe Iaso. I am a skilled force multiplier, acclaimed speaker, artist, and prolific blogger. My writing is widely viewed across 15 time zones and is one of the most viewed software blogs in the world.]

I specialize in helping people realize their latent abilities and help to unblock them when they get stuck. This creates unique value streams and lets me bring others up to my level to help create more senior engineers. I am looking for roles that allow me to build upon existing company cultures and transmute them into new and innovative ways of talking about a product I believe in. I am prioritizing remote work at companies that align with my values of transparency, honesty, equity, and equality.

If you want someone that is dedicated to their craft, a fearless innovator and a genuine force multiplier, please look no further. I'm more than willing to hear you out.    

== Experience

#titleCompany[Archmage of Infrastructure][Tailscale]\
#term[2020-12 -- 2023-10][Ottawa, CA]

At Tailscale I founded the developer relations team with the goal of inspiring people to use Tailscale in fun and interesting ways. I work with the DevRel team to write articles for #text(fill: blue)[#link("https://tailscale.dev")[tailscale.dev]] to help teach people fundamentals of computer science and networking in the process of learning about new product features and things you can do with them.

Tailscale has easily been the best job I've ever had and I've had an enormous impact on how Tailscale is percieved by developers worldwide. For a while my actions were directly attributable to MAU growth. One of the hardest projects I've worked on was making DevRel efforts more than single flashes in the pan and create a reason for people to visit the website on a regular basis.

I regularly wrote articles and lead internal talks about how to use Tailscale and other technology topics in new and interesting ways.

#coloredLine()

#titleCompany[Expert principal en fiabilité du site][Lightspeed]\
#term[2019-05 -- 2020-11][Montréal, CA]

I created and maintained the internal Kubernetes deployment (with the goal of being functionally an in-house Platform-as-a-service) and all of the assorted tooling around it, helping internal developers deploy new features to customers faster. I also helped to create custom icons and color schemes for internal tools, with the goal of having consistent brand design for knowing which tool is which at a glance.

- Created a custom Kubernetes deployment system with a focus on developer experience
- Built developer tooling on top of a custom CI/CD flow
- Managed hidden complexity in the Kubernetes deployment system

#coloredLine()

#titleCompany[Senior Software Engineer][Heroku]\
#term[2017-02 -- 2019-03][Bellevue, USA]

I maintained the subsystem for processing terabytes of customer metrics per week in real time, and tools that consumed this data, such as threshold alerting and autoscaling. We were hitting theoretical limits for Kafka performance by the time I left.

- Wrote and maintained integrations for JVM application metrics and Go runtime metrics
- Developed a FaaS platform prototype with my team
- Helped my team navigate complicated corporate politics and market conditions beyond our control

#coloredLine()

My work history before 2017 is available upon request.

== Notable Publications

#for pub in resume.notablePublications [
    #publication(pub)
    
]

== Projects

==== Xesite\
The custom blog engine that powers #text(fill: blue)[#link("https://xeiaso.net")[xeiaso.net]]. It is a handcrafted work of art written in Rust with two goals:

1. To be as fast as possible to survive traffic surges from news aggregators without flinching.
2. To be easily extensible and hackable to meet my needs.

This project has been an overwhelming success and is the backbone of a lot of my personal development. Most of the things that would otherwise be written as separate projects have become extensions and modifications to my blog engine, allowing them to get a lot of traffic and hands-on user experience as soon as possible.

==== XeDN\
The software that powers #text(font: "Iosevka Curly Iaso Extended")[#link("https://cdn.xeiaso.net")[cdn.xeiaso.net]], currently serving over 4 terabytes of traffic per month without breaking a sweat. This powers the images, video, slides, and other files that I use to enrich my posts.

==== Xeact\
My custom frontend JavaScript framework. I regularly write about things I have learned working on it #text(fill: blue)[#link("https://xeiaso.net/blog/series/xeact")[on my blog]].

Working on this project has been the catalyst for me finally understanding how to do front-end development in web browsers. Xeact also powers internal tooling at Tailscale, being used daily by the support team to understand and diagnose customer issues.

==== #link("https://github.com/Xe/waifud")[waifud]\
My custom virtual machine manager for my homelab. waifud replicates most of the advantages of using cloud computing in my basement. I write about it regularly #text(fill: blue)[#link("https://xeiaso.net/blog/series/waifud")[on my blog]].
