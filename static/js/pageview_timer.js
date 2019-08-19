/*
  Hi,

  If you are reading this, you have found this script in the referenced scripts
  for pages on this site. I know you're gonna have to take me at my word on this,
  but I'm literally using this to collect how much time people spend reading my
  webpages. See metrics here: https://christine.website/metrics

  If you have the "do not track" setting enabled in your browser, this code will
  be ineffectual.
*/

(function() {
    let dnt = navigator.doNotTrack;
    if (dnt == "1") {
        return;
    }

    let startTime = new Date();

    function logTime() {
        let stopTime = new Date();
        let message = JSON.stringify(
            {
                "path": window.location.pathname,
                "start_time": startTime.toISOString(),
                "end_time": stopTime.toISOString(),
            }
        );

        window.navigator.sendBeacon("/api/pageview-timer", message);
    }

    window.addEventListener("pagehide", logTime, false);
})();
