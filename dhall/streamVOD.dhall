let xesite = ./types/package.dhall

let VOD = xesite.StreamVOD

in  [ VOD::{
      , title = "Fixing Xesite in reader mode and RSS readers"
      , slug = "reader-mode-css"
      , description =
          ''
          When you are using reader mode in Firefox, Safari or Google Chrome, the browser rends control of the website's design and renders its own design. This is typically done in order to prevent people's bad design decisions from making webpages unreadable and also to strip away advertisements from content. As a website publisher, I rely on the ability to control the CSS of my blog a lot. This stream covers the research/implementation process for fixing some long-standing issues with the Xesite CSS and making a fix to XeDN so that the site renders acceptably in reader mode.

          This stream covers the following topics:

          * Understanding complicated CSS rules and creating fixes for issues with them
          * Using content distribution networks (CDNs) to help reduce page load time for readers
          * Implementing image resizing capabilities into an existing CDN program (XeDN)
          * Design with end-users in mind
          ''
      , date = "2022-01-21"
      , cdnPath = "talks/vod/2023/01-21-reader-mode"
      , tags = [ "css", "xedn", "imageProcessing", "scalability", "bugFix" ]
      }
    , VOD::{
      , title = "Implementing the Pronouns service in Rust and Axum"
      , slug = "pronouns-service"
      , description =
          ''
          In this stream I implemented the [pronouns](https://pronouns.within.lgbt) service and deployed it to the cloud with [fly.io](https://fly.io). This was mostly writing a bunch of data files with [Dhall](https://dhall-lang.org) and then writing a simple Rust program to query that 'database' and then show results based on the results of those queries.

          This stream covers the following topics:

          * Starting a new Rust project from scratch with Nix flakes, Axum, and Maud
          * API design for human and machine-paresable outputs
          * DevOps deployment to the cloud via [fly.io](https://fly.io)
          * Writing Terraform code for the pronouns service
          * Building Docker images with Nix flakes and `pkgs.dockerTools.buildLayeredImage`
          * Writing API documentation
          * Writing [the writeup](https://xeiaso.net/blog/pronouns-service) on the service
          ''
      , date = "2022-01-07"
      , cdnPath = "talks/vod/2023/01-07-pronouns"
      , tags = [ "rust", "axum", "terraform", "nix", "flyio", "docker" ]
      }
    , VOD::{
      , title = "Modernizing hlang with the nguh compiler"
      , slug = "hlang-nguh-compiler"
      , description =
          ''
          This stream was the last stream of 2022 and focused on modernizing the [hlang](https://xeiaso.net/blog/series/h) compiler. In this stream I reverse-engineered how WebAssembly modules work and wrote my own compiler for a trivial esoteric programming language named h. The existing compiler relied on legacy features of WebAssembly tools that don't work anymore.

          This stream covers the following topics:

          * Reverse-engineering the WebAssembly module format based on the specification and other reverse-engineering tools
          * Adapting an existing compiler to output WebAssembly directly
          * Deploying a new service to my NixOS machines in the cloud
          * Building a Nix flake and custom NixOS module to build and deploy the new hlang website
          * Terraform DNS config
          * Writing [the writeup on the new compiler](https://xeiaso.net/blog/hlang-nguh)
          ''
      , date = "2022-12-31"
      , cdnPath = "talks/vod/2022/12-31-nguh"
      , tags =
        [ "hlang"
        , "go"
        , "wasm"
        , "philosophy"
        , "devops"
        , "terraform"
        , "aws"
        , "route53"
        , "nixos"
        ]
      }
    ]
