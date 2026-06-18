package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"strings"

	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

type LibraryClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewLibraryClient() *LibraryClient {
	return &LibraryClient{
		baseURL:    os.Getenv("LIBRARY_SERVICE_URL"),
		httpClient: &http.Client{},
	}
}

func (c *LibraryClient) ValidateBook(ctx context.Context, bookID int) (*model.Book, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/books/%d", c.baseURL, bookID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrBookNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrLibraryServiceUnavailable
	}

	var book model.Book
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		return nil, err
	}

	if book.AvailableCopies == 0 {
		return nil, ErrNoCopiesAvailable
	}

	return &book, nil
}

func (c *LibraryClient) UpdateCopies(ctx context.Context, bookID int, action string) error {
	body := fmt.Sprintf(`{"action":"%s"}`, action)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/books/%d/copies", c.baseURL, bookID),
		strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ErrLibraryServiceUnavailable
	}

	return nil
}
