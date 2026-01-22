package books

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type googleBooksAPIResponse struct {
	Items []googleVolume `json:"items"`
}

type googleVolume struct {
	VolumeInfo googleVolumeInfo `json:"volumeInfo"`
}

type googleVolumeInfo struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
}

type GoogleBooksProvider struct {
	httpClient *http.Client
}

func NewGoogleBooksProvider(client *http.Client) *GoogleBooksProvider {
	return &GoogleBooksProvider{
		httpClient: client,
	}
}

func (p *GoogleBooksProvider) SearchBooks(query string, ctx context.Context) ([]externalBookModel, error) {
	googleBooks, err := p.queryAPI(query, ctx)
	if err != nil {
		return nil, err
	}
	externalBooks := mapGoogleToExternal(googleBooks)
	return externalBooks, nil
}

func (p *GoogleBooksProvider) queryAPI(query string, ctx context.Context) ([]googleVolume, error) {
	baseURL := "https://www.googleapis.com/books/v1/volumes?q="

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL+query,
		nil,
	)
	if err != nil {
		return nil, err
	}

	// TODO: consider making configurable
	q := req.URL.Query()
	q.Set("q", query)
	q.Set("maxResults", "40")
	q.Set("startIndex", "0")

	req.URL.RawQuery = q.Encode()

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("google books api returned status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse googleBooksAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, err
	}

	return apiResponse.Items, nil
}

func mapGoogleToExternal(googleBooks []googleVolume) []externalBookModel {
	books := make([]externalBookModel, 0, len(googleBooks))

	for _, g := range googleBooks {
		books = append(books, externalBookModel{
			Title:   g.VolumeInfo.Title,
			Authors: g.VolumeInfo.Authors,
		})
	}

	return books
}
