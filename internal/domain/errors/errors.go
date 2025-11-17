package errors

import "errors"

// Assessment errors
var (
	ErrInvalidAssessmentID    = errors.New("domain: invalid assessment ID")
	ErrInvalidMaterialID      = errors.New("domain: invalid material ID")
	ErrInvalidMongoDocumentID = errors.New("domain: mongo document ID must be exactly 24 characters")
	ErrEmptyTitle             = errors.New("domain: assessment title cannot be empty")
	ErrInvalidTotalQuestions  = errors.New("domain: total questions must be between 1 and 100")
	ErrInvalidPassThreshold   = errors.New("domain: pass threshold must be between 0 and 100")
	ErrInvalidMaxAttempts     = errors.New("domain: max attempts must be at least 1")
	ErrInvalidTimeLimit       = errors.New("domain: time limit must be between 1 and 180 minutes")
)

// Attempt errors
var (
	ErrInvalidAttemptID        = errors.New("domain: invalid attempt ID")
	ErrInvalidStudentID        = errors.New("domain: invalid student ID")
	ErrInvalidScore            = errors.New("domain: score must be between 0 and 100")
	ErrInvalidTimeSpent        = errors.New("domain: time spent must be positive and <= 7200 seconds")
	ErrInvalidStartTime        = errors.New("domain: invalid start time")
	ErrInvalidEndTime          = errors.New("domain: end time must be after start time")
	ErrAttemptAlreadyCompleted = errors.New("domain: attempt already completed, cannot modify")
	ErrNoAnswersProvided       = errors.New("domain: at least one answer must be provided")
)

// Answer errors
var (
	ErrInvalidAnswerID         = errors.New("domain: invalid answer ID")
	ErrInvalidQuestionID       = errors.New("domain: invalid question ID")
	ErrInvalidSelectedAnswerID = errors.New("domain: invalid selected answer ID")
)
