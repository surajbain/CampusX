package event

import (
	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/apperror"
)

// AllowedTransitions defines the state machine for event lifecycle.
//
//	DRAFT     ──► PUBLISHED | CANCELLED
//	PUBLISHED ──► ONGOING   | CANCELLED
//	ONGOING   ──► COMPLETED | CANCELLED
//	COMPLETED ──► (terminal)
//	CANCELLED ──► (terminal)
var AllowedTransitions = map[models.EventStatus][]models.EventStatus{
	models.StatusDraft:     {models.StatusPublished, models.StatusCancelled},
	models.StatusPublished: {models.StatusOngoing, models.StatusCancelled},
	models.StatusOngoing:   {models.StatusCompleted, models.StatusCancelled},
	models.StatusCompleted: {},
	models.StatusCancelled: {},
}

// CanTransition returns nil if from→to is allowed, otherwise an AppError.
func CanTransition(from, to models.EventStatus) error {
	allowed, ok := AllowedTransitions[from]
	if !ok {
		return apperror.Conflict("unknown current status: " + string(from))
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return apperror.Conflict(
		"invalid transition: cannot move from " + string(from) + " to " + string(to),
	)
}
