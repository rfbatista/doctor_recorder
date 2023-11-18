import { SignalingChannel } from "./signaling.js";

const signaling = new SignalingChannel();

let localStream;

// Put variables in global scope to make them available to the browser console.
let constraints = (window.constraints = {
  audio: true,
  video: false,
});

/**
 * Start session
 * @param {string} RTCPeerConnection
 */
async function startSession(pc) {
  console.log("starting");
  const offer = await pc.createOffer();
  pc.setLocalDescription(offer);
  const answer = await signaling.send(offer);
  pc.setRemoteDescription(answer);
  /* startSignaling(pc); */
  console.log("finished init");
}

/**
 * Create session
 * @returns {RTCPeerConnection}
 */
async function createSession() {
  const configuration = {
    iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
  };
  const pc = new RTCPeerConnection(configuration);
  pc.oniceconnectionstatechange = () =>
    console.log("state changed", pc.iceConnectionState);
  pc.onicecandidate = async (event) => {
    console.log(event);
    if (event.candidate === null) {
      document.getElementById("localSessionDescription").value = btoa(
        JSON.stringify(pc.localDescription)
      );
    }
  };
  // Listen for connectionstatechange on the local RTCPeerConnection
  pc.addEventListener("connectionstatechange", (event) => {
    console.log("state", pc.connectionState);
    if (pc.connectionState === "connected") {
      // Peers connected!
    }
  });
  return pc;
}

async function handleSuccess(stream) {
  localStream = stream;
  window.stream = stream; // make variable available to browser console
  const audioTracks = stream.getAudioTracks();
  if (audioTracks.length > 0) {
    console.log(`Using Audio device: ${audioTracks[0].label}`);
  }
  const pc = createSession();
  console.log(`Using video device: ${audioTracks[0].label}`);
  localStream.getTracks().forEach((track) => pc.addTrack(track, localStream));
  await startSession(pc);
}

function handleError(error) {
  if (error.name === "OverconstrainedError") {
    const v = constraints.video;
    errorMsg(
      `The resolution ${v.width.exact}x${v.height.exact} px is not supported by your device.`
    );
  } else if (error.name === "NotAllowedError") {
    errorMsg(
      "Permissions have not been granted to use your camera and " +
        "microphone, you need to allow the page access to your devices in " +
        "order for the demo to work."
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
    handleSuccess(stream);
    e.target.disabled = true;
  } catch (e) {
    handleError(e);
  }
}

document.querySelector("#showVideo").addEventListener("click", (e) => init(e));
