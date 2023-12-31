package repository

import "gitlab.com/gjerry134679/vigilate/pkg/models"

// DatabaseRepo is the database repository
type DatabaseRepo interface {
	// preferences
	AllPreferences() ([]models.Preference, error)
	// SetSystemPref(name, value string) error
	SetSystemPref(name models.SystemPreference, value string) error
	InsertOrUpdateSitePreferences(pm map[string]string) error

	// users and authentication
	GetUserById(id int) (models.User, error)
	InsertUser(u models.User) (int, error)
	UpdateUser(u models.User) error
	DeleteUser(id int) error
	UpdatePassword(id int, newPassword string) error
	Authenticate(email, testPassword string) (int, string, error)
	AllUsers() ([]*models.User, error)
	InsertRememberMeToken(id int, token string) error
	DeleteToken(token string) error
	CheckForToken(id int, token string) bool

	// hosts
	InsertHost(h models.Host) (int, error)
	GetAllHost() ([]models.Host, error)
	GetHostByID(id int) (models.Host, error)
	UpdateHost(h models.Host) error
	UpdateHostServiceStatusByID(hostId, serviceId, active int) error
	UpdateHostService(hs models.HostService) error
	GetAllServiceStatusCounts() (map[models.ServiceStatus]int, error)
	GetServiceByStatus(status models.ServiceStatus) ([][6]string, error)
	GetHostServiceByID(id int) (models.HostService, error)
	GetServivesToMonitor() ([]models.HostService, []string, error)
	GetHostByHostIDServiceID(hostID, serviceID int) (models.HostService, error)
	InsertEvent(e models.Event) error
	GetAllEvent() ([]models.Event, error)
	DeleteHost(id int) error
}
