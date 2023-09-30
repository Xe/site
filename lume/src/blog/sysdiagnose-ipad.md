---
title: How to run a sysdiagnose on an iPad
date: 2023-07-20
tags:
 - iPad
 - iPadOS
---

Sometimes you need to dump your system logs for a developer of an
application to understand why things are failing.
[sysdiagnose](https://it-training.apple.com/tutorials/support/sup075)
lets you have your iDevice emit a giant tarball of information so that
they can pick out what the problem is.

However, their official button pressing procedure is finicky and
doesn't give confirmation that anything is happening. Here's how you
unconditionally force your iPad to do a sysdiagnose:

* Remove your iPad from the keyboard dock
* Disable Stage Manager
* Open Settings
* Tap Accessibility
* Tap Touch
* Tap AssistiveTouch
* Enable it
* Tap Customize Top Level Menu
* Add another icon
* Tap the empty icon
* Choose Analytics
* Swipe up to the home screen
* Tap the AssistiveTouch button
* Tap Analytics
* Wait for it to finish
* Open Settings
* Tap Privacy & Security
* Tap Analytics & Improvements
* Tap Analytics Data
* Scroll all the way to the bottom
* Find something called "sysdiagnose" and tap on it
* Tap the share icon in the upper right hand corner of the screen
* Save to Files
* Save to your iCloud Desktop folder

Then you can give the developer the file they need to diagnose the
issue. Hope this helps!
