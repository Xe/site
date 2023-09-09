export default [
  {
    "cdnPath": "talks/vod/2023/03-04-cursorless",
    "date": "2023-03-04",
    "description": "This is a bit of an experimental stream where I attempted to dictate code with [cursorless](https://www.cursorless.org/). When I recorded this stream, I was at minute twenty of playing with this tool. This stream is going to sound really weird, because I am going to be rattling off voice commands that will sound weird at first.\n\nOn this stream, I decided to implement a stable diffusion feature for my CDN XeDN. It replicates the API of the service gravatar, but backed by stable diffusion based off of the hash. There is a terrible bit of code that turns a gravatar hash into a stable diffusion prompt and seed combination.\n\nThis stream covers the following topics:\n\n* Basic navigation with cursorless\n* Data transformations\n* How to execute on terrible ideas\n",
    "slug": "cursorless",
    "tags": [
      "accessibility",
      "voiceControl",
      "go",
      "stableDiffusion"
    ],
    "title": "Shouting at my editor"
  },
  {
    "cdnPath": "talks/vod/2023/02-04-emacs",
    "date": "2023-02-04",
    "description": "  This is a shorter stream where I switched my Emacs config from [Spacemacs](https://spacemacs.org) to a custom configuration I've been prototyping for a year or so that has everything managed with [home-manager](https://nixos.wiki/wiki/Home_Manager) on NixOS. This allows my configuration to be completely managed in configuration and all packages that I depend on can be precompiled at deploy time_, allowing me to run my complicated configurations on less powerful hardware without having to wait for bytecode compilation to happen. Most of the rest of the stream was just going through the motions of actually making the change, and then trying to make some ergonomics changes so that I could use it as a replacement for tmux.\n\n  This stream covers the following topics:\n\n  * Nix/NixOS configuration management\n  * Emacs Lisp programming\n  * Writing custom interactive commands in Emacs\n  * Proving chat wrong about the capabilities of Emacs\n",
    "slug": "emacs-nix",
    "tags": [
      "emacs",
      "nix",
      "lisp",
      "tmux"
    ],
    "title": "Ripping the bandaid off and using Emacs managed by Nix"
  },
  {
    "cdnPath": "talks/vod/2023/01-21-reader-mode",
    "date": "2023-01-21",
    "description": "When you are using reader mode in Firefox, Safari or Google Chrome, the browser rends control of the website's design and renders its own design. This is typically done in order to prevent people's bad design decisions from making webpages unreadable and also to strip away advertisements from content. As a website publisher, I rely on the ability to control the CSS of my blog a lot. This stream covers the research/implementation process for fixing some long-standing issues with the Xesite CSS and making a fix to XeDN so that the site renders acceptably in reader mode.\n\nThis stream covers the following topics:\n\n* Understanding complicated CSS rules and creating fixes for issues with them\n* Using content distribution networks (CDNs) to help reduce page load time for readers\n* Implementing image resizing capabilities into an existing CDN program (XeDN)\n* Design with end-users in mind\n",
    "slug": "reader-mode-css",
    "tags": [
      "css",
      "xedn",
      "imageProcessing",
      "scalability",
      "bugFix"
    ],
    "title": "Fixing Xesite in reader mode and RSS readers"
  },
  {
    "cdnPath": "talks/vod/2023/01-07-pronouns",
    "date": "2023-01-07",
    "description": "In this stream I implemented the [pronouns](https://pronouns.within.lgbt) service and deployed it to the cloud with [fly.io](https://fly.io). This was mostly writing a bunch of data files with [Dhall](https://dhall-lang.org) and then writing a simple Rust program to query that 'database' and then show results based on the results of those queries.\n\nThis stream covers the following topics:\n\n* Starting a new Rust project from scratch with Nix flakes, Axum, and Maud\n* API design for human and machine-paresable outputs\n* DevOps deployment to the cloud via [fly.io](https://fly.io)\n* Writing Terraform code for the pronouns service\n* Building Docker images with Nix flakes and `pkgs.dockerTools.buildLayeredImage`\n* Writing API documentation\n* Writing [the writeup](https://xeiaso.net/blog/pronouns-service) on the service\n",
    "slug": "pronouns-service",
    "tags": [
      "rust",
      "axum",
      "terraform",
      "nix",
      "flyio",
      "docker"
    ],
    "title": "Implementing the Pronouns service in Rust and Axum"
  },
  {
    "cdnPath": "talks/vod/2022/12-31-nguh",
    "date": "2022-12-31",
    "description": "This stream was the last stream of 2022 and focused on modernizing the [hlang](https://xeiaso.net/blog/series/h) compiler. In this stream I reverse-engineered how WebAssembly modules work and wrote my own compiler for a trivial esoteric programming language named h. The existing compiler relied on legacy features of WebAssembly tools that don't work anymore.\n\nThis stream covers the following topics:\n\n* Reverse-engineering the WebAssembly module format based on the specification and other reverse-engineering tools\n* Adapting an existing compiler to output WebAssembly directly\n* Deploying a new service to my NixOS machines in the cloud\n* Building a Nix flake and custom NixOS module to build and deploy the new hlang website\n* Terraform DNS config\n* Writing [the writeup on the new compiler](https://xeiaso.net/blog/hlang-nguh)\n",
    "slug": "hlang-nguh-compiler",
    "tags": [
      "hlang",
      "go",
      "wasm",
      "philosophy",
      "devops",
      "terraform",
      "aws",
      "route53",
      "nixos"
    ],
    "title": "Modernizing hlang with the nguh compiler"
  }
];
