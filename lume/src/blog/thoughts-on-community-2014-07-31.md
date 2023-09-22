---
title: Thoughts on Community Management
date: 2014-07-31
---

Many open source community projects lack proper management. They can put too
much of their resources in too few places. When that one person falls out of
contact or goes rogue on everyone, it can have huge effects on everyone
involved in the project. Users, Contributors and Admins.

Here, I propose an alternative management structure based on what works.

## Organization

Contributors and Project Administrators are there to take input/feedback from
Users, rectify the situation or explain why doing so is counterproductive.
Doing so will be done kindly and will be ran through at least another person
before it is posted publicly. This includes (but is not limited to) email, IRC,
forums, anything. A person involved in the project is a representative of it.
They are the face of it. If they are rude it taints the image of everyone
involved.

## Access

Project Administrators will have full, unfiltered access to anything the
project has. This includes root access, billing access, everything. There will
be no reason to hide things. Operational conversations will be shared. All
group decisions will be voted on with a simple Yes/No/Abstain process. As such
this team should be kept small.

## Contributions

Contributors will have to make pull requests, as will Administrators. There
will be review on all changes made. No commits will be pushed to master by
themselves unless there is approval. This will allow for the proper review and
testing procedures to be done to all code contributed.

Additionally, for ease of scripts scraping the commits when something is
released, a commit style should be enforced.

### Commit Style

The following section is borrowed from [Deis' commit
guidelines](https://github.com/deis/deis/blob/master/CONTRIBUTING.md).

---

We follow a rough convention for commit messages borrowed from CoreOS, who borrowed theirs
from AngularJS. This is an example of a commit:

```
feat(scripts/test-cluster): add a cluster test command

this uses tmux to setup a test cluster that you can easily kill and
start for debugging.
```

To make it more formal, it looks something like this:

```
{type}({scope}): {subject}
<BLANK LINE>
{body}
<BLANK LINE>
{footer}
```

The `{scope}` can be anything specifying place of the commit change.

The `{subject}` needs to use imperative, present tense: “change”, not “changed”
nor “changes”. The first letter should not be capitalized, and there is no dot
(.) at the end.

Just like the `{subject}`, the message `{body}` needs to be in the present tense,
and includes the motivation for the change, as well as a contrast with the
previous behavior. The first letter in a paragraph must be capitalized.

All breaking changes need to be mentioned in the `{footer}` with the description
of the change, the justification behind the change and any migration notes
required.

Any line of the commit message cannot be longer than 72 characters, with the
subject line limited to 50 characters. This allows the message to be easier to
read on github as well as in various git tools.

The allowed `{types}` are as follows:

    feat -> feature
    fix -> bug fix
    docs -> documentation
    style -> formatting
    ref -> refactoring code
    test -> adding missing tests
    chore -> maintenance

---

I believe that these guidelines would lead towards a harmonious community.
