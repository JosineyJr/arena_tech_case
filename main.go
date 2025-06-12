// main.go
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type IPRecord struct {
	LowerIPID   uint32
	UpperIPIP   uint32
	CountryCode string
	CountryName string
	Region      string
	City        string
}

type LocationResponse struct {
	Country     string `json:"country"`
	CountryCode string `json:"countryCode"`
	City        string `json:"city"`
}

var ipDatabase []IPRecord

func ipToID(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address format")
	}

	ip = ip.To4()
	if ip == nil {
		return 0, fmt.Errorf("not an IPv4 address")
	}

	return (uint32(ip[0]) << 24) | (uint32(ip[1]) << 16) | (uint32(ip[2]) << 8) | uint32(ip[3]), nil
}

func loadData(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("could not open csv file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("could not read csv data: %w", err)
	}

	ipDatabase = make([]IPRecord, len(records))
	for i, row := range records {
		lowerIP, _ := strconv.ParseUint(row[0], 10, 32)
		upperIP, _ := strconv.ParseUint(row[1], 10, 32)

		ipDatabase[i] = IPRecord{
			LowerIPID:   uint32(lowerIP),
			UpperIPIP:   uint32(upperIP),
			CountryCode: row[2],
			CountryName: row[3],
			Region:      row[4],
			City:        row[5],
		}
	}

	log.Printf("Successfully loaded %d records into memory.", len(ipDatabase))
	return nil
}

func findLocationByIPID(id uint32) *LocationResponse {
	index := sort.Search(len(ipDatabase), func(i int) bool {
		return ipDatabase[i].LowerIPID > id
	})

	if index > 0 {
		record := ipDatabase[index-1]
		if id >= record.LowerIPID && id <= record.UpperIPIP {
			return &LocationResponse{
				Country:     record.CountryName,
				CountryCode: record.CountryCode,
				City:        record.City,
			}
		}
	}

	return nil
}

func locationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ipQuery := r.URL.Query().Get("ip")
	if ipQuery == "" {
		http.Error(w, `{"error": "IP query parameter is required"}`, http.StatusBadRequest)
		return
	}

	ipID, err := ipToID(ipQuery)
	if err != nil {
		http.Error(w, `{"error": "Invalid IPv4 address format"}`, http.StatusNotFound)
		return
	}

	location := findLocationByIPID(ipID)

	w.Header().Set("Content-Type", "application/json")
	if location == nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Not Found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(location)
}

func main() {
	const datasetFile = "./backend_test/IP2LOCATION-LITE-DB11.CSV"

	log.Println("Starting IP Location API...")

	if err := loadData(datasetFile); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	http.HandleFunc("/ip/location", locationHandler)

	const port = "3000"
	log.Printf("Server listening at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
