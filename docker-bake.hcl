variable "ALPINE_VERSION" {
  default = "edge"
}

variable "GO_VERSION" {
  default = "1.24"
}

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
    GO_VERSION = null
  }
  context = "."
  dockerfile = "./docker/xesite.Dockerfile"
  platforms = [ "linux/amd64", "linux/arm64" ]
  pull = true
  tags = [
    "registry.int.xeserv.us/xe/site/bin:main"
  ]
}