export class SignalingChannel {
  constructor(log) {
    this.log = log;
    this.ws = new WebSocket("ws://" + document.location.host + "/ws");
    this.ws.addEventListener("open", (event) => {
      log("websocket connection stabelished");
    });
    this.ws.onopen = function () {
      log("signaling service connected");
    };
  }

  setIceCallback(callback) {
    this.iceCallback = callback;
  }

  setSdpCallback(callback) {
    this.sdpCallback = callback;
  }

  async sendMessage(message) {
    try {
      const raw = JSON.stringify(message);
      this.ws.send(raw);
    } catch (e) {
      console.error(e);
      return;
    }
  }
  async send(offer) {
    try {
      const rawAnswer = await fetch("/webrtc", {
        method: "POST",
        headers: {
          Accept: "application/json",
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ type: offer.type, sdp: offer.sdp }),
      });
      const answer = await rawAnswer.json();
      return new RTCSessionDescription({ type: "answer", sdp: answer.sdp });
    } catch (e) {
      console.error(e);
      return;
    }
  }
}
