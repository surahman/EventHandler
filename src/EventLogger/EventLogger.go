package EventLogger

import (
    "bufio"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    _ "github.com/gorilla/mux"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "time"
)

/*
 * Data structure for a single event.
 */
type Event struct {
    ServiceName string `json:"service_name"` // Service Name
    ServerID    string `json:"server_id"`    // Server's unique identifier
    Date        string `json:"date"`         // Date of event
    Time        string `json:"time"`         // Time of event
    Level       string `json:"level"`        // Level of event: Critical, Warn, Info,...
    EventType   string `json:"event_type"`   // Type of event is specific to service
    Description string `json:"description"`  // Details of the event
}

func (event *Event) String() string {
    return fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s\n",
        event.Date,
        event.Time,
        event.ServiceName,
        event.ServerID,
        event.Level,
        event.EventType,
        event.Description)
}

// baseLogDirectory will be used to locate the directory to read/write log files.
var baseLogDirectory string

// Format string: [TIME STAMP]-[LEVEL] MESSAGE
const baseLogStr = "[%v]-[%s] %s"

// Port number
const portNumber = 45456

// Server log file location
var serverLogFile string

// Logger
func logger(level, msg string) {
    logMessage := fmt.Sprintf(baseLogStr, time.Now(), level, msg)

    log.Print(logMessage)

    fd, err := os.OpenFile(serverLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        panic(err)
    }
    defer fd.Close()

    if _, err = fd.WriteString(logMessage + "\n"); err != nil {
        panic(err)
    }
}

// Initialize the Event Logger.
func init() {
    // Configure log directory and server log file.
    baseDirectory, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    baseLogDirectory = fmt.Sprintf("%s/%s/", baseDirectory, "logs")
    serverLogFile = filepath.Join(baseLogDirectory, filepath.Base("server_log.log"))
    if err = os.MkdirAll(baseLogDirectory, os.ModePerm); err != nil {
        panic(err)
    }
    logger("INFO", "Event Server Starting...")
    logger("INFO", "Log storage directory set to "+baseLogDirectory)
    logger("INFO", "Server logs located at: "+serverLogFile)

    // Start serving requests.
    handleRequests()

    logger("INFO", "Event Logger started")
}

// Gorilla Mux connection multiplexer.
func handleRequests() {
    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/", homePage).Methods("GET")
    router.HandleFunc("/logs", serverLog).Methods("GET")
    router.HandleFunc("/append", appendEvent).Methods("POST")
    router.HandleFunc("/logs/{service_name}/{server_id}/{date}", retrieveLog).Methods("GET")
    log.Fatal(http.ListenAndServe(":"+strconv.Itoa(portNumber), router))
}

// Default home page
func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Event Logger is Active")
    logger("INFO", "Homepage visitor")
}

// View Event Logger server logs.
func serverLog(w http.ResponseWriter, r *http.Request) {
    logger("WARN", "Served request for server logs")

    // Open server log file
    fd, err := os.Open(serverLogFile)
    if err != nil {
        error := fmt.Sprintf("No log files found for the Event Logger: %v", err)
        fmt.Fprintf(w, "ERROR:"+error)
        logger("SEVERE", error)
        return
    }
    defer fd.Close()

    // Write log file to responder.
    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        fmt.Fprintln(w, scanner.Text())
    }

    // Check for error from the scanner.
    if err := scanner.Err(); err != nil {
        error := fmt.Sprintf("No log files found for the Event Logger: %v", err)
        fmt.Fprintf(w, "ERROR:"+error)
        logger("SEVERE", error)
        return
    }
}

// Append an event to a log file.
func appendEvent(w http.ResponseWriter, r *http.Request) {
    // Unmarshall event into Event.
    reqBody, _ := ioutil.ReadAll(r.Body)
    var event Event
    json.Unmarshal(reqBody, &event)

    // TODO: Error check to make sure event fields are fully populated.

    // Generate path for log file
    path := fmt.Sprintf("%s/%s/%s/", baseLogDirectory, event.ServiceName, event.ServerID)
    filePath := fmt.Sprintf("%s%s.log", path, event.Date)

    // ECHO BACK
    //json.NewEncoder(w).Encode(event)

    // Create service log directory as required.
    if err := os.MkdirAll(path, os.ModePerm); err != nil {
        error := fmt.Sprintf("Unable to create directory for log: %v", err)
        fmt.Fprintf(w, "ERROR:"+error)
        logger("SEVERE", error)
        return
    }

    // Open log file and write entry.
    fd, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        error := fmt.Sprintf("Unable to entry to log: %v", err)
        fmt.Fprintf(w, "ERROR:"+error)
        logger("SEVERE", error)
        return
    }
    defer fd.Close()

    if _, err = fd.WriteString(event.String()); err != nil {
        panic(err)
    }

    logger("INFO", fmt.Sprintf("Appended new event to log: %s", filePath))
}

// Append an event to a log file.
func retrieveLog(w http.ResponseWriter, r *http.Request) {
    urlFields := mux.Vars(r)
    serviceName := urlFields["service_name"]
    serverID := urlFields["server_id"]
    date := urlFields["date"]

    // TODO: validate fields retrieved from URL.

    filePath := fmt.Sprintf("%s%s/%s/%s.log", baseLogDirectory, serviceName, serverID, date)

    // Attempt to write log file to requester
    if success := logFileWriter(&w, filePath); success {
        logger("INFO", fmt.Sprintf("Served event log: %s", filePath))
    }
}

// Loads and writes out a log file for an HTTP request.
func logFileWriter(w *http.ResponseWriter, logFileName string) bool {
    // Open server log file
    fd, err := os.Open(logFileName)
    if err != nil {
        err := fmt.Sprintf("No log files found: %v", err)
        fmt.Fprintf(*w, "ERROR:"+err)
        logger("SEVERE", err)
        return false
    }
    defer fd.Close()

    // Write log file to responder.
    scanner := bufio.NewScanner(fd)
    for scanner.Scan() {
        fmt.Fprintln(*w, scanner.Text())
    }

    // Check for error from the scanner.
    if err := scanner.Err(); err != nil {
        err := fmt.Sprintf("No log files found: %v", err)
        fmt.Fprintf(*w, "ERROR:"+err)
        logger("SEVERE", err)
        return false
    }

    return true
}
