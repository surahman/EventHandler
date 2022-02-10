# REST Event Logger POC

I am working on this project intermittently. I have set myself a time limit of ~3hrs which includes the time to learn and adapt to Golang. This is a very quick Proof of Concept for demonstration purposes only. It shall also serve as a refresher to Golang.

Gorilla MUX library is used for concurrent multiplexed HTTP request handling.


## TODO's and Shortcuts

* **Test suite**: Since this is not meant for production and is a very quick exercise I have not had the time to build a test suite. This is a poor practice and does not fly in production environments.
* **Code Coverage**
* **Benchmarking**
* **Profiling**
* **HTTP Codes**: Appropriate HTTP response codes should be sent for errors and success.
* **Write Chunking**: Log entries should be buffered and then written to disk to improve performance.


## Event Tuple

* `ServiceName`: Reporting service's name.
* `ServerID`: Server's unique identifier.
* `Date`: Date of event. `ddmmyyyy`
* `Time`: Time of event. `hhmmss`
* `Level`: Level of event: `Critical`, `Warn`, `Info`, etc.
* `EventType`: Type of event is specific to service.
* `Description`: Details of the event.

## Details

### Storage
Events are tail appended to files stored on disk in the file structure `service_name/server_id/date.log`.

Using a file on disk that is tail appended to will improve memory consumption but result in reduced performance due to writes to disk. Ideally, writes would be chunked in blocks and tail-appended to files on disk. Tail-appending in blocks also improves concurrency as we will not require implicitly locking of the file to append. This is a similar scheme to that which is used in Apache Kafka.

The file structure above is efficient for ETL/ELT jobs (Spark etc.) when moving the logs to a data warehouse (OLAP) for analysis. 


### Event Logger Server Logs
You may view the activity of the Event Logger by visiting `/logs`

##### Example
`http://localhost:45456/logs`


### Submission
Events are submitted to the server via the `body` of an `HTTP: put` with the details of the event structured in `JSON`. The url would be `server_address/add_event/`.

##### Example
```bash
curl -X POST http://localhost:45456/append \
   -H "Content-Type: application/json" \
   -d '{"service_name": "serviceA", "server_id": "server001", "date": "09022022", "time": "000102", "level": "INFO", "event_type": "Account Created", "description": "New user Bilbo Baggins"}'
```


### Retrieval
Events are retrieved via an `HTTP: get` on the url structure `server_address/service_name/server_id/date/`. This will return all events on the specific date for a specific service's server. In a production system the log would be read and written back to the client in chunks as they can get very large.

##### Example
`http://localhost:45456/logs/serviceA/server001/19092022`
