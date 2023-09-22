---
title: "Metaprogramming: Partial Application..."
date: 2015-08-26
---

The title of this post looks intimidating. There's a lot of words there that
look like they are very complicated and will take a long time to master. In
reality, they are really very simple things. Let's start with a mundane example
and work our way up to a real-world bit of code. Let's begin with a small
story:

---

ACMECorp has a world-renowned Python application named Itera that is known for
its superb handling of basic mathematic functions. It's so well known and
industry proven that it is used in every school and on every home computer. You
have just accepted a job there as an intermediate programmer set to do
maintenance on it. Naturally, you are very excited to peek under the hood of
this mysterious and powerful program and offer your input to make it even
better for the next release and its users.

Upon getting there, you settle in and look at your ticket queue for the day.
A user is complaining that whenever they add `3` and `5`, they get `7` instead
of `8`, which is what they expected. Your first step is to go look into the
`add3` function and see what it does:

```Python
def add1(x):
    return x + 1

def add2(x):
    return x + 2

def add3(x):
    return x + 2

def add4(x):
    return x + 4
```

You are aghast. Your company's multi-billion dollar calculator is brought to
its knees by a simple copy-paste error. You wonder, "how in Sam Hill are these
people making any money???" (The answer, of course, is that they are a big
enterprise corporation)

You let your boss know about the bad news, you are immediately given any
resource in the company that you need to get this mission-critical problem
solved for *any input*. Yesterday. Without breaking the API that the rest of
the program has hard-coded in.

---

Let's look at what is common about all these functions. The `add*` family of
functions seems to all be doing one thing consistently: adding one number to
another.

Let's define a function called `add` that adds any two numbers:

```Python
def add(x, y):
    return x + y
```

This is nice, but it won't work for the task we were given, which is to not
break the API.

Let's go over what a function is in Python. We can define a function as
something that takes some set of Python values and produces some set of Python
values:

```haskell
PythonFunction :: [PythonValue] -> [PythonValue]
```

We can read this as "a Python function takes a set of Python values and
produces a set of Python values". Now we need to define what a Python value
actually is. To keep things simple, we're only going to define the following
types of values:

- `None` -> no value
- `Int` -> any whole number (Python calls this `int`)
- `Text` -> any string value (Python calls this `str`)
- `Function` -> something that takes and produces values

