let xesite = ./types/package.dhall

let Resume = xesite.Resume

let Link = xesite.Link

in  Resume::{
    , buzzwords =
      [ "Docker"
      , "Git"
      , "Go"
      , "Rust"
      , "C"
      , "DevOps"
      , "Heroku"
      , "WebAssembly"
      , "Lua"
      , "Mindfulness"
      , "Nix"
      , "NixOS"
      , "HTTP/2"
      , "Ubuntu"
      , "Alpine Linux"
      , "GraphViz"
      , "JavaScript"
      , "TypeScript"
      , "SQLite"
      , "PostgreSQL"
      , "Dudeism"
      , "Technical writing"
      , "Emacs"
      , "Continuous Integration"
      , "Continuous Delivery"
      ]
    , notablePublications =
      [ Link::{
        , url = "https://blog.heroku.com/how-to-make-progressive-web-app"
        , title = "How to Make a Progressive Web App From Your Existing Website"
        , description =
            "An article summarizing how easy it is to make a webpage into an installable Progressive Web App using APIs available in the most commonly used browsers."
        }
      , Link::{
        , url =
            "https://web.archive.org/web/20210318102148/https://tech.lightspeedhq.com/palisade-version-bumping-at-scale-in-ci/"
        , title = "Palisade: Version Bumping at Scale in CI"
        , description =
            "The release post for Palisade, a tool to automate version bumping, release tagging and more."
        }
      , Link::{
        , url = "https://tailscale.com/blog/magicdns-why-name/"
        , title = "An epic treatise on DNS, magical and otherwise"
        , description =
            "A deep dive into all of the problems that DNS has at scale and how Tailscale makes most of those problems go away, with the rest of them being easier in comparison."
        }
      , Link::{
        , url = "https://tailscale.dev/blog/weaponizing-hyperfocus"
        , title =
            "Weaponizing hyperfocus: Becoming the first DevRel at Tailscale"
        , description =
            "A brief history of the developer relations team at Tailscale and how I found myself creating it. I cover one of my largest internal demons and how I managed to wield it as a force of empowerment rather than a limiting force of pain."
        }
      , Link::{
        , url = "https://tailscale.dev/blog/headscale-funnel"
        , title = "Using Tailscale without Using Tailscale"
        , description =
            "An award-winning April Fools Day post describing how you can use Tailscale via Headscale via Tailscale Funnel. This post is notable for demonstrating all five of the Tailscale company values at the same time."
        }
      , Link::{
        , url = "https://fly.io/blog/how-i-fly/"
        , title = "How I Fly"
        , description =
            "A post about how I use Fly.io to host my personal website and its supporting infrastructure."
        }
      ]
    }
