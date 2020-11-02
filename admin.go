package rclient

import (
	"context"
	"github.com/khorevaa/ras-client/messages"
)

func (c *Client) GetAgentVersion(ctx context.Context) (string, error) {

	req := &messages.GetAgentVersionRequest{}

	resp, err := c.sendEndpointRequest(ctx, req)

	if err != nil {
		return "", err
	}

	response := resp.(*messages.GetAgentVersionResponse)

	return response.Version, err
}
