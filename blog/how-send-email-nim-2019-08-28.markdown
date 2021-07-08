---
title: How to Send Email with Nim
date: 2019-08-28
series: howto
tags:
 - nim
 - email
---

Nim offers an [smtp][nimsmtp] module, but it is a bit annoying to use out of the
box. This blogpost hopes to be a mini-tutorial on the basics of how to use the
smtp library and give developers best practices for handling outgoing email in
ways that Google or iCloud will accept.

## SMTP in a Nutshell

[SMTP][SMTPrfc], or the Simple Mail Transfer Protocol is the backbone of how 
email works. It's a very simple line-based protocol, and there are wrappers for
it in almost every programming language. Usage is pretty simple:

- The client connects to the server
- The client authenticates itself with the server
- The client signals that it would like to create an outgoing message to the server
- The client sends the raw contents of the message to the server
- The client ends the message
- The client disconnects

Unfortunately, the devil is truly in the details here. There are a few things
that _absolutely must_ be present in your emails in order for services like 
GMail to accept them. They are:

- The `From` header specifying where the message was sent from
- The Mime-Version that your code is using (if you aren't sure, put `1.0` here)
- The Content-Type that your code is sending to users (probably `text/plain`)

For a more complete example, let's create a `Mailer` type and a constructor:

```nim
# mailer.nim
import asyncdispatch, logging, smtp, strformat, strutils

type Mailer* = object
  address: string
  port: Port
  myAddress: string
  myName: string
  username: string
  password: string
  
proc newMailer*(address, port, myAddress, myName, username, password: string): Mailer =
  result = Mailer(
    address: address,
    port: port.parseInt.Port,
    myAddress: myAddress,
    myName: myName,
    username: username,
    password: password,
  )
```

And let's write a `mail` method to send out email:

```nim
proc mail(m: Mailer, to, toName, subject, body: string) {.async.} =
  let
    toList = @[fmt"{toName} <{to}>"]
    msg = createMessage(subject, body, toList, @[], [
      ("From", fmt"{m.myName} <{m.myAddress}"),
      ("MIME-Version", "1.0"),
      ("Content-Type", "text/plain"),
    ])

  var client = newAsyncSmtp(useSsl = true)
  await client.connect(m.address, m.port)
  await client.auth(m.username, m.password)
  await client.sendMail(m.myAddress, toList, $msg)
  info "sent email to: ", to, " about: ", subject
  await client.close()
```

Breaking this down, you can clearly see the parts of the SMTP connection as I
laid out before. The `Mailer` creates a new transient SMTP connection, 
authenticates with the remote server, sends the properly formatted email to
the server and then closes the connection cleanly. 

If you want to test this code, I suggest testing it with a freely available
email provider that offers TLS/SSL-encrypted SMTP support. This also means that
you need to compile this code with `--define: ssl`, so create `config.nims` and
add the following:

```nimscript
--define: ssl
```

Here's a little wrapper using [cligen][cligen]:

```nim
when isMailModule:
  import cligen, os
  
  let
    smtpAddress = getEnv("SMTP_ADDRESS")
    smtpPort = getEnv("SMTP_PORT")
    smtpMyAddress = getEnv("SMTP_MY_ADDRESS")
    smtpMyName = getEnv("SMTP_MY_NAME")
    smtpUsername = getEnv("SMTP_USERNAME")
    smtpPassword = getEnv("SMTP_PASSWORD")
  
  proc sendAnEmail(to, toName, subject, body: string) =
    let m = newMailer(smtpAddress, smtpPort, smtpMyAddress, smtpMyName, smtpUsername, smtpPassword)
    waitFor m.mail(to, toName, subject, body)
  
  dispatch(sendAnEmail)
```

Usage is simple:

```console
$ nim c -r mailer.nim --help
Usage:
  sendAnEmail [required&optional-params]
Options(opt-arg sep :|=|spc):
  -h, --help                         print this cligen-erated help
  --help-syntax                      advanced: prepend,plurals,..
  -t=, --to=       string  REQUIRED  set to
  --toName=        string  REQUIRED  set toName
  -s=, --subject=  string  REQUIRED  set subject
  -b=, --body=     string  REQUIRED  set body
```

I hope this helps, this module is going to be used in my future post on how to
create an application using Nim's [Jester][jester] framework.

[nimsmtp]: https://nim-lang.org/docs/smtp.html
[SMTPrfc]: https://tools.ietf.org/html/rfc5321
[jester]: https://github.com/dom96/jester
[cligen]: https://github.com/c-blake/cligen
