---
title: The Cult of Kubernetes
date: 2019-09-07
series: howto
tags:
 - kubernetes
 - digitalocean
 - githubactions
---

or: How I got my blog onto it with autodeployment via GitHub Actions.

The world was once a simple place. Things used to make sense, or at least there
weren't so many layers that it became difficult to tell what the hell is going
on.

Then complexity happened. This is a tale of how I literally recreated this meme:

<center><blockquote class="twitter-tweet"><p lang="en" dir="ltr">Deployed my blog on Kubernetes <a href="https://t.co/XHXWLrmYO4">pic.twitter.com/XHXWLrmYO4</a></p>&mdash; DevOps Thought Liker (@dexhorthy) <a href="https://twitter.com/dexhorthy/status/856639005462417409?ref_src=twsrc%5Etfw">April 24, 2017</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

This is how I deployed my blog (the one you are reading right now) to Kubernetes.

## The Old State of the World

Before I deployed my blog to Kubernetes, I used [Dokku][dokku], as I had been
for years. Dokku is great. It emulates most of the Heroku "git push; don't care"
workflow, but on your own server that you can self-manage.

This is a blessing and a curse.

The real advantage of managed services like Heroku is that you literally just
_HAND OFF_ operations to Heroku's team. This is not the case with Dokku. Unless
you pay someone a lot of money, you are going to have to manage the server
yourself. My dokku server was unmanaged, and I run _many_ apps on it (this
listing was taken after I started to move apps over):

```
=====> My Apps
bsnk
cinemaquestria
fordaplot-backup
graphviz.christine.website
identicond
ilo-kesi
johaus
maison
olin
printerfacts
since
tulpaforce.tk
tulpanomicon
```

