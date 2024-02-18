# Task Tracker Project - Technical Requirement Document

## Project Overview

The Task Tracker project aims to develop a RESTful API for managing tasks, allowing users to create, update, and track tasks efficiently. The application will provide features such as task filtering and summary generation to enhance task management capabilities.

## Technical Specification

### Tech Stack

| Aspect              | Details                               |
|---------------------|---------------------------------------|
| Programming Language| Go (Golang)                           |
| RDBMS               | MySQL                                 |
| Framework           | [Mux](https://github.com/gorilla/mux) |
| API Protocol        | RESTful API                           |
| Docker              | Used for containerization             |
| Docker Compose      | Used for orchestrating containers     |

### Architecture Overview

The project follows a **layered architecture**, which is a common architectural pattern for structuring software applications. In a layered architecture, the application is divided into multiple layers, each responsible for a specific set of tasks. This separation of concerns helps in better organization, maintainability, and scalability of the application.

### Layers

1. **Presentation Layer**:
   - Located in the `handler` package.
   - Responsible for handling HTTP requests and responses.
   - Contains HTTP request handlers (controllers) that interact with the service layer.

2. **Service Layer**:
   - Located in the `service` package.
   - Implements the business logic of the application.
   - Orchestrates interactions between the presentation layer and the data access layer.
   - Contains service functions that encapsulate specific business rules and operations.

3. **Data Access Layer**:
   - Located in the `repository` package.
   - Responsible for interacting with the database or any external data source.
   - Contains repository functions for performing CRUD (Create, Read, Update, Delete) operations on data entities.
   - Abstracts away the details of data storage and retrieval from the rest of the application.

### Folder Structure

- **cmd**: Contains the main application entry point (`main.go`).
- **config**: Contains configuration settings for the application.
- **internal**: Contains the main codebase of the application, divided into various packages representing different layers and components.
  - **constant**: Contains constant values used throughout the application.
  - **handler**: Contains HTTP request handlers (presentation layer).
  - **middleware**: Contains middleware functions for request processing.
  - **model**: Contains data models/structures used in the application.
  - **repository**: Contains functions for interacting with the database (data access layer).
  - **router**: Contains router configuration for defining application routes.
  - **service**: Contains business logic and service functions (service layer).
  - **util**: Contains utility/helper functions used across the application.
- **test**: Contains unit tests for the application.
  - **unit**: Further divided into `repository` and `service` packages, containing unit tests for the repository and service layers respectively.

This layered architecture promotes separation of concerns, making the application more modular, maintainable, and easier to test. Each layer has its specific responsibilities, and interactions between layers are well-defined, contributing to a clear and structured codebase.

### Database Diagram

#### Table Spec

Table Name: tasks
| Column Name | Data Type                         | Constraints                                               |
|-------------|-----------------------------------|-----------------------------------------------------------|
| id          | int                               | Primary Key, Auto Increment                               |
| uuid        | char(36)                          | Not Null                                                  |
| title       | varchar(255)                      | Not Null                                                  |
| description | text                              |                                                           |
| due_date    | date                              |                                                           |
| status      | enum('in-progress', 'completed')  | Default: 'in-progress'                                    |
| version     | int                               | Not Null                                                  |
| created_at  | datetime                          | Default: Current Timestamp                                |
| updated_at  | datetime                          | Default: Current Timestamp, On Update: Current Timestamp  |

#### Data Definition Language (DDL)

```sql
CREATE TABLE IF NOT EXISTS tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    uuid CHAR(36) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    due_date DATE,
    status ENUM('in-progress', 'completed') DEFAULT 'in-progress',
    version INT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### API Documentation

| API Action     | Description                                              | HTTP Method | URL                                           | Headers                                     | Payload                                                                                     | CURL URL                                                              |
|----------------|----------------------------------------------------------|-------------|-----------------------------------------------|---------------------------------------------|---------------------------------------------------------------------------------------------|-----------------------------------------------------------------------|
| Create Task    | Creates a new task with the provided details            | POST        | http://localhost:8080/tasks/create           | x-api-token: HJZkQrCwgrBN23aTcnyo          | {"title": "test 1", "description": "test desc", "due_date": "2023-02-17"}                  | `curl --location 'http://localhost:8080/tasks/create' --header 'x-api-token: HJZkQrCwgrBN23aTcnyo' --header 'Content-Type: application/json' --data '{"title": "test 1", "description": "test desc", "due_date": "2023-02-17"}'` |
| Update Task    | Updates an existing task with the provided details       | PUT         | http://localhost:8080/tasks/{taskID}/update  | x-api-token: HJZkQrCwgrBN23aTcnyo          | {"title": "test 1", "description": "test desc", "due_date": "2023-02-17", "status": "completed"} | `curl --location --request PUT 'http://localhost:8080/tasks/{taskID}/update' --header 'x-api-token: HJZkQrCwgrBN23aTcnyo' --header 'Content-Type: application/json' --data '{"title": "test 1", "description": "test desc", "due_date": "2023-02-17", "status": "completed"}'` |
| Task List      | Retrieves a list of tasks with optional filters and pagination | GET         | http://localhost:8080/tasks                  | x-api-token: HJZkQrCwgrBN23aTcnyo          | query parameters: `search`, `status`, `page`, `limit`                                                                                             | `curl --location 'http://localhost:8080/tasks?search={search}&status={status}&page={page}&limit={limit}' --header 'x-api-token: HJZkQrCwgrBN23aTcnyo'` |
| Task Summary   | Retrieves a summary of tasks with optional filter by due date | GET         | http://localhost:8080/tasks/summary          | x-api-token: HJZkQrCwgrBN23aTcnyo          | query parameter: `due_date`                                                                                             | `curl --location 'http://localhost:8080/tasks/summary?due_date={due_date}' --header 'x-api-token: HJZkQrCwgrBN23aTcnyo'` |

### Run API with Postman

If you want to run with postman, you can import postman collection file from this project

To import the collection file into Postman:

1. Open Postman.
2. Click on "Import" in the top left corner.
3. Choose the collection file called ```Task Tracker.postman_collection.json``` inside this root folder project.
4. Click "Open" to import the collection into Postman.

### Running the Application

#### Clone Repository

### How to Run

Follow these steps to run the application:

1. **Clone the Repository**: Clone the project repository to your local machine using the following command:

   ```bash
   git clone https://github.com/handrixn/task-tracker.git
   ```

#### Using Docker Compose

To run the application using Docker Compose, follow these steps:

1. **Install Docker and Docker Compose**: Make sure you have Docker and Docker Compose installed on your system. If not, you can download and install them from the official Docker website.

2. **Navigate to Project Directory**: Open a terminal and navigate to the root directory of your project where the `docker-compose.yml` file is located.

3. **Build the Docker Images**: Run the following command to build the Docker images:

   ```bash
   docker-compose -f docker-compose.yml build
   ```

4. **Start the Application**: Run the following command to start the application:

   ```bash
   docker-compose -f docker-compose.yml up
   ```

   - If you want to run the application in detached mode, add the `-d` flag:

     ```bash
     docker-compose -f docker-compose.yml up -d
     ```

5. **Access the Application**: Once the application is running, you can access it in your web browser at `http://localhost:8080`.

6. **Stopping the Application**:
   - If you started the application without detached mode, press `Ctrl + C` in the terminal where Docker Compose is running to stop the application.
   - If you started the application with detached mode, run the following command to stop the application:

     ```bash
     docker-compose -f docker-compose.yml down
     ```

#### Using Makefile

If you prefer using the provided Makefile to manage the application, follow these steps:

1. **Install Docker and Docker Compose**: Make sure you have Docker and Docker Compose installed on your system. If not, you can download and install them from the official Docker website.

2. **Navigate to Project Directory**: Open a terminal and navigate to the root directory of your project where the `Makefile` is located.

3. **Build the Docker Images**: Run the following command to build the Docker images:

   ```bash
   make build
   ```

4. **Start the Application**: Run the following command to start the application without detached mode:

   ```bash
   make run
   ```

   - If you want to run the application in detached mode, run the following command:

     ```bash
     make run-d
     ```

5. **Access the Application**: Once the application is running, you can access it in your web browser at `http://localhost:8080`.

6. **Stopping the Application**:
   - If you started the application without detached mode, press `Ctrl + C` in the terminal where Docker Compose is running to stop the application.
   - If you started the application with detached mode, you can run the following command to stop the application and remove the containers:

     ```bash
     make clean
     ```
