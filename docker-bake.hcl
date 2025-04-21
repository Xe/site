variable "ALPINE_VERSION" { default = "edge" }
variable "DENO_VERSION" { default = "2.2.11" }
variable "DHALL_VERSION" { default = "1.42.2" }
variable "DHALL_JSON_VERSION" { default = "1.7.12" }
variable "GO_VERSION" { default = "1.24" }
variable "UBUNTU_VERSION" { default = "24.04" }

group "default" {
  targets = [ "patreon-saasproxy", "xesite" ]
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
    GO_VERSION = null
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