variable "ALPINE_VERSION" { default = "edge" }
variable "DENO_SHA" { default = "6ef38d16cbe99c3d610576b56aaa9ede9d988e8a2e5c1ed9c9d502e3167ef758" }
variable "DENO_VERSION" { default = "2.2.11" }
variable "DHALL_VERSION" { default = "1.42.2" }
variable "DHALL_JSON_VERSION" { default = "1.7.12" }
variable "DHALL_JSON_SHA" { default = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855" }
variable "FONTS_SHA" { default = "2d96002c16d611fe8498a71c0b44362b4a98e18023cce34e7e37f581f34def22" }
variable "FONTS_VERSION" { default = "20250421" }
variable "GO_VERSION" { default = "1.25" }
variable "TYPST_SHA" { default = "7d214bfeffc2e585dc422d1a09d2b144969421281e8c7f5d784b65fc69b5673f" }
variable "TYPST_VERSION" { default = "0.13.1" }
variable "UBUNTU_VERSION" { default = "24.04" }

group "default" {
  targets = [ "patreon-saasproxy", "xesite", "github-sponsor-webhook" ]
}

target "patreon-saasproxy" {
  args = {
    ALPINE_VERSION = null
    GO_VERSION = null
  }
  context = "."
  dockerfile = "./docker/patreon-saasproxy.Dockerfile"
  platforms = [ "linux/amd64", "linux/arm64" ]
  pull = true
  tags = [
    "registry.int.xeserv.us/xe/site/patreon-saasproxy:main"
  ]
}

target "xesite" {
  args = {
    ALPINE_VERSION = null
    DENO_VERSION = null
    DHALL_VERSION = null
    DHALL_JSON_VERSION = null
    FONTS_VERSION = null
    GO_VERSION = null
    TYPST_VERSION = null
    UBUNTU_VERSION = "24.04"
  }
  context = "."
  dockerfile = "./docker/xesite.Dockerfile"
  platforms = [
    "linux/amd64"
    #"linux/arm64",
  ]
  pull = true
  tags = [
    "registry.int.xeserv.us/xe/site/bin:main"
  ]
}

target "github-sponsor-webhook" {
  args = {
    ALPINE_VERSION = null
    GO_VERSION = null
  }
  context = "."
  dockerfile = "./docker/github-sponsor-webhook.Dockerfile"
  platforms = [ "linux/amd64", "linux/arm64" ]
  pull = true
  tags = [
    "registry.int.xeserv.us/xe/site/github-sponsor-webhook:main"
  ]
}