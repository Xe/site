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
      , name = "Astrid Bek"
      , tags =
        [ "dotnet"
        , "java"
        , "python"
        , "typescript"
        , "rust"
        , "c"
        , "cpp"
        , "full-stack"
        , "devops"
        , "linux"
        , "docker"
        , "gamedev"
        ]
      , gitLink = "https://github.com/xSke"
      , twitter = "https://twitter.com/floofstrid"
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
      , name = "Jamie Bliss"
      , tags = [ "python", "devops", "full-stack", "saltstack", "web", "linux" ]
      , gitLink = "https://github.com/astronouth7303"
      , twitter = "https://twitter.com/AstraLuma"
      }
    , Person::{
      , name = "Jaden Weiss"
      , tags =
        [ "go"
        , "wasm"
        , "docker"
        , "kubernetes"
        , "java"
        , "linux"
        , "c"
        , "hardware"
        , "compilers"
        , "machine-learning"
        ]
      , gitLink = "https://github.com/jaddr2line"
      , twitter = "https://twitter.com/jaddr2line"
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
        , "PHP"
        , "google-cloud"
        ]
      , gitLink = "https://github.com/Piyushhbhutoria"
      , twitter = "https://twitter.com/PiyushhB"
      }
    ]
