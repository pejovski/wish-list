package logger

import (
	"github.com/sirupsen/logrus"
	"io"
)

// Hook represents a Logstash hook.
// It has two fields: writer to write the entry to Logstash and
// formatter to format the entry to a Logstash format before sending.
//
// To initialize it use the `New` function.
//
type Hook struct {
	writer    io.Writer
	formatter logrus.Formatter
}

// New returns a new logrus.Hook for Logstash.
//
// To create a new hook that sends logs to `tcp://logstash.corp.io:9999`:
//
// conn, _ := net.Dial("tcp", "logstash.corp.io:9999")
// hook := logrustash.New(conn, logrustash.DefaultFormatter())
func New(w io.Writer, f logrus.Formatter) logrus.Hook {
	return Hook{
		writer:    w,
		formatter: f,
	}
}

// Fire takes, formats and sends the entry to Logstash.
// Hook's formatter is used to format the entry into Logstash format
// and Hook's writer is used to write the formatted entry to the Logstash instance.
func (h Hook) Fire(e *logrus.Entry) error {
	dataBytes, err := h.formatter.Format(e)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(dataBytes)
	return err
}

// Levels returns all logrus levels.
func (h Hook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// LogstashFormatter represents a Logstash format.
// It has logrus.Formatter which formats the entry and logrus.Fields which
// are added to the JSON message if not given in the entry data.
//
// Note: use the `DefaultFormatter` function to set a default Logstash formatter.
type LogstashFormatter struct {
	logrus.Formatter
	logrus.Fields
}

var (
	logstashFields   = logrus.Fields{"@version": "1", "type": "log"}
	logstashFieldMap = logrus.FieldMap{
		logrus.FieldKeyTime: "@timestamp",
		logrus.FieldKeyMsg:  "message",
	}
)

// DefaultFormatter returns a default Logstash formatter:
// A JSON format with "@version" set to "1" (unless set differently in `fields`,
// "type" to "log" (unless set differently in `fields`),
// "@timestamp" to the log time and "message" to the log message.
//
// Note: to set a different configuration use the `LogstashFormatter` structure.
func DefaultFormatter(fields logrus.Fields) logrus.Formatter {

	for k, v := range logstashFields {
		if _, ok := fields[k]; !ok {
			fields[k] = v
		}
	}

	return LogstashFormatter{
		Formatter: &logrus.JSONFormatter{FieldMap: logstashFieldMap},
		Fields:    fields,
	}
}

// Format formats an entry to a Logstash format according to the given Formatter and Fields.
//
// Note: the given entry is copied and not changed during the formatting process.
func (f LogstashFormatter) Format(e *logrus.Entry) ([]byte, error) {
	addFields(e, f.Fields)
	dataBytes, err := f.Formatter.Format(e)
	return dataBytes, err
}

// addFields adds all the fields in `fields` that are missing in the data.
func addFields(e *logrus.Entry, fields logrus.Fields) {
	for k, v := range fields {
		e.Data[k] = v
	}
}
