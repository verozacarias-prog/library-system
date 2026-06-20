package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/verozacarias-prog/library-system/loans-service/internal/model"
)

type LibraryClient struct {
	baseURL    string
	httpClient *http.Client
	jwtSecret  []byte
}

func NewLibraryClient() *LibraryClient {
	return &LibraryClient{
		baseURL:    os.Getenv("LIBRARY_SERVICE_URL"),
		httpClient: &http.Client{Timeout: 5 * time.Second},
		jwtSecret:  []byte(os.Getenv("JWT_SECRET")),
	}
}

func (c *LibraryClient) serviceToken() (string, error) {
	claims := jwt.MapClaims{
		"role": "service",
		"exp":  time.Now().Add(time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(c.jwtSecret)
}

func (c *LibraryClient) ValidateBook(ctx context.Context, bookID int) (*model.Book, error) {
	token, err := c.serviceToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/books/%d", c.baseURL, bookID), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrLibraryServiceUnavailable, err)
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
	token, err := c.serviceToken()
	if err != nil {
		return err
	}

	delta := -1
	if action == "increment" {
		delta = 1
	}
	body := fmt.Sprintf(`{"delta":%d}`, delta)
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch,
		fmt.Sprintf("%s/books/%d/copies", c.baseURL, bookID),
		strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrLibraryServiceUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusConflict {
		return ErrNoCopiesAvailable
	}

	if resp.StatusCode != http.StatusOK {
		return ErrLibraryServiceUnavailable
	}

	return nil
}
