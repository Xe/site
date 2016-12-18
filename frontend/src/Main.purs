module Main where

import App.Layout (Action(PageView), State, view, update)
import App.Routes (match)
import Control.Bind ((=<<))
import Control.Monad.Eff (Eff)
import DOM (DOM)
import Network.HTTP.Affjax (AJAX)
import Prelude (bind, pure)
import Pux (renderToDOM, renderToString, App, Config, CoreEffects, start)
import Pux.Devtool (Action, start) as Pux.Devtool
import Pux.Router (sampleUrl)
import Signal ((~>))

type AppEffects = (dom :: DOM, ajax :: AJAX)

-- | App configuration
config :: forall eff. State -> Eff (dom :: DOM | eff) (Config State Action AppEffects)
config state = do
  -- | Create a signal of URL changes.
  urlSignal <- sampleUrl

  -- | Map a signal of URL changes to PageView actions.
  let routeSignal = urlSignal ~> \r -> PageView (match r)

  pure
    { initialState: state
    , update: update
    , view: view
    , inputs: [routeSignal] }

-- | Entry point for the browser.
main :: State -> Eff (CoreEffects AppEffects) (App State Action)
main state = do
  app <- start =<< config state
  renderToDOM "#app" app.html
  -- | Used by hot-reloading code in support/index.js
  pure app

-- | Entry point for the browser with pux-devtool injected.
debug :: State -> Eff (CoreEffects AppEffects) (App State (Pux.Devtool.Action Action))
debug state = do
  app <- Pux.Devtool.start =<< config state
  renderToDOM "#app" app.html
  -- | Used by hot-reloading code in support/index.js
  pure app

-- | Entry point for server side rendering
ssr :: State -> Eff (CoreEffects AppEffects) String
ssr state = do
  app <- start =<< config state
  res <- renderToString app.html
  pure res
