---
date: 2018-03-30
title: When Then Zen
---

# When Then Zen

Meditation is something that is very easy to experience but very difficult to
explain in any way that is understandable. Historically, things that man could
not explain on his own get attributed to gods. As such, religious texts that
describe meditation can be very difficult to understand without context in the
religion in question. 

I would like to change this and make meditation more accessible. As such, I 
have created the [`When Then Zen`](https://github.com/Xe/when-then-zen) 
project. This project aims to divorce meditation methods from the context of
their spirituality and distill them down into what the steps to the process 
are. 

## A better way to teach meditation

At a high level, meditation is the act of practicing the separation of action
and reaction. A lot of the meditation methods that people have been publishing
over the years are the equivalent of what works for them on their PC (tm), and
as such things are generally described using whatever comparators the author of
the meditation guide is comfortable with. This can lead to confusion.

The way I am teaching meditation is simple: teach the method and have people do
it and see what happens. I've decided to teach methods using [Gherkin](https://github.com/cucumber/cucumber/wiki/Gherkin).
Gherkin can be kind of strange to read if you are not used to it, so consider
the game of baseball, specifically the act of the batter hitting a home run. 

```
Feature: home run
  Scenario: home run
    As a batter
    In order to hit a home run
    Given the pitcher has thrown the ball
    When I swing
    Then I hit the ball out of the park
```

As shown above, a Gherkin scenario clearly identifies who the feature is 
affecting, what actions they take and what things should happen to them as a
result of them taking those actions. This translates very well when trying to
explain some of the finer points of meditation, EG:

```
  # from when then zen's metta feature
  Scenario: Nature Walking
    # this is optional
    # but it helps when you're starting
    # physical fitness
    As a meditator
    In order to help me connect with the environment
    Given a short route to walk on
    When I walk down the route
    Then I should relax and enjoy the scenery
    And feel the sensations of the world around me
```



- Overall summary of the project
  - goals
  - explanation of method of teaching
  - contextual/terminology dump
- Sneak peek of in-progress features
- Potential other topics I'd like to dig into
  - kundalini
  - yoga
  - tulpamancy
  - musical synchronization
