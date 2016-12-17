module Pux.DocumentTitle where

import Pux.Html (Html, Attribute)

-- | Declaratively set `document.title`.  See [react-document-title](https://github.com/gaearon/react-document-title)
-- | for more information.
foreign import documentTitle :: forall a. Array (Attribute a) -> Array (Html a) -> Html a