This is enough apps (plus 5 more that I've already migrated) that it really
doesn't make sense paying for something like Heroku; nor does it really make
sense to use the free tier either.

So, I decided that it was time for me to properly learn how to Kubernetes, and I
set off to create a cluster via [DigitalOcean managed Kubernetes][dok8s].

## The Cluster

I decided it would be a good idea to create my cluster using
[Terraform][terraform], mostly because I wanted to learn how to use it better.
I use Terraform at work, so I figured this would also be a way to level up my
skills in a mostly sane environment.

<center><blockquote class="twitter-tweet"><p lang="en" dir="ltr">Terraform is suffering as a service</p>&mdash; Cadey Ratio 🌐 (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1165390942679048192?ref_src=twsrc%5Etfw">August 24, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

I have been creating and playing with a small Terraform wrapper tool called
[dyson][dyson]. This tool is probably overly simplistic and is written in Nim.
With the config in `~/.config/dyson/dyson.ini`, I can simplify my Terraform
usage by moving my secrets _out of_ the Terraform code directly. I also avoid
having my API tokens exposed in my shell to avoid accidental exposure of the
secrets.

Dyson is very simple to use:

```console
$ dyson
Usage:
  dyson {SUBCMD}  [sub-command options & parameters]
where {SUBCMD} is one of:
  help         print comprehensive or per-cmd help
  apply        apply Terraform code to production
  destroy      destroy resources managed by Terraform
  env          dump envvars
  init         init Terraform
  manifest     generate a somewhat sane manifest for a kubernetes app based on the arguments.
  plan         plan a future Terraform run
  slug2docker  converts a heroku/dokku slug to a docker image

dyson {-h|--help} or with no args at all prints this message.
dyson --help-syntax gives general cligen syntax help.
Run "dyson {help SUBCMD|SUBCMD --help}" to see help for just SUBCMD.
Run "dyson help" to get *comprehensive* help.
```

So I wrote up my config:

```
# main.tf
provider "digitalocean" {}

resource "digitalocean_kubernetes_cluster" "main" {
  name    = "kubermemes"
  region  = "${var.region}"
  version = "${var.kubernetes_version}"

  node_pool {
    name       = "worker-pool"
    size       = "${var.node_size}"
    node_count = 2
  }
}
```

```
# variables.tf
variable "region" {
  type    = "string"
  default = "nyc3"
}

variable "kubernetes_version" {
  type    = "string"
  default = "1.15.3-do.1"
}

variable "node_size" {
  type    = "string"
  default = "s-1vcpu-2gb"
}
```

and ran it:

```console
$ dyson plan
<... many lines of plan output ...>
$ dyson apply
<... many lines of apply output ...>
```

Then I had a working but mostly unconfigured Kubernetes cluster.

## Configuration

This is where things started to go downhill. I wanted to do a few things with
this cluster so I could consider it "ready" for me to use for deploying applications
to.

I wanted to do the following:

- setup [helm][helm] to install packages for things like DNS management and HTTP/HTTPS ingress
- setup [automatic certificate management][certmanager] with [Let's Encrypt][letsencrypt]
- setup HTTP/HTTPS request ingress with [nginx-ingress][nginxingress] (which uses [nginx](https://www.nginx.com/))
- setup [automatic DNS management][autodns] because the external IP addresses of Kubernetes nodes can and will change

After a lot of trial, error, pain, suffering and the like, I created
[this script][setupdotsh] which I am not pasting here. Look at it if you want to
get a streamlined overview of how to set these things up.

Now that all of this is set up, I can deploy an [example app][exanple] with a
manifest that looks something like [this][ingresstestdotyaml]:

```yaml
apiVersion: v1
kind: Service
metadata:
  name: hello-kubernetes-first
  annotations:
    external-dns.alpha.kubernetes.io/hostname: exanple.within.website
    external-dns.alpha.kubernetes.io/ttl: "120" #optional
    external-dns.alpha.kubernetes.io/cloudflare-proxied: "false"
spec:
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 8080
  selector:
    app: hello-kubernetes-first
    
---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-kubernetes-first
spec:
  replicas: 1
  selector:
    matchLabels:
      app: hello-kubernetes-first
  template:
    metadata:
      labels:
        app: hello-kubernetes-first
    spec:
      containers:
      - name: hello-kubernetes
        image: paulbouwer/hello-kubernetes:1.5
        ports:
        - containerPort: 8080
        env:
        - name: MESSAGE
          value: Henlo this are an exanple deployment
          
---

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: hello-kubernetes-ingress
  annotations:
    kubernetes.io/ingress.class: nginx
    certmanager.k8s.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - exanple.within.website
    secretName: prod-certs
  rules:
  - host: exanple.within.website
    http:
      paths:
      - backend:
          serviceName: hello-kubernetes-first
          servicePort: 80
```

<center><blockquote class="twitter-tweet"><p lang="en" dir="ltr">Nope, I was wrong, Kubernetes is the real suffering as a service</p>&mdash; Cadey Ratio 🌐 (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1169997202971930624?ref_src=twsrc%5Etfw">September 6, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

It was about this time when I wondered if I was making a mistake moving off of
Dokku. Dokku really does a lot to abstract almost everything involved with nginx
away from you, and it _really shows_.

However, as a side effect of everything being so declarative and Kubernetes really
not assuming anything, you have _a lot_ more freedom to do basically anything
you want. You don't have to have specially magic names for tasks like `web` or
`worker` like you do in Heroku/Dokku. You just have a deployment that belongs to
an "app" that just so happens to expose a TCP port that just so happens to have
a correlating ingress associated with it.

Lucky for me, most of the apps I write fit into that general format, and the ones
that don't can mostly use the same format without the ingress.

So I [templated][deploymenttemplateyaml] that sucker as a subcommand in dyson.
This lets me do commands like [this][exampledysonmanifestcommand]:

```console
$ dyson manifest \
      --name=hlang \
      --domain=h.christine.website \
      --dockerImage=docker.pkg.github.com/xe/x/h:v1.1.8 \
      --containerPort=5000 \
      --replicas=1 \
      --useProdLE=true | kubectl apply -f-
```

And the service gets shunted into the cloud without any extra effort on my part.
This also automatically sets up Let's Encrypt, DNS and other things that were
manual in my Dokku setup. This saves me time for when I want to go add services
in the future. All I have to do is create a docker image somehow, identify what
port should be exposed, give it a domain name and number of replicas and just
send it on its merry way.

## GitHub Actions

This does however mean that deployment is no longer as simple as
"git push; don't care". This is where [GitHub Actions][actions] come into play.
They claimed to have the ability to run full end-to-end CI/CD on my applications.

I have been using them for a while for [CI on my website][sitegoci] and have
been pleased with them, so I decided to give it a try and set up continuous
deployment with them.

As [the commit log for the deployment manifest can tell][kubernetescddotyml],
this took a lot of trial and error. One of the main sources of problems here
was that GitHub Actions had recently had _a lot_ of changes made to
configuration and usage as compared to when it was in private beta. This
included changing the configuration schema from [HCL][hcl] to [YAML][yaml].

Of course, all of the documentation (outside of GitHub's
[quite excellent documentation][githubactionsdocs]) was out of date and wrong.
I tried following a tutorial by [DigitalOcean themselves][dotutorialkube] on
how to do this exact thing I wanted to do, but it referenced the old HCL syntax
for GitHub Actions and did not work. To make things worse, examples
[in the marketplace READMEs][marketplacereadmeexample] simply DID NOT WORK because
they were written for the old GitHub Actions syntax.

This was frustrating to say the least.

After trying to make them work anyways with a combination of the "Use Latest
Version" button in the marketplace, prayer and gratuitous use of the `with.args`
field in steps I gave up and decided to manually download the tools I needed
from their upstream providers and execute them by hand.

This is how I ended up with [this monstrosity][monster]:

```yaml
- name: Configure/Deploy/Verify Kubernetes
  run: |
    curl -L https://github.com/digitalocean/doctl/releases/download/v1.30.0/doctl-1.30.0-linux-amd64.tar.gz | tar xz
    ./doctl auth init -t $DIGITALOCEAN_ACCESS_TOKEN
    ./doctl kubernetes cluster kubeconfig show kubermemes > .kubeconfig

    curl -LO https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl
    chmod +x kubectl
    ./kubectl --kubeconfig .kubeconfig apply -n apps -f deploy.yml
    sleep 2
    ./kubectl --kubeconfig .kubeconfig rollout -n apps status deployment/christinewebsite
  env:
    DIGITALOCEAN_ACCESS_TOKEN: ${{ secrets.DIGITALOCEAN_TOKEN }}
```

~~I am almost _certain_ that I am doing it wrong here, I don't know how robust this
is and I'm very sure that this can and should be done another way; but this is
the only thing I could get working (for some definition of "working").~~

EDIT: it got fixed, see below

<center><blockquote class="twitter-tweet"><p lang="hu" dir="ltr">kubernetes is a cult</p>&mdash; Andrew Kelley (@andy_kelley) <a href="https://twitter.com/andy_kelley/status/1169999209438859264?ref_src=twsrc%5Etfw">September 6, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

---

Now when I git push things to the master branch of my blog repo, it will
automatically get deployed to my Kubernetes cluster.

~~If you work at DigitalOcean and are reading this post. Please get someone to
update [this tutorial][dotutorialkube] and the README of [this repo][marketplacereadmeexample].
The examples listed _DO NOT WORK_ for me because I was not in the private beta
of GitHub Actions. It would also be nice if you had better documentation on how
to use [your premade action][doctlgithubaction] for usecases like mine. I just
wanted to download the kubernetes configuration file and run apply against a yaml
file.~~

EDIT: The above complaint has been fixed! See [here](https://github.com/Xe/site/commit/a9329bfbeffc4a290876a56795c23286c537ca94)
for the simpler way of doing things.

Thanks for reading, I hope this was entertaining. Be well.

[dokku]: https://dokku.com/
[dok8s]: https://www.digitalocean.com/products/kubernetes/
[terraform]: https://www.terraform.io
[dyson]: https://github.com/Xe/within-terraform/tree/master/dyson
[helm]: https://helm.sh
[certmanager]: https://docs.cert-manager.io/en/latest/
[letsencrypt]: https://letsencrypt.org
[nginxingress]: https://kubernetes.github.io/ingress-nginx/
[autodns]: https://github.com/kubernetes-incubator/external-dns
[setupdotsh]: https://github.com/Xe/within-terraform/blob/master/do/setup.sh
[exanple]: https://exanple.within.website
[ingresstestdotyaml]: https://github.com/Xe/within-terraform/blob/master/do/ingress_test.yaml
[deploymenttemplateyaml]: https://github.com/Xe/within-terraform/blob/master/dyson/src/dysonPkg/deployment_with_ingress.yaml
[exampledysonmanifestcommand]: https://github.com/Xe/within-terraform/blob/master/kube_manifests/h.sh
[actions]: https://github.com/features/actions
[sitegoci]: https://github.com/Xe/site/blob/e4d7c3c2691acad73d6240ff0c9b208273b95997/.github/workflows/go.yml
[githubactionsdocs]: https://help.github.com/en/articles/about-github-actions
[kubernetescddotyml]: https://github.com/Xe/site/commits/e4d7c3c2691acad73d6240ff0c9b208273b95997/.github/workflows/kubernetes-cd.yml
[hcl]: https://github.com/hashicorp/hcl
[yaml]: https://yaml.org
[marketplacereadmeexample]: https://github.com/marketplace/actions/github-action-for-digitalocean-doctl
[monster]: https://github.com/Xe/site/blob/e4d7c3c2691acad73d6240ff0c9b208273b95997/.github/workflows/kubernetes-cd.yml#L53-L65
[dotutorialkube]: https://blog.digitalocean.com/how-to-deploy-to-digitalocean-kubernetes-with-github-actions/
[doctlgithubaction]: https://github.com/digitalocean/action-doctl
