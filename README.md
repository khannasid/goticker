# 🚀 GoTicker: Real-Time Distributed Crypto Aggregator

A high-throughput, low-latency cryptocurrency price aggregator built with Go. 
This project demonstrates **concurrency patterns**, **race-condition handling**, and **real-time data broadcasting**.

## 🏗 Architecture
The system follows a **Fan-In / Fan-Out** concurrency architecture:
1.  **Ingestion (Fan-In):** Multiple independent Goroutines fetch data from **Binance** and **Coinbase** APIs concurrently.
2.  **Processing:** A central engine normalizes data and updates a thread-safe in-memory store (protected by `sync.RWMutex`).
3.  **Distribution (Fan-Out):** A WebSocket Hub broadcasts updates to connected clients in O(1) time relative to the data source.

## 🛠 Tech Stack
* **Language:** Golang (1.21+)
* **Concurrency:** Goroutines, Buffered Channels, Mutexes, WaitGroups
* **Transport:** WebSockets (Gorilla), REST API (net/http)
* **Design Patterns:** Producer-Consumer, Dependency Injection, Clean Architecture

## 🚀 How to Run
1.  **Clone the repo**
    ```bash
    git clone [https://github.com/YOUR_USERNAME/goticker.git](https://github.com/YOUR_USERNAME/goticker.git)
    cd goticker
    ```

2.  **Run the Engine**
    ```bash
    go run cmd/server/main.go
    ```

3.  **View the Dashboard**
    Open `index.html` in your browser.

## 📂 Project Structure
```text
├── cmd/server/       # Entry point (Main)
├── internal/
│   ├── api/          # REST Handlers & WebSocket Hub
│   ├── generator/    # Data Fetchers (Producers)
│   ├── model/        # Data Structs
│   ├── processor/    # Aggregation Logic (Consumer)
│   └── store/        # Thread-Safe In-Memory DB
└── index.html        # Real-time Dashboard