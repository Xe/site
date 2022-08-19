---
title: "Nix Flakes on WSL"
date: 2022-05-01
series: nix-flakes
tags:
 - nixos
 - wsl
vod:
  youtube: https://youtu.be/VzQ_NwFJObc
  twitch: https://www.twitch.tv/videos/1464781566
---

About five years ago, Microsoft released the Windows Subsystem for Linux
([WSL](https://docs.microsoft.com/en-us/windows/wsl/)). This allows you to run
Linux programs on a Windows machine. When they released WSL version 2 in 2019,
this added support for things like Docker and systemd. As a result, this is
enough to run NixOS on Windows.

<xeblog-conv name="Mara" mood="hacker">This will give you an environment to run
Nix and Nix Flakes commands with. You can use this to follow along with this
series without having to install NixOS on a VM or cloud server. This is going to
retread a bunch of ground from the first article. If you have been following
along through this entire series, once you get to the point where you convert
the install to flakes there isn't much more new material here.
</xeblog-conv>

## Installation

Head to the NixOS-WSL [releases
page](https://github.com/nix-community/NixOS-WSL/releases/) and download the
`nixos-wsl-installer-fixed.tar.gz` file to your Downloads folder.

Then open Powershell and make a folder called `WSL`:

```powershell
New-Item -Path .\WSL -ItemType Directory
```

<xeblog-conv name="Mara" mood="hacker">It's worth noting that Powershell does
have a bunch of aliases for common coreutils commands to the appropriate
Powershell CMDlets. However these aliases are <b>NOT</b> flag-compatible and use
the Powershell semantics instead of the semantics of the command it is aliasing.
This will bite you when you use commands like <code>wget</code> out of instinct
to download things. In order to avoid your muscle memory betraying you, the
Powershell CMDlets are shown here in their full overly verbose glory.
</xeblog-conv>

Then enter the directory with `Set-Location`:

```powershell
Set-Location -Path .\WSL
```

<xeblog-conv name="Mara" mood="hacker">This directory is where the NixOS root
filesystem will live. If you want to put this somewhere else, feel free to.
Somewhere in `%APPDATA%` will work, just as long as it's on an NTFS volume
somewhere.
</xeblog-conv>

Make a folder for the NixOS filesystem:

```powershell
New-Item -Path .\NixOS -ItemType Directory
```

Then install the NixOS root image with the `wsl` command:

```powershell
wsl --import NixOS .\NixOS\ ..\Downloads\nixos-wsl-installer-fixed.tar.gz --version 2
```

And start NixOS once to have it install itself:

```powershell
wsl -d NixOS
```

Once that finishes, press control-D (or use the `exit` command) to exit out of
NixOS and restart the WSL virtual machine:

```powershell
exit
wsl --shutdown
wsl -d NixOS
```

And then you have yourself a working NixOS environment! It's very barebones, but
we can use it to test the `nix run` command against our gohello command:

```console
$ nix run github:Xe/gohello
Hello reader!
```

## Local Services

We can also use this NixOS environment to run a local nginx server. Open
`/etc/nixos/configuration.nix`:

```nix
{ lib, pkgs, config, modulesPath, ... }:

with lib;
let
  nixos-wsl = import ./nixos-wsl;
in
{
  imports = [
    "${modulesPath}/profiles/minimal.nix"

    nixos-wsl.nixosModules.wsl
  ];

  wsl = {
    enable = true;
    automountPath = "/mnt";
    defaultUser = "nixos";
    startMenuLaunchers = true;

    # Enable integration with Docker Desktop (needs to be installed)
    # docker.enable = true;
  };

  # Enable nix flakes
  nix.package = pkgs.nixFlakes;
  nix.extraOptions = ''
    experimental-features = nix-command flakes
  '';
}
```

Right after the `wsl` block, add this nginx configuration to the file:

```nix
services.nginx.enable = true;
services.nginx.virtualHosts."test.local.cetacean.club" = {
  root = "/srv/http/test.local.cetacean.club";
};
```

This will create an nginx configuration that points the domain
`test.local.cetacean.club` to the contents of the folder `/srv/http/test.local.cetacean.club`.

<xeblog-conv name="Mara" mood="hacker">The <code>/srv</code> folder is set aside
for site-specific data, which is code for "do whatever you want with this
folder". In many cases people make a separate <code>/srv/http</code> folder and
put each static subdomain in its own folder under that, however I am also told
that it is idiomatic to put stuff in <code>/var/www</code>. Pick your poison.
</xeblog-conv>

Then you can test the web server with the `curl` command:

```console
$ curl http://test.local.cetacean.club
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
```

This is good! Nginx is running and since we haven't created the folder with our
website content yet, this 404 means that it can't find it! Let's create the
folder so that nginx has permission to it and we can modify things in it:

```
sudo mkdir -p /srv/http/test.local.cetacean.club
sudo chown nixos:nginx /srv/http/test.local.cetacean.club
```

Finally we can make an amazing website. Open
`/srv/http/test.local.cetacean.club/index.html` in nano:

```
nano /srv/http/test.local.cetacean.club/index.html
```

And paste in this HTML:

```html
<title>amazing website xD</title>
<h1>look at my AMAZING WEBSITE</h1>
It's so cool *twerks*
```

<xeblog-conv name="Mara" mood="hacker">This doesn't have to just be artisanal
handcrafted HTML in bespoke folders either. You can set the <code>root</code> of
a nginx virtual host to point to a Nix package as well. This will allow you to
automatically generate your website somehow and deploy it with the rest of the
system. Including being able to roll back changes.</xeblog-conv>

And then you can see it show up with `curl`:

```console
$ curl http://test.local.cetacean.club
<title>amazing website xD</title>
<h1>look at my AMAZING WEBSITE</h1>
It's so cool *twerks*
```

You can also check this out in a browser:

![a browser window titled "amazing website xD" with the header "look at my
AMAZING WEBSITE" and content of "It's so cool
\*twerks\*"](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-04-23+141937.png)

## Installing `gohello`

To install the `gohello` service, first we will need to convert this machine to
use NixOS flakes. We can do that really quick and easy by adding this file to
`/etc/nixos/flake.nix`:

<xeblog-conv name="Mara" mood="happy">Do this as root!</xeblog-conv>

```nix
{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
  };
  
  outputs = { self, nixpkgs, ... }: {
    nixosConfigurations.nixos = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        
        # add things here
      ];
    }; 
  };
}
```

Then run `nix flake check` to make sure everything is okay:

```
sudo nix flake check /etc/nixos
```

And finally activate the new configuration with flakes:

```
sudo nixos-rebuild switch
```

<xeblog-conv name="Mara" mood="hmm">Why don't you have the <code>--flake</code>
flag here? Based on what I read in the documentation, I thought you had to have
it there.</xeblog-conv>

<xeblog-conv name="Cadey" mood="enby"><code>nixos-rebuild</code> will
auomatically detect flakes in <code>/etc/nixos</code>. The only major thing it
cares about is the hostname matching. If you want to customize the hostname of
the WSL VM, change the <code>nixos</code> in
<code>nixosConfigurations.nixos</code> above and set
<code>networking.hostName</code> to the value you want to use. To use flakes
explicitly, pass <code>--flake /etc/nixos#hostname</code> to your
<code>nixos-rebuild</code> call. 
</xeblog-conv>

After it thinks for a bit, you should notice that nothing happened. This is
good, we have just converted the system over to using Nix flakes instead of the
classic `nix-channel` rebuild method.

To get `gohello` in the system, first we need to add `git` to the commands
available on the system in `configuration.nix`:

```nix
environment.systemPackages = with pkgs; [ git ];
```

Then we can add `gohello` to our system flake:

```nix
{
  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    # XXX(Xe): this URL may change for you, such as github:Xe/gohello-http
    gohello.url = "git+https://tulpa.dev/cadey/gohello-http?ref=main";
  };

  outputs = { self, nixpkgs, gohello, ... }: {
    nixosConfigurations.nixos = nixpkgs.lib.nixosSystem {
      system = "x86_64-linux";
      modules = [
        ./configuration.nix
        
        # add things here
        gohello.nixosModule
        ({ pkgs, ... }: {
          xeserv.services.gohello.enable = true;
        })
      ];
    }; 
  };
}
```

<xeblog-conv name="Mara" mood="hacker">The block of code under
<code>gohello.nixosModule</code> is an inline NixOS module. If we put
<code>gohello.nixosModule</code> before the <code>./configuration.nix</code>
reference, we could put the <code>xeserv.services.gohello.enable = true;</code>
line inside <code>./configuration.nix</code>. This is an exercise for the
reader.</xeblog-conv>

And rebuild the system with `gohello` enabled:

```
sudo nixos-rebuild switch
```

Finally, poke it with `curl`:

```console
$ curl http://gohello.local.cetacean.club
hello world :)
```

To update it, update the flake inputs in `/etc/nixos` and run `nixos-rebuild`:

```
sudo nix flake update /etc/nixos
sudo nixos-rebuild switch
```

---

And from here you can do whatever you want with NixOS. You can use
[containers](https://nixos.org/manual/nixos/stable/#ch-containers), set up
arbitrary services, or plan for world domination as normal.

<xeblog-conv name="Numa" mood="delet">I thought it was "to save the world from
devastation", not "to plan for world domination". Who needs a monopoly on
violence for world domination when you have Nix expressions?</xeblog-conv>

<xeblog-conv name="Cadey" mood="coffee">Siiiiiiiiiiiiiiiiiigh.</xeblog-conv>

I will use this setup in future posts to make this more accessible and easy to
hack at without having to have a dedicated NixOS machine laying around.
