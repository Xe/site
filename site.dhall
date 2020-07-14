let kms =
      https://xena.greedo.xeserv.us/pkg/dhall/kubermemes/k8s/package.dhall sha256:e47e95aba6a08f8ca3e38fbabc436566d6558a05a9b4ac149e8e712c8583b8f0

let kubernetes =
      https://xena.greedo.xeserv.us/pkg/dhall/dhall-kubernetes/1.15/package.dhall sha256:271494d6e3daba2a47d9d023188e35bf44c9c477a1cfbad1c589695a6b626e56

let tag = env:GITHUB_SHA as Text ? "latest"

let image = "xena/christinewebsite:${tag}"

let vars
    : List kubernetes.EnvVar.Type
    = [ kubernetes.EnvVar::{ name = "PORT", value = Some "3030" } ]

in  kms.app.make
      kms.app.Config::{
      , name = "christinewebsite"
      , appPort = 3030
      , image = image
      , domain = "christine.website"
      , leIssuer = "prod"
      , envVars = vars
      }
