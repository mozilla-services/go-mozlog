package mozlog // import "go.mozilla.org/mozlog"

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

var Logger = &MozLogger{
	Output:     os.Stdout,
	LoggerName: "Application",
}

var hostname string

func Hostname() string {
	return hostname
}

// MozLogger implements the io.Writer interface
type MozLogger struct {
	Output     io.Writer
	LoggerName string
}

func init() {
	var err error
	hostname, err = os.Hostname()
	if err != nil {
		log.Printf("Can't resolve hostname: %v", err)
	}

	log.SetOutput(Logger)
	log.SetFlags(log.Lshortfile)
}

// Write converts the log to AppLog
func (m *MozLogger) Write(l []byte) (int, error) {
	log := NewAppLog(m.LoggerName, l)

	out, err := log.ToJSON()
	if err != nil {
		// Need someway to notify that this happened.
		fmt.Fprintln(os.Stderr, err)
		return 0, err
	}

	_, err = m.Output.Write(append(out, '\n'))
	return len(l), err
}

// AppLog implements Mozilla logging standard
type AppLog struct {
	Timestamp  int64
	Time       string
	Type       string
	Logger     string
	Hostname   string `json:",omitempty"`
	EnvVersion string
	Pid        int `json:",omitempty"`
	Severity   int `json:",omitempty"`
	Fields     map[string]interface{}
}

// NewAppLog returns a loggable struct
func NewAppLog(loggerName string, msg []byte) *AppLog {
	now := time.Now().UTC()
	return &AppLog{
		Timestamp:  now.UnixNano(),
		Time:       now.Format(time.RFC3339),
		Type:       "app.log",
		Logger:     loggerName,
		Hostname:   hostname,
		EnvVersion: "2.0",
		Pid:        os.Getpid(),
		Fields: map[string]interface{}{
			"msg": string(bytes.TrimSpace(msg)),
		},
	}
}

// ToJSON converts a logline to JSON
func (a *AppLog) ToJSON() ([]byte, error) {
	return json.Marshal(a)
}

func FromJSON(log string) (a AppLog, err error) {
	err = json.Unmarshal([]byte(log), &a)
	return
}

func (a *AppLog) ToString() string {
	str := fmt.Sprintf("Timestamp=%d Time=%q Type=%q Logger=%q Hostname=%q EnvVersion=%q Pid=%d Severity=%d Fields[",
		a.Timestamp, a.Time, a.Type, a.Logger, a.Hostname, a.EnvVersion, a.Pid, a.Severity)
	var namedFields []string
	for name, value := range a.Fields {
		val, _ := json.Marshal(value)
		namedFields = append(namedFields, name+"="+string(val))
	}
	str += strings.Join(namedFields, ", ")
	str += "]"
	return str
}

func (a *AppLog) Evaluate() {
	if a.Timestamp == 0 {
		fmt.Println("error: nanosecond Timestamp is missing")
	}
	if a.Time == "" {
		fmt.Println("info: RFC3339 Time not set")
	}
	if a.Type == "" {
		fmt.Println("error: Type is missing")
	}
	if a.Logger == "" {
		fmt.Println("error: Logger is missing")
	}
	if a.Hostname == "" {
		fmt.Println("error: Hostname is missing")
	}
	if a.EnvVersion != "2.0" {
		fmt.Printf("error: EnvVersion should be 2.0, not %q", a.EnvVersion)
	}
	if a.Pid == 0 {
		fmt.Println("warning: Pid is missing")
	}
	if a.Severity == 0 {
		fmt.Println("warning: Severity is missing")
	}
	if len(a.Fields) == 0 {
		fmt.Println("error: no field was found")
	}
}
