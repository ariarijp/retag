package redash

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"golang.org/x/xerrors"
	"gopkg.in/yaml.v3"
)

type Query struct {
	Id   int      `yaml:"id"`
	Name string   `yaml:"name"`
	Tags []string `yaml:"tags,flow"`
}

type QueriesResponse struct {
	Count    int      `json:"count"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
	Results  []*Query `json:"results"`
}

var apiKey string
var baseUrl string
var client *resty.Client
var pageSize = 50

func Init(_baseUrl string, _apiKey string, _client *resty.Client) {
	apiKey = _apiKey
	baseUrl = _baseUrl
	client = _client
}

func get(url string, params map[string]string, result interface{}) (*resty.Response, error) {
	params["api_key"] = apiKey
	resp, err := client.R().
		SetQueryParams(params).
		SetResult(result).
		Get(url)
	if err != nil {
		return nil, xerrors.Errorf("Error: %+w", err)
	} else if resp.StatusCode() != 200 {
		return resp, xerrors.Errorf("Error: %s, %+w", resp.Status(), resp.String())
	}

	return resp, nil
}

func post(url string, body map[string]interface{}) (*resty.Response, error) {
	resp, err := client.R().
		SetQueryParams(map[string]string{
			"api_key": apiKey,
		}).
		SetBody(body).
		Post(url)
	if err != nil {
		return nil, xerrors.Errorf("Error: %+w", err)
	} else if resp.StatusCode() != 200 {
		return resp, xerrors.Errorf("Error: %s, %+w", resp.Status(), resp.String())
	}

	return resp, nil
}

func GetQuery(id int) (*Query, error) {
	url := fmt.Sprintf("%s/api/queries/%d", baseUrl, id)
	resp, err := get(url, map[string]string{}, Query{})
	if err != nil {
		return nil, err
	}

	return resp.Result().(*Query), nil
}

func ExportQueries() (*[]*Query, error) {
	url := fmt.Sprintf("%s/api/queries", baseUrl)
	var queries []*Query
	page := 1
	count := 0

	for {
		if count < (page-1)*pageSize {
			break
		}

		params := map[string]string{
			"page":      fmt.Sprintf("%d", page),
			"page_size": fmt.Sprintf("%d", pageSize),
			"order":     "created_at",
		}
		resp, err := get(url, params, QueriesResponse{})
		if err != nil {
			return nil, err
		}

		respResult := resp.Result().(*QueriesResponse)

		if count == 0 {
			count = respResult.Count
		}

		queries = append(queries, respResult.Results...)
		page++

		time.Sleep(1000 * time.Millisecond)
	}

	return &queries, nil
}

func ExportQueriesAsYAML() (*bytes.Buffer, error) {
	queries, err := ExportQueries()
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&out)
	yamlEncoder.SetIndent(2)
	err = yamlEncoder.Encode(map[string]*[]*Query{
		"queries": queries,
	})
	if err != nil {
		return nil, err
	}

	return &out, nil
}

func UpdateQuery(query Query) error {
	url := fmt.Sprintf("%s/api/queries/%d", baseUrl, query.Id)
	body := map[string]interface{}{
		"id":   query.Id,
		"name": query.Name,
		"tags": query.Tags,
	}
	_, err := post(url, body)
	if err != nil {
		return err
	}

	return nil
}
