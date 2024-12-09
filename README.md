# mytheresa

## Overview
This project is a REST API implemented in Golang, designed to apply discounts to a list of products and return the resulting data with optional filtering. It offers a robust and performant solution that adheres to scalable design principles, capable of handling a large dataset of over 20,000 products.

The API is built using:

* **Golang** for high performance and simplicity.
* **MongoDB** as the database for storing and querying products and discounts.
* **Docker** to ensure ease of setup and portability across environments.

## 1- How to run
### Prerequisites

To run the API and the frontend, you need to have Docker and Docker Compose installed on the machine.
- Docker (Min version: 20.10.12): [Get Docker](https://docs.docker.com/get-docker/)
- Docker Compose (Min version: 2.27.1): [Get Docker Compose](https://docs.docker.com/compose/install/)

### Build and run

In the main folder of the project, where the `docker-compose.yml` file is located, execute:

```
docker-compose up --build
```

This command builds the Dockerfiles and pulls any necessary images from DockerHub. Once docker-compose has finished building and running the images, the services will be up and running.

### Swagger Documentation

Swagger allows you to visualize and interact with the API in a user-friendly interface, providing detailed information about the various endpoints, including the requests and responses.

To access the API documentation, open your browser and navigate to the following URL:

```
http://localhost:8001
```
swagger.json -> https://github.com/carlos2380/mytheresa/blob/main/docs/swagger.json

![swagger_app](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/swagger.png)

## 2- Performace
### Client for Load Testing

A custom client has been created to test the performance of the API server. This client simulates concurrent requests to evaluate how well the application handles load.

The client is located -> https://github.com/carlos2380/mytheresa/blob/main/cmd/client/main.go

Running the Client

Build the Client After running the Docker Compose setup for the API, build the client image using the following command:

```
docker build -t client --target client .
```
Run the Client Once the client image is built:

```
docker run -it --network=mytheresa_default client sh -c "/client -c 8 -nc 5000 -url http://mytheresa:8000/api/products"
```

* **-c 8** Specifies the number of threads (set to 8 to match the 8 CPUs available).

* **-nc 50000** Specifies the number of transactions per thread (total requests = 8 * 50000 = 400,000).

* **-url http://localhost:8000/api/products** The URL to test (the API endpoint for this application).

#### Hardware Assumptions
The API server is tested on a docker shared with other dockers and programs running in the machine with 8 CPUs
and 8 GB of RAM.

#### Results
![client_test](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/client.png)

Executing concurrency 1 and 40000 transactions per thread we have a TPS (Transactions Per Second) of 1100.

Executing concurrency 4 and 10000 transactions per thread we have a TPS of 3159.

### AB Apache Benchmark Testing

Using AB (Apache Benchmark), it is straightforward to evaluate the server's performance by specifying the number of concurrent threads (-c) and the total number of requests (-n).

```
ab -k -c 1 -n 100000 http://localhost:8000/api/products
```

#### Results

![client_test](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/ab1.png)
![client_test](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/cpu1.png)

With a concurrency level of 1, the server achieved 951.91 transactions per second (TPS), utilizing approximately 3.5% CPU with an average request time of 1.051 ms.

![client_test](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/ab4.png)
![client_test](https://github.com/carlos2380/webCarlos2380/blob/master/mytheresa/cpu4.png)

With a concurrency level of 4, the server achieved 2856.62 transactions per second (TPS), utilizing approximately 8% CPU with an average request time of 0.35 ms.

### Conclusion

The results show that increasing concurrency improves transactions per second, highlighting Go's efficient concurrency handling. Since all tests were conducted on a local machine, future testing on a dedicated environment would provide more accurate insights into performance and scalability.

## 3- Documentation
### Domain-Driven Design (DDD) and Hexagonal Architecture

This project applies Domain-Driven Design (DDD) and Hexagonal Architecture principles to create a modular, scalable, and maintainable system. Here's how these concepts are implemented:

#### **Domain-Driven Design (DDD)**

##### Domain-Centric Structure
* Core business logic, such as applying discounts and filtering products, is encapsulated within the ProductApplication.
- https://github.com/carlos2380/mytheresa/blob/main/internal/application/product/product.go

* Models like Product and Price represent model concepts, reflecting business rules such as discount prioritization.
- https://github.com/carlos2380/mytheresa/blob/main/internal/models/discount.go

##### Abstraction Through Interfaces
* Interfaces like ProductStorage and DiscountStorage abstract infrastructure details, ensuring the domain logic remains decoupled from external systems.
-https://github.com/carlos2380/mytheresa/blob/main/internal/storage/product.go

#### **Hexagonal Architecture**

* **Separation of Concerns**: Core domain logic (ProductApplication) is isolated from external systems like MongoDB and the HTTP API, ensuring a clean boundary between layers.

* **Port-and-Adapter Model**: Interfaces (ProductStorage, DiscountStorage) act as ports for external interactions, while MongoDB implementations serve as adapters, enabling seamless integration without coupling.

* **Independent Business Logic**: Core logic is infrastructure-agnostic, allowing changes to databases or frameworks without affecting the domain.

* **Dependency Inversion**: The core logic relies on interfaces instead of concrete implementations, with dependencies injected in main.go to ensure flexibility and ease of testing.

* **High Testability**: Mock implementations enable independent testing of the core domain, avoiding reliance on network or database connections.

#### Key Benefits

* **Scalability**: The clear separation of concerns supports easy addition of new features or integrations.
* **Flexibility**: External systems (e.g., databases) can be replaced without altering core logic.
* **Testability**: Core functionality can be thoroughly tested without requiring live infrastructur

### MongoDB Integration

The project leverages MongoDB, a NoSQL, document-based database, to efficiently manage and query product and discount data. Its flexible schema and document structure align seamlessly with the application's hierarchical and semi-structured data requirements.

#### Database Design

##### **products.json**
* Stores product data with fields such as SKU, name, category, price, and optional SKU-specific discounts.
* Avoids redundancy by storing discounts tied to specific SKUs (if have it) within the same document.

Example product:

```
{
  "sku": "000003",
  "name": "BV Lean leather ankle boots",
  "category": "boots",
  "price": 89000,
  "sku_discount": 15
}
```
##### **discounts.json**
* Holds category-level discount information.

Example product:

```
{
  "category": "boots",
  "percent": 30
}
```

#### Automation with init.js

The init.js script initializes the MongoDB database during the docker-compose build process by importing products.json and discounts.json.

Handles MongoDB's 16MB batch size limit by using multiple smaller inserts  when necessary (batches of 10000 elements), ensuring all data is properly loaded.

#### Indices for Performance
Indices have been added to the MongoDB collections to optimize filtering and retrieval:

* **Category Filtering**: An index on the category field speeds up filtering products by category.
* **Price Filtering**: Indexing on price allows fast querying when filtering by priceLessThan.

```
db.Product.createIndex({ category: 1 });
db.Product.createIndex({ price: 1 });
db.Product.createIndex({ category: 1, price: 1 });
```

#### Conclusion
MongoDB offers flexibility, performance, and scalability. Its schema-less design adapts to changes without migrations, while indices on category and price ensure fast queries. Its document model simplifies handling hierarchical product and discount data, and horizontal scaling supports growing data and traffic efficiently.

### Testing

To ensure isolated and efficient testing, the project uses mock implementations of storage interfaces (ProductStorage and DiscountStorage). This approach removes dependencies on external services or localhost, enabling fast and reliable validation of core business logic and simulate the data_storage.

mock_product -> https://github.com/carlos2380/mytheresa/blob/main/internal/storage/mock/mock_product.go
mock_data -> https://github.com/carlos2380/mytheresa/blob/main/internal/storage/mock/mock_data.go

Tests are structured using table-driven design, making it easy to add and maintain scenarios. This method ensures coverage for filtering by category and priceLessThan, discount application, pagination, and error handling (e.g., invalid inputs).

Handler_test with mocks and table test -> https://github.com/carlos2380/mytheresa/blob/main/internal/handlers/productHandler_test.go

By combining mocks and table-driven tests, the setup is highly maintainable, fast, and aligned with the project's requirement for self-contained, dependency-free tests.

### Flags

I have added flags to start the application with different setups.

```GO
port := flag.String("port", "8000", "Port on which the server will be listening for incoming requests.")
portMongodb := flag.String("port_mongodb", "27017", "Port on which the MongoDB server will be serving connections.")
ipMongodb := flag.String("ip_mongodb", "mongodb", "IP address on which the MongoDB server will be listening.")
databaseName := flag.String("database_name", "mytheresadb", "Name of the database connection.")
pageSizeStr := flag.String("page_size_product", "5", "Number the products per page response.")
flag.Parse()
```

### Close server properly

To prevent the server from shutting down while there are pending tasks. The server captures any shutdown signal and waits for pending tasks to finish.

The server has programmed a counter that after 5 seconds of receiving the shutdown signal. The server will shut down, avoiding stuck tasks and the server never shuts down.

### CORS (Control Access HTTP)

CORS are enable to be able to connect swagger with apigo.

```GO
func setHeaders(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
}
```

### Linter

I have integrated a linter into our project to ensure code quality and consistency. The linter checks for potential issues and enforces coding standards.
- https://github.com/carlos2380/mytheresa/blob/main/golangci.yml

#### Intall

````
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
`````

#### Run

```
golangci-lint run
```

### Continuous Integration (CI)

Continuous Integration (CI) pipeline is added to automatically test and validate the codebase with each push to the repository. This ensures that our code remains reliable and that new changes do not introduce any issues.

- https://github.com/carlos2380/mytheresa/blob/main/.github/workflows/ci.yml

### Error Handling

Centralized error handling mechanism with a dedicated errors.go file. 
- https://github.com/carlos2380/mytheresa/blob/main/internal/errors/errors.go

This approach ensures consistent and clear error messages, making the application easier to maintain add new errors and debug, while adhering to best programming practices.
- https://github.com/carlos2380/mytheresa/blob/main/internal/errors/codes.go

## Areas for Improvement and Next steps

### Enhanced Test Coverage

More Test Cases: Expand the current test suite to cover a wider range of scenarios, including edge cases and error conditions. This ensures that the application handles all possible inputs and states correctly.

### Add products

Implementing a POST /products endpoint would enhance the API by allowing users to add new products dynamically. This addition would complement the existing GET /products functionality, enabling full CRUD operations and expanding the system's usability.

### Automatic Documentation Generation

Integrate automatic generation of documentation into the CI pipeline. This ensures that the documentation is always up-to-date with the latest code changes.

### Monitoring and Metrics

Prometheus Integration: Integrate Prometheus for monitoring application performance and gathering metrics. This helps in identifying performance bottlenecks and monitoring the health of the application.

Grafana Dashboards: Set up Grafana dashboards to visualize the metrics collected by Prometheus, providing a clear and accessible overview of application performance and health.

Health Check: Implement an endpoint dedicated to verifying the status of key system components, such as database connection, service availability, etc.
