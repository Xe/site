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
import Pux.Html (Html, a, code, div, h1, h3, h4, li, nav, p, pre, text, ul)
import Pux.Html.Attributes (classID, className, id_, role, href)
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
update _ state              = noEffects $ state

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

contact :: Html Action
contact =
  div
    [ className "row" ]
    [ div
      [ className "col s6" ]
      [ h3 [] [ text "Email" ]
      , div [ className "email" ] [ text "me@christine.website" ]
      , p []
        [ text "My GPG fingerprint is "
        , code [] [ text "799F 9134 8118 1111" ]
        , text ". If you get an email that appears to be from me and the signature does not match that fingerprint, it is not from me. You may download a copy of my public key "
        , a [ href "/static/gpg.pub" ] [ text "here" ]
        , text "."
        ]
      ]
    , div
      [ className "col s6" ]
      [ h3 [] [ text "Other Information" ]
      , p []
        [ text "To send me donations, my bitcoin address is "
        , code [] [ text "1Gi2ZF2C9CU9QooH8bQMB2GJ2iL6shVnVe" ]
        , text "."
        ]
      , div []
        [ h4 [] [ text "IRC" ]
        , p [] [ text "I am on many IRC networks. On Freenode I am using the nick Xe but elsewhere I will use the nick Xena or Cadey." ]
        ]
      , div []
        [ h4 [] [ text "Telegram" ]
        , a [ href "https://telegram.me/miamorecadenza" ] [ text "@miamorecadenza" ]
        ]
      , div []
        [ h4 [] [ text "Discord" ]
        , pre [] [ text "Cadey~#1932" ]
        ]
      ]
    ]

page :: Route -> State -> Html Action
page NotFound _ = h1 [] [ text "not found" ]
page Home state = map Child $ Counter.view state.count
page Resume state = h1 [] [ text "Christine Dodrill" ]
page BlogIndex state = map BIChild $ BlogIndex.view state.bistate
page (BlogPost _) state = map BEChild $ BlogEntry.view state.bestate
page ContactPage _ = contact
page _ _ = h1 [] [ text "not implemented yet" ]
