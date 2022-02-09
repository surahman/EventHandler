package EventLogger

import (
    "fmt"
    "github.com/gorilla/mux"
    _ "github.com/gorilla/mux"
    "log"
    "net/http"
    "os"
    "strconv"
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

// Port number
const portNumber = 45456

func init() {
    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Event Logger starting..."))

    baseDirectory, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    baseLogDirectory = fmt.Sprintf("%s/%s/", baseDirectory, "logs")
    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Log storage directory set to "+baseLogDirectory))

    // Start serving requests.
    handleRequests()

    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Event Logger started"))
}

func handleRequests() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homePage)
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), router))
}

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Event Logger is Active")
    log.Print(fmt.Sprintf(baseLogStr, time.Now(), "INFO", "Visit to homepage"))
}
