package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestGreetHandler tests the greet function
func TestGreetHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet?name=Alice", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expected := "Hello, Alice!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerWithoutName tests the greet function when no name is provided
func TestGreetHandlerWithoutName(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expected := "Hello, Stranger!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerWithEmptyName tests the greet function when name parameter is empty
func TestGreetHandlerWithEmptyName(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet?name=", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expected := "Hello, Stranger!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerWithSpecialCharacters tests the greet function with special characters in name
func TestGreetHandlerWithSpecialCharacters(t *testing.T) {
	// Test with URL encoded special characters
	req := httptest.NewRequest("GET", "/greeter/greet?name="+url.QueryEscape("José María"), nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expected := "Hello, José María!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerWithLongName tests the greet function with a very long name
func TestGreetHandlerWithLongName(t *testing.T) {
	longName := strings.Repeat("A", 1000)
	req := httptest.NewRequest("GET", "/greeter/greet?name="+url.QueryEscape(longName), nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	expected := "Hello, " + longName + "!\n"
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != expected {
		t.Errorf("Expected body %q, got %q", expected, string(body))
	}
}

// TestGreetHandlerDifferentMethods tests the greet function with different HTTP methods
func TestGreetHandlerDifferentMethods(t *testing.T) {
	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/greeter/greet?name=TestUser", nil)
			w := httptest.NewRecorder()

			greet(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			expected := "Hello, TestUser!\n"
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200 for %s method, got %d", method, resp.StatusCode)
			}
			if string(body) != expected {
				t.Errorf("Expected body %q for %s method, got %q", expected, method, string(body))
			}
		})
	}
}

// TestGreetHandlerTableDriven uses table-driven tests for multiple scenarios
func TestGreetHandlerTableDriven(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected string
	}{
		{"Normal name", "name=John", "Hello, John!\n"},
		{"Name with spaces", "name=" + url.QueryEscape("John Doe"), "Hello, John Doe!\n"},
		{"Single character", "name=X", "Hello, X!\n"},
		{"Numbers in name", "name=User123", "Hello, User123!\n"},
		{"Name with symbols", "name=" + url.QueryEscape("User@Domain"), "Hello, User@Domain!\n"},
		{"Unicode name", "name=" + url.QueryEscape("用户"), "Hello, 用户!\n"},
		{"Empty query", "", "Hello, Stranger!\n"},
		{"Multiple name parameters", "name=First&name=Second", "Hello, First!\n"}, // Should use first occurrence
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/greeter/greet"
			if tc.query != "" {
				url += "?" + tc.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			greet(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
			if string(body) != tc.expected {
				t.Errorf("Expected body %q, got %q", tc.expected, string(body))
			}
		})
	}
}

// TestGreetHandlerContentType tests that the response has the correct content type
func TestGreetHandlerContentType(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/greet?name=Alice", nil)
	w := httptest.NewRecorder()

	greet(w, req)

	resp := w.Result()

	// Note: Go's http.ResponseWriter sets text/plain by default for fmt.Fprintf
	contentType := resp.Header.Get("Content-Type")
	if contentType != "" && !strings.Contains(contentType, "text/plain") {
		t.Errorf("Expected Content-Type to be text/plain or empty, got %q", contentType)
	}
}

// TestFarewellHandler tests the farewell function
func TestFarewellHandler(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected string
	}{
		{"With name", "name=Alice", "Goodbye, Alice! Have a great day!\n"},
		{"Without name", "", "Goodbye, Stranger! Have a great day!\n"},
		{"Empty name", "name=", "Goodbye, Stranger! Have a great day!\n"},
		{"Special characters", "name=" + url.QueryEscape("José"), "Goodbye, José! Have a great day!\n"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/greeter/farewell"
			if tc.query != "" {
				url += "?" + tc.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			farewell(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}
			if string(body) != tc.expected {
				t.Errorf("Expected body %q, got %q", tc.expected, string(body))
			}
		})
	}
}

// TestHealthCheckHandler tests the health check endpoint
func TestHealthCheckHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/greeter/health", nil)
	w := httptest.NewRecorder()

	healthCheck(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to contain application/json, got %q", contentType)
	}

	var health HealthResponse
	if err := json.Unmarshal(body, &health); err != nil {
		t.Errorf("Failed to unmarshal health response: %v", err)
	}

	if health.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %q", health.Status)
	}
	if health.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got %q", health.Version)
	}
	if health.Timestamp.IsZero() {
		t.Error("Expected timestamp to be set")
	}
}

