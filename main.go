package main

import (
    "fmt"
    "net"
    "os"
    "sync"
    "time"
)

func scanPort(host string, port int, wg *sync.WaitGroup, results chan<- int) {
    defer wg.Done()
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, 1*time.Second)
    if err != nil {
        return
    }
    conn.Close()
    results <- port
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage: go run main.go <host> <startPort-endPort>")
        return
    }

    host := os.Args[1]
    var startPort, endPort int
    _, err := fmt.Sscanf(os.Args[2], "%d-%d", &startPort, &endPort)
    if err != nil {
        fmt.Println("Invalid port range format. Use startPort-endPort")
        return
    }

    var wg sync.WaitGroup
    ports := make(chan int, endPort-startPort+1)
    results := make(chan int, endPort-startPort+1)

    for port := startPort; port <= endPort; port++ {
        wg.Add(1)
        go scanPort(host, port, &wg, results)
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    openPorts := []int{}
    for port := range results {
        openPorts = append(openPorts, port)
    }

    fmt.Printf("Open ports on %s:\n", host)
    for _, port := range openPorts {
        fmt.Println(port)
    }
}
