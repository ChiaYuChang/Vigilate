package models

import (
	"errors"
	"time"

	"github.com/robfig/cron/v3"
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

type SystemPreference int

const (
	MonitoringLive SystemPreference = iota + 1
	CheckIntervalAmount
	CheckIntervalUnit
	NotifyViaEmail
)

func (ss SystemPreference) String() string {
	var s string = "unknown"
	switch ss {
	case CheckIntervalAmount:
		s = "check_interval_amount"
	case CheckIntervalUnit:
		s = "check_interval_unit"
	case NotifyViaEmail:
		s = "notify_via_email"
	case MonitoringLive:
		s = "monitoring_live"
	}
	return s
}

func ParseSystemPreference(s string) SystemPreference {
	var state SystemPreference
	switch s {
	case CheckIntervalAmount.String():
		state = CheckIntervalAmount
	case CheckIntervalUnit.String():
		state = CheckIntervalUnit
	case NotifyViaEmail.String():
		state = NotifyViaEmail
	case MonitoringLive.String():
		state = MonitoringLive
	}
	return state
}

type Schedule struct {
	ID            int
	EntryID       cron.EntryID
	Entry         cron.Entry
	Host          string
	Service       string
	LastRunFromHS time.Time
	HostServiceID int
	ScheduleText  string
}

type Event struct {
	ID            int
	Type          string
	HostServiceID int
	HostID        int
	HostName      string
	ServiceID     int
	ServiceName   string
	Message       string
	TimeInfo
}