// TestTimeBasedGreetHandler tests the time-based greeting function
func TestTimeBasedGreetHandler(t *testing.T) {
	testCases := []struct {
		name     string
		query    string
		expected []string // Multiple possible greetings based on time
	}{
		{"With name", "name=Alice", []string{"Good morning, Alice!\n", "Good afternoon, Alice!\n", "Good evening, Alice!\n"}},
		{"Without name", "", []string{"Good morning, Stranger!\n", "Good afternoon, Stranger!\n", "Good evening, Stranger!\n"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/greeter/time-greet"
			if tc.query != "" {
				url += "?" + tc.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			timeBasedGreet(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status 200, got %d", resp.StatusCode)
			}

			// Check if the response matches any of the expected greetings
			bodyStr := string(body)
			found := false
			for _, expected := range tc.expected {
				if bodyStr == expected {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected one of %v, got %q", tc.expected, bodyStr)
			}
		})
	}
}

// TestUserInfoHandlerGET tests the GET method of user info endpoint
func TestUserInfoHandlerGET(t *testing.T) {
	testCases := []struct {
		name           string
		query          string
		expectedStatus int
		expectedName   string
		expectedAge    int
	}{
		{"Basic user", "name=John", http.StatusOK, "John", 0},
		{"User with age", "name=John&age=25", http.StatusOK, "John", 25},
		{"User with all fields", "name=John&age=25&location=NYC&email=john@example.com", http.StatusOK, "John", 25},
		{"Missing name", "", http.StatusBadRequest, "", 0},
		{"Invalid age", "name=John&age=invalid", http.StatusOK, "John", 0},
		{"Negative age", "name=John&age=-5", http.StatusOK, "John", 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/greeter/user-info"
			if tc.query != "" {
				url += "?" + tc.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			userInfoHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectedStatus == http.StatusOK {
				var user UserInfo
				if err := json.Unmarshal(body, &user); err != nil {
					t.Errorf("Failed to unmarshal user response: %v", err)
				}

				if user.Name != tc.expectedName {
					t.Errorf("Expected name %q, got %q", tc.expectedName, user.Name)
				}
				if user.Age != tc.expectedAge {
					t.Errorf("Expected age %d, got %d", tc.expectedAge, user.Age)
				}
			}
		})
	}
}

// TestUserInfoHandlerPOST tests the POST method of user info endpoint
func TestUserInfoHandlerPOST(t *testing.T) {
	testCases := []struct {
		name           string
		payload        UserInfo
		expectedStatus int
		expectError    bool
	}{
		{"Valid user", UserInfo{Name: "John", Age: 25, Location: "NYC"}, http.StatusCreated, false},
		{"Minimal user", UserInfo{Name: "Jane"}, http.StatusCreated, false},
		{"Missing name", UserInfo{Age: 25}, http.StatusBadRequest, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tc.payload)
			req := httptest.NewRequest("POST", "/greeter/user-info", bytes.NewBuffer(payloadBytes))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			userInfoHandler(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if !tc.expectError && tc.expectedStatus == http.StatusCreated {
				var response map[string]interface{}
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if message, ok := response["message"].(string); !ok || !strings.Contains(message, tc.payload.Name) {
					t.Errorf("Expected message to contain user name %q", tc.payload.Name)
				}
			}
		})
	}
}

// TestUserInfoHandlerInvalidMethod tests unsupported HTTP methods
func TestUserInfoHandlerInvalidMethod(t *testing.T) {
	methods := []string{"PUT", "DELETE", "PATCH"}

	for _, method := range methods {
		t.Run("Method_"+method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/greeter/user-info", nil)
			w := httptest.NewRecorder()

			userInfoHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != http.StatusMethodNotAllowed {
				t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, resp.StatusCode)
			}
		})
	}
}

// TestBulkGreetHandler tests the bulk greeting endpoint
func TestBulkGreetHandler(t *testing.T) {
	testCases := []struct {
		name           string
		query          string
		expectedStatus int
		expectedCount  int
	}{
		{"Single name", "names=Alice", http.StatusOK, 1},
		{"Multiple names", "names=Alice,Bob,Charlie", http.StatusOK, 3},
		{"Names with spaces", "names=" + url.QueryEscape("Alice,Bob Smith,Charlie"), http.StatusOK, 3},
		{"Empty names", "names=,,", http.StatusOK, 1}, // Should return "Hello, Stranger!"
		{"Missing names parameter", "", http.StatusBadRequest, 0},
		{"Names with extra spaces", "names=" + url.QueryEscape(" Alice , Bob , Charlie "), http.StatusOK, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "/greeter/bulk-greet"
			if tc.query != "" {
				url += "?" + tc.query
			}

			req := httptest.NewRequest("GET", url, nil)
			w := httptest.NewRecorder()

			bulkGreet(w, req)

			resp := w.Result()
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, resp.StatusCode)
			}

			if tc.expectedStatus == http.StatusOK {
				var response map[string][]string
				if err := json.Unmarshal(body, &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				greetings, ok := response["greetings"]
				if !ok {
					t.Error("Expected 'greetings' field in response")
				}

				if len(greetings) != tc.expectedCount {
					t.Errorf("Expected %d greetings, got %d", tc.expectedCount, len(greetings))
				}

				// Verify greeting format
				for _, greeting := range greetings {
					if !strings.HasPrefix(greeting, "Hello, ") || !strings.HasSuffix(greeting, "!") {
						t.Errorf("Invalid greeting format: %q", greeting)
					}
				}
			}
		})
	}
}

// BenchmarkGreetHandler benchmarks the greet function performance
func BenchmarkGreetHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/greeter/greet?name=BenchmarkUser", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		greet(w, req)
	}
}

// BenchmarkGreetHandlerNoName benchmarks the greet function with no name parameter
func BenchmarkGreetHandlerNoName(b *testing.B) {
	req := httptest.NewRequest("GET", "/greeter/greet", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		greet(w, req)
	}
}

// BenchmarkFarewellHandler benchmarks the farewell function
func BenchmarkFarewellHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/greeter/farewell?name=BenchmarkUser", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		farewell(w, req)
	}
}

// BenchmarkHealthCheckHandler benchmarks the health check function
func BenchmarkHealthCheckHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/greeter/health", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		healthCheck(w, req)
	}
}

// BenchmarkBulkGreetHandler benchmarks the bulk greet function
func BenchmarkBulkGreetHandler(b *testing.B) {
	req := httptest.NewRequest("GET", "/greeter/bulk-greet?names=Alice,Bob,Charlie", nil)

	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		bulkGreet(w, req)
	}
}
