module App.Routes where

import Control.Alt ((<|>))
import Control.Apply ((<*), (*>))
import Data.Functor ((<$))
import Data.Maybe (fromMaybe)
import Prelude (($), (<$>))
import Pux.Router (param, router, lit, str, end)

data Route = Home
           | Resume
           | StaticPage String
           | BlogIndex
           | BlogPost String
           | NotFound

match :: String -> Route
match url = fromMaybe NotFound $ router url $
    Home <$ end
  <|>
    BlogIndex <$ lit "blog" <* end
  <|>
    BlogPost <$> (lit "blog" *> str) <* end
