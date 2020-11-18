#!/usr/bin/env nix-shell
#! nix-shell -p doctl -p kubectl -p curl -i bash
#! nix-shell -I nixpkgs=https://nixos.org/channels/nixpkgs-unstable

nix-env -if ./nix/dhall-yaml.nix
doctl kubernetes cluster kubeconfig save kubermemes
dhall-to-yaml-ng < ./site.dhall | kubectl apply -n apps -f -
kubectl rollout status -n apps deployment/christinewebsite
kubectl apply -f ./k8s/job.yml
sleep 10
kubectl delete -f ./k8s/job.yml
curl --http1.1 -H "Authorization: $MI_TOKEN" https://mi.within.website/api/blog/refresh
