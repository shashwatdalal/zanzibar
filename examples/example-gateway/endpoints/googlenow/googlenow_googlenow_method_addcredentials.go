// Code generated by zanzibar
// @generated

package googlenow

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/uber-go/zap"
	"github.com/uber/zanzibar/examples/example-gateway/clients"
	zanzibar "github.com/uber/zanzibar/runtime"

	"github.com/uber/zanzibar/examples/example-gateway/clients/googlenow"
)

// HandleAddCredentialsRequest handles "/googlenow/add-credentials".
func HandleAddCredentialsRequest(
	ctx context.Context,
	inc *zanzibar.IncomingMessage,
	g *zanzibar.Gateway,
	clients *clients.Clients,
) {
	// Handle request headers.
	h := http.Header{}
	for _, header := range []string{"x-uuid", "x-token"} {
		h.Set(header, inc.Header.Get(header))
	}

	// Handle request body.
	rawBody, ok := inc.ReadAll()
	if !ok {
		return
	}
	var body AddCredentialsHTTPRequest
	if ok := inc.UnmarshalBody(&body, rawBody); !ok {
		return
	}
	clientRequest := convertToAddCredentialsClientRequest(&body)
	clientResp, err := clients.GoogleNow.AddCredentials(ctx, clientRequest, h)
	if err != nil {
		g.Logger.Error("Could not make client request",
			zap.String("error", err.Error()),
		)
		inc.SendError(500, errors.Wrap(err, "could not make client request:"))
		return
	}

	defer func() {
		if cerr := clientResp.Body.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Handle client respnse.
	if !inc.IsOKResponse(clientResp.StatusCode, []int{200, 202}) {
		g.Logger.Warn("Unknown response status code",
			zap.Int("status code", clientResp.StatusCode),
		)
	}
	inc.WriteJSONBytes(clientResp.StatusCode, nil)
}

func convertToAddCredentialsClientRequest(body *AddCredentialsHTTPRequest) *googlenowClient.AddCredentialsHTTPRequest {
	clientRequest := &googlenowClient.AddCredentialsHTTPRequest{}

	clientRequest.AuthCode = string(body.AuthCode)

	return clientRequest
}
