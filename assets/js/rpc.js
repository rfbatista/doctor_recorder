import { SignalingChannel } from "./signaling.js";
import { Message } from "./message.js";

const configuration = {
  iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
};

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
 * Handle success of media access
 * @param {MediaStream} stream
 */
async function startRPC(stream) {
  window.stream = stream; // make variable available to browser console
  const tracks = stream.getTracks();
  if (tracks.length > 0) {
    log(`Using Audio device: ${tracks[0].label}`);
  }
  const pc = startSession(stream);
  for (const track of stream.getAudioTracks()) {
    log("adding track to Peer Connection");
    pc.addTrack(track, stream);
  }
}

/**
 * Start session,
 * @param {MediaStream} stream
 */
async function startSession(stream) {
  log("starting session");
  const pc = new RTCPeerConnection(configuration);
  pc.onicecandidate = async (event) => {
    if (!event.candidate) {
      return;
    }
    const message = new Message();
    message.setICE(event.candidate.candidate);
    await signaling.sendMessage(message.getJSON());
  };
  pc.addEventListener("connectionstatechange", () => {
    log(`connection state change ${pc.connectionState}`);
    if (pc.connectionState === "connected") {
      log("peer connected!");
    }
  });
  signaling.ws.addEventListener("message", (evt) => {
    const data = JSON.parse(evt.data);
    if (data.type == "sdp") {
      if (!pc.currentRemoteDescription) {
        log("received a sdp answer");
        const answer = JSON.parse(data.sdp);
        const answerDescription = new RTCSessionDescription(answer);
        pc.setRemoteDescription(answerDescription);
      }
    }
    if (data.type == "ice") {
      log("received an ice candidate");
      const candidate = new RTCIceCandidate(data.ice);
      pc.addIceCandidate(candidate);
    }
  });
  stream.getTracks().forEach((track) => pc.addTrack(track, stream));
  const offer = await pc.createOffer();
  pc.setLocalDescription(offer);
  const message = new Message();
  message.setSDP(offer);
  signaling.sendMessage(message.getJSON()).catch(console.error);
}

function handleError(error) {
  if (error.name === "OverconstrainedError") {
    const v = constraints.video;
    log(
      `The resolution ${v.width.exact}x${v.height.exact} px is not supported by your device.`
    );
  } else if (error.name === "NotAllowedError") {
    log(
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
    const video = document.querySelector("#video1");
    video.srcObject = stream;
    video.onloadedmetadata = () => {
      video.play();
    };
    startRPC(stream);
    e.target.disabled = true;
  } catch (e) {
    handleError(e);
  }
}

document.querySelector("#showVideo").addEventListener("click", (e) => init(e));
