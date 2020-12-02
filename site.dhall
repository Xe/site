let kms = https://tulpa.dev/cadey/kubermemes/raw/branch/master/k8s/package.dhall

let kubernetes =
      https://raw.githubusercontent.com/dhall-lang/dhall-kubernetes/master/1.15/package.dhall

let tag = env:GITHUB_SHA as Text ? "latest"

let image = "ghcr.io/xe/site:${tag}"

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
      , kubernetes.EnvVar::{
        , name = "MI_TOKEN"
        , value = Some env:MI_TOKEN as Text
        }
      ]

in  kms.app.make
      kms.app.Config::{
      , name = "christinewebsite"
      , appPort = 3030
      , image
      , replicas = 2
      , domain = "christine.website"
      , leIssuer = "prod"
      , envVars = vars
      }
