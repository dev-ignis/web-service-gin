# Go Gin Application

This is a simple Go application that provides a RESTful API to manage a list of music albums. The application is built using the Gin web framework.

## Features
- List all albums
- Retrieve an album by ID
- Add a new album

## Prerequisites
- Go (for running the app outside of a container)
- Docker (for running the app inside a container)

## Running the Application

### Running the Application Outside of a Container

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-gin-app.git
    cd go-gin-app
    ```

2. **Build the application:**

    ```bash
    go build -o main .
    ```

3. **Run the application:**

    ```bash
    ./main
    ```

4. **Access the API:**

   Open your browser or use `curl` to interact with the API:

    ```bash
    curl http://localhost:8080/albums
    ```

### Running the Application Inside a Docker Container

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-gin-app.git
    cd go-gin-app
    ```

2. **Build the Docker image:**

    ```bash
    docker build -t go-gin-app .
    ```

3. **Run the Docker container:**

    ```bash
    docker run -p 8080:8080 go-gin-app
    ```

4. **Access the API:**

   Open your browser or use `curl` to interact with the API:

    ```bash
    curl http://localhost:8080/albums
    ```

### Accessing the API

Once the application is running (either inside or outside of Docker), you can access the following endpoints:

- **List all albums:**

    ```bash
    curl http://localhost:8080/albums
    ```

- **Retrieve an album by ID:**

    ```bash
    curl http://localhost:8080/albums/{id}
    ```

- **Add a new album:**

    ```bash
    curl -X POST http://localhost:8080/albums -d '{"id":"4","title":"New Album","artist":"New Artist","price":29.99}' -H "Content-Type: application/json"
    ```

### Notes

- Make sure that the port `8080` is available on your host machine. If the port is in use, you can modify the port in the commands by replacing `8080` with an available port number.
- If running in Docker, ensure that Docker is correctly configured and running on your system.
