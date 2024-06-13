
# Fiber Options Analysis

This project analyzes options contracts using the Fiber web framework in Go. The analysis includes calculating profit/loss values, maximum profit, maximum loss, and break-even points for a given set of options contracts.

## Prerequisites

- Go (version 1.16 or higher)
- Fiber framework
- Go modules

## Installation

1. Install Go from the official [website](https://golang.org/dl/).

2. Set up your Go workspace and ensure your `GOPATH` and `GOBIN` environment variables are correctly configured.

3. Install the Fiber framework:
    ```sh
    go get -u github.com/gofiber/fiber/v2
    ```

4. Install the testify package for testing:
    ```sh
    go get github.com/stretchr/testify
    ```

5. Clone the repository or copy the project files into your workspace.

## Project Structure

```
.
├── controllers
│   └── analysisContoller.go
├── models
│   └── optionsContract.go
├── routes
│   └── routes.go
├── testdata
│   └── testdata.json
├── tests
│   └── analysis_test.go
└── main.go
```

- `controllers/analysisController.go`: Contains the logic for analyzing options contracts.
- `models/optionsContract.go`: Defines the data structure for an options contract.
- `routes/routes.go`: Sets up the routes for the application.
- `testdata/testdata.json`: Sample test data for options contracts.
- `tests/analysis_test.go`: Contains test cases for the application.
- `main.go`: The entry point of the application.

## Running the Application

1. Navigate to the project directory.

2. Run the application:
    ```sh
    go run main.go
    ```

3. The server will start at `http://localhost:8080`. You can test the `/analyze` endpoint using the sample data in `testdata/testdata.json`.

## Testing the Application

1. Navigate to the project directory.

2. Run the tests:
    ```sh
    go test ./tests
    ```

3. The tests will validate the functionality of the model and the `/analyze` endpoint.

## Testing the Endpoint

You can use `curl` or any API testing tool like Postman to send a POST request to the `/analyze` endpoint with the sample data.

Example `curl` command:

```sh
curl -X POST http://localhost:8080/analyze      -H "Content-Type: application/json"      -d @testdata/testdata.json
```
