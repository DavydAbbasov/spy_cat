package catapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

func (c *Client) IsValid(ctx context.Context, breed string) (bool, error) {
	q := normBreed(breed)
	if q == "" {
		return false, nil
	}

	breeds, status, err := c.SearchBreeds(ctx, q)
	if err != nil {
		return false, err
	}
	if status != http.StatusOK {
		return false, fmt.Errorf("catapi status=%d", status)
	}
	return len(breeds) > 0, nil
}
func normBreed(b string) string {
	b = strings.TrimSpace(strings.ToLower(b))
	return strings.ReplaceAll(b, "_", " ")
}
