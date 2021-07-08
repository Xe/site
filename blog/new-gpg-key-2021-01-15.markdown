---
title: New PGP Key Fingerprint
date: 2021-01-15
---

EDIT(Xe): M02 11 2020 I fucked it up, this key is broken. Don't bother trying to
use it, it won't work. Use [age](https://github.com/FiloSottile/age) to send
messages with the key `ssh-ed25519
AAAAC3NzaC1lZDI1NTE5AAAAIPg9gYKVglnO2HQodSJt4z4mNrUSUiyJQ7b+J798bwD9` instead.

This morning I got an encrypted email, and in the process of trying to decrypt
it I discovered that I had _lost_ my PGP key. I have no idea how I lost it. As
such, I have created a new PGP key and replaced the one on my website with it.
I did the replacement in [this
commit](https://github.com/Xe/site/commit/66233bcd40155cf71e221edf08851db39dbd421c),
which you can see is verified with a subkey of my new key.

My new PGP key ID is `803C 935A E118 A224`. The key with the ID `799F 9134 8118
1111` should not be used anymore. Here are all the subkey fingerprints:

```
Signature key ....: 378E BFC6 3D79 B49D 8C36  448C 803C 935A E118 A224
      created ....: 2021-01-15 13:04:28
Encryption key....: 8C61 7F30 F331 D21B 5517  6478 8C5C 9BC7 0FC2 511E
      created ....: 2021-01-15 13:04:28
Authentication key: 7BF7 E531 ABA3 7F77 FD17  8F72 CE17 781B F55D E945
      created ....: 2021-01-15 13:06:20
General key info..: pub  rsa2048/803C935AE118A224 2021-01-15 Christine Dodrill (Yubikey) <me@christine.website>
sec>  rsa2048/803C935AE118A224  created: 2021-01-15  expires: 2031-01-13
                                card-no: 0006 03646872
ssb>  rsa2048/8C5C9BC70FC2511E  created: 2021-01-15  expires: 2031-01-13
                                card-no: 0006 03646872
ssb>  rsa2048/CE17781BF55DE945  created: 2021-01-15  expires: 2031-01-13
                                card-no: 0006 03646872
```

I don't really know what the proper way is to go about revoking an old PGP key.
It probably doesn't help that I don't use PGP very often. I think this is the
first encrypted email I've gotten in a year.

Let's hope that I don't lose this key as easily!
