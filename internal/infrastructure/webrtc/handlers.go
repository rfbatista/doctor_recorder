package webrtc

import "github.com/pion/webrtc/v3"

type OnNewICECandidateCallback func(ice *webrtc.ICECandidate)
type OnSendOfferCallback func(sd *webrtc.SessionDescription)

type TrackHandler func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)
