package scopes

import (
	"gorm.io/gorm"
)

// College restricts a query to a single college's data.
func College(collegeID string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if collegeID == "" {
			return db.Where("1 = 0")
		}
		return db.Where("college_id = ?", collegeID)
	}
}

// ActiveCollege also filters to active colleges.
func ActiveCollege(collegeID string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if collegeID == "" {
			return db.Where("1 = 0")
		}
		return db.Where("college_id = ? AND is_active = TRUE", collegeID)
	}
}

// Published filters events that are publicly visible.
func Published() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = 'PUBLISHED'")
	}
}

// Featured filters featured events.
func Featured() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("is_featured = TRUE")
	}
}

// EventOwner filters events created by a specific user.
func EventOwner(userID string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if userID == "" {
			return db.Where("1 = 0")
		}
		return db.Where("created_by = ?", userID)
	}
}

// EventVisible returns published events OR events owned by user.
// Used for "my events" where drafts should be visible.
func EventVisible(userID string) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if userID == "" {
			return Published()(db)
		}
		return db.Where("status = 'PUBLISHED' OR created_by = ?", userID)
	}
}

// NotDeleted — explicit (GORM does this automatically for Base models).
func NotDeleted() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("deleted_at IS NULL")
	}
}