Python [itself has a lot more types that any value can be](https://docs.Python.org/3.4/library/stdtypes.html),
but for the scope of this blog post, this will do just fine.

Now, since a function can return a value and a function is a value, let's see
what happens if you return a function:

```python
def outer():
    def inner():
        return "Hello!"
    return inner
```

And in the repl:

```
>>> type(outer)
<type 'function'>
```

So `outer` is a function as we expect. It takes `None` (in Python, a function
without arguments has `None` for the type of them) and returns a function that
takes `None` and that function returns `Text` containing `"Hello!"`.
Let's make sure of this:

```
>>> outer()()
'Hello!'
>>> type(outer()())
<type 'str'>
```

Yay! When nothing is applied to the result of applying nothing to `outer`, it
returns the `Text` value `"Hello!"`. We can define the type of `outer` as the
following:

```haskell
outer :: None -> None -> Text
```

Now, let's use this for addition:

```python
# add :: Int -> Int -> Int
def add(x):
    def inner(y):
        return x + y

    return inner
```

And in the repl:

```
>>> add(4)(5)
9
```

A cool feature about this is that now we can dip into something called Partial
Application. Partial application lets you apply part of the arguments of
a function and you get another function out of it. Let's trace the type of the
`inner` function inside the `add` function, as well as the final computation
for clarity:

```python
# add :: Int -> Int -> Int
def add(x):
    # inner :: Int -> Int
    def inner(y):
        return x + y # :: Int

    return inner
```

Starting from the inside, we can see how the core computation here is `x + y`,
which returns an `Int`. Then we can see that `y` is passed in and in the scope
also as an `Int`. Then we can also see that `x` is passed in the outermost
layer as an int, giving it the type `Int -> Int -> Int`. Since `inner` is
a value, and a Python variable can contain any Python value, let's make
a function called `increment` using the `add` function:

```python
# increment :: Int -> Int
increment = add(1)
```

And in the repl:

```
>>> increment(50)
51
```

`increment` takes the integer given and increases it by 1, it is the same thing
as defining:

```python
def increment50():
    return 51
```

Or even `51` directly.

Now, let's see how we can use this for the `add*` family of function mentioned
above:

```python
# add :: Int -> Int -> Int
def add(x):
    def inner(y):
        return x + y

    return inner

# add1 :: Int -> Int
add1 = add(1)

# add2 :: Int -> Int
add2 = add(2)

# add3 :: Int -> Int
add3 = add(3)

# add4 :: Int -> Int
add4 = add(4)
```

And all we need to do from here is a few simple tests to prove it will work:

```python
if __name__ == "__main__":
    assert add(1)(1) == 2 # 1 + 1
    assert add(1)(2) == add(2)(1) # 1+2 == 2+1
    print("all tests passed")
```

```console
$ python addn.py
all tests passed
```

Bam. The `add*` family of functions is now a set of partial applications. It is
just a set of half-filled out forms.

---

You easily mechanically rewrite all of the `add*` family of functions to use
the metaprogramming style you learned on your own. Your patch goes in for
consideration to the code review team. Meanwhile your teammates are frantically
going through every function in the 200,000 line file that defines the `add*`
family of functions. They are estimating months of fixing is needed not to
mention millions of lines of test code. They are also estimating an additional
budget of contractors being brought in to speed all this up. Your code has made
all of this unneeded.

Your single commit was one of the biggest in company history. Billboards that
were red are now beaming a bright green. Your code fixed 5,000 other copy-paste
errors that have existed in the product for years. You immediately get a raise
and live happily ever after, a master in your craft.

---

For fun, let's rewrite the `add` function in Haskell.

```haskell
add :: Int -> Int -> Int
add x y = x + y
```

And then we can create a partial application with only:

```haskell
add1 :: Int -> Int
add1 = (add 1)
```

And use it in the repl:

```
Prelude> add1 3
4
```

Experienced haskellers would probably gawk at this. Because functions are the
base data type in Haskell, and partial application means that you can make
functions out of functions, we can define `add` as literally the addition
operator `(+)`:

```haskell
add :: Int -> Int -> Int
add = (+)
```

And because operators are just functions, we can further simplify the `add1`
function by partially applying the addition operation:

```haskell
add1 :: Int -> Int
add1 = (+1)
```

And that will give us the same thing.

```
Prelude> let add1 = (+1)
Prelude> add1 3
4
```

---

Now, real world example time. I recently wrote a simple JSON api based off of
a lot of data that has been marginally useful to some people. This api has
a series of HTTP endpoints that return data about My Little Pony: Friendship is
Magic episodes. Its code is [here](https://github.com/Xe/ponyapi) and its
endpoint is `http://ponyapi.apps.xeserv.us`.

One of the challenges when implementing it was how to avoid a massive amount of
copy-pasted code when doing so. I had started with a bunch of functions like:

```python
# all_episodes :: IO [Episode]
def all_episodes():
    r = requests.get(API_ENDPOINT + "/all")

    if r.status_code != 200:
        raise Exception("Not found or server error")

    return r.json()["episodes"]
```

Which was great and all, but there was so much code duplication involved to
just get one result for all the endpoints. My first step was to write something
that just automated the getting of json from an endpoint in the same way
I automated addition above:

```python
# _base_get :: Text -> None -> IO (Either Episode [Episode])
def _base_get(endpoint):
    def doer():
        r = requests.get(API_ENDPOINT + endpoint)

        if r.status_code != 200:
            raise Exception("Not found or server error")

    try:
        return r.json()["episodes"]
    except:
        return r.json()["episode"]

# all_episodes :: IO [Episode]
all_episodes = _base_get("/all")
```

Where `_base_get` returned the function that satisfied the request.

This didn't end up working so well with the endpoints
[that take parameters](https://github.com/Xe/ponyapi#seasonnumber), so I had to
account for that in my code:

```python
# _base_get :: Text -> Maybe [Text] -> (Maybe [Text] -> IO (Either Episode [Episode]))
# _base_get takes a text, a splatted list of texts and returns a function such that
#     the function takes a splatted list of texts and returns either an Episode or
#     a list of Episode as an IO action.
def _base_get(endpoint, *fragments):
    def doer(*args):
        r = None

        assert len(fragments) == len(args)

        if len(fragments) == 0:
            r = requests.get(API_ENDPOINT + endpoint)
        else:
            url = API_ENDPOINT + endpoint

            for i in range(len(fragments)):
                url = url + "/" + fragments[i] + "/" + str(args[i])

            r = requests.get(url)

        if r.status_code != 200:
            raise Exception("Not found or server error")

        try:
            return r.json()["episodes"]
        except:
            return r.json()["episode"]

    return doer

# all_episodes :: IO [Episode]
all_episodes = _base_get("/all")

# newest :: IO Episode
newest = _base_get("/newest")

# last_aired :: IO Episode
last_aired = _base_get("/last_aired")

# random :: IO Episode
random = _base_get("/random")

# get_season :: Int -> IO [Episode]
get_season = _base_get("", "season")

# get_episode :: Int -> Int -> IO Episode
get_episode = _base_get("", "season", "episode")
```

And that was it, save the `/search` route, which was acceptable to implement by
hand:

```python
# search :: Text -> IO [Episode]
def search(query):
    params = {"q": query}
    r = requests.get(API_ENDPOINT + "/search", params=params)

    if r.status_code != 200:
        raise Exception("Not found or server error")

    return r.json()["episodes"]
```

---

Months later you have been promoted as high as you can go. You've been teaching
the other engineers at ACMECorp metaprogramming and even convinced management
to let the next big project be in Haskell.

You are set for life. You have won.

---

For comments on this article, please feel free to email me, poke me in `#geek`
on `irc.ponychat.net` (my nick is Xena), or leave thoughts at one of the below
places this article has been posted.

Comments:

- [Hacker News](https://news.ycombinator.com/item?id=10127389)
- [Reddit](https://www.reddit.com/r/programming/comments/3ijyyz/metaprogramming_partial_application_and_currying/)
