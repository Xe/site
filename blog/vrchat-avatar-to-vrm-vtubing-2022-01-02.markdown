---
title: Converting a VRChat Avatar to VRM Format for VTubing
date: 2022-01-02
series: vtuber
tags:
 - vrchat
 - vrm
 - tutorial
---

So you want to be a VTuber, but you already have a custom avatar that you really
like in VRChat. Surely there's a way to make it work, right? This tutorial is
going to cover the problem, how I do it and everything I learned along the way.

To say that VR is in its infancy is an understatement. This early on in the
game, there are many different ecosystems and each of them really came up with
their own avatar format out of necessity. As such, there's a lot of vastly
different and incompatible ecosystems for this kind of avatar data. However, a
common avatar data interchange format is starting to crop up organically:
[VRM](https://vrm.dev). This format mostly came from the efforts of Japanese
VTubers, and as such most 3d VTubing software supports it.

However, you still have your avatar stuck in VRChat format, thus this
article exists.

[Please note that this REQUIRES you to have access to your avatar data in the
Unity editor. Please support VRChat creators, making this kind of content is
hard and expensive.](conversation://Mara/hacker)
  
So at a high level, here are the steps for doing this:

* Add the VRM plugins to your VRChat avatar project
* Export a basic unfinished VRM
* Make a new Unity project with the VRM plugins
* Import the unfinished VRM into it
* Fix the shaders
* Fix the blendshapes
* Test in VSeeFace

## First Export

First, you need to get your avatar into a "dirty" VRM file. This is going to
become the basis of your VTuber model. In order to do this, you need to import
the VRM plugins for Unity into your VRChat avatar project.

There are two ways to do this, one of them is to import the Unity packages
directly and the other is to customize your Unity package manifest to include
the VRM plugins. I'm going to cover both of these below:

<div id="installvrm"></div>

With either of these methods, you will need the VRM plugins [from
GitHub](https://github.com/vrm-c/UniVRM).

### UnityPackage Install Method

The simplest way to get the VRM plugins installed is to use the UnityPackage
files that the VRM plugin team releases. Download the UniVRM and VRM packages
from the latest release:

And then import them into Unity by going to the "Assets" menu, clicking on the
"Import Package" menu and then click on "Custom Package" and select each of the
UnityPackage files you downloaded.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+153934.png)

### `manifest.json` Install Method

On the release page for UniVRM, there's a JSON fragment labeled `manifest.json`.
This corresponds to a file in your Unity project where you can tell Unity to
automatically download and manage packages and updates to them from GitHub and
other sources.

In your Unity project folder, open the "Packages" folder and open
`manifest.json` in your favorite text editor (such as emacs). There will be a
JSON object that looks like this:

```json
{
  "dependencies": {
    "com.bar.foo": "1.0.0"
  }
}
```

This will be full of other things by default, but you can mostly ignore it
because those are things Unity puts in there by itself. Go to the end of the
object and add a comma to the end of the last value in `dependencies`. Then
paste in the VRM `manifest.json` fragment in the release page. You should end up
with something like this:

```json
{
  "dependencies": {
    "com.bar.foo": "1.0.0",
    "com.vrmc.vrmshaders": "https://github.com/vrm-c/UniVRM.git?path=/Assets/VRMShaders#v0.92.0",
    "com.vrmc.gltf": "https://github.com/vrm-c/UniVRM.git?path=/Assets/UniGLTF#v0.92.0",
    "com.vrmc.univrm": "https://github.com/vrm-c/UniVRM.git?path=/Assets/VRM#v0.92.0",
    "com.vrmc.vrm": "https://github.com/vrm-c/UniVRM.git?path=/Assets/VRM10#v0.92.0"
  }
}
```

The exact release of the VRM plugin based on when you read this article, but
overall it should look something like that.

### Make the first VRM

When you install the Unity VRM plugins, you should have a new "VRM0" menu in
your Unity menu bar. Click on it and select "Export UniVRM-0.80.0".

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+112733.png)

A new window will open titled "VRM Exporter".

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+113005.png)

Change the language to "En" (English) so that you can read the error and warning
messages. Then fill out the title, version and author information and export
your avatar to a VRM file on your desktop folder.

[Strictly speaking, you don't NEED to use your desktop, however that's where I
put intermediate files for things like this.](conversation://Mara/hacker)

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+113040.png)

## Actually Making Your VRM Shine

Strictly speaking, I don't know if you _need_ to do this, but I've found that
things are the most reliable when you make a separate Unity project for
manipulating your VRM avatar. Open the Unity Hub and create a separate project,
then re-import the [VRM plugin](#installvrm).

Once you have that ready, import your raw VRM into your new Unity project by
clicking on the "VRM0" menu and selecting "Import from VRM 0.x". Select your raw
VRM from before and then save it as a prefab in your folder. Unity will lock up
for a moment while everything gets crunched and you will have a shiny new prefab
in your project.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+113833.png)

Drag that prefab into the scene and then click on it to open its inspector.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+114246.png)

Reset its position to 0,0,0 so that it's in the center of the world.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+114319.png)

### Shaders

