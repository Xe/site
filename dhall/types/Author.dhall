{ Type =
    { name : Text
    , handle : Text
    , image : Optional Text
    , url : Optional Text
    , sameAs : List Text
    , jobTitle : Text
    , inSystem : Bool
    }
, default =
  { name = ""
  , handle = ""
  , image = None Text
  , url = None Text
  , sameAs = [] : List Text
  , jobTitle = ""
  , inSystem = False
  }
}
