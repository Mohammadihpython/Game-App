# Game App
## Table of Contents
- [About](#about)
- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
- [Technologies Used](#technologies-used)
- [Contributing](#contributing)
- [License](#license)

## About
This project built with go and echo framework

## Features

### 1. Authentication & Authorization System,

- Implemented Authentication system with JWT.
- Implement a Role Base Authorization
## Scheduled Task With Redis
- Match users from Waiting list and send to Event
## Message Driven system for the Matching Service and Play service
## Real Time communication with Websocket
## GRPC System for presence service and Matching Service
- Use protobuf protocol
## UintTest & Integration Test
### Implement Integration Test
- use dokcertest package to run services
- Run in memory GRPC Server & Client for Tests
- Run an isolate mysql service and migrate data and ready for run test



## Installation

### Prerequisites

- Docker installed on your system

1. Clone the repository:
```
https://github.com/Mohammadihpython/Game-App.git
```

2. Set up environment variables:
configure your environments in .env  and yaml file in config folder like:
- SECRET_KEY
- DATABASE_DB
- .......

3. Start Docker containers:
```
run: docker compose  up --build
go run main.go
```




5. Access the application:
Visit http://localhost:8080/ in your browser to access the application

## Usage

1. **Login & Register**:
- Login or register user with phone number and get access token



2. **Asynchronous Tasks**:
- Experience improved performance due to the asynchronous handling of tasks

3. **CI with GitHub Actions**:
- Utilize the predefined GitHub Actions workflows for automated testing, linting, and deployment.
- View CI status and check build, test, and deployment logs directly on GitHub.

4. **Customization**:
- Explore and modify the codebase to customize the platform according to your specific requirements.



## Technologies Used

- Go
- Echo framework
- Docker
- GRPC
- GitHub Actions
- Message Queue
- Other relevant technologies and libraries used in the project
- dockerTest
- Mysql




## Contributing
Explain how others can contribute to your project. Include guidelines for pull requests and issue reporting.

## License
This project is licensed under the [License Name] License - see the [LICENSE](LICENSE) file for details.
