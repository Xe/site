module App.BlogIndex where

import Control.Monad.Aff (attempt)
import Data.Argonaut (class DecodeJson, decodeJson, (.?))
import Data.Either (Either(Left, Right), either)
import DOM (DOM)
import Network.HTTP.Affjax (AJAX, get)
import Prelude (($), bind, map, const, show, (<>), pure, (<<<))
import Pux (EffModel, noEffects)
import Pux.Html (Html, br, div, h1, ol, li, button, text, span, p)
import Pux.Html.Attributes (key, className, id_)
import Pux.Html.Events (onClick)
import Pux.Router (link)

data Action = RequestPosts
            | ReceivePosts (Either String Posts)

type State =
  { posts  :: Posts
  , status :: String }

data Post = Post
  { title   :: String
  , link    :: String
  , summary :: String
  , date    :: String }

type Posts = Array Post

instance decodeJsonPost :: DecodeJson Post where
  decodeJson json = do
    obj <- decodeJson json
    title <- obj .? "title"
    link <- obj .? "link"
    summ <- obj .? "summary"
    date <- obj .? "date"
    pure $ Post { title: title, link: link, summary: summ, date: date }

init :: State
init =
  { posts: []
  , status: "Loading..." }

update :: Action -> State -> EffModel State Action (ajax :: AJAX, dom :: DOM)
update (ReceivePosts (Left err)) state =
  noEffects $ state { status = ("error: " <> err) }
update (ReceivePosts (Right posts)) state =
  noEffects $ state { posts = posts, status = "Posts" }
update RequestPosts state =
  { state: state { status = "Fetching posts..." }
  , effects: [ do
      res <- attempt $ get "/api/blog/posts"
      let decode r = decodeJson r.response :: Either String Posts
      let posts = either (Left <<< show) decode res
      pure $ ReceivePosts posts
    ]
  }

post :: Post -> Html Action
post (Post state) =
  div
    [ className "col s6" ]
    [ div
      [ className "card pink lighten-4" ]
      [ div
        [ className "card-content black-text" ]
        [ span [ className "card-title" ] [ text state.title ]
        , br [] []
        , p [] [ text ("Posted on: " <> state.date) ]
        , span [] [ text state.summary ]
        ]
      , div
        [ className "card-action pink" ]
        [ link state.link [] [ text "Read More" ] ]
      ]
    ]

view :: State -> Html Action
view state =
  div
    []
    [ h1 [] [ text state.status ]
    , button [ onClick (const RequestPosts), id_ "requestbutton", className "hidden" ] [ text "Fetch posts" ]
    , div [ className "row" ] $ map post state.posts ]
