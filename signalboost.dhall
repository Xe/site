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
      , name = "Jamie Bliss"
      , tags = [ "python", "devops", "full-stack", "saltstack", "web", "linux" ]
      , gitLink = "https://github.com/astronouth7303"
      , twitter = "https://twitter.com/AstraLuma"
      }
    ]
