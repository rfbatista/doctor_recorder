package webrtc

import (
	"fmt"
)

var (
	FailedToCreateAnswer            = &FailedToCreateAnswerError{}
	FailedToSetLocalDescription     = &FailedToSetLocalDescriptionError{}
	FailedToCreateNewPeerConnection = &FailedToCreateNewPeerConnectionError{}
)

type FailedToCreateAnswerError struct {
	Err error
}

func (e *FailedToCreateAnswerError) Error() string {
	return fmt.Sprintf("FailedToCreateAnswerError: %v", e.Err)
}

type FailedToSetLocalDescriptionError struct {
	Err error
}

func (e *FailedToSetLocalDescriptionError) Error() string {
	return fmt.Sprintf("FailedToSetLocalDescriptionError: %v", e.Err)
}

type FailedToCreateNewPeerConnectionError struct {
	Err error
}

func (e *FailedToCreateNewPeerConnectionError) Error() string {
	return fmt.Sprintf("FailedToCreateNewPeerConnectionError: %v", e.Err)
}
