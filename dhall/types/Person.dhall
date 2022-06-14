{ Type =
    { name : Text
    , tags : List Text
    , gitLink : Optional Text
    , twitter : Optional Text
    }
, default =
  { name = "", tags = [] : List Text, gitLink = None Text, twitter = None Text }
}
