let Link = ./Link.dhall

in  { Type =
        { title : Text
        , slug : Text
        , date : Text
        , description : Text
        , cdnPath : Text
        , tags : List Text
        }
    , default =
      { title = ""
      , slug = ""
      , date = ""
      , description = ""
      , cdnPath = ""
      , tags = [] : List Text
      }
    }
