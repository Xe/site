---
title: "Funnel 101: sharing your local developer preview with the world"
date: 2023-03-30
redirect_to: https://tailscale.dev/blog/funnel-101
---

🚀 Do you want to share your web server with the world without exposing your computer to the world? 🚀

If you’re like me, you love using Tailscale to create a secure and private network for your devices. But sometimes, you need to let the outside world access your web server, whether it’s for testing, hosting, or collaborating.

That’s why I’m super excited about Tailscale Funnel, a new feature that lets you route traffic from the internet to your Tailscale node. You can think of it as publicly sharing a node for anyone to access, even if they don’t have Tailscale themselves.

Tailscale Funnel is easy to set up and use. You just need to enable it in the admin console and on your node, and you’ll get a public DNS name for your node that points to Tailscale’s Funnel servers. These servers will proxy the incoming requests over Tailscale to your node, where you can terminate the TLS and serve your content.

The best part is that Tailscale Funnel is secure and private. The Funnel servers don’t see any information about your traffic or what you’re serving. They only see the source IP and port, the SNI name, and the number of bytes passing through. And they can’t connect to your nodes directly. They only offer a TCP connection, which your nodes can accept or reject.

Tailscale Funnel is currently in beta and available for all users. I’ve been using it for a while now and I’m blown away by how simple and powerful it is. It’s like having your own personal cloud service without any hassle or cost.
