# Toronto Time API with MySQL Database

This is a simple Go-based API that provides the current time in Toronto (converted from UTC), logs the time into a MySQL database.
## Features

- **Current Time in Toronto**: The `/current-time` endpoint returns the current time in Toronto.
- **Logging**: Every time the `/current-time` endpoint is hit, the current time is logged into a MySQL database.
- **Log Retrieval**: The `/logs` endpoint retrieves all logged times stored in the database.
- **Dockerization**: The Go application and MySQL database are Dockerized for easy deployment and consistent development environments.

---



## Project Setup

### Prerequisites

- **Go**: Version 1.20 or higher installed on your machine.
- **Docker**: Docker and Docker Compose must be installed to run the application in containers.
- **MySQL**: The project uses MySQL to store the logged times.



**The timestamp is displayed in JSON format in the browser and can also be seen in the database after querying.**

