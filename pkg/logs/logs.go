package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

type LogEntry struct {
	Time      time.Time   `json:"time"`
	Level     string      `json:"level"`
	Msg       string      `json:"msg"`
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Status    int         `json:"status"`
	Latency   Duration    `json:"latency"`
	ExtraData interface{} `json:"extraData"`
}

func (l *LogEntry) CreateLog() {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		fmt.Printf("failed to create logs dir: %v\n", err)
		return
	}

	filename := filepath.Join(logDir, time.Now().Format("2006-01-02")+".log")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	l.Time = time.Now().UTC()

	data, err := json.Marshal(l)
	if err != nil {
		fmt.Printf("failed to marshal log entry: %v\n", err)
		return
	}

	if _, err := file.Write(append(data, '\n')); err != nil {
		fmt.Printf("failed to write log entry: %v\n", err)
	}
}
