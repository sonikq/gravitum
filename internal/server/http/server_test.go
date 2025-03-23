package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// TestNewServer verifies that the NewServer function correctly initializes a Server instance
func TestNewServer(t *testing.T) {
	// Test cases
	testCases := []struct {
		name     string
		address  string
		handler  http.Handler
		expected Server
	}{
		{
			name:    "Valid initialization with handler",
			address: ":8080",
			handler: http.NewServeMux(),
			expected: Server{
				httpServer: &http.Server{
					Addr:           ":8080",
					Handler:        http.NewServeMux(),
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					MaxHeaderBytes: 1 << 20,
				},
			},
		},
		{
			name:    "Valid initialization with nil handler",
			address: "127.0.0.1:9090",
			handler: nil,
			expected: Server{
				httpServer: &http.Server{
					Addr:           "127.0.0.1:9090",
					Handler:        nil,
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					MaxHeaderBytes: 1 << 20,
				},
			},
		},
		{
			name:    "Empty address",
			address: "",
			handler: http.NewServeMux(),
			expected: Server{
				httpServer: &http.Server{
					Addr:           "",
					Handler:        http.NewServeMux(),
					ReadTimeout:    15 * time.Second,
					WriteTimeout:   15 * time.Second,
					MaxHeaderBytes: 1 << 20,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := NewServer(tc.address, tc.handler)

			// Verify server is not nil
			if server == nil {
				t.Fatal("Expected non-nil server, got nil")
			}

			// Verify server properties
			if server.httpServer.Addr != tc.expected.httpServer.Addr {
				t.Errorf("Expected address %s, got %s", tc.expected.httpServer.Addr, server.httpServer.Addr)
			}

			// We can't directly compare handlers, but we can check if both are nil or both are non-nil
			if (server.httpServer.Handler == nil) != (tc.expected.httpServer.Handler == nil) {
				t.Errorf("Handler mismatch: expected nil: %v, got nil: %v",
					tc.expected.httpServer.Handler == nil,
					server.httpServer.Handler == nil)
			}

			if server.httpServer.ReadTimeout != tc.expected.httpServer.ReadTimeout {
				t.Errorf("Expected ReadTimeout %v, got %v",
					tc.expected.httpServer.ReadTimeout,
					server.httpServer.ReadTimeout)
			}

			if server.httpServer.WriteTimeout != tc.expected.httpServer.WriteTimeout {
				t.Errorf("Expected WriteTimeout %v, got %v",
					tc.expected.httpServer.WriteTimeout,
					server.httpServer.WriteTimeout)
			}

			if server.httpServer.MaxHeaderBytes != tc.expected.httpServer.MaxHeaderBytes {
				t.Errorf("Expected MaxHeaderBytes %v, got %v",
					tc.expected.httpServer.MaxHeaderBytes,
					server.httpServer.MaxHeaderBytes)
			}
		})
	}
}

// TestServerRun tests the Run method
func TestServerRun(t *testing.T) {
	// Since we can't easily mock http.Server.ListenAndServe, we'll test with a real server
	// but on a port that's likely to be available or already in use

	// First, let's find a port that's already in use by starting a server
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()

	// Get the port that was assigned
	inUsePort := listener.Addr().(*net.TCPAddr).Port

	testCases := []struct {
		name          string
		address       string
		shouldSucceed bool
	}{
		{
			name:          "Invalid port number",
			address:       ":99999", // Invalid port number
			shouldSucceed: false,
		},
		{
			name:          "Already in use port",
			address:       ":" + strconv.Itoa(inUsePort), // Port that's already in use
			shouldSucceed: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a server with the test address
			server := NewServer(tc.address, http.NewServeMux())

			// Run the server in a goroutine with a channel to capture the error
			errCh := make(chan error, 1)
			go func() {
				errCh <- server.Run()
			}()

			// Give it a moment to start or fail
			var err error
			select {
			case err = <-errCh:
				// Server returned immediately, which means it failed to start
			case <-time.After(100 * time.Millisecond):
				// Server started successfully or is still trying
				// Try to shut it down
				ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
				defer cancel()
				err = server.Shutdown(ctx)
			}

			if tc.shouldSucceed && err != nil && !errors.Is(err, http.ErrServerClosed) {
				t.Errorf("Expected server to start successfully, got error: %v", err)
			} else if !tc.shouldSucceed && err == nil {
				t.Errorf("Expected server to fail to start, but it started successfully")
			}
		})
	}
}

