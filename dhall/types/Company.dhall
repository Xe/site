let Location = ./Location.dhall

in  { Type =
        { name : Text
        , url : Optional Text
        , tagline : Text
        , location : Location.Type
        , defunct : Bool
        }
    , default =
      { name = ""
      , url = None Text
      , tagline = ""
      , location = Location::{=}
      , defunct = False
      }
    }
