package jira

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type editTestServer struct{ code int }

func (e *editTestServer) serve(t *testing.T, expectedBody string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/rest/api/2/issue/TEST-1", r.URL.Path)
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Accept"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		actualBody := new(strings.Builder)
		_, _ = io.Copy(actualBody, r.Body)

		assert.JSONEq(t, expectedBody, actualBody.String())

		w.WriteHeader(e.code)
	}))
}

func TestEditWithADFCustomField(t *testing.T) {
	// ADF fields in edit are wrapped in set operations: [{"set": "converted markdown"}]
	expectedBody := `{"update":{"customfield_10042":[{"set":"*Bold text*"}]},"fields":{"parent":{}}}`
	testServer := editTestServer{code: 204}
	server := testServer.serve(t, expectedBody)
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))

	err := client.Edit("TEST-1", &EditRequest{
		ADFFields: map[string]string{
			"customfield_10042": "**Bold text**",
		},
	})
	assert.NoError(t, err)
}

func TestEditWithADFCustomFieldFromFile(t *testing.T) {
	tmpDir := t.TempDir()
	mdFile := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(mdFile, []byte("**Bold text**"), 0644)
	assert.NoError(t, err)

	expectedBody := `{"update":{"customfield_10042":[{"set":"*Bold text*"}]},"fields":{"parent":{}}}`
	testServer := editTestServer{code: 204}
	server := testServer.serve(t, expectedBody)
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))

	err = client.Edit("TEST-1", &EditRequest{
		ADFFields: map[string]string{
			"customfield_10042": "@" + mdFile,
		},
	})
	assert.NoError(t, err)
}

func TestEditWithSummaryAndADFField(t *testing.T) {
	expectedBody := `{"update":{"summary":[{"set":"New title"}],"customfield_10042":[{"set":"*Bold text*"}]},"fields":{"parent":{}}}`
	testServer := editTestServer{code: 204}
	server := testServer.serve(t, expectedBody)
	defer server.Close()

	client := NewClient(Config{Server: server.URL}, WithTimeout(3*time.Second))

	err := client.Edit("TEST-1", &EditRequest{
		Summary: "New title",
		ADFFields: map[string]string{
			"customfield_10042": "**Bold text**",
		},
	})
	assert.NoError(t, err)
}
