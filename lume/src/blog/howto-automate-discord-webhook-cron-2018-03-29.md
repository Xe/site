---
title: How to Automate Discord Message Posting With Webhooks and Cron
date: 2018-03-29
series: howto
---

Most Linux systems have [`cron`](https://en.wikipedia.org/wiki/Cron) installed to run programs at given intervals. An example usecase would be to install package updates every Monday at 9 am (keep the sysadmins awake!).

Discord lets us post things using [webhooks](https://discordapp.com/developers/docs/resources/webhook). Combining this with cron lets us create automated message posting bots at arbitrary intervals.

## The message posting script

Somewhere on disk, copy down the following script:

```sh
#!/bin/sh
# msgpost.sh
# change MESSAGE, WEBHOOK and USERNAME as makes sense
# This code is trivial, and not covered by any license or warranty.

# explode on errors
set -e

MESSAGE='haha memes are funny xD'
WEBHOOK=https://discordapp.com/api/webhooks/0892379892092/AFkljAoiuj098oKA_98kjlA85jds
USERNAME=KRONK

curl -X POST \
     -F "content=${MESSAGE}" \
     -F "username=${USERNAME}" \
     "${WEBHOOK}"
```

Test run it and get a message like this:

![example discord message](https://i.imgur.com/dtjXcei.png)

## How to automate it

To automate it, first open your [`crontab(5)`](https://man7.org/linux/man-pages/man5/crontab.5.html) file:

```console
$ crontab -e
```

Then add a crontab entry as such:

```crontab
# Post this funny message every hour, on the hour
0 * * * *  sh /path/to/msgpost.sh

# Also valid with some implementations of cron (non-standard)
@hourly    sh /path/to/msgpost.sh
```

Then save this with your editor and it will be loaded into the cron daemon. For more information on crontab formats, see [here](https://crontab.guru/).

To run multiple copies of this, create multiple copies of `msgpost.sh` on your drive with multiple crontab entries.

Have fun :)
