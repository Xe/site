---
title: Encrypted Secrets with NixOS
date: 2021-01-20
series: nixos
tags:
 - age
 - ed25519
---

# Encrypted Secrets with NixOS

One of the best things about NixOS is the fact that it's so easy to do
configuration management using it. The Nix store (where all your packages live)
has a huge flaw for secret management though: everything in the Nix store is
globally readable. This means that anyone logged into or running code on the
system could read any secret in the Nix store without any limits. This is
sub-optimal if your goal is to keep secret values secret. There have been a few
approaches to this over the years, but I want to describe how I'm doing it. 
Here are my goals and implementation for this setup and how a few other secret
management strategies don't quite pan out.

At a high level I have these goals:

* It should be trivial to declare new secrets
* Secrets should never be globally readable in any useful form
* If I restart the machine, I should not need to take manual human action to
  ensure all of the services come back online
* GPG should be avoided at all costs

As a side goal being able to roll back secret changes would also be nice.

The two biggest tools that offer a way to help with secret management on NixOS
that come to mind are NixOps and Morph.

[NixOps](https://github.com/NixOS/nixops) is a tool that helps administrators
operate NixOS across multiple servers at once. I use NixOps extensively in my
own setup. It calls deployment secrets "keys" and they are documented
[here](https://hydra.nixos.org/build/115931128/download/1/manual/manual.html#idm140737322649152).
At a high level they are declared like this:

```nix
deployment.keys.example = {
  text = "this is a super sekrit value :)";
  user = "example";
  group = "keys";
  permissions = "0400";
};
```

This will create a new secret in `/run/keys` that will contain our super secret
value.

[Wait, isn't `/run` an ephemeral filesystem? What happens when the system
reboots?](conversation://Mara/hmm)

Let's make an example system and find out! So let's say we have that `example`
secret from earlier and want to use it in a job. The job definition could look
something like this:

```nix
# create a service-specific user
users.users.example.isSystemUser = true;

# without this group the secret can't be read
users.users.example.extraGroups = [ "keys" ]; 

systemd.services.example = {
  wantedBy = [ "multi-user.target" ];
  after = [ "example-key.service" ];
  wants = [ "example-key.service" ];
  
  serviceConfig.User = "example";
  serviceConfig.Type = "oneshot";
  
  script = ''
    stat /run/keys/example
  '';
};
```

This creates a user called `example` and gives it permission to read deployment
keys. It also creates a systemd service called `example.service` and runs
[`stat(1)`](https://linux.die.net/man/1/stat) to show the permissions of the
service and the key file. It also runs as our `example` user. To avoid systemd
thinking our service failed, we're also going to mark it as a
[oneshot](https://www.digitalocean.com/community/tutorials/understanding-systemd-units-and-unit-files#the-service-section).

Altogether it could look something like
[this](https://gist.github.com/Xe/4a71d7741e508d9002be91b62248144a). Let's see
what `systemctl` has to report:

```console
$ nixops ssh -d blog-example pa -- systemctl status example
● example.service
     Loaded: loaded (/nix/store/j4a8f6mnaw3v4sz7dqlnz95psh72xglw-unit-example.service/example.service; enabled; vendor preset: enabled)
     Active: inactive (dead) since Wed 2021-01-20 20:53:54 UTC; 37s ago
    Process: 2230 ExecStart=/nix/store/1yg89z4dsdp1axacqk07iq5jqv58q169-unit-script-example-start/bin/example-start (code=exited, status=0/SUCCESS)
   Main PID: 2230 (code=exited, status=0/SUCCESS)
         IP: 0B in, 0B out
        CPU: 3ms

Jan 20 20:53:54 pa example-start[2235]:   File: /run/keys/example
Jan 20 20:53:54 pa example-start[2235]:   Size: 31                Blocks: 8          IO Block: 4096   regular file
Jan 20 20:53:54 pa example-start[2235]: Device: 18h/24d        Inode: 37428       Links: 1
Jan 20 20:53:54 pa example-start[2235]: Access: (0400/-r--------)  Uid: (  998/ example)   Gid: (   96/    keys)
Jan 20 20:53:54 pa example-start[2235]: Access: 2021-01-20 20:53:54.010554201 +0000
Jan 20 20:53:54 pa example-start[2235]: Modify: 2021-01-20 20:53:54.010554201 +0000
Jan 20 20:53:54 pa example-start[2235]: Change: 2021-01-20 20:53:54.398103181 +0000
Jan 20 20:53:54 pa example-start[2235]:  Birth: -
Jan 20 20:53:54 pa systemd[1]: example.service: Succeeded.
Jan 20 20:53:54 pa systemd[1]: Finished example.service.
```

So what happens when we reboot? I'll force a reboot in my hypervisor and we'll
find out:

```console
$ nixops ssh -d blog-example pa -- systemctl status example
● example.service
     Loaded: loaded (/nix/store/j4a8f6mnaw3v4sz7dqlnz95psh72xglw-unit-example.service/example.service; enabled; vendor preset: enabled)
     Active: inactive (dead)
```

The service is inactive. Let's see what the status of `example-key.service` is:

```console
$ nixops ssh -d blog-example pa -- systemctl status example-key
● example-key.service
     Loaded: loaded (/nix/store/ikqn64cjq8pspkf3ma1jmx8qzpyrckpb-unit-example-key.service/example-key.service; linked; vendor preset: enabled)
     Active: activating (start-pre) since Wed 2021-01-20 20:56:05 UTC; 3min 1s ago
Cntrl PID: 610 (example-key-pre)
         IP: 0B in, 0B out
         IO: 116.0K read, 0B written
      Tasks: 4 (limit: 2374)
     Memory: 1.6M
        CPU: 3ms
     CGroup: /system.slice/example-key.service
             ├─610 /nix/store/kl6lr3czkbnr6m5crcy8ffwfzbj8a22i-bash-4.4-p23/bin/bash -e /nix/store/awx1zrics3cal8kd9c5d05xzp5ikazlk-unit-script-example-key-pre-start/bin/example-key-pre-start
             ├─619 /nix/store/kl6lr3czkbnr6m5crcy8ffwfzbj8a22i-bash-4.4-p23/bin/bash -e /nix/store/awx1zrics3cal8kd9c5d05xzp5ikazlk-unit-script-example-key-pre-start/bin/example-key-pre-start
             ├─620 /nix/store/kl6lr3czkbnr6m5crcy8ffwfzbj8a22i-bash-4.4-p23/bin/bash -e /nix/store/awx1zrics3cal8kd9c5d05xzp5ikazlk-unit-script-example-key-pre-start/bin/example-key-pre-start
             └─621 inotifywait -qm --format %f -e create,move /run/keys

Jan 20 20:56:05 pa systemd[1]: Starting example-key.service...
```

The service is blocked waiting for the keys to exist. We have to populate the
keys with `nixops send-keys`:

```console
$ nixops send-keys -d blog-example
pa> uploading key ‘example’...
```

Now when we check on `example.service`, we get the following:

```console
$ nixops ssh -d blog-example pa -- systemctl status example
● example.service
     Loaded: loaded (/nix/store/j4a8f6mnaw3v4sz7dqlnz95psh72xglw-unit-example.service/example.service; enabled; vendor preset: enabled)
     Active: inactive (dead) since Wed 2021-01-20 21:00:24 UTC; 32s ago
    Process: 954 ExecStart=/nix/store/1yg89z4dsdp1axacqk07iq5jqv58q169-unit-script-example-start/bin/example-start (code=exited, status=0/SUCCESS)
   Main PID: 954 (code=exited, status=0/SUCCESS)
         IP: 0B in, 0B out
        CPU: 3ms

Jan 20 21:00:24 pa example-start[957]:   File: /run/keys/example
Jan 20 21:00:24 pa example-start[957]:   Size: 31                Blocks: 8          IO Block: 4096   regular file
Jan 20 21:00:24 pa example-start[957]: Device: 18h/24d        Inode: 27774       Links: 1
Jan 20 21:00:24 pa example-start[957]: Access: (0400/-r--------)  Uid: (  998/ example)   Gid: (   96/    keys)
Jan 20 21:00:24 pa example-start[957]: Access: 2021-01-20 21:00:24.588494730 +0000
Jan 20 21:00:24 pa example-start[957]: Modify: 2021-01-20 21:00:24.588494730 +0000
Jan 20 21:00:24 pa example-start[957]: Change: 2021-01-20 21:00:24.606495751 +0000
Jan 20 21:00:24 pa example-start[957]:  Birth: -
Jan 20 21:00:24 pa systemd[1]: example.service: Succeeded.
Jan 20 21:00:24 pa systemd[1]: Finished example.service.
```

This means that NixOps secrets require _manual human intervention_ in order to
repopulate them on server boot. If your server went offline overnight due to an
unexpected issue, your services using those keys could be stuck offline until
morning. This is undesirable for a number of reasons. This plus the requirement
for the `keys` group (which at time of writing was undocumented) to be added to
service user accounts means that while they do work, they are not very
ergonomic.

[You can read secrets from files using something like
`deployment.keys.example.text = "${builtins.readFile ./secrets/example.env}"`,
but it is kind of a pain to have to do that. It would be better to just
reference the secrets by filesystem paths in the first
place.](conversation://Mara/hacker)

On the other hand [Morph](https://github.com/DBCDK/morph) gets this a bit
better. It is sadly even less documented than NixOps is, but it offers a similar
experience via [deployment
secrets](https://github.com/DBCDK/morph/blob/master/examples/secrets.nix). The
main differences that Morph brings to the table are taking paths to secrets and
allowing you to run an arbitrary command on the secret being uploaded. Secrets
are also able to be put anywhere on the disk, meaning that when a host reboots it
will come back up with the most recent secrets uploaded to it.

However, like NixOps, Morph secrets don't have the ability to be rolled back.
This means that if you mess up a secret value you better hope you have the old
information somewhere. This violates what you'd expect from a NixOS machine.

So given these examples, I thought it would be interesting to explore what the
middle path could look like. I chose to use
[age](https://github.com/FiloSottile/age) for encrypting secrets in the Nix
store as well as using SSH host keys to ensure that every secret is decryptable
at runtime by _that machine only_. If you get your hands on the secret
cyphertext, it should be unusable to you.

One of the harder things here will be keeping a list of all of the server host
keys. Recently I added a
[hosts.toml](https://github.com/Xe/nixos-configs/blob/master/ops/metadata/hosts.toml)
file to my config repo for autoconfiguring my WireGuard overlay network. It was
easy enough to add all the SSH host keys for each machine using a command like
this to get them:

[We will cover how this WireGuard overlay works in a future post.](conversation://Mara/hacker)

```console
$ nixops ssh-for-each -d hexagone -- cat /etc/ssh/ssh_host_ed25519_key.pub 
firgu....> ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIB8+mCR+MEsv0XYi7ohvdKLbDecBtb3uKGQOPfIhdj3C root@nixos
chrysalis> ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGDA5iXvkKyvAiMEd/5IruwKwoymC8WxH4tLcLWOSYJ1 root@chrysalis
lufta....> ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIMADhGV0hKt3ZY+uBjgOXX08txBS6MmHZcSL61KAd3df root@lufta
keanu....> ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGDZUmuhfjEIROo2hog2c8J53taRuPJLNOtdaT8Nt69W root@nixos
```

age lets you use SSH keys for decryption, so I added these keys to my
`hosts.toml` and ended up with something like
[this](https://github.com/Xe/nixos-configs/commit/14726e982001e794cd72afa1ece209eed58d3f38#diff-61d1d8dddd71be624c0d718be22072c950ec31c72fded8a25094ea53d94c8185).

Now we can encrypt secrets on the host machine and safely put them in the Nix
store because they will be readable to each target machine with a command like
this:

```shell
age -d -i /etc/ssh/ssh_host_ed25519_key -o $dest $src
```

From here it's easy to make a function that we can use for generating new
encrypted secrets in the Nix store. First we need to import the host metadata
from the toml file:

```nix
let
  cfg = config.within.secrets;
  metadata = lib.importTOML ../../ops/metadata/hosts.toml;

  mkSecretOnDisk = name:
    { source, ... }:
    pkgs.stdenv.mkDerivation {
      name = "${name}-secret";
      phases = "installPhase";
      buildInputs = [ pkgs.age ];
      installPhase =
        let key = metadata.hosts."${config.networking.hostName}".ssh_pubkey;
        in ''
          age -a -r "${key}" -o $out ${source}
        '';
    };
```

And then we can generate systemd oneshot jobs with something like this:

```nix
  mkService = name:
    { source, dest, owner, group, permissions, ... }: {
      description = "decrypt secret for ${name}";
      wantedBy = [ "multi-user.target" ];

      serviceConfig.Type = "oneshot";

      script = with pkgs; ''
        rm -rf ${dest}
        ${age}/bin/age -d -i /etc/ssh/ssh_host_ed25519_key -o ${dest} ${
          mkSecretOnDisk name { inherit source; }
        }

        chown ${owner}:${group} ${dest}
        chmod ${permissions} ${dest}
      '';
    };
```

And from there we just need some [boring
boilerplate](https://github.com/Xe/nixos-configs/blob/master/common/crypto/default.nix#L8-L38)
to define a secret type. Then we declare the secret type and its invocation:

```nix
in {
  options.within.secrets = mkOption {
    type = types.attrsOf secret;
    description = "secret configuration";
    default = { };
  };

  config.systemd.services = let
    units = mapAttrs' (name: info: {
      name = "${name}-key";
      value = (mkService name info);
    }) cfg;
  in units;
}
```

And we have ourself a NixOS module that allows us to:

* Trivially declare new secrets
* Make secrets in the Nix store useless without the key
* Make every secret be transparently decrypted on startup
* Avoid the use of GPG
* Roll back secrets like any other configuration change

Declaring new secrets works like this (as stolen from [the service definition
for the website you are reading right now](https://github.com/Xe/nixos-configs/blob/master/common/services/xesite.nix#L35-L41)):

```nix
within.secrets.example = {
  source = ./secrets/example.env;
  dest = "/var/lib/example/.env";
  owner = "example";
  group = "nogroup";
  permissions = "0400";
};
```

Barring some kind of cryptographic attack against age, this should allow the
secrets to be stored securely. I am working on a way to make this more generic.
This overall approach was inspired by [agenix](https://github.com/ryantm/agenix)
but made more specific for my needs. I hope this approach will make it easy for
me to manage these secrets in the future.
