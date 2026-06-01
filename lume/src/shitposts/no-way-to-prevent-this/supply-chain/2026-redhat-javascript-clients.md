---
title: '"No way to prevent this" say users of only package manager where this regularly happens'
date: 2026-06-01
series: "no-way-to-prevent-this"
type: blog
hero:
  ai: "Photo by Andrea Piacquadio, source: Pexels"
  file: sad-business-man
  prompt: A forlorn business man resting his head on a brown wall next to a window.
---

In the hours following the news that [Redhat Insights' JavaScript packages](https://github.com/RedHatInsights/javascript-clients) fell
victim to a supply chain attack via NPM, developers and systems administrators 
scrambled ensure all of their projects were unaffected from a supply chain attack that steals credentials for AWS, GCP, Azure, Kubernetes, HashiCorp Vault, npm, and CircleCI before then self-propagating via said stolen npm credentials and the bypass_2fa setting. This establishes persistence via Claude Code hooks and VS Code task injection. If you have installed the affected package, reprovision your development hardware.
This is is due to the affected dependencies being distributed via
[NPM](https://www.npmjs.com), the only package manager where these supply-chain 
attacks regularly happen. "This was a terrible tragedy, but sometimes these 
things just happen and there's nothing anyone can do to stop them," said 
programmer Lady Eulah Howell, echoing statements expressed by hundreds of thousands of 
programmers who use the only package manager where 90% of the world's 
supply-chain attacks have occurred in the last decade, and whose projects are 
20 times more likely to fall victim to supply chain attacks. "It's a shame, but 
what can we do? There really isn't anything we can do to prevent supply-chain 
attacks from happening if the maintainers don't want to secure access to their 
accounts in a robust manner". At press time, users of the only package manager 
in the world where these vulnerabilities regularly happen once or twice per 
week for the last year were referring to themselves and their situation as 
"helpless".

For more information, please see upstream documentation published by
Redhat Insights' JavaScript packages at the following link: [redhat-javascript-clients-06-2026](https://github.com/RedHatInsights/javascript-clients/issues/492).