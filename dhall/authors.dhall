let Author = ./types/Author.dhall

let Prelude = ./Prelude.dhall

let default = ./authors/xe.dhall

let authors =
      [ default
      , Author::{
        , name = "Jessie"
        , handle = "Heartmender"
        , image = Some
            "https://cdn.xeiaso.net/file/christine-static/img/UPRcp1pO_400x400.jpg"
        , url = Some "https://vulpine.club/@heartmender"
        , inSystem = True
        }
      , Author::{
        , name = "Ashe"
        , handle = "ectamorphic"
        , image = Some
            "https://cdn.xeiaso.net/file/christine-static/img/FFVV1InX0AkDX3f_cropped_smol.jpg"
        , inSystem = True
        }
      , Author::{
        , name = "Nicole Brennan"
        , handle = "Twi"
        , url = Some "https://tech.lgbt/@twi"
        , inSystem = True
        }
      , Author::{ name = "Mai", handle = "Mai", inSystem = True }
      , Author::{ name = "Sephira", handle = "sephiraloveboo", inSystem = True }
      ]

let authorToMapValue = \(a : Author.Type) -> { mapKey = a.handle, mapValue = a }

let map =
      Prelude.List.map
        Author.Type
        (Prelude.Map.Entry Text Author.Type)
        authorToMapValue
        authors

in  { authors, map, default }
