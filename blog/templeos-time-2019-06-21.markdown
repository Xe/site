---
title: "TempleOS: Date/Time Code"
date: 2019-06-22
---

# TempleOS: Date/Time Code

## Introduction

In the [last post](https://christine.website/blog/templeos-2-god-the-rng-2019-05-30)
we discussed `god`, the random number generator of TempleOS. This post is going to 
cover the time handling subsystem of the TempleOS kernel. As mentioned [before](https://christine.website/blog/sorting-time-2018-08-26), 
time is an unfortunately complicated thing, and legacy systems stapled onto 
legacy systems complicate this a lot. It also doesn't help that we consider 
one day 24 hours, even though it's more accurately slightly less than that. 
Most people consider a year to be 365 days, but the actual movement of the ball
Earth around the Sun is closer to 365 days and about [six hours](https://www.scienceabc.com/eyeopeners/why-do-we-have-a-leap-day-leap-year-after-every-4-years.html). 
To compensate for this, every 4 years we add a day to the calendar called a 
Leap Day. This complicates time calculations. And this isn't even going to 
cover time zones because I want you readers to have some shred of sanity left. 

## Problem space description

Time is an unfortunately complicated thing to handle. The goal of this part of 
the TempleOS kernel is to allow users to programmatically get the current "real 
time", or the time of the day. This is sometimes called "wall clock time", 
referring to the offices where Unix was developed having literal clocks on the 
wall to tell you the time. When you read a clock off the wall, you get the 
current time in hours, minutes and sometimes seconds. See this example below:

![an analog clock](/static/img/analog_clock.jpg)

This clock either says it is 2:51:59 or 14:51:59, as clocks only show a 12 hour
view of 24 hour days for historical reasons.

## How x86 hardware handles time

The x86 real-time clock keeps track of time with a similar kind of API. The OS
can query the real-time clock and get the following fields individually:

* Seconds
* Minutes
* Hours
* Weekday
* Day of month
* Month
* Last two digits of the year
* Sometimes the current century (19 or 20)

Note the last two digits of the year. This is returned because computers 
historically did not keep track of the current century in order to save memory. 
This famously caused the Y2K bug that [almost caused chaos globally](https://en.wikipedia.org/wiki/Year_2000_problem). Thankfully this problem was averted.

## How Linux handles time

Linux handles time by keeping track of [the number of seconds since January 1, 1970](https://en.wikipedia.org/wiki/Unix_time). 
Traditionally programs use a signed 32 bit integer value for this. You might 
notice that this is a drastically different time format than what the real-time
clock provides. When Linux reads from the real time clock, it does some math to
convert these real time values to a unix timestamp.

### The Rockchip calendar hack

Of course, this assumes that the hardware is behaving correctly. There have been
cases in the past where [buggy hardware required patches](https://lore.kernel.org/patchwork/patch/628142/). 
If the hardware calendar and the software understanding of that calendar get out 
of sync, things get really confused.

## How TempleOS handles time

TempleOS has two main time types, [CDateStruct](https://github.com/Xe/TempleOS/blob/master/Kernel/KernelA.HH#L193-L198) 
and [CDate](https://github.com/Xe/TempleOS/blob/master/Kernel/KernelA.HH#L186-L190). 
Let's look at each of these in more detail:

```c++
public class CDateStruct
{
  U8	sec10000,sec100,sec,min,hour,
	day_of_week,day_of_mon,mon;
  I32	year;
};

public I64 class CDate
{
  U32	time;
  I32	date;
};
```

The `I64 class CDate` part of the definition makes the HolyC compiler align
instances of it across the size of a int64 value, or 8 bytes. This also makes
the compiler treat instances of it like an integer value.

Here are approximate Zig equivalents:

```zig
pub const CDateStruct = packed struct {
  sec10000:    u8,
  sec100:      u8,
  sec:         u8,
  min:         u8,
  hour:        u8,
  day_of_week: u8,
  day_of_mon:  u8,
  mon:         u8,
  year:        i32,
};

pub const CDate = packed struct {
  time: u32,
  date: i32,
}
```

The [TempleOS documentation](https://templeos.holyc.xyz/Wb/Doc/TimeDate.html) 
describes `CDate` like this:

```
TempleOS uses a 64-bit value, CDate, for date/time.  The upper 32-bits are the 
days since Christ.  The lower 32-bits store time of day divided by 4 billion 
which works out to 49710ths of a second.  You can subtract two CDate's to get a 
time span.

Use CDATE_FREQ to convert.
```

And it [describes the `date` field](https://templeos.holyc.xyz/Wb/Doc/Date.html) 
of `CDate` like this:

```
Dates are 32-bit signed ints representing the number of days since the birth of 
Christ.  Negative values represent B.C.E. dates.
```

B.C.E. usually means "Before Common Era", which is a way to say "the time before
Jesus Christ was born". Christianity ended up being such a core part of global 
culture that the calendar was reset because of it. 



* CMOS real time clock API
    * NowDateTimeStruct();
        * How to read from the CMOS clock
        * The bug in the time code
* Live Demo (Zig/wasm+JS)

