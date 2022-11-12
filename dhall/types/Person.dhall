{ Type =
    { name : Text
    , tags : List Text
    , gitLink : Optional Text
    , twitter : Optional Text
    , linkedin : Optional Text
    , fediverse : Optional Text
    , cover_letter : Optional Text
    , website: Optional Text
    }
, default =
  { name = "", tags = [] : List Text, gitLink = None Text, twitter = None Text, linkedin : None Text, fediverse: None Text, cover_letter : None Text, website: None Text }
}
