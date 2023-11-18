export class SignalingChannel {
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
