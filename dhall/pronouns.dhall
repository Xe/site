let xesite = ./types/package.dhall

let Pronouns = xesite.PronounSet

in  [ Pronouns::{
      , nominative = "xe"
      , accusative = "xer"
      , possessiveDeterminer = "xer"
      , possessive = "xers"
      , reflexive = "xerself"
      , singular = True
      }
    , ./pronouns/they.dhall
    , ./pronouns/she.dhall
    ]
