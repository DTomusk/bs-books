package extraction

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
	Title      string   `json:"title"`
	Authors    []string `json:"authors"`
	ImageLinks struct {
		Thumbnail string `json:"thumbnail"`
	} `json:"imageLinks"`
	Description string `json:"description"`
}

type GoogleBooksProvider struct {
	httpClient *http.Client
	apiKey     string
}

func NewGoogleBooksProvider(client *http.Client, apiKey string) *GoogleBooksProvider {
	return &GoogleBooksProvider{
		httpClient: client,
		apiKey:     apiKey,
	}
}

func (p *GoogleBooksProvider) SearchBooks(query string, maxResults int, ctx context.Context) ([]externalBookModel, error) {
	googleBooks, err := p.queryAPI(query, maxResults, ctx)
	if err != nil {
		return nil, err
	}
	externalBooks := mapGoogleToExternal(googleBooks)
	return externalBooks, nil
}

func (p *GoogleBooksProvider) queryAPI(query string, maxResults int, ctx context.Context) ([]googleVolume, error) {
	baseURL := "https://www.googleapis.com/books/v1/volumes?q="

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		baseURL,
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "bs-books-api/1.0")

	q := req.URL.Query()
	q.Set("q", query)
	q.Set("maxResults", fmt.Sprintf("%d", maxResults))
	q.Set("startIndex", "0")
	q.Set("key", p.apiKey)

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
			Title:    g.VolumeInfo.Title,
			Authors:  g.VolumeInfo.Authors,
			ImageURL: g.VolumeInfo.ImageLinks.Thumbnail,
			Synopsis: g.VolumeInfo.Description,
		})
	}

	return books
}
