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
)

type HostService struct {
	ID             string
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
