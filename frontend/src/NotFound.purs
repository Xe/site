module App.NotFound where

import Pux.Html (Html, (#), div, h2, text)

view :: forall state action. state -> Html action
view state =
  div # do
    h2 # text "404 Not Found"
