{ Type =
    { name : Text
    , handle : Text
    , picUrl : Optional Text
    , link : Optional Text
    , twitter : Optional Text
    , default : Bool
    , inSystem : Bool
    }
, default =
  { name = ""
  , handle = ""
  , picUrl = None Text
  , link = None Text
  , twitter = None Text
  , default = False
  , inSystem = False
  }
}
