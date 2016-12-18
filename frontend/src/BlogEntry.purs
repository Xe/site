module App.BlogEntry where

import App.Utils (mdify)
import Control.Monad.Aff (attempt)
import DOM (DOM)
import Data.Argonaut (class DecodeJson, decodeJson, (.?))
import Data.Either (Either(..), either)
import Data.Maybe (Maybe(..))
import Network.HTTP.Affjax (AJAX, get)
import Prelude (bind, pure, show, ($), (<>), (<<<))
import Pux (noEffects, EffModel)
import Pux.DocumentTitle (documentTitle)
import Pux.Html (Html, div, h1, p, text)
import Pux.Html.Attributes (dangerouslySetInnerHTML, className, id_, title)

data Action = RequestPost
            | ReceivePost (Either String Post)

type State =
  { status :: String
  , id     :: Maybe Int
  , post   :: Post
  , name   :: String }

data Post = Post
  { title :: String
  , body  :: String
  , date  :: String }

instance decodeJsonPost :: DecodeJson Post where
  decodeJson json = do
    obj <- decodeJson json
    title <- obj .? "title"
    body <- obj .? "body"
    date <- obj .? "date"
    pure $ Post { title: title, body: body, date: date }

init :: State
init =
  { status: "Loading..."
  , post: Post
    { title: ""
    , body: ""
    , date: "" }
  , name: ""
  , id: Nothing }

update :: Action -> State -> EffModel State Action (ajax :: AJAX, dom :: DOM)
update (ReceivePost (Left err)) state =
  noEffects $ state { id = Nothing, status = err }
update (ReceivePost (Right post)) state = noEffects $ state { status = "", id = Just 1, post = post }
update RequestPost state =
  { state: state
  , effects: [ do
      res <- attempt $ get ("/api/blog/post?name=" <> state.name)
      let decode r = decodeJson r.response :: Either String Post
      let post = either (Left <<< show) decode res
      pure $ ReceivePost post
    ]
  }

view :: State -> Html Action
view { id: id, status: status, post: (Post post) } =
      case id of
        Nothing -> div [] []
        (Just _) ->
          div [ className "row" ]
          [ h1 [] [ text status ]
          , documentTitle [ title $ post.title <> " - Christine Dodrill" ] []
          , div [ className "col s8 offset-s2" ]
            [ p [ id_ "blogpost", dangerouslySetInnerHTML $ mdify post.body ] [] ]
          ]
