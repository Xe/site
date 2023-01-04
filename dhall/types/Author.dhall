let PronounSet = ./PronounSet.dhall

in  { Type =
        { name : Text
        , handle : Text
        , image : Optional Text
        , url : Optional Text
        , sameAs : List Text
        , jobTitle : Text
        , inSystem : Bool
        , pronouns : PronounSet.Type
        }
    , default =
      { name = ""
      , handle = ""
      , image = None Text
      , url = None Text
      , sameAs = [] : List Text
      , jobTitle = ""
      , inSystem = False
      , pronouns = ../pronouns/she.dhall
      }
    }
