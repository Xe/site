module App.Layout where

import App.BlogEntry as BlogEntry
import App.BlogIndex as BlogIndex
import App.Counter as Counter
import App.Routes (Route(..))
import Control.Monad.RWS (state)
import DOM (DOM)
import Network.HTTP.Affjax (AJAX)
import Prelude (($), (#), map, pure)
import Pux (EffModel, noEffects, mapEffects, mapState)
import Pux.Html (Html, div, h1, li, nav, text, ul)
import Pux.Html.Attributes (classID, className, id_, role)
import Pux.Router (link)

data Action
  = Child (Counter.Action)
  | BIChild (BlogIndex.Action)
  | BEChild (BlogEntry.Action)
  | PageView Route

type State =
  { route   :: Route
  , count   :: Counter.State
  , bistate :: BlogIndex.State
  , bestate :: BlogEntry.State }

init :: State
init =
  { route: NotFound
  , count: Counter.init
  , bistate: BlogIndex.init
  , bestate: BlogEntry.init }

update :: Action -> State -> EffModel State Action (ajax :: AJAX, dom :: DOM)
update (PageView route) state = routeEffects route $ state { route = route }
update (BIChild action) state = BlogIndex.update action state.bistate
                              # mapState (state { bistate = _ })
                              # mapEffects BIChild
update (BEChild action) state = BlogEntry.update action state.bestate
                              # mapState (state { bestate = _ })
                              # mapEffects BEChild
update (Child action) state = noEffects $ state { count = Counter.update action state.count }

routeEffects :: Route -> State -> EffModel State Action (dom :: DOM, ajax :: AJAX)
routeEffects BlogIndex state = { state: state
                               , effects: [ pure BlogIndex.RequestPosts ] } # mapEffects BIChild
routeEffects (BlogPost page) state = { state: state { bestate = BlogEntry.init { name = page } }
                              , effects: [ pure BlogEntry.RequestPost ] } # mapEffects BEChild
routeEffects _ state = noEffects $ state

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
    [ className "pink lighten-1", role "navigation" ]
    [ div
      [ className "nav-wrapper container" ]
      [ link "/" [ className "brand-logo", id_ "logo-container" ] [ text "Christine Dodrill" ]
      , ul
        [ className "right hide-on-med-and-down" ]
        [ li [] [ link "/blog" [] [ text "Blog" ] ]
        , li [] [ link "/projects" [] [ text "Projects" ] ]
        , li [] [ link "/resume" [] [ text "Resume" ] ]
        , li [] [ link "/contact" [] [ text "Contact" ] ]
        ]
      ]
    ]

page :: Route -> State -> Html Action
page NotFound _ = h1 [] [ text "not found" ]
page Home state = map Child $ Counter.view state.count
page Resume state = h1 [] [ text "Christine Dodrill" ]
page BlogIndex state = map BIChild $ BlogIndex.view state.bistate
page (BlogPost _) state = map BEChild $ BlogEntry.view state.bestate
page _ _ = h1 [] [ text "not implemented yet" ]
