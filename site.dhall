let kms =
      https://xena.greedo.xeserv.us/pkg/dhall/kubermemes/k8s/package.dhall sha256:e47e95aba6a08f8ca3e38fbabc436566d6558a05a9b4ac149e8e712c8583b8f0

let kubernetes =
      https://xena.greedo.xeserv.us/pkg/dhall/dhall-kubernetes/1.15/package.dhall sha256:271494d6e3daba2a47d9d023188e35bf44c9c477a1cfbad1c589695a6b626e56

let tag = env:GITHUB_SHA as Text ? "latest"

let image = "xena/christinewebsite:${tag}"

let vars
    : List kubernetes.EnvVar.Type
    = [ kubernetes.EnvVar::{ name = "PORT", value = Some "3030" }
      , kubernetes.EnvVar::{ name = "RUST_LOG", value = Some "info" }
      , kubernetes.EnvVar::{
        , name = "PATREON_CLIENT_ID"
        , value = Some env:PATREON_CLIENT_ID as Text
        }
      , kubernetes.EnvVar::{
        , name = "PATREON_CLIENT_SECRET"
        , value = Some env:PATREON_CLIENT_SECRET as Text
        }
      , kubernetes.EnvVar::{
        , name = "PATREON_ACCESS_TOKEN"
        , value = Some env:PATREON_ACCESS_TOKEN as Text
        }
      , kubernetes.EnvVar::{
        , name = "PATREON_REFRESH_TOKEN"
        , value = Some env:PATREON_REFRESH_TOKEN as Text
        }
      ]

in  kms.app.make
      kms.app.Config::{
      , name = "christinewebsite"
      , appPort = 3030
      , image = image
      , replicas = 2
      , domain = "christine.website"
      , leIssuer = "prod"
      , envVars = vars
      }
