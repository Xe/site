import HLS from "npm:hls.js";

export default function execFor(id, path) {
    const video = document.getElementById(id);
    if (HLS.isSupported()) {
        const hls = new HLS();
        hls.loadSource(path);
        hls.attachMedia(video);
    } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
        video.src = path;
    }
}
