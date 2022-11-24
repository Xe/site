let Link = ./types/Link.dhall

let Person = ./types/Person.dhall

in  [ Person::{
      , name = "nicoo"
      , tags =
        [ "cryptography"
        , "Debian"
        , "distributed systems"
        , "embedded"
        , "nix"
        , "rust"
        , "privacy"
        , "security"
        , "SDR"
        ]
      , links =
        [ Link::{ url = "https://github.com/nbraud", title = "GitHub" } ]
      }
    , Person::{
      , name = "Prajjwal Singh"
      , tags =
        [ "full-stack javascript"
        , "ruby"
        , "rails"
        , "vuejs"
        , "emberjs"
        , "golang"
        , "linux"
        , "docker"
        , "google-cloud"
        , "typescript"
        ]
      , links =
        [ Link::{ url = "https://github.com/Prajjwal", title = "GitHub" }
        , Link::{ url = "https://twitter.com/prajjwalsin", title = "Twitter" }
        ]
      }
    , Person::{
      , name = "Piyushh Bhutoria"
      , tags =
        [ "golang"
        , "react-native"
        , "full-stack developer"
        , "javascript"
        , "php"
        , "google-cloud"
        ]
      , links =
        [ Link::{ url = "https://github.com/Piyushhbhutoria", title = "GitHub" }
        , Link::{ url = "https://twitter.com/PiyushhB", title = "twitter" }
        ]
      }
    , Person::{
      , name = "Jeremy White"
      , tags =
        [ "kubernetes"
        , "golang"
        , "devops"
        , "python"
        , "rust"
        , "csharp"
        , "angular"
        , "react"
        , "javascript"
        , "saltstack"
        , "aws"
        , "google-cloud"
        , "azure"
        ]
      , links =
        [ Link::{ url = "https://github.com/dudymas", title = "GitHub" }
        , Link::{ url = "https://twitter.com/dudymas", title = "Twitter" }
        ]
      }
    , Person::{
      , name = "Violet White"
      , tags =
        [ "c++"
        , "linux"
        , "python"
        , "javascript"
        , "sql"
        , "lisps"
        , "rust"
        , "backend"
        ]
      , links =
        [ Link::{ url = "https://github.com/epsilon-phase", title = "GitHub" } ]
      }
    , Person::{
      , name = "Henri Shustak"
      , tags =
        [ "backend"
        , "generalist"
        , "documentation"
        , "support"
        , "electronics"
        , "javascript"
        , "python"
        , "ruby"
        , "bash"
        , "sh"
        , "fish"
        , "zsh"
        , "tsch"
        , "software"
        , "full-stack"
        , "linux"
        , "R&D"
        , "SRE / system adminsitration"
        ]
      , links =
        [ Link::{ url = "https://github.com/henri", title = "GitHub" }
        , Link::{ url = "https://twitter.com/henri_shustak", title = "Twitter" }
        ]
      }
    , Person::{
      , name = "Andrei Jiroh Halili"
      , tags = [ "backend", "bash", "nodejs", "deno", "alpinelinux", "linux" ]
      , links =
        [ Link::{ url = "https://github.com/ajhalili2006", title = "GitHub" }
        , Link::{ url = "https://twitter.com/Kuys_Potpot", title = "Twitter" }
        , Link::{
          , url = "https://tilde.zone/@ajhalili2006"
          , title = "Fediverse"
          }
        , Link::{ url = "https://ajhalili2006.bio.link", title = "Website" }
        ]
      }
    ]
