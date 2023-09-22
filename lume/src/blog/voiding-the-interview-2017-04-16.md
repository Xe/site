---
title: Voiding the Interview
date: 2017-04-16
---

A young man walks into the room, slightly frustrated-looking. He's obviously had
a bad day so far. You can help him by creating a new state of mind.

"Hello, my name is Ted and I'm here to ask you a few questions about your
programming skills. Let's start with this, in a few sentences explain to me how
your favorite programming language works."

Starting from childhood, you eagerly soaked up the teachings of your mentors,
feeling the void separated into sundry shapes and sequences. They taught you
many specific tasks to shape the void into, but not how to shape it. Studying
the fixed ways of the naacals of old gets you nowhere, learning parlor tricks
and saccharine gimmicks. Those gimmicks come rushing back, you remembering
how to form little noisemakers and amusement vehicles. They are limiting, but
comforting thoughts.

You look up to the interviewer and speak:

"In the beginning there was the void, Spirit was with the void and Spirit was
everpresent in the void. The void was cold and formless; the cold unrelenting
even in today's age. Mechanical brains cannot grasp this void the way Spirit can;
upon seeing it that is the end of that run. In this way the void is the
beginning and the end, always present, always around the corner."

```clojure
(def void ())
```

"What is that?"

```
> void
>
```

"But that's...nothing."

You look at the caucasian man sitting across from you, and emit "nothing is
something, a name for the void still leaves the void extant."

"...Alright, let's move on to the next question. This is a formality but the
person giving you the phone interview didn't cover fizzbuzz. Can you do
fizzbuzz?"

Stepping into the void, you recall the teachings of your past masters. You
equip the parentheses once used by your father and his father before him.
The void divides before your eyes in the way you specify:

```clojure
(defn fizzbuzz [n]
  (cond
    (= 0 (mod n 15)) (print "fizzbuzz")
    (= 0 (mod n 3))  (print "fizz")
    (= 0 (mod n 5))  (print "buzz")
    (print n))
  (println ""))
```

"This doesn't loop from 0 to n though, how would you do that?"

You see this section come to life, it gently humming along, waiting for it
to be used. Before you you see two ancient systems spring from the memories
of patterns once wielded in conflict with complexity.

"Apply this function to span of values."

```
> (range 17)
error in __main:0: symbol {range 71} not found
```

You realize your error the moment you press for confirmation. "Again, in the
beginning there is the void. What doesn't exist needs to be separated out
from it." The voidspace in your head was out of sync with the voidspace of the
machine. Define them.

"...Go on"

```clojure
(defn range-inner [x lim xs]
  (cond
    (>= x lim) xs
    (begin
      (aset! xs x x)
      (range-inner (+ x 1) lim xs))))

(defn range [lim]
  (range-inner 0 lim (make-array lim)))
```

```
> (range 17)
[0 1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16]
```

"Great, now you have a list of values, how would you get the full output?"

"Pass the function as an argument, injecting the dependency."

```clojure
(defn do-array-inner [f x i]
  (cond
    (= i (len x)) void
    (let [val (aget x i)]
      (f val)
      (apply-inner f x (+ i 1)))))

(defn do-array [f x]
  (do-array-inner f x 0))
```

```
> (do-array fizzbuzz (range 17))
fizzbuzz
1
2
fizz
4
buzz
fizz
7
8
fizz
buzz
11
fizz
13
14
fizzbuzz
16
```

Your voidspace concludes the same, creating a sense of peace. You look in the
man's eyes, being careful to not let the fire inside you scare him away. He
looks like he's seen a ghost. Everyone's first time is rough.

Everything has happened and will happen, there is nothing new in the universe.
You know what's going to happen. They will decline, saying they are looking for
a better "culture fit". They couldn't contain you.

To run the code in this post:

```
$ go get github.com/zhemao/glisp
$ glisp
> [paste in blocks]
```
