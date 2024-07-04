---
title: "Setting up Gitea on Kubernetes from scratch"
date: 2024-06-30
tags:
  - kubernetes
  - gitea
  - helm
  - pain
vod:
  path: video/2024/gitea-k8s-vod
---

I set up Gitea on Kubernetes with only the documentation, a dream, and existential suffering.

## High Level Topics Covered

- Setting up Gitea, a self-hosted Git service, on a Kubernetes cluster
- Using Helm, a package manager for Kubernetes, to simplify the deployment process
- Troubleshooting common issues that arise during the setup process
- Configuring and managing persistent storage for Gitea
- Integrating Gitea with other services, such as Tigris and S3

## Interesting Lessons Learned

- Helm can be a useful tool for simplifying Kubernetes deployments, but it can also introduce complexity and potential security risks.
- Persistent storage is essential for ensuring that data is preserved across pod restarts.
- Careful configuration of Gitea's storage settings is important to avoid potential errors and data loss.
- Integrating Gitea with other services can enhance its functionality and make it a more powerful tool for managing Git repositories.
- It is important to thoroughly test and troubleshoot any changes made to a Kubernetes deployment to ensure that it is functioning properly.
