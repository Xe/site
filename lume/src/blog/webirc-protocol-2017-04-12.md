---
title: "IRCv3.2 `webirc` Extension"
date: 2017-04-12
---

This document does not describe a new IRCv3 standard. It is designed to 
document how the existing `WEBIRC` mechanism works so there is a specification 
to test things against. This is known to be implemented by all major IRC 
daemons as of the time of this writing.

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT", "SHOULD",
"SHOULD NOT", "RECOMMENDED",  "MAY", and "OPTIONAL" in this document are to be
interpreted as described in RFC 2119.

Summary
-------

The `WEBIRC` verb allows a connecting IRC client to spoof its origin IP address 
so that a user connecting via a gateway of some kind may have accountability 
for their actions and bans against them do not affect unintended users of said 
gateway.

This protocol verb must be sent before the initial `NICK` and `USER` handshake 
and may be advertised as the client capability `webirc`. The remote server may 
send a pre-connection `NOTICE` clarifying that the user has their specified IP 
address and reverse DNS. Gateway implementors must not let the user set their 
own IP address as part of connection negotiations.

Formatting
----------

The `WEBIRC` verb must be used as such:

```
WEBIRC <password> <client ident> <client reverse DNS> <client IP address>
```

Access to `WEBIRC` must be protected by a password to prevent abuse. If the 
password the client gives fails, the IRC daemon should disconnect the client 
with an appropriate error message. IRC daemon authors should also restruct the
use of the `WEBIRC` verb to a specific IP address and may force the use of
a specific identd reply.

Example Session
---------------

```
>> WEBIRC snowflower Mibbit anonyhash.mibbit.com 127.0.0.1
>> NICK mib_4002
>> USER Mibbit x x :http://mibbit.com AJAX IRC Client
<< :hostname.domain.tld 001 mib_4002 :Welcome to ShadowNET mib_4002!
```

Limitations
-----------

In order for this to be secure, the relay server must be trusted by the IRC
server. A remote server may kill off clients that fail the password and host
check, but this is not required.

---

This was recovered from an old backup of my site data on 2019-04-12.
