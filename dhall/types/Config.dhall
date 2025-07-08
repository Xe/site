let Person = ./Person.dhall

let Character = ./Character.dhall

let Job = ./Job.dhall

let Link = ./Link.dhall

let SeriesDescription = ./SeriesDescription.dhall

let PronounSet = ./PronounSet.dhall

let Resume = ./Resume.dhall

let defaultPort = env:PORT ? 3030

let defaultWebMentionEndpoint =
        env:WEBMENTION_ENDPOINT
      ? "https://mi.within.website/api/webmention/accept"

in  { Type =
        { signalboost : List Person.Type
        , port : Natural
        , clackSet : List Text
        , webMentionEndpoint : Text
        , miToken : Text
        , jobHistory : List Job.Type
        , seriesDescriptions : List SeriesDescription.Type
        , notableProjects : List Link.Type
        , contactLinks : List Link.Type
        , pronouns : List PronounSet.Type
        , characters : List Character.Type
        , resume : Resume.Type
        }
    , default =
      { signalboost = [] : List Person.Type
      , port = defaultPort
      , clackSet = [ "Ashlynn" ]
      , webMentionEndpoint = defaultWebMentionEndpoint
      , miToken = "${env:MI_TOKEN as Text ? ""}"
      , jobHistory = [] : List Job.Type
      , seriesDescriptions = [] : List SeriesDescription.Type
      , notableProjects = [] : List Link.Type
      , contactLinks = [] : List Link.Type
      , pronouns = [] : List PronounSet.Type
      , characters = [] : List Character.Type
      , resume = Resume::{=}
      }
    }
