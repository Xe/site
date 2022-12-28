let Step = ./StoryStep.dhall

in  { Type = { name : Text, steps : List Step.Type }
    , default = { name = "", steps = [] : List Step.Type }
    }
