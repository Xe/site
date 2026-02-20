# Picture Component Guide

Use `<Picture>` and `<XeblogPicture>` to embed images in posts. The `path` value should point to the uploaded asset location (see Uploading below). Always include a clear `desc` when the component supports it.

## Components

```jsx
<Picture
  path="blog/2025/squandered-holy-grail/berlin-tv-tower"
  desc="An AI-generated illustration of the East Berlin TV tower at sunset."
/>

<XeblogPicture
  path="blog/2023/xeact/star-history-202342"
/>
```

## Path Rules

- The `path` value is the canonical storage path for the asset.
- The `path` value should match the upload path used with `uploud`.

## Uploading Images

Upload images using the `uploud` command in the `Xe/x` repo on GitHub.

Syntax:

```bash
$ uploud <local path> <folder path in tigris>
```

Or for a folder of files:

```bash
$ uploud <folder> <path in tigris>
```

Notes:

- The path supplied should be the `path` argument passed to the `Picture` component.
- If uploading a folder, each filename becomes the final element in the path.
