let xesite = ./types/package.dhall

let Config = xesite.Config

let Link = xesite.Link

let authors = ./authors.dhall

let desc = ./seriesDescriptions.dhall

in  Config::{
    , signalboost = ./signalboost.dhall
    , authors = authors.map
    , defaultAuthor = authors.default
    , clackSet =
      [ "Ashlynn", "Terry Davis", "Dennis Ritchie", "Steven Hawking" ]
    , jobHistory = ./jobHistory.dhall
    , seriesDescriptions = desc.descriptions
    , seriesDescMap = desc.map
    , notableProjects =
      [ Link::{
        , url = "https://github.com/PonyvilleFM/aura"
        , title = "Aura"
        , description = "PonyvilleFM live DJ recording bot"
        }
      , Link::{
        , url = "https://h.within.lgbt"
        , title = "The h Programming Language"
        , description =
            "An esoteric programming language that compiles to WebAssembly"
        }
      , Link::{
        , url = "https://github.com/Xe/olin"
        , title = "Olin"
        , description = "WebAssembly on the server"
        }
      , Link::{
        , url = "https://printerfacts.cetacean.club/"
        , title = "Printer Facts"
        , description = "Useful facts about printers"
        }
      , Link::{
        , url = "https://github.com/Xe/waifud"
        , title = "waifud"
        , description = "A VM manager for my homelab cluster"
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
        , url = "https://github.com/Xe/Xeact"
        , title = "Xeact"
        , description =
            "My personal JavaScript femtoframework for high productivity development"
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
      , Link::{ url = "https://keybase.io/xena", title = "Keybase" }
      , Link::{ url = "https://www.patreon.com/cadey", title = "Patreon" }
      , Link::{ url = "https://www.twitch.tv/princessxen", title = "Twitch" }
      , Link::{ url = "https://pony.social/@cadey", title = "Fediverse" }
      , Link::{ url = "https://t.me/miamorecadenza", title = "Telegram" }
      , Link::{ url = "irc://irc.libera.chat/#xeserv", title = "IRC" }
      , Link::{ url = "https://staging.bsky.app/profile/xeiaso.net", title = "Bluesky" }
      ]
    , pronouns = ./pronouns.dhall
    , characters = ./characters.dhall
    , vods = ./streamVOD.dhall
    }
