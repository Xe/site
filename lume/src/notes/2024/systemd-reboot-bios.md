---
title: How to reboot a Linux system into the BIOS
date: 2024-01-11
---

Usually to reboot your system into the BIOS, you have to press a
button or button combination on boot. This is usually one of the
following keys:

- F2
- F10
- Escape

Sometimes you just can't get into the BIOS, but you can boot into a
Linux system. In order to force the system to reboot into the BIOS,
run this command:

```
sudo systemctl reboot --firmware
```

This will force it to go into the BIOS so you can change settings or
update the BIOS.

On Windows, search for "Advanced Startup Options" in the start menu
and click through the "yes I really need to do something advanced"
boxes until you get something on the lines of "UEFI firmware
settings". That will get you into the BIOS.
