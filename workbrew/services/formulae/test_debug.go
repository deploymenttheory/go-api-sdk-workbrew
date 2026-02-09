package formulae

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/client"
	"github.com/deploymenttheory/go-api-sdk-workbrew/workbrew/services/formulae/mocks"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestDebugURLs(t *testing.T) {
	logger := zap.NewNop()

	httpClient, err := client.NewClient("test-api-key", "test-workspace",
		client.WithLogger(logger),
		client.WithBaseURL("https://console.workbrew.com"),
	)
	require.NoError(t, err)

	httpmock.ActivateNonDefault(httpClient.GetHTTPClient().Client())
	defer httpmock.DeactivateAndReset()

	// Register a catch-all to see what URL is being called
	httpmock.RegisterNoResponder(func(req *http.Request) (*http.Response, error) {
		fmt.Printf("UNMATCHED URL: %s\n", req.URL.String())
		return httpmock.NewStringResponse(404, "not found"), nil
	})

	mockHandler := &mocks.FormulaeMock{}
	mockHandler.RegisterMocks("https://console.workbrew.com/workspaces/test-workspace")

	service := NewService(httpClient)
	ctx := context.Background()
	_, err = service.GetFormulae(ctx)
	
	fmt.Printf("Error: %v\n", err)
	fmt.Printf("Total calls: %d\n", httpmock.GetTotalCallCount())
}
