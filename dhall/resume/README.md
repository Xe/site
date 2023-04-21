# Alta Typst

A simple Typst CV template, inspired by [AltaCV by LianTze Lim](https://github.com/liantze/AltaCV). [LaurenzV's simplecv](https://github.com/LaurenzV/simplecv) was used as a Typst code reference. See [`example.pdf`](example.pdf) for the rendered PDF output.

<img src="screenshot.png" width="500"></img>

## Usage

### On [typst.app](https://typst.app/)

Upload the `.typ` and `.svg` files to your Typst project, then see `example.typ`.

### With [Typst CLI](https://github.com/typst/typst)

Fork and clone this repo, then run `typst --watch example.typ`.

Note that the template is intended for use with the IBM Plex Sans font, which isn't currenly bundled with the Typst CLI — so the rendered output on your machine may differ slightly. You can run `typst --fonts` to see which fonts can be used instead.

### Icons

Add extra icons by uploading more `.svg` files — the existing ones are free icons from [Font Awsome](https://fontawesome.com/search?o=r&m=free). You can then reference these as name values in the links array.

## Licence

[MIT](./LICENSE)

Icons are from Font Awesome, subject to [their terms](https://github.com/FortAwesome/Font-Awesome/blob/6.x/LICENSE.txt).
