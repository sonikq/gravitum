package reader

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

// mockReadCloser implements io.ReadCloser for testing
type mockReadCloser struct {
	reader    io.Reader
	closeFunc func() error
}

func (m mockReadCloser) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m mockReadCloser) Close() error {
	if m.closeFunc != nil {
		return m.closeFunc()
	}
	return nil
}

func TestGetBody_Success(t *testing.T) {
	// Test data
	expected := []byte("test data")

	// Create a mock ReadCloser with our test data
	mockRC := mockReadCloser{
		reader: bytes.NewReader(expected),
	}

	// Call the function
	result, err := GetBody(mockRC)

	// Check results
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestGetBody_ReadError(t *testing.T) {
	// Create a mock ReadCloser that returns an error on Read
	expectedErr := errors.New("read error")
	mockRC := mockReadCloser{
		reader: &errorReader{err: expectedErr},
	}

	// Call the function
	_, err := GetBody(mockRC)

	// Check results
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestGetBody_CloseError(t *testing.T) {
	// Create a mock ReadCloser that returns an error on Close
	expectedErr := errors.New("close error")
	mockRC := mockReadCloser{
		reader:    bytes.NewReader([]byte("test data")),
		closeFunc: func() error { return expectedErr },
	}

	// Call the function
	_, err := GetBody(mockRC)

	// Check results
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error %v, got %v", expectedErr, err)
	}
}

func TestGetBody_EmptyReader(t *testing.T) {
	// Create a mock ReadCloser with empty data
	mockRC := mockReadCloser{
		reader: bytes.NewReader([]byte{}),
	}

	// Call the function
	result, err := GetBody(mockRC)

	// Check results
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %v bytes", len(result))
	}
}

func TestGetBody_LargeData(t *testing.T) {
	// Create a large test data (1MB)
	size := 1024 * 1024
	expected := make([]byte, size)
	for i := 0; i < size; i++ {
		expected[i] = byte(i % 256)
	}

	// Create a mock ReadCloser with our large test data
	mockRC := mockReadCloser{
		reader: bytes.NewReader(expected),
	}

	// Call the function
	result, err := GetBody(mockRC)

	// Check results
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !bytes.Equal(result, expected) {
		t.Errorf("Data mismatch for large input")
	}

	if len(result) != size {
		t.Errorf("Expected %d bytes, got %d bytes", size, len(result))
	}
}

// errorReader is a helper that always returns an error when Read is called
type errorReader struct {
	err error
}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, e.err
}
