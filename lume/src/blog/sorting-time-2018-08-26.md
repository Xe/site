---
title: Sorting Time
date: 2018-08-26
tags:
 - time
 - javascript
 - tale-of-woe
---

Computers have a very interesting relationship with time. Time is how we keep track of many things, but mainly we use time to keep track of how far along in a day cycle we are. These daily sunrise/sunset cycles take about 24 hours on average, and the periodicity of them runs just about everything. Computers use time to keep track of just about everything on the board, usually measured in tiny fractions of seconds. (The common rating of gigahertz for computer processors actually measures how much time it takes for the processor to execute one instruction. A processor with a clock of 3.4 gigahertz means that the processor executes, best case, 3.4 billion instructions per second.)  Computer programmers have several popular methods of storing time with computers, the number of time intervals since a fixed date (usually the number of seconds since January 1st 1970) or as a human-readable string. These intervals are normally ever only added to and read from, almost never updated by human hands after being initially set by the network time service.

Pulling things back into the real world, let's consider storing time in Javascript. Let's say we're using Javascript in the browser and have a date object like so:

```
var date = new DateTime();
```

Say this is for Thursday, August 23rd 2018 at midnight UTC. If we turn it into a string using the toString method:

```
date.toString(); // -> "Thu Aug 23 2018 00:00:00 GMT+0000 (UTC)"
```

We get the date and time as a string.  The application in question uses a data store that has an interesting problem: it will automatically coerce things to a string type without alerting developers.

```
typeof date === 'object' // -> true
```

We expect date to be a normal object after we add it to the store. Let's add it to the store and see what happens to it.

```
const record = store.createRecord("widget", { createdAt: date });
typeof record.get("createdAt"); // -> string
```

Oh boy. It's suddenly a string now. That's not good. 

```
console.log(record.get("createdAt")); // -> "Thu Aug 23 2018 00:00:00 GMT+0000 (UTC)"
```

This works all fine and well, but sometimes a few lists of things can get bizarrely out of order in the UI. Things created or updated right at a midnight UTC barrier would sometimes cause lists of things to show the newest elements at the bottom of the list. This confused us, sorting data really fitting it into the order it belongs to, and time doesn't usually advance out of order; so something being sorted wrongly by time is very intuitively confusing.

Consider a function like this at the given date above:

```
function minutesAgo(minutes) {
  return moment().subtract(minutes, "minute").toDate();
}

const date1 = minutesAgo(0);
const date2 = minutesAgo(1);
const date3 = minutesAgo(30);
```

If we were to sort date1, date2 and date3 with the current time being Thursday August 23 2018 at midnight UTC, it would make sense for the objects to sort ascendingly in the following order: date3, date2, date1. Not as strings however. As strings:

```
date1.toString(); // -> "Thu Aug 23 2018 00:00:00 GMT+0000 (UTC)"
date2.toString(); // -> "Wed Aug 22 2018 23:59:00 GMT+0000 (UTC)"
date3.toString(); // -> "Wed Aug 22 2018 23:30:00 GMT+0000 (UTC)"
```

Since T comes before W in the alphabet, the actual sort order is: date1, date3, date2. This causes an assertion failure in both humans and machines. This caused test failures, but only at about midnight UTC on Mondays, Thursdays, Fridays and Saturdays at 00:00 UTC through 00:30 UTC. How did we fix this? Turns out the time data from the API we get this information from is already properly sortable; this is because the API uses IS08601 timestamps.

```
const thursday = '2018-08-23T00:00:00.000Z';
const wednesday = '2018-08-22T23:30:00.000Z';

thursday > wednesday // true
```

This time data is also easy to convert back into a native DateTime object should we need it. The fix was to only ever store time as strings unless you need to actively do something with them, then you coerce them back into a native DateTime like it never happened. This is not an ideal fix, but given the larger complexity of the problem, it's what we're gonna have to live with for the time being. This solution at the very least seems to be less bad than the original problem, as things get sorted properly in the UI now. Yay computers! 

---

This is an adaptation of a pull request made by a coworker to work around an annoying to track down bug that caused flaky tests. It's not my story, but it just goes to show how many moving parts truly are at play with computers. Even when you think you have all of the moving parts kept track of, complicated systems interface in unpredictable ways. Increasingly complicated systems interface in increasingly unpredictable ways too, which makes finding problems like these more of a hunt.

Happy hunting and be well to you all.
