---
title: Dhall for Kubernetes
date: 2020-01-25
tags:
  - dhall
  - kubernetes
  - witchcraft
---

Kubernetes is a surprisingly complicated software package. Arguably, it has to
be that complicated as a result of the problems it solves being complicated; but
managing yaml configuration files for Kubernetes is a complicated task. [YAML][yaml]
doesn't have support for variables or type metadata. This means that the
validity (or sensibility) of a given Kubernetes configuration file (or files)
isn't easy to figure out without using a Kubernetes server.

[yaml]: https://yaml.org

In my [last post][cultk8s] about Kubernetes, I mentioned I had developed a tool
named [dyson][dyson] in order to help me manage Terraform as well as create
Kubernetes manifests from [a template][template]. This works for the majority of
my apps, but it is difficult to extend at this point for a few reasons:

[cultk8s]: https://xeiaso.net/blog/the-cult-of-kubernetes-2019-09-07
[dyson]: https://github.com/Xe/within-terraform/tree/master/dyson
[template]: https://github.com/Xe/within-terraform/blob/master/dyson/src/dysonPkg/deployment_with_ingress.yaml

- It assumes that everything passed to it are already valid yaml terms
- It doesn't assert the type of any values passed to it
- It is difficult to add another container to a given deployment
- Environment variables implicitly depend on the presence of a private git repo
- It depends on the template being correct more than the output being correct

So, this won't scale. People in the community have created other solutions for
this like [Helm][helm], but a lot of them have some of the same basic problems.
Helm also assumes that your template is correct. [Kustomize][kustomize] does
help with a lot of the type-safe variable replacements, but it doesn't have the
ability to ensure your manifest is valid.

[helm]: https://helm.sh
[kustomize]: https://kustomize.io

I looked around for alternate solutions for a while and eventually found
[Dhall][dhall] thanks to a friend. Dhall is a _statically typed_ configuration
language. This means that you can ensure that inputs are _always_ the correct
type or the configuration file won't load. There's also a built-in
[dhall-to-yaml][dhallyaml] tool that can be used with the [Kubernetes
package][dhallk8s] in order to declare Kubernetes manifests in a type-safe way.

[dhall]: https://dhall-lang.org
[dhallyaml]: https://github.com/dhall-lang/dhall-haskell/tree/master/dhall-yaml#dhall-yaml
[dhallk8s]: https://github.com/dhall-lang/dhall-kubernetes

Here's a small example of Dhall and the yaml it generates:

```dhall
-- Mastodon usernames
[ { name = "Cadey", mastodon = "@cadey@mst3k.interlinked.me" }
, { name = "Nicole", mastodon = "@sharkgirl@mst3k.interlinked.me" }
]
```

Which produces:

```yaml
- mastodon: "@cadey@mst3k.interlinked.me"
  name: Cadey
- mastodon: "@sharkgirl@mst3k.interlinked.me"
  name: Nicole
```

Which is fine, but we still have the type-safety problem that you would have in
normal yaml. Dhall lets us define [record types][dhallrecord] for this data like
this:

[dhallrecord]: https://www.haskellforall.com/2020/01/dhall-year-in-review-2019-2020.html

```dhall
let User =
      { Type = { name : Text, mastodon : Optional Text }
      , default = { name = "", mastodon = None }
      }

let users =
      [ User::{ name = "Cadey", mastodon = Some "@cadey@mst3k.interlinked.me" }
      , User::{
        , name = "Nicole"
        , mastodon = Some "@sharkgirl@mst3k.interlinked.me"
        }
      ]

in  users
```

Which produces:

```yaml
- mastodon: "@cadey@mst3k.interlinked.me"
  name: Cadey
- mastodon: "@sharkgirl@mst3k.interlinked.me"
  name: Nicole
```

This is type-safe because you cannot add arbitrary fields to User instances
without the compiler rejecting it. Let's add an invalid "preferred_language"
field to Cadey's instance:

```
-- ...
let users =
      [ User::{
        , name = "Cadey"
        , mastodon = Some "@cadey@mst3k.interlinked.me"
        , preferred_language = "en-US"
        }
      -- ...
      ]
```

