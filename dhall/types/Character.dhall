let PronounSet = ./PronounSet.dhall

in  { Type =
        { name : Text
        , stickerName : Text
        , defaultPose : Text
        , description : Text
        , pronouns : PronounSet.Type
        , stickers : List Text
        }
    , default =
      { name = ""
      , stickerName = ""
      , defaultPose = ""
      , description = ""
      , pronouns = ../pronouns/she.dhall
      , stickers = [] : List Text
      }
    }
