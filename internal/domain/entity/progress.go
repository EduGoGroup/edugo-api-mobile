package entity

import (
	"time"

	"github.com/EduGoGroup/edugo-api-mobile/internal/domain/valueobject"
	"github.com/EduGoGroup/edugo-shared/pkg/errors"
	"github.com/EduGoGroup/edugo-shared/pkg/types/enum"
)

// Progress representa el progreso de lectura de un material
type Progress struct {
	materialID     valueobject.MaterialID
	userID         valueobject.UserID
	percentage     int
	lastPage       int
	status         enum.ProgressStatus
	lastAccessedAt time.Time
	createdAt      time.Time
	updatedAt      time.Time
}

func NewProgress(materialID valueobject.MaterialID, userID valueobject.UserID) *Progress {
	now := time.Now()
	return &Progress{
		materialID:     materialID,
		userID:         userID,
		percentage:     0,
		lastPage:       0,
		status:         enum.ProgressStatusNotStarted,
		lastAccessedAt: now,
		createdAt:      now,
		updatedAt:      now,
	}
}

func ReconstructProgress(
	materialID valueobject.MaterialID,
	userID valueobject.UserID,
	percentage, lastPage int,
	status enum.ProgressStatus,
	lastAccessedAt, createdAt, updatedAt time.Time,
) *Progress {
	return &Progress{
		materialID:     materialID,
		userID:         userID,
		percentage:     percentage,
		lastPage:       lastPage,
		status:         status,
		lastAccessedAt: lastAccessedAt,
		createdAt:      createdAt,
		updatedAt:      updatedAt,
	}
}

func (p *Progress) MaterialID() valueobject.MaterialID { return p.materialID }
func (p *Progress) UserID() valueobject.UserID         { return p.userID }
func (p *Progress) Percentage() int                    { return p.percentage }
func (p *Progress) LastPage() int                      { return p.lastPage }
func (p *Progress) Status() enum.ProgressStatus        { return p.status }
func (p *Progress) LastAccessedAt() time.Time          { return p.lastAccessedAt }
func (p *Progress) CreatedAt() time.Time               { return p.createdAt }
func (p *Progress) UpdatedAt() time.Time               { return p.updatedAt }

func (p *Progress) UpdateProgress(percentage, lastPage int) error {
	if percentage < 0 || percentage > 100 {
		return errors.NewValidationError("percentage must be between 0 and 100")
	}

	p.percentage = percentage
	p.lastPage = lastPage
	p.lastAccessedAt = time.Now()
	p.updatedAt = time.Now()

	if percentage == 0 {
		p.status = enum.ProgressStatusNotStarted
	} else if percentage >= 100 {
		p.status = enum.ProgressStatusCompleted
	} else {
		p.status = enum.ProgressStatusInProgress
	}

	return nil
}
