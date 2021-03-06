package main

import (
    "fmt"
    "net"
    "os"
    "flag"
    "log"
	"encoding/csv"
	"encoding/json"
    "bufio"
	"io"
	"strings"
)

type CovidPatient struct {
    Positive    string      `json:"Covid_Positive"`
    Performed   string      `json:"Coivd_Performed"`
	Date        string      `json:"Covid_Date"`
	Discharged  string      `json:"Covid_Discharged"`
	Expired     string      `json:"Covid_Expired"`
	Region      string      `json:"Covid_Region"`
	Admitted    string      `json:"Covid_Admitted"`
}

type DataRequest struct {   
	Get string `json:"get"`
}

type DataError struct {     
	Error string `json:"Covid_error"`
}

func Load(path string) []CovidPatient {
	table := make([]CovidPatient, 0)
	var patient CovidPatient
	file, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
    defer file.Close()

	reader := csv.NewReader(file)
	csvData, err := reader.ReadAll()
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
	for _, row := range csvData{
		patient.Positive =  row[0]
		patient.Performed =  row[1]
		patient.Date =       row[2]
        patient.Discharged = row[3]
        patient.Expired =    row[4]
        patient.Region =     row[5]
        patient.Admitted =   row[6]
		table = append(table, patient)
	}
	return table
}

func Find(table []CovidPatient, filter string) []CovidPatient {
	if filter == "" || filter == "*" {
		return table
	}
	result := make([]CovidPatient, 0)
	filter = strings.ToUpper(filter)
	for _, cp := range table {
		if cp.Date == filter ||
			cp.Region == filter ||
			strings.Contains(strings.ToUpper(cp.Positive), filter) 		||
            strings.Contains(strings.ToUpper(cp.Performed), filter) 	||
            strings.Contains(strings.ToUpper(cp.Date), filter) 			||
            strings.Contains(strings.ToUpper(cp.Discharged), filter) 	||
            strings.Contains(strings.ToUpper(cp.Expired), filter) 		||
            strings.Contains(strings.ToUpper(cp.Region), filter) 		||
            strings.Contains(strings.ToUpper(cp.Admitted), filter){
			result = append(result, cp)
		}
	}
	return result
}

var (
	patientsDetail = Load("./covid_final_data.csv")
)

func main(){
    var addr string
	var network string
	flag.StringVar(&addr, "e", ":4040", "service endpoint [ip addr or socket path]")
	flag.StringVar(&network, "n", "tcp", "network protocol [tcp,unix]")
	flag.Parse()

	switch network {
	case "tcp", "tcp4", "tcp6", "unix":
	default:
		fmt.Println("unsupported network protocol")
		os.Exit(1)
	}

	ln, err := net.Listen(network, addr)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer ln.Close()
	log.Println("Covid19 Condition in Pakistan")
    log.Printf("Service started: (%s) %s\n", network, addr)
    
    for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			conn.Close()
			continue
		}
		log.Println("Connected to ", conn.RemoteAddr())
		go handleConnection(conn)
	}
}
func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println("error closing connection:", err)
		}
	}()

	reader := bufio.NewReaderSize(conn, 4)

	for {
		buf, err := reader.ReadSlice('}')
		if err != nil {
			if err != io.EOF {
				log.Println("connection read error:", err)
				return
			}
		}
        reader.Reset(conn)
        
        var req DataRequest
		if err := json.Unmarshal(buf, &req); err != nil {
			log.Println("failed to unmarshal request:", err)
			cerr, jerr := json.Marshal(DataError{Error: err.Error()})
			if jerr != nil {
				log.Println("failed to marshal DataError:", jerr)
				continue
			}
			if _, werr := conn.Write(cerr); werr != nil {
				log.Println("failed to write to DataError:", werr)
				return
			}
			continue
		}

        result := Find(patientsDetail, req.Get)

        rsp, err := json.Marshal(&result)
		if err != nil {
			log.Println("failed to marshal data:", err)
			if _, err := fmt.Fprintf(conn, `{"data_error":"internal error"}`); err != nil {
				log.Printf("failed to write to client: %v", err)
				return
			}
			continue
		}
		if _, err := conn.Write(rsp); err != nil {
			log.Println("failed to write response:", err)
			return
		}
	}
}


