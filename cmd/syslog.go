package cmd

import (
	"fmt"
	"net"

	"github.com/crewjam/rfc5424"
	"github.com/sirupsen/logrus"
)

// TODO move this to it's own package
// reconnect?
// filtering? limiting output? maybe? probably leave it for syslog-ng/rsyslog and co ?
// benchmark it
// buffer messages before sending them

// loosely based on https://godoc.org/github.com/sirupsen/logrus/hooks/syslog
type syslogHook struct {
	Writer        net.Conn
	SyslogNetwork string
	SyslogRaddr   string
}

func newSyslogHook(network, raddr string) (*syslogHook, error) {
	w, err := net.Dial(network, raddr)
	return &syslogHook{w, network, raddr}, err
}

func (hook *syslogHook) Fire(entry *logrus.Entry) error {
	sdParams := []rfc5424.SDParam{
		{Name: "level", Value: entry.Level.String()},
	}
	for name, value := range entry.Data {
		// TODO maybe do it only for some?
		// TODO have custom logic for different things ?
		sdParams = append(sdParams, rfc5424.SDParam{Name: name, Value: fmt.Sprint(value)})
	}

	m := rfc5424.Message{
		Priority:  rfc5424.Daemon | rfc5424.Info, // TODO figure this out
		Timestamp: entry.Time,
		Message:   []byte(entry.Message),
		StructuredData: []rfc5424.StructuredData{
			{
				ID:         "k6",
				Parameters: sdParams,
			},
		},
	}
	_, err := m.WriteTo(hook.Writer)
	return err
}

func (hook *syslogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
