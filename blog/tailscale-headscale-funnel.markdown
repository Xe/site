---
title: "Using Tailscale without using Tailscale"
date: 2023-04-01
redirect_to: https://tailscale.dev/blog/headscale-funnel
---

I’m always amazed by the power of technology and the creativity of people who use it. Today I want to share with you a project that I did for fun and learning. It’s about how to use Tailscale Funnel to host a Headscale server from behind NAT. Sounds crazy, right? Let me explain.

Tailscale Funnel is a tool that lets you share a web service on your private network with the public internet. It’s like having your own personal cloud without the hassle of setting up servers and domains. Headscale is a tool that lets you create your own Tailscale control plane. It’s like having your own private VPN without relying on a third-party service.

Now, what if you could combine these two tools and create a network that uses Tailscale without using Tailscale? That’s exactly what I did. I used waifud, NixOS, and Funnel to create a virtual machine that runs Headscale and exposes it to the internet. Then I used Funnel to connect other devices to this network and enjoy the benefits of Tailscale.

Why did I do this? Because I love learning new things and challenging myself. Because I believe in the power of open source and decentralized solutions. Because I wanted to have some fun and make some jokes along the way.

This project taught me a lot about networking, security, and automation. It also showed me how much potential there is in tools like Tailscale and Headscale. And it made me laugh at the absurdity of using Tailscale without using Tailscale.

If you’re curious about how I did this, you can check out my blog post where I explain everything in detail. It’s not a serious tutorial, but rather a playful experiment.

I hope this post inspires you to try new things and have fun with technology. Remember, nothing is impossible if you put your mind to it. And don’t forget to laugh at yourself sometimes.
