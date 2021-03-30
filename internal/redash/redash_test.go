package redash

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getServer(body string, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}))
}

func TestGetQuery(t *testing.T) {
	body := `{"id": 42, "name": "New Query", "tags": ["foo", "bar", "baz"]}`
	ts := getServer(body, 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q, err := GetQuery(1)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 42, q.Id)
	assert.Equal(t, "New Query", q.Name)
	assert.Equal(t, []string{"foo", "bar", "baz"}, q.Tags)
}

func TestGetQuery_EmptyResponse(t *testing.T) {
	ts := getServer("", 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q, err := GetQuery(1)

	assert.Nil(t, q)
	assert.Error(t, err)
}

func TestGetQuery_NotFound(t *testing.T) {
	ts := getServer(`{"message": "The requested URL was not found on the server."}`, 404)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q, err := GetQuery(1)

	assert.Nil(t, q)
	assert.Error(t, err)
}

func TestExportQueries(t *testing.T) {
	body := `{"results": [
		{"id": 42, "name": "New Query", "tags": ["foo", "bar", "baz"]},
		{"id": 43, "name": "New Query 2", "tags": []}
	]}`
	ts := getServer(body, 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	queries, err := ExportQueries()
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, 2, len(*queries))
	assert.Equal(t, 43, (*queries)[1].Id)
	assert.Equal(t, "New Query 2", (*queries)[1].Name)
	assert.Empty(t, (*queries)[1].Tags)
}

func TestExportQueries_EmptyResponse(t *testing.T) {
	ts := getServer("", 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	queries, err := ExportQueries()

	assert.Nil(t, queries)
	assert.Error(t, err)
}

func TestExportQueriesAsYAML(t *testing.T) {
	body := `{"results": [
		{"id": 42, "name": "New Query", "tags": ["foo", "bar", "baz"]},
		{"id": 43, "name": "New Query 2", "tags": []}
	]}`
	ts := getServer(body, 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	b, err := ExportQueriesAsYAML()
	if err != nil {
		t.Error(err)
	}

	expected := `queries:
  - id: 42
    name: New Query
    tags: [foo, bar, baz]
  - id: 43
    name: New Query 2
    tags: []
`
	assert.Equal(t, expected, b.String())
}

func TestExportQueriesAsYAML_EmptyResponse(t *testing.T) {
	ts := getServer("", 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	b, err := ExportQueriesAsYAML()

	assert.Nil(t, b)
	assert.Error(t, err)
}

func TestUpdateQuery(t *testing.T) {
	ts := getServer("", 200)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q := Query{Id: 42, Name: "New Query", Tags: []string{"foo", "bar", "baz"}}
	err := UpdateQuery(q)
	if err != nil {
		t.Error(err)
	}

	assert.Nil(t, err)
}

func TestUpdateQuery_ConnectionRefused(t *testing.T) {
	ts := getServer(`{"message": "Bad Request"}`, 400)
	ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q := Query{Id: 42, Name: "New Query", Tags: []string{"foo", "bar", "baz"}}
	err := UpdateQuery(q)

	assert.NotNil(t, err)
}

func TestUpdateQuery_BadRequest(t *testing.T) {
	ts := getServer(`{"message": "Bad Request"}`, 400)
	defer ts.Close()

	client = resty.New()
	Init(ts.URL, "supersecret", client)
	q := Query{Id: 42, Name: "New Query", Tags: []string{"foo", "bar", "baz"}}
	err := UpdateQuery(q)

	assert.NotNil(t, err)
}
