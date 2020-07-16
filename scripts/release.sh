#!/usr/bin/env nix-shell
#! nix-shell -p doctl -p kubectl -i bash
nix-env -if ./nix/dhall-yaml.nix
doctl kubernetes cluster kubeconfig save kubermemes
dhall-to-yaml-ng < ./site.dhall | kubectl apply -n apps -f -
kubectl rollout status -n apps deployment/christinewebsite
