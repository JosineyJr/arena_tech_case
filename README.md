# IP Location API - Go Implementation

This project is a Go implementation of the Backend Developer Challenge. It provides a high-performance REST API to find the geographic location of an IPv4 address.

## Architectural Decisions

The primary goal was to meet the performance requirement of **<100ms response time** and handle **100+ concurrent users**, while adhering to the constraint of using **no third-party libraries** where possible.

### 1. In-Memory Data Storage

- **Problem**: The dataset is a 330Mb CSV file. Reading from and searching this file on disk for every API request would be incredibly slow and would never meet the performance target.
- **Solution**: The entire CSV file is parsed **once** at server startup and loaded into memory (RAM). RAM access is thousands of times faster than disk access. This creates an in-memory database that all incoming requests can query instantly.

### 2. Binary Search for High-Performance Lookups

- **Problem**: Even in memory, searching through 3 million records one-by-one (a linear search) for every request is too slow. In the worst case, it would have to check all 2,979,950 rows.
- **Solution**: The data is stored in a slice, sorted by the lower IP ID. This allows us to use a **binary search** (`sort.Search` in Go). A binary search is logarithmically fast (`O(log n)`).


### 3. Concurrency with Go's Standard HTTP Server

- **Problem**: The API must handle many users at once.
- **Solution**: We use Go's built-in `net/http` server. By default, it automatically handles each incoming request in a separate, lightweight thread called a **goroutine**.

### 4. Go Standard Library

- To satisfy the "no third-party libraries" preference, this implementation relies exclusively on Go's robust standard library for the HTTP server (`net/http`), CSV parsing (`encoding/csv`), string manipulation (`strings`), and searching (`sort`).

## How to Run

1.  **Prerequisites**:

    - Go (version 1.18 or higher) installed.

2.  **Run the Server**:

    ```bash
    go run main.go
    ```

    The server will start and listen on `http://localhost:3000`.

3.  **Run Tests**:
    ```bash
    go test
    ```
