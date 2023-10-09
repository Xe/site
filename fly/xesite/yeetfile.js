nix.build(".#docker")
docker.load("./result")
docker.push("ghcr.io/xe/site/bin")
fly.deploy()
