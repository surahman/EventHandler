# REST Event Logger POC

I am working on this project intermittently. I have set myself a time limit of ~3hrs which includes the time to learn and adapt to Golang.

Gorilla MUX library is used for concurrent multiplexed server.

## Event Tuple

* `ServiceName`: Reporting service's name.
* `ServerID`: Server's unique identifier.
* `Date`: Date of event.
* `Time`: Time of event.
* `Level`: Level of event: `Critical`, `Warn`, `Info`, etc.
* `EventType`: Type of event is specific to service.
* `Description`: Details of the event.

## Procedure

### Storage
Events are tail appended to files stored on disk in the file structure `service_name/server_id/date/time_stamp.log`.

Using a file on disk that is tail appended to will improve memory consumption but result in reduced performance due to writes to disk. Ideally, writes would be chunked in blocks and tail-appended to files on disk. Tail-appending in blocks also improves concurrency as we will not require implicitly locking of the file to append. 

The file structure above is efficient for ETL/ELT jobs (Spark etc.) when moving the logs to a data warehouse (OLAP) for analysis. 


### Submission
Events are submitted to the server via the `body` of an `HTTP: put` with the details of the event structured in `JSON`. The url would be `server_address/add_event/`.


### Retrieval
Events are retrieved via an `HTTP: get` on the url structure `server_address/service_name/server_id/date/`. This will return all events on the specific date for a specific service's server. In a production system the log would be read and written back to the client in chunks as they can get very large.
