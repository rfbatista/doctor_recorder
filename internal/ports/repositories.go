package ports

import (
	"context"
	"doctor_recorder/internal/entities"
)

type SDPRepository interface {
	Save(ctx context.Context, sdp entities.SDP) (entities.SDP, error)
	GetByPeerId(ctx context.Context, peerId entities.PeerId) (entities.PeerId, error)
	DeleteByPeerId(ctx context.Context, peerId entities.PeerId) (entities.PeerId, error)
}
