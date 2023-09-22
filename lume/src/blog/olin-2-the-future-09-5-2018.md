---
title: "Olin: 2: The Future"
date: 2018-09-05
series: olin
---

This post is a continuation of [this post](https://xeiaso.net/blog/olin-1-why-09-1-2018).

Suppose you are given the chance to throw out the world and start from scratch
in a minimal environment. You can then work up from nothing and build the world
from there.

How would you do this?

One of the most common ways is to pick a model that they are Stockholmed into 
after years of badness and then replicate it, with all of the flaws of the model
along with it. Dagger is a direct example of this. I had been stockholmed into
thinking that everything was a file stream and replicated Dagger's design based
on it. There was a really [brilliant](https://write.as/excerpts/conversation-with-_wmd-on-hacker-news)
Hacker News comment that inspired a bit of a rabbit hole internally, and I think
we have settled on an idea for a primitive that would be easy to implement and
use from multiple languages.

So, let's stop and ask ourselves a question that is going to sound really simple
or basic, but really will define a lot of what we do here.

What do we want to do with a computer that could be exposed to a WebAssembly 
module? What are the basic operations that we can expose that would be primitive
enough to be universally useful but also simple to understand from an implementation
standpoint from multiple languages?

Well, what are the programs actually doing with the interfaces? How can we use
that normal semantic behavior and provide a more useful primitive?

## The Parable of the Poison Arrow

When designing things such as these, it is very easy to get lost in the 
philosophical weeds. I mean, we are getting the chance to redefine the basic 
things that we will get angry at. There's a lot of pain and passion that goes
into our work and it shows.

As such, consider the following Buddhist parable:

> It's just as if a man were wounded with an arrow thickly smeared with poison. 
> 
> His friends & companions, kinsmen & relatives would provide him with a surgeon, and the man would say, 'I won't have this arrow removed until I know whether the man who wounded me was a noble warrior, a priest, a merchant, or a worker.' 
> 
> He would say, 'I won't have this arrow removed until I know whether the shaft with which I was wounded was that of a common arrow, a curved arrow, a barbed, a calf-toothed, or an oleander arrow.'
> 
> The man would die and those things would still remain unknown to him.

[Source](https://en.wikipedia.org/wiki/Parable_of_the_Poisoned_Arrow)

At some point, we are going to have to just try something and see what it is 
like. Let's not get lost too deep into what the bowstring of the person who shot
us with the poison arrow is made out of and focus more on the task at hand right
now, designing the ground floor.

## Core Operations

Let's try a new primitive. Let's call this primitive the interface. An interface
is a collection of types and methods that allows a WebAssembly module to perform
some action that it otherwise would be unable to do. As such, the only functions
we really need are a `require` function to introduce the dependency into the
environment, a `close` function to remove dependencies from the environment, and
an `invoke` function to call methods of the dependent interfaces. These can be
expressed in the following C-style types:

```c
// require loads the dependency by package into the environment. The int64 value
// returned by this function is effectively random and should be treated as
// opaque.
//
// If this returns less than zero, the value times negative 1 is the error code.
//
// Anything created by this function is to be considered initialized but
// unconfigured.
extern int64 require(const char* package);

// close removes a given dependency from the environment. If this returns less
// than zero, the value times negative 1 is the error code.
extern int64 close(int64 handle);

// invoke calls the given method with an input and output structure. This allows
// the protocol buffer generators to more easily build the world for us.
// 
// The resulting int64 value is zero if everything succeeded, otherwise it is the
// error code (if any) times negative 1.
//
// The in and out pointers must be to a C-like representation of the protocol
// buffer definition of the interface method argument. If this ends up being an
// issue, I guess there's gonna be some kinda hacky reader thing involved. No
// biggie though, that can be codegenned.
extern int64 invoke(int64 handle, int64 method, void* in, void* out);
```

(Yes, I know I made a lot of fuss about not just blindly following the design
decisions of the past and then just suggested returning a negative value from a
function to indicate the presence of an error. I just don't know of a better and
more portable mechanism for errors yet. If you have one, please suggest it to me.)

You may have noticed that the `invoke` function takes void pointers. This is
intentional. This will require additional code generation on the server side to
support copying the values out of WebAssembly memory. This may serve to be 
completely problematic, but I bet we can at least get Rust working with this.

Using these basic primitives, we can actually model way more than you think would
be possible. Let's do a simple example.

## Example: Logging

Consider logging. It is usually implemented as a stream of logging messages containing 
unstructured text that usually only has meaning to the development team and the
regular expressions that trigger the pager. Knowing this, we can expose a logging
interface like this:

```proto
syntax = "proto3";

package us.xeserv.olin.dagger.logging.v1;
option go_package = "logging";

// Writer is a log message writer. This is append-only. All text in log messages
// may be read by scripts and humans.
service Writer {
  // method 0
  rpc Log(LogMessage) returns (Nil) {};
}

// When nothing remains, everything is equally possible.
// TODO(Xe): standardize this somehow.
message Nil {}

// LogMessage is an individual log message. This will get added to as it gets
// propagated up through the layers of the program and out into the world, but 
// those don't matter right now.
message LogMessage {
  bytes message = 1;
}
```

And at a low level, this would be used like this:

```c
extern int64 require(const char* package);
extern int64 close(int64 handle);
extern int64 invoke(int64 handle, int64 method, void* in, void* out);

// This exposes logging_LogMessage, logging_Nil, 
// int64 logging_Log(int64 handle, void* in, void* out)
// assume this is magically generated from the protobuf file above.
#include <services/us.xeserv.olin.dagger.logging.v1.h> 

int64 main() {
  int64 logHdl = require("us.xeserv.olin.dagger.logging.v1");
  logging_LogMessage msg;
  logging_Nil none;
  msg.message = "Hello, world!";
  
  // The following two calls are equivalent:
  assert(logging_Log(logHdl, &msg, &none));
  assert(invoke(logHdl, logging_Writer_method_Log, &msg, &none));
  
  assert(close(logHdl));
}
```

This is really great to codegen, audit, validate, and not to mention we can easily
verify what logging interface the user actually wants from which vendor. This 
allows people who install Olin to their own cluster to potentially define their
own custom interfaces. This actually gives us the chance to make this a primitive.

Some problems that probably are going to come up pretty quickly is that every
language under the sun has their own idea of how to arrange memory. This may make
directly scraping the values out of ram inviable in the future. 

If reading values out of memory does become inviable, I suggest the following
changes:

```c
extern int64 require(const char* package);
extern int64 close(int64 handle);
extern int64 invoke(int64 handle, int64 method, char* in, int32 inlen, char* out int32 outlen);
```

(I don't know how to describe "pointer to bytes" in C, so I am using a C string 
here to fill in that gap.)
In this case, the arguments to `invoke()` would be pointers to protocol 
buffer-encoded ram. This may prove to be a huge burden in terms of deserializing
and serializing the protocol buffers over and over every time a syscall has to
be made, but it may actually be enough of a performance penalty that it prevents
spurious syscalls, given the "cost" of them. Code generators should remove most
of the pain when it comes to actually using this interface though, the 
automatically generated code should automatically coax things into protocol
buffers without user interaction.

For fun, let's take this basic model and then map Dagger's concept of file I/O to
it:

```proto
syntax = "proto3";

package us.xeserv.olin.dagger.files.v1;
option go_package = "files";

// When nothing remains, everything is equally possible.
// TODO(Xe): standardize this somehow.
message Nil {}

service Files {
  rpc Open(OpenRequest) returns (FID) {};
  rpc Read(ReadRequest) returns (ReadResponse) {};
  rpc Write(WriteRequest) returns (N) {};
  rpc Close(FID) returns (Nil) {};
  rpc Sync(FID) returns (Nil) {};
}

message FID {
  int64 opaque_id;
}

message OpenRequest {
  string identifier = 1;
  int64 flags = 2;
}

message N {
  int64 count
}

message ReadRequest {
  FID fid = 1;
  int64 max_length = 2;
}

message ReadResponse {
  bytes data = 1;
  N n = 2;
}

message WriteRequest {
  FID fid = 1;
  bytes data = 2;
}
```

Using these methods, we can rebuild (most of) the original API:

```c
extern int64 require(const char* package);
extern int64 close(int64 handle);
extern int64 invoke(int64 handle, int64 method, void* in, void* out);

#include <services/us.xeserv.olin.dagger.files.v1.h>

int64 filesystem_service_id;

void setup_filesystem() {
  filesystem_service_id = require("us.xeserv.olin.dagger.files")
}

int64 open(char *furl, int64 flags) {
  files_OpenRequest req;
  files_FID resp;
  int64 err;
  
  req.identifier = char*(furl);
  req.flags = flags;
  
  // could also be err = file_Files_Open(filesystem_service_id, &req, &resp);
  err = invoke(filesystem_service_id, files_Files_method_Open, &req, &resp);
  if (err != 0) {
    return err;
  }
  
  return resp.opaque_id;
}

int64 d_close(int64 fd) {
  files_FID req;
  files_Nil resp;
  int64 err;
  
  req.opaque_id = fd;
  
  err = invoke(filesystem_service_id, files_Files_method_Close, &req, &resp);
  if (err != 0) {
    return err;
  }
  
  return 0;
}

int64 read(int64 fd, void* buf, int64 nbyte) {
  files_FID fid;
  files_ReadRequest req;
  files_ReadResponse resp;
  int64 err;
  int i;
  
  fid.opaque_id = fd;
  req.fid = fid;
  req.max_length = nbyte;
  
  err = invoke(filesystem_service_id, file_Files_method_Read, &req, &resp);
  if (err != 0) {
    return err;
  }
  
  // TODO(Xe): replace with memcpy once we have libc or something
  for (i = 0; i < resp.n.count; i++) {
    buf[i] = resp.data[i]
  }
  
  return 0;
}

int64 write(int64 fd, void* buf, int64 nbyte) {
  files_FID fid;
  files_WriteRequest req;
  files_N resp;
  int64 err;
  
  fid.opaque_id = fd;
  req.fid = fid;
  req.data = buf; // let's pretend this works, okay?
  
  err = invoke(filesystem_service_id, files_Files_method_Write, &req, &resp);
  if (err != 0) {
    return err;
  }
  
  return resp.count;
}

int64 sync(int64 fd) {
  files_FID req;
  files_Nil resp;
  int64 err;
  
  req.opaque_id = fd;
  
  err = invoke(filesystem_service_id, files_Files_method_Sync, &req, &resp);
  if (err != 0) {
    return err;
  }
  
  return 0;
}
```

And with that we should have the same interface as Dagger's, save the fact that
the name `close` is now shadowed by the global close function. On the server side
we could implement this like so:

```go
package files

import (
  "context"
  "errors"
  "math/rand"
  
  "github.com/Xe/olin/internal/abi/dagger"
)

func init() {
  rand.Seed(time.Now().UnixNano())
}

type FilesImpl struct {
  *dagger.Process
}

func (FilesImpl) getRandomNumber() int64 {
  return rand.Int63()
}

func daggerError(respValue int64, err error) error {
  if err == nil {
    err = errors.New("")
  }
  
  return dagger.Error{Errno: dagger.Errno(respValue * -1), Underlying: err}
}

func (fs *FilesImpl) Open(ctx context.Context, op *OpenRequest) (*FID, error) {
  fd := fs.Process.OpenFD(op.Identifier, uint32(op.Flags))
  if fd < 0 {
    return nil, daggerError(fd, nil)
  
  return &FID{OpaqueId: fd}, nil
}


func (fs *FilesImpl) Read(ctx context.Context, rr *ReadRequest) (*ReadResponse, error) {
  fd := rr.Fid.OpaqueId
  data := make([]byte, rr.MaxLength)
  
  n := fs.Process.ReadFD(fd, data)
  if n < 0 {
    return nil, daggerError(n, nil)
  }
  
  result := &ReadResponse{
    Data: data,
    N: N{
      Count: n
    },
  }
  
  return result, nil
}

func (fs *FilesImpl) Write(ctx context.Context, wr *WriteRequest) (*N, error) {
  fd := wr.Fid.OpaqueId
  
  n := fs.Process.WriteFD(fd, wr.Data)
  if n < 0 {
    return nil, daggerError(n, nil)
  }
  
  return &N{Count: n}, nil
}

func (fs *FilesImpl) Close(ctx context.Context, fid *Fid) (*Nil, error) {
  return &Nil{}, daggerError(fs.Process.CloseFD(fid.OpaqueId), nil)
}

func (fs *FilesImpl) Sync(ctx context.Context, fid *Fid) (*Nil, error) {
  return &Nil{}, daggerError(fs.Process.SyncFD(fid.OpaqueId), nil)
}
```

And then we have all of these arbitrary methods bound to WebAssembly modules,
where they are free to use them how they want. I think that initially there is 
going to be support for this interface from Go WebAssembly modules as we can 
make a lot more assumptions about how Go handles its memory management, making
it a lot easier for us to code generate reading Go structures/pointers/whatever
out of Go WebAssembly memory than we can code generate reading C structures
(recursively with pointers and C-style strings galore too).
The really cool part is that this is all powered by those three basic functions:
`require`, `invoke` and `close`. The rest is literally just stuff we can treat
as a black box for now and code generate.

As before, I would love any comments that people have on this article. Please
contact me somehow to let me know what you think. This design is probably wrong.
