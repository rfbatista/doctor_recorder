import { SignalingChannel } from "./signaling.js";


let localStream;

var log = (msg) => {
  const logs = document.querySelector("#logs");
  logs.innerHTML += msg + "<br>";
};

const signaling = new SignalingChannel(log);
// Put variables in global scope to make them available to the browser console.
let constraints = (window.constraints = {
  audio: true,
  video: true,
});

/**
 * Start session
 * @param {string} RTCPeerConnection
 */
async function startSession(pc) {
  log("starting session");
  const offer = await pc.createOffer();
  pc.setLocalDescription(offer);
  const answer = await signaling.sendOffer(offer);
  pc.setRemoteDescription(answer);
  /* startSignaling(pc); */
  log("finished to start session");
}

/**
 * Create session
 * @returns {RTCPeerConnection}
 */
function createSession() {
  const configuration = {
    iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  };
  const pc = new RTCPeerConnection(configuration);
  pc.oniceconnectionstatechange = () =>
    log(`state changed", ${pc.iceConnectionState}"`);
  pc.onicecandidate = async (event) => {
    log("ice candidate identified")
    console.log(event.candidate)
    if (event.candidate === null) {
      /* document.getElementById("localSessionDescription").value = btoa( */
      /*   JSON.stringify(pc.localDescription), */
      /* ); */
    }
  };
  // Listen for connectionstatechange on the local RTCPeerConnection
  pc.addEventListener("connectionstatechange", (event) => {
    console.log("state", pc.connectionState);
    if (pc.connectionState === "connected") {
      log("peer connected!")
    }
  });
  return pc;
}

/**
 * Handle success of media access
 * @param {MediaStream} stream
 */
async function handleSuccess(stream) {
  localStream = stream;
  window.stream = stream; // make variable available to browser console
  const tracks = stream.getTracks();
  if (tracks.length > 0) {
    log(`Using Audio device: ${tracks[0].label}`);
  }
  const pc = createSession();
  for (const track of stream.getAudioTracks()) {
    log("adding track to Peer Connection");
    pc.addTrack(track, stream);
  }
  await startSession(pc);
}

function handleError(error) {
  if (error.name === "OverconstrainedError") {
    const v = constraints.video;
    log(
      `The resolution ${v.width.exact}x${v.height.exact} px is not supported by your device.`,
    );
  } else if (error.name === "NotAllowedError") {
    log(
      "Permissions have not been granted to use your camera and " +
        "microphone, you need to allow the page access to your devices in " +
        "order for the demo to work.",
    );
  }
  errorMsg(`getUserMedia error: ${error.name}`, error);
}

function errorMsg(msg, error) {
  if (typeof error !== "undefined") {
    console.error(error);
  }
  console.error(msg);
}

async function init(e) {
  try {
    const stream = await navigator.mediaDevices.getUserMedia(constraints);
    const video = document.querySelector("#video1");
    video.srcObject = stream;
    video.onloadedmetadata = () => {
      video.play();
    };
    handleSuccess(stream);
    e.target.disabled = true;
  } catch (e) {
    handleError(e);
  }
}

document.querySelector("#showVideo").addEventListener("click", (e) => init(e));
