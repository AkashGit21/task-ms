# task-ms
A simple Task Management System using Go. The system allow users to create, read, update, and delete tasks. The system should be designed as a microservice with clear separation of concerns and should demonstrate your understanding of microservices architecture. 

Tasks Listing should include:

- **Pagination**: Implement pagination for the GET /tasks endpoint.
- **Filtering**: Allow filtering tasks by status (e.g., GET /tasks?status=Completed).
---
### Problem breakdown and design decisions.


---
### Pre-requisites and dependencies
1. Docker setup 

    OR

1. Golang v1.22.1 or above
2. MySQL database (v8.0 or above)
3. make (used to run Makefile)


### Steps to run the backend server via Docker

#### Start both task-service and auth-service:
```sh
docker-compose up --build -d
```

#### Stop the server(s) completely
```sh
docker-compose down
```

**NOTE**: Docker compose is running `mysql` internally on port 3306, and may require local Mysql instance to stop.
```sh
sudo systemctl stop mysql
```


### Steps to run the backend server without Docker
1. Navigate to the project directory in terminal and fetch all the dependencies using following command:
    ```sh
    go mod tidy
    ```
1. Add the required table(s) to your RDBMS system by using commands from `init.sql`.
    
    **Note**: This is a one-time step only. Taking this step again will clear all previous tasks data from local setup.
1. Create and update the .env file in the project directory by using `.env.example` file. Fill in all the necessary values to make connection with the pre-requisites defined above.
1. Run the task service using terminal by typing: 
    ```sh
    $ make run-task
    ```
1. Open a new terminal in parallel. Navigate to the project directory in new terminal as well. Run the user service:
    ```sh
    $ make run-authn
    ```
---

### API Documentation
- [X] **POST**    `/api/v1/tasks` Allows user to create a new Task

    **Sample Request Body**
    ```json
    {
        "title":" Task A", // [Optional]
        "content": "Test task",
        "stylized_content": "<div>Test <br/> task</div>" // [Optional]
        // Allowed values - "todo", "in-progress", "backlog", "on-hold", "completed"
        "status": "in-progress" // [Optional] 
    }
    ```

    **Sample Response Body**
    ```json
    {
        "id": "4ea8e62a-162a-4b00-8064-1ae0e21d7b48",
        "title": " Task A",
        "content": "Test task",
        "stylized_content": "",
        "status": "in-progress",
        "created_at": "2025-03-23T16:30:46.575179643Z",
        "created_by": "SYSTEM",
        "modified_at": "2025-03-23T16:30:46.575179643Z",
        "modified_by": "SYSTEM"
    }
    ```

- [X] **GET**     `/api/v1/tasks/{taskID}` Retrieve a specific task based on a unique taskID.

    **Sample Response Body**
    ```json
    {
        "id": "4ea8e62a-162a-4b00-8064-1ae0e21d7b48",
        "title": " Task A",
        "content": "Test task",
        "stylized_content": "",
        "status": "in-progress",
        "created_at": "2025-03-23T16:30:46.575179643Z",
        "created_by": "SYSTEM",
        "modified_at": "2025-03-23T16:30:46.575179643Z",
        "modified_by": "SYSTEM"
    }
    ```

- [X] **PATCH**     `/api/v1/tasks/{taskID}` Update an existing task with limited set of fields.

    Any other except the following listed below are ignored:
    - title
    - content
    - stylized_content
    - status

    **Sample Request Body**
    ```json
    {
        "status": "on-hold", // [Optional]
        "title": "Task Patched", // [Optional]
        "content": "Patching the content of Task", // [Optional]
        "stylized_content": "Patching the stylized content of Task" // [Optional]
    }
    ```

    **Sample Response Body**
    ```json
    {
        "id": "4ea8e62a-162a-4b00-8064-1ae0e21d7b48",
        "title": " Task A",
        "content": "Test task",
        "stylized_content": "",
        "status": "on-hold",
        "created_at": "2025-03-23T16:30:47Z",
        "created_by": "SYSTEM",
        "modified_at": "2025-03-23T16:30:47Z",
        "modified_by": "SYSTEM"
    }
    ```

- [X] **DELETE**  `/api/v1/tasks/{taskID}` Delete a specific task based on a unique taskID. 

    **Note**: Applied a soft delete for Task, instead of hard delete to allow data recovery from System.


- [X] **GET**     `/api/v1/tasks` List all available tasks and their details.
    - *Pagination*: Using cursor-based pagination to avoid running into similar results in subsequent pages on write-heavy systems  
    - *Filtering*: Applied filtering by status value.
    NOTE: status value is internally stored as int rather than string to reduce memory

    **Available Query Params**
    - *status*: Allows filtering tasks by status value

        // Allowed values are - "todo", "in-progress", "backlog", "on-hold", "completed"
    - *cursor*: A cursor built using last Task's modifiedAt and ID
    - *limit*: Maximum number of tasks to be shown

    **Sample Response Body**
    ```json
    {
        "tasks": [
            {
                "ID": "b9f4725d-9f0c-4712-9209-52e1132a7646",
                "title": " Task A",
                "content": "Test task",
                "stylized_content": "",
                "status": 1,
                "created_at": "2025-03-23T16:30:46Z",
                "created_by": "SYSTEM",
                "modified_at": "2025-03-23T16:30:46Z",
                "modified_by": "SYSTEM"
            },
            {
                "ID": "9146da0d-21cb-4d6e-bae0-e75f023ca47b",
                "title": " Task A",
                "content": "Test task",
                "stylized_content": "",
                "status": 1,
                "created_at": "2025-03-23T16:30:43Z",
                "created_by": "SYSTEM",
                "modified_at": "2025-03-23T16:30:43Z",
                "modified_by": "SYSTEM"
            }
        ],
        "next_cursor": "eyJtb2RpZmllZF9hdCI6IjIwMjUtMDMtMjNUMTY6MzA6NDNaIiwiaWQiOiI5MTQ2ZGEwZC0yMWNiLTRkNmUtYmFlMC1lNzVmMDIzY2E0N2IifQ=="
    }
    ```

- [ ] **POST** `/api/login` Verify username / password and generate token for the user

---

Explanation of how the service demonstrates microservices concepts.

---
### Additional Points considered for Application Server
1. **Graceful shutdown**: This avoids any side effects on conflicts that may occur on closing the server and the new deployment can be started without any kind of difficulty.
1. **Logging**: For debugging and monitoring the application on remote servers, it is recommended to log the application functionality.
1. **Panic Handler**: Prevents the application from crashing, in case of any runtime errors or application malfunctioning.

### Improvements that can be done
1. Unit tests
1. Caching the API response for short periods, to save DB queries, reduce server stress and latency
