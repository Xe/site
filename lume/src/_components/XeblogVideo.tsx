import Hls from "npm:hls.js";
import { sha256 } from "npm:js-sha256";

export interface VideoProps {
  path: string;
  vertical?: boolean;
}

export default function Video({ path, vertical }: VideoProps) {
  const streamURL = `https://cdn.xeiaso.net/file/christine-static/${path}/index.m3u8`;
  const id = sha256(streamURL);
  const video = (
    <video id={id} className="not-prose sm:max-h-[50vh] mx-auto" controls>
      <source src={streamURL} type="application/vnd.apple.mpegurl" />
      <source
        src="https://cdn.xeiaso.net/file/christine-static/blog/HLSBROKE.mp4"
        type="video/mp4"
      />
    </video>
  );

  const script = (
    <script type="module">
      {`import execFor from ${`'`}/js/hls.js${`'`};
      
      execFor(${`'`}${id}${`'`}, ${`'`}${streamURL}${`'`});`}
    </script>
  );

  return (
    <>
      {video}
      {script}
      <div className="text-center">
        Want to watch this in your video player of choice? Take this:
        <br />
        <a href={streamURL} rel="noreferrer">
          {streamURL}
        </a>
      </div>
    </>
  );
}
