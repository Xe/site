---
title: "How to delete a ZFS zvol when it claims to be busy but you're sure it isn't"
date: 2024-01-05
---

Sometimes when you delete ZFS datasets you get this error:

```
$ sudo zfs destroy -rf arsene/vms
cannot destroy 'arsene/vms/oracle-linux-9': dataset is busy
cannot destroy 'arsene/vms/rocky-linux-9': dataset is busy
cannot destroy 'arsene/vms': dataset already exists
```

The "dataset is busy" error is thrown when deleting a zvol if
something has the file open. This can happen with zvols for linux
systems when your kernel's LVM stack is autodetecting the LVM groups
in the zvols. You can confirm this with the `lvdisplay` command:

```
$ lvdisplay

--- Logical volume ---
LV Path                /dev/vg_main/lv_swap
LV Name                lv_swap
VG Name                vg_main
LV UUID                0phfZx-OYpN-EHUe-11Gs-8fzl-iuwf-jgh9Nl
LV Write Access        read/write
LV Creation host, time localhost.localdomain, 2023-05-22 05:50:13 -0400
LV Status              NOT available
LV Size                4.00 GiB
Current LE             1024
Segments               1
Allocation             inherit
Read ahead sectors     auto
```

You can disable them with the `dmsetup` command:

```
$ sudo dmsetup remove -f /dev/vg_main/lv_swap
```

This will let you delete the offending zvols.
