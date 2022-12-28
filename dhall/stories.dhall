let xesite = ./types/package.dhall

let Prelude = ./Prelude.dhall

let Story = xesite.Story

let Step = xesite.StoryStep

let stories =
      [ Story::{
        , name = "unwrapped-2022"
        , steps =
          [ Step::{
            , file = "title"
            , title = "Xesite 2022"
            , text = "The last year all wrapped up!"
            }
          , Step::{
            , file = "posts"
            , title = "Blogposts"
            , text = "Xe wrote \$NUMBER posts this year!"
            }
          , Step::{
            , file = "talks"
            , title = "Talks"
            , text = "Xe gave 4 talks this year! Which one was your favorite?"
            }
          , Step::{
            , file = "commits"
            , title = "340 commits"
            , text = "There were 340 commits to Xesite this year!"
            }
          ]
        }
      ]

let storyToMapValue =
      \(story : Story.Type) -> { mapKey = story.name, mapValue = story }

let map =
      Prelude.List.map
        Story.Type
        (Prelude.Map.Entry Text Story.Type)
        storyToMapValue
        stories

in  { stories, map }
