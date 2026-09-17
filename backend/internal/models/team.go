package models

import "time"

type Team struct {
	Base
	EventID   string `gorm:"type:uuid;not null;index" json:"event_id"`
	CollegeID string `gorm:"type:uuid;not null" json:"college_id"`
	Name      string `gorm:"size:150;not null" json:"name"`
	Code      string `gorm:"size:20;not null" json:"code"`
	LeaderID  string `gorm:"type:uuid;not null;index" json:"leader_id"`
	IsLocked  bool   `gorm:"not null;default:false" json:"is_locked"`

	Members []TeamMember `gorm:"foreignKey:TeamID" json:"members,omitempty"`
}

func (Team) TableName() string { return "teams" }

type TeamMemberRole string

const (
	TeamRoleLeader TeamMemberRole = "LEADER"
	TeamRoleMember TeamMemberRole = "MEMBER"
)

type TeamMember struct {
	BaseNoSoftDelete
	TeamID   string         `gorm:"type:uuid;not null;index" json:"team_id"`
	UserID   string         `gorm:"type:uuid;not null;index" json:"user_id"`
	Role     TeamMemberRole `gorm:"size:20;not null;default:'MEMBER'" json:"role"`
	JoinedAt time.Time      `json:"joined_at"`
}

func (TeamMember) TableName() string { return "team_members" }