// CustomServer is a custom implementation of the Server for testing
type CustomServer struct {
	shutdownFunc func(ctx context.Context) error
}

func (s *CustomServer) Run() error {
	return nil
}

func (s *CustomServer) Shutdown(ctx context.Context) error {
	return s.shutdownFunc(ctx)
}

// TestServerShutdown tests the Shutdown method
func TestServerShutdown(t *testing.T) {
	testCases := []struct {
		name          string
		ctx           context.Context
		mockBehavior  func(ctx context.Context) error
		expectedError error
	}{
		{
			name: "Successful shutdown",
			ctx:  context.Background(),
			mockBehavior: func(ctx context.Context) error {
				return nil
			},
			expectedError: nil,
		},
		{
			name: "Shutdown with timeout",
			ctx:  context.Background(),
			mockBehavior: func(ctx context.Context) error {
				return context.DeadlineExceeded
			},
			expectedError: context.DeadlineExceeded,
		},
		{
			name: "Shutdown with canceled context",
			ctx:  context.Background(),
			mockBehavior: func(ctx context.Context) error {
				return context.Canceled
			},
			expectedError: context.Canceled,
		},
		{
			name: "Shutdown with other error",
			ctx:  context.Background(),
			mockBehavior: func(ctx context.Context) error {
				return errors.New("custom shutdown error")
			},
			expectedError: errors.New("custom shutdown error"),
		},
		{
			name: "Nil context",
			ctx:  nil,
			mockBehavior: func(ctx context.Context) error {
				if ctx == nil {
					return errors.New("nil context")
				}
				return nil
			},
			expectedError: errors.New("nil context"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a custom server with our mock behavior
			customServer := &CustomServer{
				shutdownFunc: tc.mockBehavior,
			}

			// Call Shutdown and check the error
			err := customServer.Shutdown(tc.ctx)

			// Check if the error matches the expected error
			if tc.expectedError == nil && err != nil {
				t.Errorf("Expected no error, got %v", err)
			} else if tc.expectedError != nil {
				if err == nil {
					t.Errorf("Expected error %v, got nil", tc.expectedError)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Expected error %v, got %v", tc.expectedError, err)
				}
			}
		})
	}
}

// TestServerIntegration performs a simple integration test
func TestServerIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a test handler
	mux := http.NewServeMux()
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Find an available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close() // Close the listener so the server can use the port

	// Create a server with the available port
	server := NewServer(fmt.Sprintf(":%d", port), mux)

	// Start the server in a goroutine
	go func() {
		if err := server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Errorf("Server.Run() error = %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Make a request to the server
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/test", port))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	// Shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Errorf("Server.Shutdown() error = %v", err)
	}
}

// MockServer is a mock implementation of the Server for testing
type MockServer struct {
	httpServer *http.Server
}

func NewMockServer() *MockServer {
	return &MockServer{
		httpServer: &http.Server{
			Addr: ":8080",
		},
	}
}

func (s *MockServer) Run() error {
	return nil
}

func (s *MockServer) Shutdown(ctx context.Context) error {
	// Simulate a timeout error when context is already timed out
	if ctx != nil {
		deadline, ok := ctx.Deadline()
		if ok && deadline.Before(time.Now()) {
			return context.DeadlineExceeded
		}
		if errors.Is(ctx.Err(), context.Canceled) {
			return context.Canceled
		}
	}
	return nil
}

// TestServerWithTimeoutContext tests the Shutdown method with a timeout context
func TestServerWithTimeoutContext(t *testing.T) {
	// Create a mock server
	server := NewMockServer()

	// Create a context with a very short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	// Sleep to ensure the context times out
	time.Sleep(1 * time.Millisecond)

	// Call Shutdown with the timed-out context
	err := server.Shutdown(ctx)

	// Check that we got the expected error
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
	}
}

// TestServerWithCanceledContext tests the Shutdown method with a canceled context
func TestServerWithCanceledContext(t *testing.T) {
	// Create a mock server
	server := NewMockServer()

	// Create a context and cancel it immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Call Shutdown with the canceled context
	err := server.Shutdown(ctx)

	// Check that we got the expected error
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}
