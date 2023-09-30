---
title: textile-conversion Main
date: 2017-02-08
---

Author's Note: this was intended to be documentation for a service that never ended
up being implemented. It was going to help [Derpibooru](https://derpibooru.org)
convert its existing markup to [Markdown](https://github.github.com/gfm/). This
never happened.

This program listens on port 5000 and serves an unchecked-path web handler that
converts Derpibooru Textile via HTML into Markdown, using a two-step process.

The first step is to have SimpleTextile emit a HTML AST of the comment.
The second is to have Pandoc turn that HTML into Markdown.

This is intended to be helpful during Derpi's migration from Textile.

## Pragmas

The following pragma tells the compiler to automagically tease string literals
into whatever type they need to be. For more information on this, see [this page][hs-ovs].

```haskell
{-# LANGUAGE OverloadedStrings #-}
module Main where
```

## Imports

In order to accomplish our task, we need to import some libraries.

```haskell
import Data.String.Conv (toS)
import Network.Wai
import Network.HTTP.Simple
import Network.HTTP.Types
import Network.Wai.Handler.Warp (run)
import System.Environment (lookupEnv)
import Text.Pandoc
import Text.Pandoc.Error (PandocError, handleError)
```

## Helper Functions

getEnvDefault queries an environment variable, returning a default value if it
is unset.

```haskell
getEnvDefault :: String -> String -> IO String
getEnvDefault name default' = do
    envvar <- lookupEnv name
    case envvar of
      Nothing -> pure default'
      Just x  -> pure x
```

---

htmlToMarkdown uses Pandoc to convert a HTML input string into the equivalent
Markdown. The `Either` type is used here in place of raising an exception.

```haskell
htmlToMarkdown :: String -> Either PandocError String
htmlToMarkdown inp = do
    let
        corpus = readHtml def inp

    case corpus of
        Left x ->  Left x
        Right x -> pure $ writeMarkdown def x
```

## Web Application

Now we are getting into the meat of the situation. This is the main
[Application][wai-application].

```haskell
toMarkdown :: Application
```

First, let's use a [guard][guards] to ensure that we are only accepting `POST`
requests. If the request is not a `POST` request, return [HTTP error code 405][http-4xx].

```haskell
toMarkdown req respond
    | requestMethod req /= methodPost =
        respond $ responseLBS
            status405
            [("Content-Type", "text/plain")]
            "Not allowed"
```

Otherwise, this is a `POST` request, so we should:

1. Unpack the data from the post body of the HTTP request
2. Send the data to the Sinatra app for conversion from Textile to HTML
3. Take the resulting HTML and feed it to `htmlToMarkdown`
4. Respond with the resulting Markdown.

We use [http-conduit][http-conduit] to contact the Sinatra app.

```haskell
    | otherwise = do
        body <- requestBody req
        targetHost <- getEnvDefault "TARGET_SERVER" "http://127.0.0.1:9292"
        remoteRequest' <- parseRequest ("POST " ++ targetHost ++ "/textile/html")
```

The `($)` operator is a synonym for calling functions. It is defined in the [Prelude][dolla]
as `f $ x = f x` and is mainly used for omitting parentheses. Here it is used
to combine HTTP request settings into one big request.

Additionally we use a custom [Manager][manager] to avoid any issues with
request timeouts, as those are not important for the scope of this tool.

```haskell
        let settings = defaultManagerSettings { managerResponseTimeout = Nothing }
        manager <- newManager settings

        let remoteRequest = setRequestBodyLBS (toS body)
                          $ setRequestManager manager
                          $ remoteRequest'
```

Now it is time to send off the request and unpack the response.

```haskell
        response <- httpLBS remoteRequest
```

If the sinatra app failed to deal with this properly for some reason, report
its error as `text/plain` and return `400`.

```haskell
        if getResponseStatusCode response /= 200
        then respond $ responseLBS
            status400
            [("Content-Type", "text/plain")]
            $ toS $ getResponseBody response
        else do
            let rbody = toS $ getResponseBody response
```

Convert the result body into Markdown. If there is an error, respond with a `400`
and the contents of that error.

```haskell
            let mbody = htmlToMarkdown rbody

            case mbody of
                Left x ->
                    respond $ responseLBS
                        status400
                        [("Content-Type", "text/plain")]
                        $ toS $ show x
                Right x -> do
                    respond $ responseLBS
                        status200
                        [("Content-Type", "text/markdown")]
                        $ toS x
```

Now we bootstrap it all by running the `toMarkdown` Application on port `5000`.
No other code is needed.

```haskell
main :: IO ()
main =
    run 5000 toMarkdown
```

[hs-ovs]: https://ocharles.org.uk/blog/posts/2014-12-17-overloaded-strings.html
[wai-application]: https://hackage.haskell.org/package/wai
[guards]: https://en.wikibooks.org/wiki/Haskell/Control_structures
[http-4xx]: https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#4xx_Client_Error
[http-conduit]: https://www.stackage.org/haddock/lts-6.6/http-conduit-2.1.11/Network-HTTP-Simple.html
[dolla]: https://hackage.haskell.org/package/base-4.9.0.0/docs/Prelude.html#v:-36-
[manger]: https://www.stackage.org/haddock/lts-6.5/http-client-0.4.29/Network-HTTP-Client.html#g:3