Now you will need to go in and fix all the shaders for your avatar. They should
either be the MToon shader or the UniUnlit shader. MToon should be your default
choice as it gives you the anime cartoony feeling that you see in a lot of
Japanese VTubers on YouTube. This is going to require some trial and error in
that annoying "draw the rest of the owl" kind of way.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+114723.png)

### Blendshapes

Now you will have a basic avatar, but if you import it into things like VSeeFace
you will notice that the avatar doesn't move its mouth when you speak things.
This is because you need to fix the avatar's blendshapes. Blendshapes are the
cues that the VRM format uses to know how to manipulate your avatar model's
bones in order to make it _look_ like the avatar is speaking.

Unfortunately, the avatar format is limited to having blendshapes for vowel
sounds only. In practice, this is not a huge deal (keep in mind this format was
originally designed for the Japanese market) but it can mean that you need to
enunciate your vowels a bit more to be sure the software picks up what you are
saying and adjusts your avatar lips appropriately. There are ways to get better
tracking, such as using an iPhone's FaceID sensor for
[Perfect Sync](https://malaybaku.github.io/VMagicMirror/en/tips/perfect_sync),
but for now vowels only will be good enough.

When you import a VRM, it creates a blendshapes folder in your project folder.
Open that folder and you will see a bunch of files that are named after the
different VRM animations that you can customize. The big ones are the speech
parts `A`, `I`, `U`, `E` and `O`. These will form the basis of your VTuber vibe.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+120248.png)

Start with `A`. This correlates to the
[IPA](https://en.wikipedia.org/wiki/International_Phonetic_Alphabet) vowel /a/,
which is made by having a fairly open mouth and with your tongue resting on the
bottom of your mouth. Open the `A` blendshape and expand the "Body" dropdown.
The slider `vrc.v_aa` correlates to this sound, so set it to 100.

![](https://cdn.xeiaso.net/file/christine-static/blog/Screenshot+2022-01-02+120745.png)

The rest of the vowels should be customized in the same way, here is a small
table of what they are called in the blendshapes folder, the IPA for them and
the slider that I use:

| Name  | IPA                                                                 | Slider     |
| :---- | :---                                                                | :------    |
| A     | [`/a/`](https://en.wikipedia.org/wiki/Open_central_unrounded_vowel) | `vrc.v_aa` |
| E     | [`/e/`](https://en.wikipedia.org/wiki/Mid_front_unrounded_vowel)    | `vrc.v_e`  |
| I     | [`/i/`](https://en.wikipedia.org/wiki/Close_front_unrounded_vowel)  | `vrc.v_ih` |
| O     | [`/o/`](https://en.wikipedia.org/wiki/Mid_back_rounded_vowel)       | `vrc.v_oh` |
| U     | [`/u/`](https://en.wikipedia.org/wiki/Close_back_unrounded_vowel)   | `vrc.v_ou` |

These should get you an _approximation_ of the vowel sounds, which should be
good enough for VTubing starting out.

Another good thing to care about here is blinking. Having your avatar stare
endlessly into the eyes of your audience is a great way to make people very
uncomfortable and click away. This is not a good way to grow a fanbase around
your content.

There are three blendshapes for blinking. One for both eyes blinking, one for
the left eye blinking and one for the right eye blinking. Customize these the
same way you customized the other blendshapes, but look for things that make
your avatar's eyes move. Some models may not support having only one eye
blinking. If this is the case, just make all the blendshapes have both eyes
blink.

You can also customize the moods for expressions like being angry, but I'm going
to skip this for now. Play around with it and you should be able to get
something that looks decent.

## Testing

Now we should have a semi-decent VRM file. Export it by going to the "VRM0" menu
and choosing "Export to VRM 0.x" like you did the last time. Save that to your
desktop as a VRM file.

Open your favorite VTubing software such as
[VSeeFace](https://www.vseeface.icu/) and import your avatar. Then look into
your webcam and say a sentence like this: "Alice shut the door with earnest".
Watch the avatar's lips move as you say it. There will be a very slight delay as
you say things, but overall it should roughly match up.

If you don't like how things turned out, go in and mess with the sliders. This
will take a lot of experimentation to make things look nice. This is normal.

Once you get it to a state where you are happy, save your VRM to several places
to make sure it's properly backed up. I'd suggest using DropBox, Google Drive,
iCloud Drive or anything else, as long as you have at least three copies. You
can recreate your avatar from the Unity project should you need to (or worst
case go back through this tutorial from scratch), but I'd suggest avoiding
having to do that by having proper backups. Consult the [3-2-1
rule](https://www.backblaze.com/blog/the-3-2-1-backup-strategy/) for ideas on
how/where to back things up.

Most importantly, have fun in your new life as a VTuber! It's an exciting new
frontier out there.

---

Special thanks to LithiumFox and lesliepone for reviewing this article for
content before publishing.

Want to watch these posts get written live? Check me out on
[twitch.tv](https://www.twitch.tv/princessxen)! The VOD for this post will be live
on Twitch for the next two weeks [here](https://www.twitch.tv/videos/1250763144)
and on YouTube permanently [here](https://youtu.be/ZTnFjBsm5Rs) (if you are
reading this the day of this post, the YouTube link will not be live).
