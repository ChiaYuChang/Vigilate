package models

import (
	"errors"
	"time"
)

var (
	// ErrNoRecord no record found in database error
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials invalid username/password error
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail duplicate email error
	ErrDuplicateEmail = errors.New("models: duplicate email")
	// ErrInactiveAccount inactive account error
	ErrInactiveAccount = errors.New("models: Inactive Account")
)

// User model
type User struct {
	ID          int
	FirstName   string
	LastName    string
	UserActive  int
	AccessLevel int
	Email       string
	Password    []byte
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   time.Time
	Preferences map[string]string
}

type TimeInfo struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Preference model
type Preference struct {
	ID         int
	Name       string
	Preference []byte
	TimeInfo
}

type Host struct {
	ID            int
	HostName      string
	CanonicalName string
	URL           string
	IP            string
	IPv6          string
	Location      string
	OS            string
	Active        int
	HostService   []HostService
	TimeInfo
}

type Services struct {
	ID          int
	ServiceName string
	Active      int
	Icon        string
	TimeInfo
}

type ServiceStatus int

const (
	ServiceStatusHealthy ServiceStatus = iota
	ServiceStatusProblem
	ServiceStatusPending
	ServiceStatusWarning
	ServiceStatusUnknown
)

func (ss ServiceStatus) String() string {
	var s string = "unknown"
	switch ss {
	case ServiceStatusHealthy:
		s = "healthy"
	case ServiceStatusProblem:
		s = "problem"
	case ServiceStatusPending:
		s = "pending"
	case ServiceStatusWarning:
		s = "warning"
	}
	return s
}

func NewServiceStatus(s string) ServiceStatus {
	var state ServiceStatus
	switch s {
	case ServiceStatusHealthy.String():
		state = ServiceStatusHealthy
	case ServiceStatusProblem.String():
		state = ServiceStatusProblem
	case ServiceStatusPending.String():
		state = ServiceStatusPending
	case ServiceStatusWarning.String():
		state = ServiceStatusWarning
	default:
		state = ServiceStatusUnknown
	}
	return state
}

type HostService struct {
	ID             int
	HostID         int
	ServiceID      int
	Active         int
	ScheduleNumber int
	ScheduleUnit   string
	Status         ServiceStatus
	LastCheck      time.Time
	Service        Services
	TimeInfo
}
