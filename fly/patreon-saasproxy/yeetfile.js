nix.build(".#patreon-docker")
docker.load("./result")
docker.push("ghcr.io/xe/site/patreon")
fly.deploy()
