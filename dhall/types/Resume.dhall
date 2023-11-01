let Location = ./Location.dhall

let Link = ./Link.dhall

let Job = ./Job.dhall

in  { Type =
        { name : Text
        , tagline : Text
        , location : Location.Type
        , buzzwords : List Text
        , notablePublications : List Link.Type
        }
    , default =
      { name = "Xe Iaso"
      , tagline = "Senior Technophilosopher"
      , location = Location::{
        , city = "Ottawa"
        , stateOrProvince = "ON"
        , country = "CAN"
        }
      , buzzwords = [] : List Text
      , notablePublications = [] : List Link.Type
      }
    }
