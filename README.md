# Ecommerce Golang Typescript Application

This project showcases an ecommerce website built using Typescript and Golang.

## Table of Contents

- [Features](#features)
- [Technologies Used](#technologies-used)
- [Installation](#installation)

## Features

- User authentication with JWT
- Product management
- Shopping cart functionality
- Order processing
- Scalable architecture using Docker and Kubernetes

## Technologies Used

This project employs a variety of technologies to ensure performance and scalability:

- **Frontend**: Typescript
- **Backend**: Golang with Fiber
- **Database**: PostgreSQL/Amazon RDS
- **Caching**: -
- **Storage**: -
- **Web Server**: -
- **Monitoring**: -
- **Containerization**: Docker

## Installation

To get started with the ecommerce app project, follow these steps:

1. Clone the repository:

   ```sh
   git clone https://github.com/Alero-Awani/ecommerce-app .
   ```

1. Set up the backend and install the necessary dependencies:

     ```sh
     go mod tidy
     ```

1. Set up the frontend:

   - Navigate to the frontend folder:

     ```sh
     cd frontend
     ```

   - Install the necessary dependencies:

     ```sh
     yarn dev
     ```

1. Create a `.env` file in both backend and frontend folders to configure your environment variables. You can refer to the `.env.example` files for guidance.

1. Start the backend server:

    ```sh
    go run main.go
    ```

1. Start the frontend server:

    ```bash
    npm start
    ```

Now you can access the application at `http://localhost:3000`.

## Start the Application in a Docker Container

### Dockerfile

**Reference:** [Docker Up and Running](https://www.oreilly.com/library/view/docker-up/9781098131814/)

- **Tip 1**: It is not recommended that you run commands like `app=get -y update` or `yum -y update` in you application's dockerfile. This is because it requires crawling the repository index each time you run a build, and means that your build is not guaranteed to be repeatable since package versions might change between builds. Instead, consider basing your application image on another image that already has these updates applied to it and where the versions are in a known state. It will be faster and more repeatable.

- **Tip 2**: By default, Docker runs all process as root within the container. Even though containers provide some isolation from the underlying operating system, they still run on the host kernel. Due to potential risks, production containers should almost always be run under the context of a nonprivileged user.

### Build and Run Golang Dockerfile

Build and Tag the Backend Docker Image

```sh
docker build -t golang-backend:0.1.0 .
```

Start the Docker Container

```sh
docker run -p 9000:9000 golang-backend:0.1.0 
```