Which gives us:

```
$ dhall-to-yaml --file example.dhall
Error: Expression doesn't match annotation

{ + preferred_language : …
, …
}

4│         User::{ name = "Cadey", mastodon = Some "@cadey@mst3k.interlinked.me",
5│       preferred_language = "en-US" }

example.dhall:4:9
```

Or [this more detailed explanation][explanation] if you add the `--explain` flag
to the `dhall-to-yaml` call.

[explanation]: https://clbin.com/JtVWT

We tried to do something that violated the contract that the type specified.
This means that it's an invalid configuration and is therefore rejected and no
yaml file is created.

The Dhall Kubernetes package specifies record types for _every_ object available
by default in Kubernetes. This does mean that the package is incredibly large,
but it also makes sure that _everything_ you could possibly want to do in
Kubernetes matches what it expects. In the [package
documentation][k8sdhalldocs], they give an example where a
[Deployment][k8sdeployment] is created.

[k8sdhalldocs]: https://github.com/dhall-lang/dhall-kubernetes/tree/master/1.15#quickstart---a-simple-deployment
[k8sdeployment]: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/

``` dhall
-- examples/deploymentSimple.dhall

-- Importing other files is done by specifying the HTTPS URL/disk location of
-- the file. Attaching a sha256 hash (obtained with `dhall freeze`) allows
-- the Dhall compiler to cache these files and speed up configuration loads
-- drastically.
let kubernetes =
      https://raw.githubusercontent.com/dhall-lang/dhall-kubernetes/1.15/master/package.dhall
      sha256:4bd5939adb0a5fc83d76e0d69aa3c5a30bc1a5af8f9df515f44b6fc59a0a4815
      
let deployment =
      kubernetes.Deployment::{
      , metadata = kubernetes.ObjectMeta::{ name = "nginx" }
      , spec =
          Some
            kubernetes.DeploymentSpec::{
            , replicas = Some 2
            , template =
                kubernetes.PodTemplateSpec::{
                , metadata = kubernetes.ObjectMeta::{ name = "nginx" }
                , spec =
                    Some
                      kubernetes.PodSpec::{
                      , containers =
                          [ kubernetes.Container::{
                            , name = "nginx"
                            , image = Some "nginx:1.15.3"
                            , ports =
                                [ kubernetes.ContainerPort::{
                                  , containerPort = 80
                                  }
                                ]
                            }
                          ]
                      }
                }
            }
      }

in  deployment
```

Which creates the following yaml:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 2
  template:
    metadata:
      name: nginx
    spec:
      containers:
        - image: nginx:1.15.3
          name: nginx
          ports:
            - containerPort: 80
```

Dhall's lambda functions can help you break this into manageable chunks. For
example, here's a Dhall function that helps create a docker image reference:

```
let formatImage
    : Text -> Text -> Text
    = \(repository : Text) -> \(tag : Text) ->
    "${repository}:${tag}"

in formatImage "xena/christinewebsite" "latest"
```

Which outputs `xena/christinewebsite:latest` when passed to `dhall text`.

All of this adds up into a powerful toolset that lets you express Kubernetes
configuration in a way that does what you want without as many headaches.

Most of my apps on Kubernetes need only a few generic bits of configuration:

- Their name
- What port should be exposed
- The domain that this service should be exposed on
- How many replicas of the service are needed
- Which Let's Encrypt Issuer to use (currently only `"prod"` or `"staging"`)
- The [configuration variables of the service][12factorconfig]
- Any other containers that may be needed for the service

[12factorconfig]: https://12factor.net/config

From here, I defined all of the [bits and pieces][kubermemeshttp] for the
Kubernetes manifests that Dyson produces and then created a `Config` type that
helps to template them out. Here's my [`Config` type
definition][configdefinition]:

[kubermemeshttp]: https://tulpa.dev/cadey/kubermemes/src/branch/master/k8s/http
[configdefinition]: https://tulpa.dev/cadey/kubermemes/src/branch/master/k8s/app/config.dhall

```dhall
let kubernetes = ../kubernetes.dhall

