let Author = ./types/Author.dhall

in  [ Author::{
      , name = "Xe Iaso"
      , handle = "xe"
      , picUrl = Some "/static/img/avatar.png"
      , link = Some "https://christine.website"
      , twitter = Some "theprincessxena"
      , default = True
      , inSystem = True
      }
    , Author::{
      , name = "Jessie"
      , handle = "Heartmender"
      , picUrl = Some
          "https://cdn.christine.website/file/christine-static/img/UPRcp1pO_400x400.jpg"
      , link = Some "https://heartmender.writeas.com"
      , twitter = Some "BeJustFine"
      , inSystem = True
      }
    , Author::{
      , name = "Ashe"
      , handle = "ectamorphic"
      , picUrl = Some
          "https://cdn.christine.website/file/christine-static/img/FFVV1InX0AkDX3f_cropped_smol.jpg"
      , inSystem = True
      }
    , Author::{ name = "Nicole", handle = "Twi", inSystem = True }
    , Author::{ name = "Mai", handle = "Mai", inSystem = True }
    , Author::{ name = "Sephira", handle = "Sephie", inSystem = True }
    ]
