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
      , tags = [ "backend", "bash", "nodejs", "deno", "alpinelinux", "linux", "actuallyautistic" ]
      , links =
        [ Link::{ url = "https://github.com/ajhalili2006", title = "GitHub" }
        , Link::{ url = "https://mau.dev/ajhalili2006", title = "GitLab" }
        , Link::{ url = "https://sr.ht/~ajhalili2006", title = "sourcehut" }
        , Link::{
          , url = "https://tilde.zone/@ajhalili2006"
          , title = "Fediverse"
          }
        , Link::{
          , url = "https://substack.com/@ajhalili2006"
          , title = "Substack"
        }
        , Link::{ url = "https://andreijiroh.eu.org", title = "Website" }
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
    , Person::{
      , name = "Connor Edwards"
      , tags =
        [ "golang"
        , "javascript"
        , "python"
        , "kubernetes"
        , "devops"
        , "sysadmin"
        , "google-cloud"
        , "aws"
        , "terraform"
        , "pulumi"
        , "ansible"
        , "puppet"
        ]
      , links =
        [ Link::{ url = "https://github.com/cedws", title = "GitHub" }
        , Link::{ url = "https://cedwards.xyz", title = "Website" }
        ]
      }
    , Person::{
      , name = "Dillon Baird"
      , tags =
        [ "javascript"
        , "nodejs"
        , "react"
        , "angular"
        , "vue"
        , "python"
        , "redis"
        , "sql"
        , "devops"
        , "sysadmin"
        , "docker"
        , "terraform"
        , "ansible"
        , "ui/ux"
        ]
      , links =
        [ Link::{ url = "https://dillonbaird.io", title = "Website" }
        ]
      }
    , Person::{
      , name = "antlers"
      , tags =
        [ "scheme"
        , "python"
        , "ansible"
        , "devops"
        , "sysadmin"
        , "build-systems"
        , "CI/CD"
        ]
      , links =
        [ Link::{ url = "https://illucid.net", title = "Website" }
        , Link::{ url = "https://oldbytes.space/", title = "Fedi" }
        , Link::{ url = "https://github.com/AutumnalAntlers", title = "GitHub" }
        ]
      }
    , Person::{
      , name = "Redowan Delowar"
      , tags =
        [ "python"
        , "go"
        , "javascript"
        , "nodejs"
        , "ansible"
        , "docker"
        , "aws"
        , "linux"
        , "build-systems"
        , "CI/CD"
        , "data science"
        ]
      , links =
        [ Link::{ url = "https://rednafi.com", title = "Website" }
        , Link::{ url = "https://twitter.com/rednafi", title = "Twitter" }
        , Link::{ url = "https://github.com/rednafi", title = "GitHub" }
        , Link::{ url = "https://www.linkedin.com/in/redowan/", title = "LinkedIn" }
        ]
      }
    ]