in  { Type =
        { name : Text
        , appPort : Natural
        , image : Text
        , domain : Text
        , replicas : Natural
        , leIssuer : Text
        , envVars : List kubernetes.EnvVar.Type
        , otherContainers : List kubernetes.Container.Type
        }
    , default =
        { name = ""
        , appPort = 5000
        , image = ""
        , domain = ""
        , replicas = 1
        , leIssuer = "staging"
        , envVars = [] : List kubernetes.EnvVar.Type
        , otherContainers = [] : List kubernetes.Container.Type
        }
    }
```

Then I defined a `makeApp` function that creates everything I need to deploy my
stuff on Kubernetes:

```dhall
let Prelude = ../Prelude.dhall

let kubernetes = ../kubernetes.dhall

let typesUnion = ../typesUnion.dhall

let deployment = ../http/deployment.dhall

let ingress = ../http/ingress.dhall

let service = ../http/service.dhall

let Config = ../app/config.dhall

let K8sList = ../app/list.dhall

let buildService =
        \(config : Config.Type)
      -> let myService = service config

         let myDeployment = deployment config

         let myIngress = ingress config

         in  K8sList::{
             , items =
               [ typesUnion.Service myService
               , typesUnion.Deployment myDeployment
               , typesUnion.Ingress myIngress
               ]
             }

in  buildService
```

And used it to deploy the [h language website][hlang]:

[hlang]: https://h.christine.website

```dhall
let makeApp = ../app/make.dhall

let Config = ../app/config.dhall

let cfg =
      Config::{
      , name = "hlang"
      , appPort = 5000
      , image = "xena/hlang:latest"
      , domain = "h.christine.website"
      , leIssuer = "prod"
      }

in  makeApp cfg
```

Which produces the following Kubernetes config:

```yaml
apiVersion: v1
items:
  - apiVersion: v1
    kind: Service
    metadata:
      annotations:
        external-dns.alpha.kubernetes.io/cloudflare-proxied: "false"
        external-dns.alpha.kubernetes.io/hostname: h.christine.website
        external-dns.alpha.kubernetes.io/ttl: "120"
      labels:
        app: hlang
      name: hlang
      namespace: apps
    spec:
      ports:
        - port: 5000
          targetPort: 5000
      selector:
        app: hlang
      type: ClusterIP
  - apiVersion: apps/v1
    kind: Deployment
    metadata:
      name: hlang
      namespace: apps
    spec:
      replicas: 1
      selector:
        matchLabels:
          app: hlang
      template:
        metadata:
          labels:
            app: hlang
          name: hlang
        spec:
          containers:
            - image: xena/hlang:latest
              imagePullPolicy: Always
              name: web
              ports:
                - containerPort: 5000
          imagePullSecrets:
            - name: regcred
  - apiVersion: networking.k8s.io/v1beta1
    kind: Ingress
    metadata:
      annotations:
        certmanager.k8s.io/cluster-issuer: letsencrypt-prod
        kubernetes.io/ingress.class: nginx
      labels:
        app: hlang
      name: hlang
      namespace: apps
    spec:
      rules:
        - host: h.christine.website
          http:
            paths:
              - backend:
                  serviceName: hlang
                  servicePort: 5000
      tls:
        - hosts:
            - h.christine.website
          secretName: prod-certs-hlang
kind: List
```

And when I applied it on my Kubernetes cluster, it worked the first time and had
absolutely no effect on the existing configuration.

In the future, I hope to expand this to allow for multiple deployments (IE: a
chatbot running in a separate deployment than a web API the chatbot depends on
or non-web projects in general) as well as supporting multiple Kubernetes
namespaces. 

Dhall is probably the most viable replacement to Helm or other Kubernetes
templating tools I have found in recent memory. I hope that it will be used by
more people to help with configuration management, but I can understand that
that may not happen. At least it works for me. 

If you want to learn more about Dhall, I suggest checking out the following
links:

- [The Dhall Language homepage](https://dhall-lang.org)
- [Learn Dhall in Y Minutes](https://learnxinyminutes.com/docs/dhall/)
- [The Dhall Language GitHub Organization](https://github.com/dhall-lang)

I hope this was helpful and interesting. Be well.
