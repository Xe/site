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
      , "HTTP/2"
      , "Ubuntu"
      , "Alpine Linux"
      , "Gentoo"
      , "Bot filtering"
      , "AI"
      , "Large Language Models"
      , "ChatGPT"
      , "GPT4"
      , "Stable Diffusion"
      , "SDXL"
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
        , url = "https://xeiaso.net/talks/2025/ai-chatbot-friends/"
        , title = "Affording your AI chatbot friends"
        , description =
            "A presentation focused at technical experts about the moving parts involved with AI model hosting and how it's not actually expensive to run if you manage the scope and scale correctly."
        }
      , Link::{
        , url = "https://blog.heroku.com/how-to-make-progressive-web-app"
        , title = "How to Make a Progressive Web App From Your Existing Website"
        , description =
            "An article summarizing how easy it is to make a webpage into an installable Progressive Web App using APIs available in the most commonly used browsers."
        }
      , Link::{
        , url = "https://www.tigrisdata.com/blog/training-any-cloud/"
        , title = "Training with Big Data on Any Cloud"
        , description =
            "A detailed breakdown of how to use Object Storage to help train Large Language Models / AI models on any cloud provider that has an internet connection."
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
      ]
    }
