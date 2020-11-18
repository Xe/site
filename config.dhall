let Person =
      { Type = { name : Text, tags : List Text, gitLink : Text, twitter : Text }
      , default =
        { name = "", tags = [] : List Text, gitLink = "", twitter = "" }
      }

let defaultPort = env:PORT ? 3030

let defaultWebMentionEndpoint =
        env:WEBMENTION_ENDPOINT
      ? "https://mi.within.website/api/webmention/accept"

let Config =
      { Type =
          { signalboost : List Person.Type
          , port : Natural
          , clackSet : List Text
          , resumeFname : Text
          , webMentionEndpoint : Text
          }
      , default =
        { signalboost = [] : List Person.Type
        , port = defaultPort
        , clackSet = [ "Ashlynn" ]
        , resumeFname = "./static/resume/resume.md"
        , webMentionEndpoint = defaultWebMentionEndpoint
        }
      }

in  Config::{
    , signalboost = ./signalboost.dhall
    , clackSet = [ "Ashlynn", "Terry Davis", "Dennis Ritchie" ]
    }
