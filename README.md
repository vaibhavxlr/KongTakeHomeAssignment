# Konnect API

## Overview

Konnect API is a test project designed to provide a simple interface for listing services and retrieving service details. The project includes two ReST APIs, both of which interact with a MongoDB database. This project avoids using external packages as much as possible.

## Endpoints

### 1. List Services

**Endpoint:** `GET localhost:8080/konnectAPI/v1/listServices`

**Query Parameters:**
- `curr` (optional): Current page number. Default is `1`.
- `count` (optional): Number of services per page. Default is `5`.
- `sortOrder` (optional): Sort order for the services. Default is `0` (A-Z).
  - `0` for A-Z (name of service)
  - `1` for Z-A (name of service)
- `search` (optional): Search query string. If non-empty, a fuzzy search will be performed.

**Example Request:**
```http
GET localhost:8080/konnectAPI/v1/listServices/?curr=1&count=5&sortOrder=1&search=konnect
```
**Example Response:**
```
{
  "services": [
    {
      "id": "123i2i",
      "name": "Locate Us",
      "info": "lorem ipsum lorem ipsum",
      "versionsCount": 3
    }
  ],
  "sortOrder": {
    "A-Z": 0,
    "Z-A": 1
  },
  "pageDetails": {
    "curr": 1,
    "total": 13,
    "count": 5
  }
}
```
### 2. Service Details and Versions

**Endpoint:** `GET localhost:8080/konnectAPI/v1/service/<id>`

**Path Parameters:**
- `id` : Unique identifier of the service.

**Example Request:**
```http
GET localhost:8080/konnectAPI/v1/service/<id>
```
**Example Response:**
```
{
  "id": "12324",
  "name": "konnect",
  "info": "loremipsum",
  "versions": [
    {
      "verName": "version 1",
      "verInfo": "lorem ipsum",
      "changes": "supports x, y, z"
    }
  ]
}

```
## Installation and Setup
To run the Konnect API project, you can either pull the pre-built Docker image from Docker Hub or build the Docker image locally using the provided Dockerfile.

### Using Docker Hub
1. **Pull the Docker image from Docker Hub:**
   ```sh
   docker pull vaibhavxlr/konnect-img
   ```
2. **Run the container:**
   ```sh
   docker run -d -p 8080:8080 vaibhavxlr/konnect-img
   ```

### Using Docker file
1. **Clone the repository and switch to the project dir**
2. **Build the image using the Dockerfile**
   ```sh
   docker build -t konnect-img .
   ```
3. **Run the container**
  ```sh
  docker run -d -p 8080:8080 konnect-img
  ```
# Notes
1. The DB schema
![image](https://github.com/vaibhavxlr/KongTakeHomeAssignment/assets/36249617/dc2191d0-64d9-42b6-bca8-ad7bdbeea2ec)

2. Multi stage docker build to combine application binary and DB within a same container. Check 'Dockerfile' and 'startup.sh'. Not a recommended practice, but a trade off to ship a single image with everything included, making it easily testable by evaluators.

3. Data gets imported while building the image and  so APIs can be tested without any extra steps.

4. As part of additional considerations, wrote unit test for one of the handler functions. Check 'handler_test.go' file.

5. TODO - Planned implementing user authentication using JWT.

6. TODO - CRUD functionality could be added utilizing the existing DB schema and connection, by adding extra handling at the application level.

7. Repository branching strategy could have been improved with a separate 'develop', feature branches and main. 
