let xesite = ./types/package.dhall

let Resume = xesite.Resume

let Link = xesite.Link

in  Resume::{
    , hnLinks =
      [ Link::{
        , url = "https://news.ycombinator.com/item?id=29522941"
        , title = "'Open Source' is Broken"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=29167560"
        , title = "The Surreal Horror of PAM"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=27175960"
        , title = "Systemd: The Good Parts"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=26845355"
        , title = "I Implemented /dev/printerfact in Rust"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=25978511"
        , title = "A Model for Identity in Software"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=31390506"
        , title = "Fly.io: The reclaimer of Heroku's magic"
        }
      , Link::{
        , url = "https://news.ycombinator.com/item?id=31149801"
        , title = "Crimes with Go Generics"
        }
      ]
    }
