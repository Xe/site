---
title: Discord Webhooks via NixOS and Systemd Timers
date: 2020-11-30
series: howto
tags:
  - nixos
  - discord
  - systemd
---

# Discord Webhooks via NixOS and Systemd Timers

Recently I needed to set up a Discord message on a cronjob as a part of
moderating a guild I've been in for years. I've done this before using
[cronjobs](/blog/howto-automate-discord-webhook-cron-2018-03-29), however this
time we will be using [NixOS](https://nixos.org/) and [systemd
timers](https://wiki.archlinux.org/index.php/Systemd/Timers). Here's what you
will need to follow along:

- A machine running NixOS
- A [Discord](https://discord.com/) account
- A
  [webhook](https://support.discord.com/hc/en-us/articles/228383668-Intro-to-Webhooks)
  configured for a channel
- A message you want to send to Discord
  
[If you don't have moderation permissions in any guilds, make your own for
testing! You will need the "Manage Webhooks" permission to create a
webhook.](conversation://Mara/hacker)

## Setting Up Timers

systemd timers are like cronjobs, except they trigger systemd services instead
of shell commands. For this example, let's create a daily webhook reminder to
check on your Animal Crossing island at 9 am.

Let's create the systemd service at the end of the machine's
`configuration.nix`:

```nix
systemd.services.acnh-island-check-reminder = {
  serviceConfig.Type = "oneshot";
  script = ''
    MESSAGE="It's time to check on your island! Check those stonks!"
    WEBHOOK="${builtins.readFile /home/cadey/prefix/secrets/acnh-webhook-secret}"
    USERNAME="Domo"
    
    ${pkgs.curl}/bin/curl \
      -X POST \
      -F "content=$MESSAGE" \
      -F "username=$USERNAME" \
      "$WEBHOOK"
  '';
};
```

[This service is a <a href="https://stackoverflow.com/a/39050387">oneshot</a>
unit, meaning systemd will launch this once and not expect it to always stay
running.](conversation://Mara/hacker)

Now let's create a timer for this service. We need to do the following:

- Associate the timer with that service
- Assign a schedule to the timer

Add this to the end of your `configuration.nix`:

```nix
systemd.timers.acnh-island-check-reminder = {
  wantedBy = [ "timers.target" ];
  partOf = [ "acnh-island-check-reminder.service" ];
  timerConfig.OnCalendar = "TODO(Xe): this";
};
```

Before we mentioned that we want to trigger this reminder every morning at 9 am.
systemd timers specify their calendar config in the following format:

```
DayOfWeek Year-Month-Day Hour:Minute:Second
```

So for something that triggers every day at 9 AM, it would look like this:

```
*-*-* 8:00:00
```

[You can ignore the day of the week if it's not
relevant!](conversation://Mara/hacker)

So our final timer definition would look like this:

```nix
systemd.timers.acnh-island-check-reminder = {
  wantedBy = [ "timers.target" ];
  partOf = [ "acnh-island-check-reminder.service" ];
  timerConfig.OnCalendar = "*-*-* 8:00:00";
};
```

## Deployment and Testing

Now we can deploy this with `nixos-rebuild`:

```console
$ sudo nixos-rebuild switch
```

You should see a line that says something like this in the `nixos-rebuild`
output:

```
starting the following units: acnh-island-check-reminder.timer
```

Let's test the service out using `systemctl`:

```console
$ sudo systemctl start acnh-island-check-reminder.service
```

And you should then see a message on Discord. If you don't see a message, check
the logs using `journalctl`:

```console
$ journalctl -u acnh-island-check-reminder.service
```

If you see an error that looks like this:

```
curl: (26) Failed to open/read local data from file/application
```

This usually means that you tried to do a role or user mention at the beginning
of the message and curl tried to interpret that as a file input. Add a word like
"hey" at the beginning of the line to disable this behavior. See
[here](https://stackoverflow.com/questions/6408904/send-request-to-curl-with-post-data-sourced-from-a-file)
for more information.

---

Also happy December! My site has the [snow
CSS](https://christine.website/blog/let-it-snow-2018-12-17) loaded for the
month. Enjoy!
