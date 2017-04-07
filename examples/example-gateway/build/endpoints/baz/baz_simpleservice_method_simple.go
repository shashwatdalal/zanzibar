// Code generated by zanzibar
// @generated

package baz

import (
	"context"

	"github.com/uber-go/zap"
	"github.com/uber/zanzibar/examples/example-gateway/build/clients"
	zanzibar "github.com/uber/zanzibar/runtime"

	customBaz "github.com/uber/zanzibar/examples/example-gateway/endpoints/baz"
)

// HandleSimpleRequest handles "/baz/simple-path".
func HandleSimpleRequest(
	ctx context.Context,
	req *zanzibar.ServerHTTPRequest,
	res *zanzibar.ServerHTTPResponse,
	clients *clients.Clients,
) {

	headers := map[string]string{}

	workflow := customBaz.SimpleEndpoint{
		Clients: clients,
		Logger:  req.Logger,
		Request: req,
	}

	_, err := workflow.Handle(ctx, headers)
	if err != nil {
		req.Logger.Warn("Workflow for endpoint returned error",
			zap.String("error", err.Error()),
		)
		res.SendErrorString(500, "Unexpected server error")
		return
	}

	res.WriteJSONBytes(204, nil)
}
