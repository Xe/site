let xesite = ./types/package.dhall

let Config = xesite.Config

let Link = xesite.Link

let desc = ./seriesDescriptions.dhall

in  Config::{
    , signalboost = ./signalboost.dhall
    , clackSet =
      [ "Ashlynn"
      , "Terry Davis"
      , "Dennis Ritchie"
      , "Stephen Hawking"
      , "John Conway"
      , "Ruth Bader Ginsburg"
      , "Bram Moolenaar"
      , "Grant Imahara"
      , "David Bowie"
      , "Sir Terry Pratchett"
      , "Satoru Iwata"
      , "Kris Nóva"
      , "Joe Armstrong"
      , "Paul Allen"
      , "Kevin Mitnick"
      , "Sir Clive Sinclair"
      , "Matt Trout"
      ]
    , jobHistory = ./jobHistory.dhall
    , seriesDescriptions = desc.descriptions
    , notableProjects =
      [ Link::{
        , url = "https://h.within.lgbt"
        , title = "The h Programming Language"
        , description =
            "An esoteric programming language that compiles to WebAssembly"
        }
      , Link::{
        , url = "https://anubis.techaro.lol"
        , title = "Anubis"
        , description =
            "A anti-scraper bot filter that protects git servers for GNOME, UNESCO, and many other small communities"
        }
      , Link::{
        , url = "https://when-then-zen.christine.website/"
        , title = "When Then Zen"
        , description = "Meditation instructions in plain English"
        }
      , Link::{
        , url = "https://github.com/Xe/x"
        , title = "x"
        , description =
            "A monorepo of my experiments, toy programs and other interesting things of that nature."
        }
      , Link::{
        , url = "https://github.com/Xe/site"
        , title = "Xesite"
        , description = "The backend and templates for this website"
        }
      , Link::{
        , url = "https://github.com/Xe/Xess"
        , title = "Xess"
        , description = "My personal CSS framework"
        }
      ]
    , contactLinks =
      [ Link::{ url = "https://github.com/Xe", title = "GitHub" }
      , Link::{ url = "https://bsky.app/profile/xeiaso.net", title = "Bluesky" }
      , Link::{ url = "https://www.tiktok.com/@xeiaso.1337", title = "TikTok" }
      , Link::{ url = "https://keybase.io/xena", title = "Keybase" }
      , Link::{ url = "https://www.patreon.com/cadey", title = "Patreon" }
      , Link::{ url = "https://www.twitch.tv/princessxen", title = "Twitch" }
      , Link::{ url = "https://pony.social/@cadey", title = "Fediverse" }
      , Link::{ url = "https://t.me/miamorecadenza", title = "Telegram" }
      , Link::{ url = "irc://irc.libera.chat/#xeserv", title = "IRC" }
      , Link::{
        , url =
            "https://signal.me/#eu/Nphi3UKYkj4lgn_HPVFR6wS4VPJ7GRX3htnyHVe8m6XqOPwj8CBJmKnDfTN4mdoX"
        , title = "Signal"
        }
      , Link::{
        , url = "https://www.linkedin.com/in/xe-iaso"
        , title = "LinkedIn"
        }
      ]
    , pronouns = ./pronouns.dhall
    , characters = ./characters.dhall
    , resume = ./resume.dhall
    }
