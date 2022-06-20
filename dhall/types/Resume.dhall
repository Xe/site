let Location = ./Location.dhall

let Link = ./Link.dhall

in  { Type =
        { name : Text
        , tagline : Text
        , location : Location.Type
        , hnLinks : List Link.Type
        }
    , default =
      { name = "Xe Iaso"
      , tagline = "Archmage of Infrastructure"
      , location = Location::{
        , city = "Ottawa"
        , stateOrProvince = "ON"
        , country = "CAN"
        }
      , hnLinks = [] : List Link.Type
      }
    }
