let Person =
      { Type = { name : Text, tags : List Text, gitLink : Text, twitter : Text }
      , default =
        { name = "", tags = [] : List Text, gitLink = "", twitter = "" }
      }

in  [ Person::{
      , name = "Aisling Fae"
      , tags =
        [ "python", "bash", "kubernetes", "google-cloud", "aws", "devops" ]
      , gitLink = "https://github.com/aislingfae"
      , twitter = "https://twitter.com/aisstern"
      }
    , Person::{
      , name = "Christian Sullivan"
      , tags =
        [ "go"
        , "wasm"
        , "react"
        , "rust"
        , "react-native"
        , "swift"
        , "google-cloud"
        , "aws"
        , "docker"
        , "kubernetes"
        , "istio"
        , "typescript"
        ]
      , gitLink = "https://github.com/euforic"
      , twitter = "https://twitter.com/euforic"
      }
    , Person::{
      , name = "David Roberts"
      , tags =
        [ "ux"
        , "ui"
        , "documentation"
        , "web"
        , "html5"
        , "javascript"
        , "python"
        , "qt"
        , "bash"
        , "front-end"
        , "full-stack"
        , "linux"
        , "embedded"
        , "sql"
        ]
      , gitLink = "https://github.com/ddr0"
      , twitter = "https://twitter.com/DDR_4"
      }
    , Person::{
      , name = "Faizan Jamil"
      , tags =
        [ "java"
        , "c#"
        , "python"
        , "javascript"
        , "typescript"
        , "html"
        , "css"
        , "vue.js"
        , "express.js"
        , "flask"
        , "asp.net core"
        , "razor pages"
        , "ef core"
        , "front-end"
        , "back-end"
        , "full-stack"
        , "linux"
        ]
      , gitLink = "https://github.com/faizjamil"
      , twitter = "N/A"
      }
    , Person::{
      , name = "Joseph Crawley"
      , tags =
        [ "javascript"
        , "react"
        , "csharp"
        , "python"
        , "full-stack"
        , "web"
        , "bash"
        , "linux"
        ]
      , gitLink = "https://github.com/espe-on"
      , twitter = "https://twitter.com/espe_on_"
      }
    , Person::{
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
      , gitLink = "https://github.com/nbraud"
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
      , gitLink = "https://github.com/Prajjwal"
      , twitter = "https://twitter.com/prajjwalsin"
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
      , gitLink = "https://github.com/Piyushhbhutoria"
      , twitter = "https://twitter.com/PiyushhB"
      }
    , Person::{
      , name = "Ryan Casalino"
      , tags =
        [ "golang"
        , "react"
        , "python"
        , "javascript"
        , "aws"
        , "vue"
        , "sql"
        , "ruby"
        , "rails"
        , "flask"
        , "unix"
        ]
      , gitLink = "https://github.com/rjpcasalino"
      , twitter = "N/A"
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
      , gitLink = "https://github.com/dudymas"
      , twitter = "https://twitter.com/dudymas"
      }
    , Person::{
      , name = "Zachary McKee"
      , tags =
        [ "javascript"
        , "django"
        , "react"
        , "postgresql"
        , "firebase"
        , "aws"
        , "python"
        , "csharp"
        , "java"
        , "nginx"
        , "gunicorn"
        ]
      , gitLink = "https://github.com/ZacharyRMcKee"
      , twitter = "N/A"
      }
    , Person::{
      , name = "Muazzam Kazmi"
      , tags = [ "Rust", "C++", "x86assembly", "WinAPI", "Node.js", "React.js" ]
      , gitLink = "https://github.com/muazzamalikazmi"
      , twitter = "N/A"
      }
    , Person::{
      , name = "Jeffin Mathew"
      , tags =
        [ "Python"
        , "routing&switching"
        , "django"
        , "vue"
        , "ansible"
        , "aws"
        , "javascript"
        , "iot"
        ]
      , gitLink = "https://github.com/mjeffin"
      , twitter = "https://twitter.com/mpjeffin"
      }
    , Person::{
      , name = "Nasir Hussain"
      , tags =
        [ "python"
        , "linux"
        , "javascript"
        , "ansible"
        , "nix"
        , "docker&podman"
        , "django"
        , "golang"
        , "rpm packaging"
        ]
      , gitLink = "https://github.com/nasirhm"
      , twitter = "https://twitter.com/_nasirhm_"
      }
    , Person::{
      , name = "Eliot Partridge"
      , tags =
        [ "python"
        , "linux"
        , "typescript"
        , "javascript"
        , "docker"
        , "c#"
        , "dotnet"
        , "php"
        ]
      , gitLink = "https://github.com/BytewaveMLP"
      }
    , Person::{
      , name = "Avi Parshan"
      , tags =
        [ "python"
        , "windows"
        , "javascript"
        , "html"
        , "android"
        , "java"
        , "C#"
        ]
      , gitLink = "https://github.com/avipars"
      , twitter = "https://twitter.com/aviinfinity"
      }
    , Person:: {
      , name = "Tommy Nguyen"
      , tags =
        [ "c++"
        , "linux"
        , "cybersecurity"
        , "privacy"
        , "technical-writing"
        , "web"
        , "google-cloud-platform"
        ]
      , gitLink = "https://github.com/remyabel"
      }
    ]
