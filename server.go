package main

import (
    "fmt"
    "net"
	"os"
	"encoding/csv"
)

const (
    CONN_HOST = "localhost"
    CONN_PORT = "4040"
    CONN_TYPE = "tcp"
)
type covidData struct {
    cum_test_positive string
    cum_test_performed string
	date string
	discharged string
	expired string
	region string
	admitted string
}

func main() {

	csvFile, err := os.Open("covid_final_data.csv")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened CSV file")
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
    if err != nil {
        fmt.Println(err)
    }    
    for _, line := range csvLines {
        patients := covidData{
            cum_test_positive: line[0],
            cum_test_performed: line[1],
			date: line[2],
			discharged: line[3],
			expired: line[4],
			region: line[5],
			admitted: line[6],
        }
        fmt.Println(patients.cum_test_performed + " " + patients.cum_test_positive + " " + patients.region)
    }

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }
        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
  // Make a buffer to hold incoming data.
  buf := make([]byte, 1024)
  // Read the incoming connection into the buffer.
  _ , err := conn.Read(buf)
  if err != nil {
    fmt.Println("Error reading:", err.Error())
  }
  // Send a response back to person contacting us.
  conn.Write([]byte("Message received."))
  // Close the connection when you're done with it.
  conn.Close()
}