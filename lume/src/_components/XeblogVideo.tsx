// @jsxImportSource xeact
// @jsxRuntime automatic

import Hls from "npm:hls.js";

function uuidv4() {
  return ([1e7]+-1e3+-4e3+-8e3+-1e11).replace(/[018]/g, c =>
    (c ^ crypto.getRandomValues(new Uint8Array(1))[0] & 15 >> c / 4).toString(16)
  );
}

export interface VideoProps {
  path: string;
  vertical?: boolean;
}

export default function Video({ path, vertical }: VideoProps) {
  const streamURL =
    `https://cdn.xeiaso.net/file/christine-static/${path}/index.m3u8`;
  const id = uuidv4();
  const video = (
      <video id={id} className="not-prose sm:max-h-[50vh]" controls>
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
  )

  return <>
    {video}
    {script}
  </>;
}
