let Person = ./dhall/types/Person.dhall

let Author = ./dhall/types/Author.dhall

let Job = ./dhall/types/Job.dhall

let defaultPort = env:PORT ? 3030

let defaultWebMentionEndpoint =
        env:WEBMENTION_ENDPOINT
      ? "https://mi.within.website/api/webmention/accept"

let Config =
      { Type =
          { signalboost : List Person.Type
          , authors : List Author.Type
          , port : Natural
          , clackSet : List Text
          , resumeFname : Text
          , webMentionEndpoint : Text
          , miToken : Text
          , jobHistory : List Job.Type
          }
      , default =
        { signalboost = [] : List Person.Type
        , authors = [] : List Author.Type
        , port = defaultPort
        , clackSet = [ "Ashlynn" ]
        , resumeFname = "./static/resume/resume.md"
        , webMentionEndpoint = defaultWebMentionEndpoint
        , miToken = "${env:MI_TOKEN as Text ? ""}"
        , jobHistory = [] : List Job.Type
        }
      }

in  Config::{
    , signalboost = ./dhall/signalboost.dhall
    , authors = ./dhall/authors.dhall
    , clackSet =
      [ "Ashlynn", "Terry Davis", "Dennis Ritchie", "Steven Hawking" ]
    , jobHistory = ./dhall/jobHistory.dhall
    }
