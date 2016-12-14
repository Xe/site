module App.Counter where

import Prelude ((+), (-), const, show)
import Pux.Html (Html, a, br, div, span, text)
import Pux.Html.Attributes (className, href)
import Pux.Html.Events (onClick)

data Action = Increment | Decrement

type State = Int

init :: State
init = 0

update :: Action -> State -> State
update Increment state = state + 1
update Decrement state = state - 1

view :: State -> Html Action
view state =
  div
    [ className "row" ]
    [ div
      [ className "col s4 offset-s4" ]
      [ div
        [ className "card blue-grey darken-1" ]
        [ div
          [ className "card-content white-text" ]
          [ span [ className "card-title" ] [ text "Counter" ]
          , br [] []
          , span [] [ text (show state) ]
          ]
        , div
          [ className "card-action" ]
          [ a [ onClick (const Increment), href "#" ] [ text "Increment" ]
          , a [ onClick (const Decrement), href "#" ] [ text "Decrement" ]
          ]
        ]
      ]
    ]
