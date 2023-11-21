export class Message {
  constructor() {}

  setSDP(sdp) {
    if (sdp == null) return;
    this.type = "sdp";
    this.sdp = sdp;
  }

  setICE(ice) {
    if (ice == null) return;
    this.type = "ice";
    this.ice = ice;
  }

  getJSON() {
    return {
      type: this.type,
      ice: this.ice,
      sdp: this.sdp,
    };
  }
}
