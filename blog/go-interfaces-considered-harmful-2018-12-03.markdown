---
title: Go Interfaces Considered Harmful
date: 2018-12-03
---

# Go Interfaces Considered Harmful

A group of blind men heard that a strange animal had been brought to the town function, but none of them were aware of its type.

```
package blindmen

type Animal interface{}

func Town(strangeAnimal Animal) {
```

Out of curiosity, they said: “We must inspect and know it by type switches and touch, of which we are capable”.

```
type Toucher interface {
  Touch() interface{}
}
```

So, they sought it out, and when they found it they groped about it.

```
for man := range make([]struct{}, 6) {
   go grope(man, strangeAnimal.(Toucher).Touch())
}
```

In the case of the first person, whose hand landed on the trunk, said “This being is like a thick snake”.

```
type Snaker interface {
  Snake()
}

func grope(id int, thing interface{}) {
  switch thing.(type) {
  case Snaker:
    log.Printf("man %d: this thing is like a thick snake", id)
```

For another one whose hand reached its ear, it seemed like a kind of fan.

```
type Fanner interface {
  Fan()
}

// in grope switch block
case Fanner:
  log.Printf("man %d: this thing is like a kind of fan", id)
```

As for another person, whose hand was upon its leg, said, the it is a pillar like a tree-trunk.

```
type TreeTrunker interface {
  TreeTrunk()
}

// in grope switch block
case TreeTrunker:
  log.Printf("man %d: this thing is like a tree trunk", id)
```

The blind man who placed his hand upon its side said, “it is a wall”.

```
type Waller interface {
  Wall()
}

// in grope switch block
case Waller:
  log.Printf("man %d: this thing is like a wall", id)
```

Another who felt its tail, described it as a rope.

```
type Roper interface {
  Rope()
}

// in grope switch block
case Roper:
  log.Printf("man %d: this thing is like a rope", id)
```

The last felt its tusk, stating the thing is that which is hard, smooth and like a spear.

```
type Tusker interface {
  Tusk()
}

// in grope switch block
case Tusker:
  log.Printf("man %d: this thing is hard, smooth and like a spear", id)
```

All of the men spoke fact about the thing, but none of them spoke the truth of what it was.

```
// after grope switch block
log.Printf("%T", thing) // prints Elephant
```

---

```
  switch thing.(type) {
  case Trunker:
    log.Printf("man %d: this thing is like a thick snake", id)
  case Fanner:
    log.Printf("man %d: this thing is like a kind of fan", id)
  case TreeTrunker:
    log.Printf("man %d: this thing is like a tree trunk", id)
  case Waller:
    log.Printf("man %d: this thing is like a wall", id)
  case Roper:
    log.Printf("man %d: this thing is like a rope", id)
  case Tusker:
    log.Printf("man %d: this thing is hard, smooth and like a spear", id)
  }
```
