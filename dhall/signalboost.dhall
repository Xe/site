let Link = ./types/Link.dhall

let Person = ./types/Person.dhall

in  [ Person::{
      , name = "bri recchia"
      , tags = 
        [ "virtualization"
        , "linux"
        , "generalist"
        , "containers"
        , "networking"
        , "bgp"
        , "dns"
        , "bash"
        , "python"
        , "rust"
        , "devops"
        , "systems administration"
        ]
      , links = [ Link::{ url = "https://github.com/b-/", title = "Github" } ]
      }
    , Person::{
      , name = "Evan Pratten"
      , tags = 
        [ "rust"
        , "linux"
        , "docker"
        , "full-stack"
        , "ipv6"
        , "bgp"
        , "computer-graphics"
        ]
      , links = [ Link::{ url = "https://ewpratten.com", title = "Website" } ]
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
        , "SRE / system administration"
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
    , Person::{
      , name = "Ryan Heywood"
      , tags =
        [ "backend"
        , "rust"
        , "linux"
        , "docker"
        , "kubernetes"
        , "rook+ceph"
        , "calico"
        , "aws"
        , "terraform"
        , "ansible"
        , "nodejs"
        , "react"
        , "python"
        ]
      , links =
        [ Link::{ url = "https://github.com/RyanSquared", title = "GitHub" }
        , Link::{ url = "https://tilde.zone/@ryan", title = "Fediverse" }
        , Link::{ url = "https://ryansquared.pub", title = "Website" }
        ]
      }
    , Person::{
      , name = "Aurorah Harmony"
      , tags =
        [ "full-stack"
        , "typescript"
        , "javascript"
        , "vue"
        , "nodejs"
        , "php"
        , "laravel"
        , "docker"
        , "linux"
        ]
      , links =
        [ Link::{ url = "https://github.com/aurorahHarmony", title = "GitHub" }
        , Link::{ url = "https://pony.social/@aurorahHarmony", title = "Fediverse" }
        , Link::{ url = "https://itsaury.net", title = "Website" }
        ]
      }
   , Person::{
      , name = "Caramel Drop"
      , tags =
        [ "full-stack"
        , "c#"
        , "javascript"
        , "docker"
        , "nodejs"
        , "linux"
        , "system administration"
        , "bash"
        ]
      , links =
        [ Link::{ url = "https://github.com/caramelpony", title = "GitHub" }
        , Link::{ url = "https://caramel.horse", title = "Website" }
        ]
      }
    , Person::{
      , name = "Satvik Jagannath"
      , tags =
        [ "backend"
        , "nodejs"
        , "linux"
        , "javascript"
        , "nextjs"
        , "aws"
        , "nodejs"
        , "react"
        , "python"
        ]
      , links =
        [ Link::{ url = "https://github.com/bozzmob", title = "GitHub" }
        , Link::{ url = "https://twitter.com/@bozzmob", title = "Twitter" }
        , Link::{ url = "https://debugpointer.com", title = "Website" }
        ]
      }
    ]
