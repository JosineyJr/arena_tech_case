// main_test.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIpToID(t *testing.T) {
	testCases := []struct {
		name        string
		ip          string
		expectedID  uint32
		expectError bool
	}{
		{"Valid Public IP", "8.8.8.8", 134744072, false},
		{"Valid Private IP", "192.168.1.1", 3232235777, false},
		{"Loopback Address", "127.0.0.1", 2130706433, false},
		{"Zero IP", "0.0.0.0", 0, false},
		{"Max IP", "255.255.255.255", 4294967295, false},
		{"Invalid String", "not an ip", 0, true},
		{"Invalid Format", "1.2.3", 0, true},
		{"Out of Range", "256.1.1.1", 0, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := ipToID(tc.ip)
			if (err != nil) != tc.expectError {
				t.Errorf("ipToID() error = %v, expectError %v", err, tc.expectError)
				return
			}
			if !tc.expectError && id != tc.expectedID {
				t.Errorf("ipToID() = %v, want %v", id, tc.expectedID)
			}
		})
	}
}

func setupMockDatabase() {
	ipDatabase = []IPRecord{
		{
			LowerIPID:   16777216,
			UpperIPIP:   16777471,
			CountryCode: "US",
			CountryName: "United States",
			City:        "Mountain View",
		},
		{
			LowerIPID:   134744072,
			UpperIPIP:   134744072,
			CountryCode: "US",
			CountryName: "United States",
			City:        "New York",
		},
		{
			LowerIPID:   3232235520,
			UpperIPIP:   3232236031,
			CountryCode: "DE",
			CountryName: "Germany",
			City:        "Berlin",
		},
	}
}

func TestFindLocationByIPID(t *testing.T) {
	setupMockDatabase()

	testCases := []struct {
		name       string
		ipID       uint32
		expectFind bool
		expected   *LocationResponse
	}{
		{
			name:       "Finds IP at start of range",
			ipID:       16777216,
			expectFind: true,
			expected: &LocationResponse{
				Country:     "United States",
				CountryCode: "US",
				City:        "Mountain View",
			},
		},
		{
			name:       "Finds IP in middle of range",
			ipID:       3232235777,
			expectFind: true,
			expected:   &LocationResponse{Country: "Germany", CountryCode: "DE", City: "Berlin"},
		},
		{
			name:       "IP in a gap between ranges",
			ipID:       100000000,
			expectFind: false,
			expected:   nil,
		},
		{
			name:       "IP below all ranges",
			ipID:       1,
			expectFind: false,
			expected:   nil,
		},
		{
			name:       "IP above all ranges",
			ipID:       4000000000,
			expectFind: false,
			expected:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := findLocationByIPID(tc.ipID)
			if (result != nil) != tc.expectFind {
				t.Errorf("findLocationByIPID() found = %v, want %v", result != nil, tc.expectFind)
			}
			if tc.expectFind {
				if result.Country != tc.expected.Country || result.City != tc.expected.City {
					t.Errorf("findLocationByIPID() got %v, want %v", result, tc.expected)
				}
			}
		})
	}
}

func TestLocationHandler(t *testing.T) {
	setupMockDatabase()

	testCases := []struct {
		name               string
		url                string
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:               "Successful Lookup",
			url:                "/ip/location?ip=8.8.8.8",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"country":"United States","countryCode":"US","city":"New York"}`,
		},
		{
			name:               "IP Not Found",
			url:                "/ip/location?ip=127.0.0.1",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"error":"Not Found"}`,
		},
		{
			name:               "Invalid IP Format",
			url:                "/ip/location?ip=999.999.999",
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"error": "Invalid IPv4 address format"}`,
		},
		{
			name:               "Missing IP Parameter",
			url:                "/ip/location?ip=",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "IP query parameter is required"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			rr := httptest.NewRecorder()

			locationHandler(rr, req)

			if status := rr.Code; status != tc.expectedStatusCode {
				t.Errorf(
					"handler returned wrong status code: got %v want %v",
					status,
					tc.expectedStatusCode,
				)
			}

			actualBody := strings.TrimSpace(rr.Body.String())

			var expectedMap, actualMap map[string]interface{}
			json.Unmarshal([]byte(tc.expectedBody), &expectedMap)
			json.Unmarshal([]byte(actualBody), &actualMap)

			if fmt.Sprintf("%v", actualMap) != fmt.Sprintf("%v", expectedMap) {
				t.Errorf(
					"handler returned unexpected body: got %v want %v",
					actualBody,
					tc.expectedBody,
				)
			}
		})
	}
}
