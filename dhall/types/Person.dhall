let Link = ./Link.dhall

in  { Type = { name : Text, tags : List Text, links : List Link.Type }
    , default =
      { name = "", tags = [] : List Text, links = [] : List Link.Type }
    }
