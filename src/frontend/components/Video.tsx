// @jsxImportSource xeact
// @jsxRuntime automatic

import Hls from "@hls.js";

export interface VideoProps {
  path: string;
}

export default function Video({ path }: VideoProps) {
  const streamURL =
    `https://cdn.xeiaso.net/file/christine-static/${path}/index.m3u8`;
  const video = (
    <video style="width:100%" controls>
      <source src={streamURL} type="application/vnd.apple.mpegurl" />
      <source
        src="https://cdn.xeiaso.net/file/christine-static/blog/HLSBROKE.mp4"
        type="video/mp4"
      />
    </video>
  );

  if (Hls.isSupported()) {
    const hls = new Hls();
    hls.on(Hls.Events.MEDIA_ATTACHED, () => {
      console.log("video and hls.js are now bound together !");
    });
    hls.on(Hls.Events.MANIFEST_PARSED, (event, data) => {
      console.log(
        "manifest loaded, found " + data.levels.length + " quality level",
      );
    });
    hls.loadSource(streamURL);
    hls.attachMedia(video);
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    video.src = streamURL;
  }

  return video;
}
