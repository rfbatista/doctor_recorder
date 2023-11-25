package webrtc

import "github.com/pion/webrtc/v4"

type OnNewICECandidateCallback func(ice *webrtc.ICECandidate)
type OnSendOfferCallback func(sd *webrtc.SessionDescription)

type TrackHandler func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver)
