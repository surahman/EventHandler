package EventLogger

import (
    "fmt"
    "log"
    "os"
    "time"
)

/*
 * Data structure for a single event.
 */
type eventLogger struct {
    ServiceName string `json:"service_name"` // Service Name
    ServerID    string `json:"server_id"`    // Server's unique identifier
    Date        string `json:"date"`         // Date of event
    Time        string `json:"time"`         // Time of event
    Level       string `json:"level"`        // Level of event: Critical, Warn, Info,...
    EventType   string `json:"event_type"`   // Type of event is specific to service
    Description string `json:"description"`  // Details of the event
}

// baseLogDirectory will be used to locate the directory to read/write log files.
var baseLogDirectory string

// Format string: [TIME STAMP]-[LEVEL] MESSAGE
const baseLogStr = "[%v]-[%s] %s"

func init() {
    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Event Logger starting..."))

    baseDirectory, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    baseLogDirectory = fmt.Sprintf("%s/%s/", baseDirectory, "logs")
    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Log storage directory set to "+baseLogDirectory))

    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Event Logger started"))
}
