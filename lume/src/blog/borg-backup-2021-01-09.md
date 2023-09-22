---
title: "How to Set Up Borg Backup on NixOS"
date: 2021-01-09
series: howto
tags:
  - nixos
  - borgbackup
---

[Borg Backup](https://www.borgbackup.org/) is a encrypted, compressed,
deduplicated backup program for multiple platforms including Linux. This
combined with the [NixOS options for configuring
Borg Backup](https://search.nixos.org/options?channel=20.09&show=services.borgbackup.jobs.%3Cname%3E.paths&from=0&size=30&sort=relevance&query=services.borgbackup.jobs)
allows you to backup on a schedule and restore from those backups when you need
to.

Borg Backup works with local files, remote servers and there are even [cloud
hosts](https://www.borgbackup.org/support/commercial.html) that specialize in
hosting your backups. In this post we will cover how to set up a backup job on a
server using [BorgBase](https://www.borgbase.com/)'s free tier to host the
backup files.

## Setup

You will need a few things:

- A free BorgBase account
- A server running NixOS
- A list of folders to back up
- A list of folders to NOT back up

First, we will need to create a SSH key for root to use when connecting to
BorgBase. Open a shell as root on the server and make a `borgbackup` folder in
root's home directory:

```shell
mkdir borgbackup
cd borgbackup
```

Then create a SSH key that will be used to connect to BorgBase:

```shell
ssh-keygen -f ssh_key -t ed25519 -C "Borg Backup"
```

Ignore the SSH key password because at this time the automated Borg Backup job
doesn't allow the use of password-protected SSH keys.

Now we need to create an encryption passphrase for the backup repository. Run
this command to generate one using [xkcdpass](https://pypi.org/project/xkcdpass/):

```shell
nix-shell -p python39Packages.xkcdpass --run 'xkcdpass -n 12' > passphrase
```

## BorgBase Setup

Now that we have the basic requirements out of the way, let's configure BorgBase
to use that SSH key. In the BorgBase UI click on the Account tab in the upper
right and open the SSH key management window. Click on Add Key and paste in the
contents of `./ssh_key.pub`. Name it after the hostname of the server you are
working on. Click Add Key and then go back to the Repositories tab in the upper
right.

Click New Repo and name it after the hostname of the server you are working on.
Select the key you just created to have full access. Choose the region of the
backup volume and then click Add Repository.

On the main page copy the repository path with the copy icon next to your
repository in the list. You will need this below. Attempt to SSH into the backup
repo in order to have ssh recognize the server's host key:

```shell
ssh -i ./ssh_key o6h6zl22@o6h6zl22.repo.borgbase.com
```

Then accept the host key and press control-c to terminate the SSH connection.

## NixOS Configuration

In your `configuration.nix` file, add the following block:

```nix
services.borgbackup.jobs."borgbase" = {
  paths = [
    "/var/lib"
    "/srv"
    "/home"
  ];
  exclude = [
    # very large paths
    "/var/lib/docker"
    "/var/lib/systemd"
    "/var/lib/libvirt"
    
    # temporary files created by cargo and `go build`
    "**/target"
    "/home/*/go/bin"
    "/home/*/go/pkg"
  ];
  repo = "o6h6zl22@o6h6zl22.repo.borgbase.com:repo";
  encryption = {
    mode = "repokey-blake2";
    passCommand = "cat /root/borgbackup/passphrase";
  };
  environment.BORG_RSH = "ssh -i /root/borgbackup/ssh_key";
  compression = "auto,lzma";
  startAt = "daily";
};
```

Customize the paths and exclude lists to your needs. Once you are satisfied,
rebuild your NixOS system using `nixos-rebuild`:

```shell
nixos-rebuild switch
```

And then you can fire off an initial backup job with this command:

```shell
systemctl start borgbackup-job-borgbase.service
```

Monitor the job with this command:

```shell
journalctl -fu borgbackup-job-borgbase.service
```

The first backup job will always take the longest to run. Every incremental
backup after that will get smaller and smaller. By default, the system will
create new backup snapshots every night at midnight local time.

## Restoring Files

To restore files, first figure out when you want to restore the files from.
NixOS includes a wrapper script for each Borg job you define. you can mount your
backup archive using this command:

```
mkdir mount
borg-job-borgbase mount o6h6zl22@o6h6zl22.repo.borgbase.com:repo ./mount
```

Then you can explore the backup (and with it each incremental snapshot) to
your heart's content and copy files out manually. You can look through each
folder and copy out what you need.

When you are done you can unmount it with this command:

```
borg-job-borgbase umount /root/borgbase/mount
```

---

And that's it! You can get more fancy with nixops using a setup [like
this](https://github.com/Xe/nixos-configs/blob/master/common/services/backup.nix).
In general though, you can get away with this setup. It may be a good idea to
copy down the encryption passphrase onto paper and put it in a safe space like a
safety deposit box.

For more information about Borg Backup on NixOS, see [the relevant chapter of
the NixOS
manual](https://nixos.org/manual/nixos/stable/index.html#module-borgbase) or
[the list of borgbackup
options](https://search.nixos.org/options?channel=20.09&query=services.borgbackup.jobs)
that you can pick from.

I hope this is able to help.
