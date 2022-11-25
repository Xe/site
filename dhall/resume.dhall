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
    , jobs = ./jobHistory.dhall
    , notablePublications =
      [ Link::{
        , url = "https://blog.heroku.com/how-to-make-progressive-web-app"
        , title = "How to Make a Progressive Web App From Your Existing Website"
        , description =
            "An article summarizing how easy it is to make a webpage into an installable Progressive Web App."
        }
      , Link::{
        , url =
            "https://web.archive.org/web/20210318102148/https://tech.lightspeedhq.com/palisade-version-bumping-at-scale-in-ci/"
        , title = "Palisade: Version Bumping at Scale in CI"
        , description =
            "The release post for Palisade, a tool to automate version bumping, release tagging and more."
        }
      , Link::{
        , url = "https://tailscale.com/blog/grafana-auth/"
        , title = "How To Seamlessly Authenticate to Grafana using Tailscale"
        , description =
            "The release post for grafana-auth, a tool that lets Grafana users automagically authenticate to Grafana using Tailscale."
        }
      , Link::{
        , url = "https://tailscale.com/blog/tailscale-auth-minecraft/"
        , title = "Tailscale Authentication for Minecraft"
        , description =
            "A post explaining how Tailscale as an authentication mechanism can be used in absurd places, such as making authentication for Minecraft servers."
        }
      , Link::{
        , url = "https://tailscale.com/blog/tailscale-auth-nginx/"
        , title = "Tailscale Authentication for NGINX"
        , description =
            "The release post for nginx-auth, a tool that uses Tailscale's knowledge of IP address to person mappings to provide a weak authentication factor."
        }
      , Link::{
        , url = "https://tailscale.com/blog/steam-deck/"
        , title = "Putting Tailscale on the Steam Deck"
        , description =
            "An engineering log of all the steps taken to run Tailscale on the Valve Steam Deck and the tradeoffs between the various methods you could use to do this."
        }
      , Link::{
        , url = "https://tailscale.com/blog/gitops-acls/"
        , title = "GitOps for Tailscale ACLs"
        , description =
            "The release post of the Sync Tailscale ACLs GitHub Action, allowing administrators to automatically sync Tailscale ACLs from a GitHub repository."
        }
      , Link::{
        , url = "https://tailscale.com/blog/hamachi/"
        , title = "Tailscale: A modern replacement for Hamachi"
        , description =
            "A nostalgic piece recalling the magic of Hamachi (a product that did a similar thing to Tailscale), and how Tailscale builds on top of that to do even better things."
        }
      , Link::{
        , url = "https://tailscale.com/blog/magicdns-why-name/"
        , title = "An epic treatise on DNS, magical and otherwise"
        , description =
            "A deep dive into all of the problems that DNS has at scale and how Tailscale makes most of those problems go away, with the rest of them being easier in comparison."
        }
      , Link::{
        , url = "https://tailscale.com/blog/tsnet-virtual-private-services/"
        , title = "Virtual private services with tsnet"
        , description =
            "Tailscale lets you connect to your network from anywhere, but you have to set it up on individual computers for it to work. In this article Xe covers how to use tsnet to get all of the goodness of Tailscale in userspace so that you can have your services join your tailnet like they were separate computers."
        }
      ]
    }
