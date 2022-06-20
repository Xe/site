let xesite = ./dhall/types/package.dhall

let Config = xesite.Config

in  Config::{
    , signalboost = ./dhall/signalboost.dhall
    , authors = ./dhall/authors.dhall
    , clackSet =
      [ "Ashlynn", "Terry Davis", "Dennis Ritchie", "Steven Hawking" ]
    , jobHistory = ./dhall/jobHistory.dhall
    }
