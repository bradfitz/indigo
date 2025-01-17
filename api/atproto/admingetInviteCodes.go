package atproto

import (
	"context"

	"github.com/bluesky-social/indigo/xrpc"
)

// schema: com.atproto.admin.getInviteCodes

func init() {
}

type AdminGetInviteCodes_Output struct {
	Codes  []*ServerDefs_InviteCode `json:"codes" cborgen:"codes"`
	Cursor *string                  `json:"cursor,omitempty" cborgen:"cursor,omitempty"`
}

func AdminGetInviteCodes(ctx context.Context, c *xrpc.Client, cursor string, limit int64, sort string) (*AdminGetInviteCodes_Output, error) {
	var out AdminGetInviteCodes_Output

	params := map[string]interface{}{
		"cursor": cursor,
		"limit":  limit,
		"sort":   sort,
	}
	if err := c.Do(ctx, xrpc.Query, "", "com.atproto.admin.getInviteCodes", params, nil, &out); err != nil {
		return nil, err
	}

	return &out, nil
}
