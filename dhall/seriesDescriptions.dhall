let xesite = ./types/package.dhall

let Prelude = ./Prelude.dhall

let Desc = xesite.SeriesDescription

let descriptions
    : List Desc.Type
    = [ Desc::{
        , name = "colemak"
        , details =
            "My efforts at learning to type with colemak instead of qwerty."
        }
      , Desc::{
        , name = "conlangs"
        , details =
            "Information about constructed languages I've attempted to make."
        }
      , Desc::{
        , name = "CVE"
        , details =
            "Vulnerability information and my responsible disclosures of said vulnerabilities."
        }
      , Desc::{
        , name = "dreams"
        , details = "My attempts to write about my dreams"
        }
      , Desc::{
        , name = "ethereum"
        , details = "Why Ethereum doesn't work in the real world."
        }
      , Desc::{
        , name = "flightJournal"
        , details = "My Gemini posts, copied over to my blog for safe-keeping."
        }
      , Desc::{
        , name = "freenode"
        , details =
            "My lamentations about the collapse of the IRC network freenode."
        }
      , Desc::{
        , name = "get-going"
        , details = "Tutorials for the Go programming language."
        }
      , Desc::{
        , name = "h"
        , details = "Evolution of the h human/programming language."
        }
      , Desc::{
        , name = "homelabV2"
        , details = "My second attempt at a homelab via Fedora CoreOS and k3s."
      }
      , Desc::{
        , name = "howto"
        , details = "Instructions on how to do various things."
        }
      , Desc::{ name = "keeb", details = "Keyboard reviews." }
      , Desc::{
        , name = "magick"
        , details =
            "Writeups on things I've learned on my trip through chaos magick."
        }
      , Desc::{
        , name = "malto"
        , details = "Stories from my constructed world Malto."
        }
      , Desc::{
        , name = "medium-archive"
        , details = "Articles from my attempt at making a Medium blog long ago."
        }
      , Desc::{
        , name = "nix-flakes"
        , details =
            "Instructions on how to use Nix flakes, a new way to use Nix in a more reproducible way."
        }
        , Desc::{
          name = "no-way-to-prevent-this",
          details = "Articles about the futility of preventing memory safety vulnerabilities."
        }
      , Desc::{ name = "nixos", details = "Nix." }
      , Desc::{
        , name = "olin"
        , details = "My attempts at running WebAssembly on the server."
        }
      , Desc::{ name = "plt", details = "The saga of plt." }
      , Desc::{
        , name = "recipes"
        , details =
            "Recipes for use in your kitchen presented in a no-bullshit fashion."
        }
      , Desc::{
        , name = "reconlangmo"
        , details =
            "More details on how I tried to make a language named L'ewa."
        }
      , Desc::{
        , name = "reviews"
        , details = "Reviews on various tech or media properties."
        }
      , Desc::{
        , name = "revueBackup"
        , details = "My revue posts converted to xesite posts."
        }
      , Desc::{ name = "rust", details = "Rust." }
      , Desc::{
        , name = "short-story"
        , details = "Flash fiction stories I've written over the years"
        }
      , Desc::{
        , name = "site-update"
        , details =
            "Updates on this website. These articles will contain details on how this website is changed, new things I'm cooking up with it or more."
        }
      , Desc::{
        , name = "site-to-site-wireguard"
        , details = "Instructions on setting up your own VPN with WireGuard."
        }
      , Desc::{
        , name = "spellblade"
        , details = "Sections of my web novel Spellblade."
        }
      , Desc::{
        , name = "techaro"
        , details = "Stories about Techaro, the imaginary technology startup."
      }
      , Desc::{
        , name = "templeos"
        , details =
            "Articles about TempleOS, a public domain operating system for AMD64 computers."
        }
      , Desc::{
        , name = "thesource"
        , details = "Expansions for my TTRPG The Source."
        }
      , Desc::{
        , name = "twitter"
        , details =
            "Lamentations on the death of Twitter, a microblogging community."
        }
      , Desc::{ name = "v", details = "The V programming language." }
      , Desc::{ name = "vtuber", details = "My experience as a VTuber." }
      , Desc::{
        , name = "waifud"
        , details = "Information about my VM manager waifud."
        }
      , Desc::{
        , name = "when-then-zen"
        , details = "Meditation information sans bullshit."
        }
      , Desc::{
        , name = "xeact"
        , details = "Xeact, the best frontend femtoframework in the galaxy"
        }
      ]

let descToMapValue =
      \(desc : Desc.Type) -> { mapKey = desc.name, mapValue = desc.details }

let map =
      Prelude.List.map
        Desc.Type
        (Prelude.Map.Entry Text Text)
        descToMapValue
        descriptions

in  { descriptions, map }
