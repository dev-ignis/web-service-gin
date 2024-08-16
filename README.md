# Go Gin Application with AWS S3 Integration

This is a simple Go application that provides a RESTful API to manage a list of music albums. The application stores album data in AWS S3 as JSON files.

## Features
- List all albums
- Retrieve an album by ID
- Add a new album

## Prerequisites

- **Go**: For running the app outside of a container.
- **Docker**: For running the app inside a container.
- **AWS CLI**: To configure AWS credentials.

## AWS S3 Configuration

Make sure you have an S3 bucket created where the album data will be stored. Update the `bucketName` variable in the code with your S3 bucket name.

### AWS Credentials

Ensure your AWS credentials are set up on your local machine or environment where the application runs:

```bash
aws configure
```

# Running the Application

## Running the Application Outside of a Container

1. **Clone the repository:**

    ```bash
    git clone https://github.com/your-username/go-gin-app.git
    cd go-gin-app
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Build the application:**

    ```bash
    go build -o main .
    ```

4. **Run the application:**

    ```bash
    ./main
    ```

5. **Access the API:**

   Open your browser or use `curl` to interact with the API:

    ```bash
    curl http://localhost:8080/albums
    ```

## Running the Application Inside a Docker Container

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

## Accessing the API

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

- **AWS S3 Integration**: The application uses AWS S3 to store album data as JSON files. Ensure that your S3 bucket is properly configured, and your AWS credentials are correctly set up.
- **Environment Variables**: You may need to configure environment variables or pass AWS credentials to the Docker container, depending on your setup.

### Troubleshooting

- **Connection Issues**: If you encounter issues when connecting to AWS S3, ensure that your credentials are correct and that you have network connectivity to AWS.
- **Docker Networking**: If running in Docker, ensure that the Docker container has access to the network and that your AWS credentials are accessible within the container.

### License

This project is licensed under the MIT License.

