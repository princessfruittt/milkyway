package getrole

import (
	"context"
	"github.com/google/go-github/v33/github"
)
func start() {
	client := github.NewClient(nil)
	ctx := context.Background()

	// list all organizations for user "willnorris"
	opt, _, _ := client.Organizations.List(context.Background(), "willnorris", nil)
}
