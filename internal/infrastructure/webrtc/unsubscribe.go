package webrtc

func (w *WebRTCServer) Unsubscribe(peerId PeerID) {
	if _, exist := w.Peers[peerId]; exist {
		conn := w.Peers[peerId]
		conn.Close()
		delete(w.Peers, peerId)
	}
}
