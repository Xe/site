let Person = ./Person.dhall

let Author = ./Author.dhall

let Job = ./Job.dhall

let defaultPort = env:PORT ? 3030

let defaultWebMentionEndpoint =
        env:WEBMENTION_ENDPOINT
      ? "https://mi.within.website/api/webmention/accept"

in  { Type =
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
