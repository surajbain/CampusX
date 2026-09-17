package audit

import (
	"context"
	"encoding/json"

	"gorm.io/gorm"

	"github.com/campusx/api/internal/models"
	"github.com/campusx/api/pkg/logger"
)

type Logger struct {
	db *gorm.DB
}

func NewLogger(db *gorm.DB) *Logger {
	return &Logger{db: db}
}

type Entry struct {
	ActorID    *string
	CollegeID  *string
	Action     string
	EntityType string
	EntityID   *string
	Severity   string
	Metadata   map[string]any
	IPAddress  string
	UserAgent  string
}

// Log writes an audit entry. Never returns an error to caller — logging
// must not break the request. Errors are themselves logged.
func (l *Logger) Log(ctx context.Context, e Entry) {
	severity := e.Severity
	if severity == "" {
		severity = string(models.SeverityInfo)
	}

	var metaJSON []byte
	if len(e.Metadata) > 0 {
		var err error
		metaJSON, err = json.Marshal(e.Metadata)
		if err != nil {
			logger.From().Error().Err(err).Str("action", e.Action).Msg("audit: marshal metadata")
			metaJSON = []byte(`{"error":"marshal_failed"}`)
		}
	}

	entry := models.AuditLog{
		ActorID:    e.ActorID,
		CollegeID:  e.CollegeID,
		Action:     e.Action,
		EntityType: e.EntityType,
		EntityID:   e.EntityID,
		Severity:   models.AuditSeverity(severity),
		IPAddress:  e.IPAddress,
		UserAgent:  e.UserAgent,
	}
	if metaJSON != nil {
		entry.Metadata = metaJSON
	}

	if err := l.db.WithContext(ctx).Create(&entry).Error; err != nil {
		logger.From().Error().
			Err(err).
			Str("action", e.Action).
			Msg("audit: write failed")
	}
}
