---
title: How I Converted my Brain fMRI to a 3D Model
date: 2019-08-23
series: howto
tags:
 - python
 - blender
---

AUTHOR'S NOTE: I just want to start this out by saying I am not an expert, and
nothing in this blogpost should be construed as medical advice. I just wanted 
to see what kind of pretty pictures I could get out of an fMRI data file.

So this week I flew out to Stanford to participate in a study that involved a
fMRI of my brain while I was doing some things. I asked for (and recieved) a
data file from the fMRI so I could play with it and possibly 3D print it. This
blogpost is the record of my journey through various software to get a fully 
usable 3D model out of the fMRI data file.

## The Data File

I was given [christine_brain.nii.gz][firstniifile] by the researcher who was
operating the fMRI. I looked around for some software to convert it to a 3D
model and [/r/3dprinting][r3dprinting] suggested the use of [FreeSurfer][freesurfer]
to generate a 3D model. I downloaded and installed the software then started
to look for something I could do in the meantime, as this was going to take
something on the order of 8 hours to process.

### An Animated GIF

I started looking for the file format on the internet by googling "nii.gz brain image"
and I stumbled across a program called [gif\_your\_nifti][gyn]. It looked to be
mostly pure python so I created a virtualenv and installed it in there:

```
$ git clone https://github.com/miykael/gif_your_nifti
$ cd gif_your_nifti
$ virtualenv -p python3 env
$ source env/bin/activate
(env) $ pip3 install -r requirements.txt
(env) $ python3 setup.py install
```

Then I ran it with the following settings to get [this first result][firstgif]:

```
(env) $ gif_your_nifti christine_brain.nii.gz --mode pseudocolor --cmap plasma
```

<center><video controls> <source src="https://xena.greedo.xeserv.us/files/christine-fmri-raw.mp4" type="video/mp4" />A sideways view of the brain</video></center>

<small>(sorry the video embed isn't working in safari)</small>

It looked weird though, that's because the fMRI scanner I used has a different
rotation to what's considered "normal". The gif\_your\_nifti repo mentioned a
program called `fslreorient2std` to reorient the fMRI image, so I set out to
install and run it.

### FSL

After some googling, I found [FSL's website][fsl] which included an installer
script and required registration.

37 gigabytes of downloads and data later, I had the entire FSL suite installed
to a server of mine and ran the conversion command:

```
$ fslreorient2std christine_brain.nii.gz christine_brain_reoriented.nii.gz
```

This produced a slightly smaller [reoriented file][secondniifile].

I reran gif\_your\_nifti on this reoriented file and got [this result][secondgif]
which looked a _lot_ better:

<center><video controls> <source src="https://xena.greedo.xeserv.us/files/christine-fmri-reoriented.mp4" />A properly reoriented brain</video></center>

<small>(sorry again the video embed isn't working in safari)</small>

### FreeSurfer

By this time I had gotten back home and [FreeSurfer][freesurfer] was done installing, 
so I registered for it (god bless the institution of None) and put its license key
in the place it expected. I copied the reoriented data file to my Mac and then
set up a `SUBJECTS_DIR` and had it start running the numbers and extracting the
brain surfaces:

```
$ cd ~/tmp
$ mkdir -p brain/subjects
$ cd brain
$ export SUBJECTS_DIR=$(pwd)/subjects
$ recon-all -i /path/to/christine_brain_reoriented.nii.gz -s christine -all
```

This step took 8 hours. Once I was done I had a bunch of data in 
`$SUBJECTS_DIR/christine`. I opened my shell to that folder and went into the
`surf` subfolder:

```
$ mris_convert lh.pial lh.pial.stl
$ mris_convert rh.pial rh.pial.stl
```

Now I had standard stl files that I could stick into [Blender][blender].

### Blender

Importing the stl files was really easy. I clicked on File, then Import, then
Stl. After guiding the browser to the subjects directory and finding the STL
files, I got a view that looked something like this:

<center><blockquote class="twitter-tweet"><p lang="en" dir="ltr">BRAIN <a href="https://t.co/kGSrPj0kgP">pic.twitter.com/kGSrPj0kgP</a></p>&mdash; Cadey Ratio 🌐 (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1164526098526478336?ref_src=twsrc%5Etfw">August 22, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

I had absolutely no idea what to do from here in Blender, so I exported the
whole thing to a stl file and sent it to a coworker for 3D printing (he said
it was going to be "the coolest thing he's ever printed").

I also exported an Unreal Engine 4 compatible model and sent it to a friend of
mine that does hobbyist game development. A few hours later I got this back:

<center><blockquote class="twitter-tweet"><p lang="und" dir="ltr"><a href="https://t.co/fXnwnSpMry">pic.twitter.com/fXnwnSpMry</a></p>&mdash; Cadey Ratio 🌐 (@theprincessxena) <a href="https://twitter.com/theprincessxena/status/1164714830630203393?ref_src=twsrc%5Etfw">August 23, 2019</a></blockquote> <script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script></center>

(Hint: it is a take on the famous [galaxy brain memes][galaxybrain])

## Conclusion

Overall, this was fun! I got to play with many gigabytes of software that ran
my most powerful machine at full blast for 8 hours, I made a fully printable 3D
model out of it and I have some future plans for importing this data into
Minecraft (the NIFTI `.nii.gz` format has a limit of _256 layers_). 

I'll be sure to write more about this in the future!

## Citations

Here are my citations in [BibTex format][citations].

Special thanks goes to Michael Lifshitz for organizing the study that I 
participated in that got me this fMRI data file. It was one of the coolest
things I've ever done (if not the coolest) and I'm going to be able to get a
3D printed model of my brain out of it.

[firstniifile]: https://xena.greedo.xeserv.us/files/christine_brain.nii.gz
[secondniifile]: https://xena.greedo.xeserv.us/files/christine_brain_reoriented.nii.gz
[r3dprinting]: https://www.reddit.com/r/3Dprinting/comments/2w0zxx/magnetic_resonance_image_nii_to_stl/
[freesurfer]: https://surfer.nmr.mgh.harvard.edu/fswiki/FreeSurferWiki
[gyn]: https://github.com/miykael/gif_your_nifti
[firstgif]: /static/blog/christine-fmri-raw.mp4
[secondgif]: /static/blog/christine-fmri-reoriented.mp4
[fsl]: https://fsl.fmrib.ox.ac.uk/fsl/fslwiki/
[blender]: https://www.blender.org
[galaxybrain]: https://knowyourmeme.com/memes/expanding-brain
[citations]: /static/blog/brainfmri-to-3d-model.bib
