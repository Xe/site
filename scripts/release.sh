#!/usr/bin/env nix-shell
#! nix-shell -p doctl -p kubectl -p curl -i bash
nix-env -if ./nix/dhall-yaml.nix
doctl kubernetes cluster kubeconfig save kubermemes
dhall-to-yaml-ng < ./site.dhall | kubectl apply -n apps -f -
kubectl rollout status -n apps deployment/christinewebsite
kubectl apply -f ./k8s/job.yml
sleep 10
kubectl delete -f ./k8s/job.yml
curl -H "Authorization: $MI_TOKEN" --data "https://christine.website/blog.json" https://mi.within.website/blog/refresh
