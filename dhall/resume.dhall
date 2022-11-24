let xesite = ./types/package.dhall

let Resume = xesite.Resume

let Link = xesite.Link

in  Resume::{
    , buzzwords =
      [ "Docker"
      , "Git"
      , "Go"
      , "Rust"
      , "C"
      , "DevOps"
      , "Heroku"
      , "WebAssembly"
      , "Lua"
      , "Mindfulness"
      , "Nix"
      , "NixOS"
      , "HTTP/2"
      , "Ubuntu"
      , "Alpine Linux"
      , "GraphViz"
      , "JavaScript"
      , "TypeScript"
      , "SQLite"
      , "PostgreSQL"
      , "Dudeism"
      , "Technical writing"
      , "Emacs"
      , "Continuous Integration"
      , "Continuous Delivery"
      ]
    , jobs = ./jobHistory.dhall
    , notablePublications = [] : List Link.Type
    }
