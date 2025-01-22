# Task Management Application

## Overview

This is a simple task management web application built using Go, the Gorilla Mux router, MySQL for the database, and HTML templates for the front end. The application uses HTMX to create a dynamic, interactive user interface for task management. Users can create, read, update, and delete (CRUD) tasks asynchronously using HTMX, allowing for smoother interactions without full-page reloads.

### Features:
- Display a list of tasks
- Add new tasks dynamically with HTMX
- Update tasks in real-time
- Delete tasks without page reload
- Use Bootstrap for UI styling
- Use HTML templates for rendering forms and task lists

## Installation

### Prerequisites:
- Go 1.16+ installed

### Steps:

1. **Clone the repository:**

    ```bash
    git clone https://github.com/jake-t-dev/go-htmx-crud-site.git
    cd go-htmx-crud-site
    ```

2. **Install dependencies:**

    ```bash
    go get -u github.com/gorilla/mux
    go get -u github.com/go-sql-driver/mysql
    ```

3. **Set up the database:**

    Create a MySQL database and table for storing tasks:
    ```sql
    CREATE DATABASE testdb;

    USE testdb;

    CREATE TABLE tasks (
        id INT AUTO_INCREMENT PRIMARY KEY,
        task VARCHAR(255) NOT NULL,
        done INT DEFAULT 0
    );
    ```

4. **Modify database connection:**

    In the `initDB` function, modify the MySQL connection string to match your MySQL credentials and configuration.

    ```go
    db, err = sql.Open("mysql", "username:password@(localhost:port)/dbname")
    ```

5. **Run the application:**

    ```bash
    go run main.go
    ```

6. **Access the app:**

    Visit [http://localhost:3000](http://localhost:3000) in your browser.

## Routes

- `GET /` - Home page
- `GET /tasks` - Get a list of all tasks (loaded dynamically with HTMX)
- `GET /getnewtaskform` - Display the form to add a new task
- `POST /tasks` - Add a new task (submitted dynamically with HTMX)
- `GET /gettaskupdateform/{id}` - Display the form to update a specific task
- `PUT /tasks/{id}` - Update a task
- `DELETE /tasks/{id}` - Delete a task (removes task from the UI without reloading)