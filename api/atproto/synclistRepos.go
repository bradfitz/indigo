package atproto

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// schema: com.atproto.sync.listRepos

func init() {
}

type SyncListRepos_Output struct {
	Cursor *string               `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
	Repos  []*SyncListRepos_Repo `json:"repos" cborgen:"repos"`
}

type SyncListRepos_Repo struct {
	Did  string `json:"did" cborgen:"did"`
	Head string `json:"head" cborgen:"head"`
}

func SyncListRepos(ctx context.Context, c *xrpc.Client, cursor string, limit int64) (*SyncListRepos_Output, error) {
	var out SyncListRepos_Output

	params := map[string]interface{}{
		"cursor": cursor,
		"limit":  limit,
	}
	if err := c.Do(ctx, xrpc.Query, "", "com.atproto.sync.listRepos", params, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
