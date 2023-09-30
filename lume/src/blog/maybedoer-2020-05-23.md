---
title: "maybedoer: the Maybe Monoid for Go"
date: 2020-05-23
tags:
 - go
 - golang
 - monoid
---

I recently posted (a variant of) this image of some Go source code to Twitter
and it spawned some interesting conversations about what it does, how it works
and why it needs to exist in the first place:

![the source code of package maybedoer](/static/blog/maybedoer.png)

This file is used to sequence functions that could fail together, allowing you
to avoid doing an `if err != nil` check on every single fallible function call.
There are two major usage patterns for it.

The first one is the imperative pattern, where you call it like this:

```go
md := new(maybedoer.Impl)

var data []byte

md.Maybe(func(context.Context) error {
  var err error
 
  data, err = ioutil.ReadFile("/proc/cpuinfo")
 
  return err
})

// add a few more maybe calls?

if err := md.Error(); err != nil {
  ln.Error(ctx, err, ln.Fmt("cannot munge data in /proc/cpuinfo"))
}
```

The second one is the iterative pattern, where you call it like this:

```go
func gitPush(repoPath, branch, to string) maybedoer.Doer {
  return func(ctx context.Context) error {
    // the repoPath, branch and to variables are available here
    return nil
  }
}

func repush(ctx context.Context) error {
  repoPath, err := ioutil.TempDir("", "")
  if err != nil {
    return fmt.Errorf("error making checkout: %v", err)
  }

  md := maybedoer.Impl{
    Doers: []maybedoer.Doer{
      gitConfig, // assume this is implemented
      gitClone(repoPath, os.Getenv("HEROKU_APP_GIT_REPO")), // and this too
      gitPush(repoPath, "master", os.Getenv("HEROKU_GIT_REMOTE")),
    },
  }
  
  err = md.Do(ctx)
  if err != nil {
    return fmt.Errorf("error repushing Heroku app: %v", err)
  }
  
  return nil
}
```

Both of these ways allow you to sequence fallible actions without having to
write `if err != nil` after each of them, making this easily scale out to
arbitrary numbers of steps. The design of this is inspired by a package used at
a previous job where we used it to handle a lot of fiddly fallible actions that
need to happen one after the other.

However, this version differs because of the `Doers` element of
`maybedoer.Impl`. This allows you to specify an entire process of steps as long
as those steps don't return any values. This is very similar to how Haskell's
[`Data.Monoid.First`](https://hackage.haskell.org/package/base-4.14.0.0/docs/Data-Monoid.html#t:First)
type works, except in Go this is locked to the `error` type (due to the language
not letting you describe things as precisely as you would need to get an analog
to `Data.Monoid.First`). This is also similar to Rust's `and_then` combinator.

If we could return values from these functions, this would make `maybedoer`
closer to being a monad in the Haskell sense. However we can't so we are locked
to one specific instance of a monoid. I would love to use this for a pointer (or
pointer-like) reference to any particular bit of data, but `interface{}` doesn't
allow this because `interface{}` matches _literally everything_:

```go
var foo = []interface{
  1,
  3.4,
  "hi there",
  context.Background(),
  errors.New("this works too!"),
}
```

This could mean that if we changed the type of a Doer to be:

```go
type Doer func(context.Context) interface{}
```

Then it would be difficult to know how to handle returns from the function.
Arguably we could write some mechanism to check if it is an error:

```go
result := do(ctx)
if result != nil {
  switch result.(type) {
  case error:
    return result // result is of type error magically
  default:
    md.return = result
  }
}
```

But then it would be difficult to know how to pipe the result into the next
function, unless we change Doer's type to be:

```go
type Doer func(context.Context, interface{}) interface{}
```

Which would require code that looks like this:

```go
func getNumber(ctx context.Context, _ interface{}) interface{} {
  return 2
}

func double(ctx context.Context, num interface{}) interface{} {
  switch num.(type) {
  case int:
    return 2+2
  default:
    return fmt.Errorf("wanted num to be an int, got: %T", num)
  }
  
  return nil
}
```

But this kind of repetition would be required for _every function_. I don't
really know what the best way to solve this in a generic way would be, but I'm
fairly sure that these fundamental limitations in Go prevent this package from
being genericized to handle function outputs and inputs beyond what you can do
with currying (and maybe clever pointer usage).

I would love to be proven wrong though. If anyone can take this [source code
under the MIT license](/static/blog/maybedoer.go) and prove me wrong, I will
stand corrected and update this blogpost with the solution. 

This kind of thing is more easy to solve in Rust with its
[Result](https://doc.rust-lang.org/std/result/) type; and arguably this entire
problem solved in the Go package is irrelevant in Rust because this solution is
in the standard library of Rust.
