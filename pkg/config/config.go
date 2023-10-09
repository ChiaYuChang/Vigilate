package config

import (
	"crypto/x509"
	"html/template"

	"github.com/alexedwards/scs/v2"
	"github.com/pusher/pusher-http-go/v5"
	"github.com/robfig/cron/v3"
	"gitlab.com/gjerry134679/vigilate/pkg/channeldata"
	"gitlab.com/gjerry134679/vigilate/pkg/driver"
)

// AppConfig holds application configuration
type AppConfig struct {
	DB            *driver.DB
	Session       *scs.SessionManager
	InProduction  bool
	Domain        string
	MonitorMap    map[int]cron.EntryID
	PreferenceMap map[string]string
	Scheduler     *cron.Cron
	WsClient      pusher.Client
	PusherSecret  string
	TemplateCache map[string]*template.Template
	MailQueue     chan channeldata.MailJob
	Version       string
	Identifier    string
	CertPool      *x509.CertPool
	SSL           SSL
}

type SSL struct {
	CertificateFile string
	PrivateKeyFile  string
}
