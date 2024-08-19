package anilist

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

// Implement OAuth: https://anilist.gitbook.io/anilist-apiv2-docs/overview/oauth/getting-started
// type authedTransport struct {
// 	wrapped http.RoundTripper
// }

// func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
// 	key := "1"
// 	req.Header.Set("Authorization", "bearer "+key)
// 	return t.wrapped.RoundTrip(req)
// }

// func TestAuthRun() {
// 	ctx := context.Background()
// 	client := graphql.NewClient("https://api.github.com/graphql",
// 		&http.Client{Transport: &authedTransport{wrapped: http.DefaultTransport}})
// 	resp, err := getUserID(ctx, client, "blackrize")
// 	fmt.Println(resp, err)
// }

func TestRun() {
	ctx := context.Background()
	client := graphql.NewClient("https://graphql.anilist.co", http.DefaultClient)
	resp, err := getUserID(ctx, client, "blackrize")
	fmt.Println(resp, err)
}
