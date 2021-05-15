#!/usr/bin/env nix-shell
#! nix-shell -p python3 -i python3

import os
import time

pid = os.fork()
if pid == 0:
    for fd in {0, 1, 2}:
        os.close(fd)
    time.sleep(1)
else:
    print(pid)
