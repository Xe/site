---
title: Nixops Services on Your Home Network
date: 2020-11-09
series: howto
tags:
 - nixos
 - systemd
---

My homelab has a few NixOS machines. Right now they mostly run services inside
Docker, because that has been what I have done for years. This works fine, but
persistent state gets annoying*. NixOS has a tool called
[Nixops](https://releases.nixos.org/nixops/nixops-1.7/manual/manual.html) that
allows you to push configurations to remote machines. I use this for managing my
fleet of machines, and today I'm going to show you how to create service
deployments with Nixops and push them to your servers.

[Pedantically, Docker offers <a
href="https://releases.nixos.org/nixops/nixops-1.7/manual/manual.html">volumes</a>
to simplify this, but it is very easy to accidentally delete Docker volumes.
Plain disk files like we are going to use today are a bit simpler than docker
volumes, and thusly a bit harder to mess up.](conversation://Mara/hacker)

## Parts of a Service

For this example, let's deploy a chatbot. To make things easier, let's assume
the following about this chatbot:

- The chatbot has a git repo somewhere
- The chatbot's git repo has a `default.nix` that builds the service and
  includes any supporting files it might need
- The chatbot reads its configuration from environment variables which may
  contain secret values (API keys, etc.)
- The chatbot stores any temporary files in its current working directory
- The chatbot is "well-behaved" (for some definition of "well-behaved")

I will also need to assume that you have a git repo (or at least a folder) with
all of your configuration similar to [mine](https://github.com/Xe/nixos-configs).

For this example I'm going to use [withinbot](https://github.com/Xe/withinbot)
as the service we will deploy via Nixops. withinbot is a chatbot that I use on
my own Discord guild that does a number of vital functions including supplying
amusing facts about printers:

```
     <Cadey~> ~printerfact
<Within[BOT]> @Cadey~ Printers, especially older printers, do get cancer. Many
              times this disease can be treated successfully
```

[To get your own amusing facts about printers, see <a
href="https://printerfacts.cetacean.club">here</a> or for using its API, call <a
href="https://printerfacts.cetacean.club/fact">`/fact`</a>. This API has no
practical rate limits, but please don't test that.](conversation://Mara/hacker)

## Service Definition

We will need to do a few major things for defining this service:

1. Add the bot code as a package
1. Create a "services" folder for the service modules
1. Create a user account for the service
1. Set up a systemd unit for the service
1. Configure the secrets using [Nixops
   keys](https://releases.nixos.org/nixops/nixops-1.7/manual/manual.html#idm140737322342384)

### Add the Code as a Package

In order for the program to be installed to the remote system, you need to tell
the system how to import it. There's many ways to do this, but the cheezy way is
to add the packages to
[`nixpkgs.config.packageOverrides`](https://nixos.org/manual/nixos/stable/#sec-customising-packages)
like this:

```nix
nixpkgs.config = {
  packageOverrides = pkgs: {
    within = {
      withinbot = import (builtins.fetchTarball 
        "https://github.com/Xe/withinbot/archive/main.tar.gz") { };
    };
  };
};
```

And now we can access it as `pkgs.within.withinbot` in the rest of our config.

[In production circumstances you should probably use <a
href="https://nixos.org/manual/nixpkgs/stable/#chap-pkgs-fetchers">a fetcher
that locks to a specific version</a> using unique URLs and hashing, but this
will work enough to get us off the ground in this
example.](conversation://Mara/hacker)

### Create a "services" Folder

In your configuration folder, create a folder that you will use for these
service definitions. I made mine in `common/services`. In that folder, create a
`default.nix` with the following contents:

```nix
{ config, lib, ... }:

{
  imports = [ ./withinbot.nix ];

  users.groups.within = {};
}
```

The group listed here is optional, but I find that having a group like that can
help you better share resources and files between services.

Now we need a folder for storing secrets. Let's create that under the services
folder:

```console
$ mkdir secrets
```

And let's also add a gitignore file so that we don't accidentally commit these
secrets to the repo:

```gitignore
# common/services/secrets/.gitignore
*
```

Now we can put any secrets we want in the secrets folder without the risk of
committing them to the git repo.

### Service Manifest

Let's create `withinbot.nix` and set it up:

```nix
{ config, lib, pkgs, ... }:
with lib; {
  options.within.services.withinbot.enable =
    mkEnableOption "Activates Withinbot (the furryhole chatbot)";

  config = mkIf config.within.services.withinbot.enable {
    
  };
}
```

This sets up an option called `within.services.withinbot.enable` which will only
add the service configuration if that option is set to `true`. This will allow
us to define a lot of services that are available, but none of their config will
be active unless they are explicitly enabled.

Now, let's create a user account for the service:

```nix
# ...
  config = ... {
    users.users.withinbot = {
      createHome = true;
      description = "github.com/Xe/withinbot";
      isSystemUser = true;
      group = "within";
      home = "/srv/within/withinbot";
      extraGroups = [ "keys" ];
    };
  };
# ...
```

This will create a user named `withinbot` with the home directory
`/srv/within/withinbot`, the group `within` and also in the group `keys` so the
withinbot user can read deployment secrets. 

Now let's add the deployment secrets to the configuration:

```nix
# ...
  config = ... {
    users.users.withinbot = { ... };
    
    deployment.keys.withinbot = {
      text = builtins.readFile ./secrets/withinbot.env;
      user = "withinbot";
      group = "within";
      permissions = "0640";
    };
  };
# ...
```

Assuming you have the configuration at `./secrets/withinbot.env`, this will
register the secrets into `/run/keys/withinbot` and also create a systemd
oneshot service named `withinbot-key`. This allows you to add the secret's
existence as a condition for withinbot to run. However, Nixops puts these keys
in `/run`, which by default is mounted using a temporary memory-only filesystem,
meaning these keys will need to be re-added to machines when they are rebooted.
Fortunately, `nixops reboot` will automatically add the keys back after the
reboot succeeds.

Now that we have everything else we need, let's add the service configuration:

```nix
# ...
  config = ... {
    users.users.withinbot = { ... };
    deployment.keys.withinbot = { ... };
    
    systemd.services.withinbot = {
      wantedBy = [ "multi-user.target" ];
      after = [ "withinbot-key.service" ];
      wants = [ "withinbot-key.service" ];

      serviceConfig = {
        User = "withinbot";
        Group = "within";
        Restart = "on-failure"; # automatically restart the bot when it dies
        WorkingDirectory = "/srv/within/withinbot";
        RestartSec = "30s";
      };

      script = let withinbot = pkgs.within.withinbot;
      in ''
        # load the environment variables from /run/keys/withinbot
        export $(grep -v '^#' /run/keys/withinbot | xargs)
        # service-specific configuration
        export CAMPAIGN_FOLDER=${withinbot}/campaigns
        # kick off the chatbot
        exec ${withinbot}/bin/withinbot
      '';
    };
  };
# ...
```

This will create the systemd configuration for the service so that it starts on
boot, waits to start until the secrets have been loaded into it, runs withinbot
as its own user and in the `within` group, and throttles the service restart so
that it doesn't incur Discord rate limits as easily. This will also put all
withinbot logs in journald, meaning that you can manage and monitor this service
like you would any other systemd service.

## Deploying the Service

In your target server's `configuration.nix` file, add an import of your services
directory:

```nix
{
  # ...
  imports = [
    # ...
    /home/cadey/code/nixos-configs/common/services
  ];
  # ...
}
```

And then enable the withinbot service:

```nix
{
  # ...
  within.services = {
    withinbot.enable = true;
  };
  # ...
}
```

[Make that a block so you can enable multiple services at once like <a
href="https://github.com/Xe/nixos-configs/blob/e111413e8b895f5a117dea534b17fc9d0b38d268/hosts/chrysalis/configuration.nix#L93-L96">this</a>!](conversation://Mara/hacker)

Now you are free to deploy it to your network with `nixops deploy`:

```console
$ nixops deploy -d hexagone
```

<video controls width="100%">
    <source src="https://cdn.xeiaso.net/file/christine-static/img/nixops/tmp.Tr7HTFFd2c.webm"
            type="video/webm">
    <source src="https://cdn.xeiaso.net/file/christine-static/img/nixops/tmp.Tr7HTFFd2c.mp4"
            type="video/mp4">
    Sorry, your browser doesn't support embedded videos.
</video>


And then you can verify the service is up with `systemctl status`:

```console
$ nixops ssh -d hexagone chrysalis -- systemctl status withinbot
● withinbot.service
     Loaded: loaded (/nix/store/7ab7jzycpcci4f5wjwhjx3al7xy85ka7-unit-withinbot.service/withinbot.service; enabled; vendor preset: enabled)
     Active: active (running) since Mon 2020-11-09 09:51:51 EST; 2h 29min ago
   Main PID: 12295 (withinbot)
         IP: 0B in, 0B out
      Tasks: 13 (limit: 4915)
     Memory: 7.9M
        CPU: 4.456s
     CGroup: /system.slice/withinbot.service
             └─12295 /nix/store/qpq281hcb1grh4k5fm6ksky6w0981arp-withinbot-0.1.0/bin/withinbot

Nov 09 09:51:51 chrysalis systemd[1]: Started withinbot.service.
```

---

This basic template is enough to expand out to anything you would need and is
what I am using for my own network. This should be generic enough for most of
your needs. Check out the [NixOS manual](https://nixos.org/manual/nixos/stable/)
for more examples and things you can do with this. The [Nixops
manual](https://releases.nixos.org/nixops/nixops-1.7/manual/manual.html) is also
a good read. It can also set up deployments with VirtualBox, libvirtd, AWS,
Digital Ocean, and even Google Cloud.

The cloud is the limit! Be well.
