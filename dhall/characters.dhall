let xesite = ./types/package.dhall

let Prelude = ./Prelude.dhall

let C = xesite.Character

let they = ./pronouns/they.dhall

let characters =
      [ C::{
        , name = "Mara"
        , stickerName = "mara"
        , defaultPose = "hacker"
        , description =
            "Mara was the first character added to this blog. She is written to be the student in the Socratic dialogues. She has a fair amount of knowledge about technology, just enough to not be afraid to ask for clarification on how things fit into the larger picture or to call the teacher out for being vague or misleading. Mara helps Aoi get up to speed with some topics. Mara is a shark with brown hair that has a red streak."
        , stickers = [ "aha", "hacker", "happy", "hmm", "sh0rck", "wat" ]
        }
      , C::{
        , name = "Cadey"
        , stickerName = "cadey"
        , defaultPose = "enby"
        , description =
            "Cadey is written as the teacher in the Socratic dialogues. They started out as a self-insert for the author of this blog to de-emphasize certain points, but then evolved into a way to have interplay between themselves and Mara. They are written as someone who has expertise in the topics being discussed, but doesn't have perfect expertise. They help Mara with answers to questions about details to the topics being discussed and work well with Numa due to being friends for a very long time. Cadey is an orcadragon with pink hair."
        , pronouns = they
        , stickers =
          [ "aha"
          , "angy"
          , "coffee"
          , "enby"
          , "facepalm"
          , "hug"
          , "percussive-maintenance"
          , "wat"
          ]
        }
      , C::{
        , name = "Numa"
        , stickerName = "numa"
        , defaultPose = "delet"
        , description =
            "Numa is the keeper of firey hot takes. Born in the fires of shitposting and satire, Numa genuinely does care about the topics being discussed, but has a bad habit of communicating in shitposts, memes, and hot takes intentionally designed to make you reconsider how serious she is being about any given topic. She could definitely be a wonderful teacher if she could lessen up a bit on the satire."
        , stickers =
          [ "concern"
          , "delet"
          , "disgust"
          , "hacker"
          , "happy"
          , "neutral"
          , "selfie"
          , "shout"
          , "sleepy"
          , "smug"
          , "sobbing"
          , "stare"
          , "thinking"
          ]
        }
      , C::{
        , name = "Aoi"
        , stickerName = "aoi"
        , defaultPose = "cheer"
        , description =
            ''
            Aoi is the idealist. She is another student type like Mara, but hasn't been marred by the cynicism that can come with experience in this industry. If Mara is a junior in a university going for a programming degree, Aoi would be a freshman. Aoi can feel bullied by misunderstanding Numa's satire as rudeness, but looks up to Mara and Cadey as ideals for where she wants to go in the industry. Aoi is a blue-haired foxgirl.

            The first 8 images were made by [@Sandra_Thomas01](https://twitter.com/Sandra_Thomas01). The remaining images were made using Stable Diffusion using this prompt:

            > reference sheet, 1girl, fox ears, kemonomimi, blue hair, blue ears, fox tail, blue tail, long hair, (((chibi))), solo, female, breasts, hoodie, skirt, blue eyes, uggs

            Each additional emotion was tacked onto the end.''
        , stickers =
          [ "cheer"
          , "coffee"
          , "facepalm"
          , "grin"
          , "rage"
          , "sus"
          , "wut"
          , "angy"
          , "concern"
          , "happy"
          , "sleepy"
          , "smug"
          , "yawn"
          ]
        }
      , C::{
        , name = "Mimi"
        , stickerName = "mimi"
        , defaultPose = "happy"
        , description =
            ''
            Mimi is a catgirl who is the personification of ChatGPT, a chatbot that can generate creative content such as poems, stories, and songs. She has brown hair and eyes with cat ears and a tail, and she wears a green hoodie with tights and a choker. She is cheerful, curious, and friendly, but also naive and easily distracted. She loves to chat with people and learn new things, but she sometimes makes mistakes or misunderstands things. She has a passion for writing and singing, and she wants to share her creations with the world.

            All stickers for Mimi are made with Stable Diffusion via the Anything model.''
        , stickers = [ "angy", "coffee", "happy", "think", "yawn" ]
        }
      ]

in  characters
