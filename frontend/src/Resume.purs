module App.Resume where

import App.Utils (mdify)
import Control.Monad.Aff (attempt)
import DOM (DOM)
import Data.Argonaut (class DecodeJson, decodeJson, (.?))
import Data.Either (Either(..), either)
import Data.Maybe (Maybe(..))
import Network.HTTP.Affjax (AJAX, get)
import Prelude (Unit, bind, pure, show, unit, ($), (<>), (<<<))
import Pux (noEffects, EffModel)
import Pux.DocumentTitle (documentTitle)
import Pux.Html (Html, a, div, h1, p, text)
import Pux.Html.Attributes (href, dangerouslySetInnerHTML, className, id_, title)

data Action = RequestResume
            | ReceiveResume (Either String Resume)

type State =
  { status :: String
  , err    :: String
  , resume :: Maybe Resume }

data Resume = Resume
  { body :: String }

instance decodeJsonResume :: DecodeJson Resume where
  decodeJson json = do
    obj <- decodeJson json
    body <- obj .? "body"
    pure $ Resume { body: body }

init :: State
init =
  { status: "Loading..."
  , err: ""
  , resume: Nothing }

update :: Action -> State -> EffModel State Action (ajax :: AJAX, dom :: DOM)
update (ReceiveResume (Left err)) state =
  noEffects $ state { resume = Nothing, status = "Error in fetching resume, please use the plain text link below.", err = err }
update (ReceiveResume (Right body)) state =
  noEffects $ state { status = "", err = "", resume = Just body }
    where
      got' = Just unit
update RequestResume state =
  { state: state
  , effects: [ do
      res <- attempt $ get "/api/resume"
      let decode r = decodeJson r.response :: Either String Resume
      let resume = either (Left <<< show) decode res
      pure $ ReceiveResume resume
    ]
  }

view :: State -> Html Action
view { status: status, err: err, resume: resume } =
      case resume of
        Nothing -> div [] [ text status, p [] [ text err ] ]
        (Just (Resume resume')) ->
          div [ className "row" ]
          [ documentTitle [ title "Resume - Christine Dodrill" ] []
          , div [ className "col s8 offset-s2" ]
            [ p [ className "browser-default", dangerouslySetInnerHTML $ mdify resume'.body ] []
            , a [ href "/static/resume/resume.md" ] [ text "Plain-text version of this resume here" ], text "." ]
          ]
