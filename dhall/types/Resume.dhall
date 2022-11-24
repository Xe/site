let Location = ./Location.dhall

let Link = ./Link.dhall

let Job = ./Job.dhall

in  { Type =
        { name : Text
        , tagline : Text
        , location : Location.Type
        , buzzwords : List Text
        , jobs : List Job.Type
        , notablePublications : List Link.Type
        }
    , default =
      { name = "Xe Iaso"
      , tagline = "Archmage of Infrastructure"
      , location = Location::{
        , city = "Ottawa"
        , stateOrProvince = "ON"
        , country = "CAN"
        }
      , buzzwords = [] : List Text
      , jobs = [] : List Job.Type
      , notablePublications = [] : List Link.Type
      }
    }
