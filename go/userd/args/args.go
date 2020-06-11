package args

import (
	"errors"
	"strings"

	"github.com/GDVFox/tenjin/userd/database"
)

var (
	errSkillDifficulty = errors.New("wrong value, expected one of (low, medium, hard, unsolvable)")
	errUpdateType      = errors.New("wrong value, expected one of (add, update, delete)")
	errVacancyPriority = errors.New("wrong value, expected one of (low, medium, high)")
)

// SkillDifficultyFilter task difficulty
type SkillDifficultyFilter struct {
	Value database.SkillDifficulty
}

// Parse parses a string to a PersonType
func (a *SkillDifficultyFilter) Parse(s string) error {
	switch s {
	case "low":
		a.Value = database.LowSkillDifficulty
	case "medium":
		a.Value = database.MediumSkillDifficulty
	case "hard":
		a.Value = database.HardSkillDifficulty
	case "unsolvable":
		a.Value = database.UnsolvableSkillDifficulty
	default:
		return errSkillDifficulty
	}
	return nil
}

// UnmarshalJSON parses a []byte to a SkillDifficulty
func (a *SkillDifficultyFilter) UnmarshalJSON(s []byte) error {
	return a.Parse(strings.Trim(string(s), "\""))
}

// UpdateType represents type for list update
type UpdateType uint8

const (
	// AddUpdateType add new element
	AddUpdateType SkillUpdateType = iota
	// DeleteUpdateType remove elemetn
	DeleteUpdateType
)

// UnmarshalJSON parses a []byte to a SkillUpdateType
func (t *UpdateType) UnmarshalJSON(s []byte) error {
	representation := strings.Trim(string(s), "\"")
	switch representation {
	case "add":
		*t = AddUpdateType
	case "delete":
		*t = DeleteUpdateType
	default:
		return errUpdateType
	}

	return nil
}

// SkillUpdateType represents type for update skills
type SkillUpdateType = UpdateType

// AttachmentUpdateType represents type for update attachments
type AttachmentUpdateType = UpdateType

// RequirementCheckUpdateType represents type for update requirement check
type RequirementCheckUpdateType = UpdateType

// VacacnyPriority vacancy priority
type VacacnyPriority struct {
	Value database.VacancyPriority
}

// GetValue returns database.VacancyPriority value
func (a *VacacnyPriority) GetValue() *database.VacancyPriority {
	if a == nil {
		return nil
	}

	return &a.Value
}

// Parse parses a string to a PersonType
func (a *VacacnyPriority) Parse(s string) error {
	switch s {
	case "low":
		a.Value = database.VacancyPriorityLow
	case "medium":
		a.Value = database.VacancyPriorityMedium
	case "high":
		a.Value = database.VacancyPriorityHigh
	default:
		return errVacancyPriority
	}
	return nil
}

// UnmarshalJSON parses a []byte to a VacacnyPriority
func (a *VacacnyPriority) UnmarshalJSON(s []byte) error {
	return a.Parse(strings.Trim(string(s), "\""))
}
