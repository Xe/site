module App.Layout where

import App.BlogIndex as BlogIndex
import App.Counter as Counter
import App.Routes (Route(..))
import DOM (DOM)
import Network.HTTP.Affjax (AJAX)
import Prelude (($), (#), map, pure)
import Pux (EffModel, noEffects, mapEffects, mapState)
import Pux.Html (Html, div, h1, nav, text)
import Pux.Html.Attributes (className, id_, role)
import Pux.Router (link)

data Action
  = Child (Counter.Action)
  | BIChild (BlogIndex.Action)
  | PageView Route

type State =
  { route   :: Route
  , count   :: Counter.State
  , bistate :: BlogIndex.State }

init :: State
init =
  { route: NotFound
  , count: Counter.init
  , bistate: BlogIndex.init }

update :: Action -> State -> EffModel State Action (ajax :: AJAX, dom :: DOM)
update (PageView route) state = noEffects $ state { route = route }
update (BIChild action) state = BlogIndex.update action state.bistate
                              # mapState (state { bistate = _ })
                              # mapEffects BIChild
update (Child action) state = noEffects $ state { count = Counter.update action state.count }

view :: State -> Html Action
view state =
  div
    []
    [ navbar state
    , div
      [ className "container" ]
      [ page state.route state ]
    ]

navbar :: State -> Html Action
navbar state =
  nav
    [ className "light-blue lighten-1", role "navigation" ]
    [ div
      [ className "nav-wrapper container" ]
      [ link "/" [ className "brand-logo", id_ "logo-container" ] [ text "Christine Dodrill" ] ]
    ]

page :: Route -> State -> Html Action
page NotFound _ = h1 [] [ text "not found" ]
page Home state = map Child $ Counter.view state.count
page Resume state = h1 [] [ text "Christine Dodrill" ]
page BlogIndex state = map BIChild $ BlogIndex.view state.bistate
page _ _ = h1 [] [ text "not implemented yet" ]
