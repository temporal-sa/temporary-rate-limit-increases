package capacity

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	defaultBaseURL = "https://saas-api.tmprl.cloud"
)

// ProvisionInput holds the parameters for a capacity provisioning request.
type ProvisionInput struct {
	Namespace string
	APSLimit  int32
}

// Activities holds the dependencies required by the provisioning activities.
type Activities struct {
	HTTPClient *http.Client
	BaseURL    string
	APIKey     string
}

type getNamespaceResponse struct {
	Namespace struct {
		Spec            json.RawMessage `json:"spec"`
		ResourceVersion string          `json:"resourceVersion"`
	} `json:"namespace"`
}

type updateNamespaceRequest struct {
	Spec            json.RawMessage `json:"spec"`
	ResourceVersion string          `json:"resourceVersion"`
}

func (a *Activities) baseURL() string {
	if a.BaseURL != "" {
		return a.BaseURL
	}
	return defaultBaseURL
}

func (a *Activities) httpClient() *http.Client {
	if a.HTTPClient != nil {
		return a.HTTPClient
	}
	return http.DefaultClient
}

func (a *Activities) getNamespace(ctx context.Context, namespace string) (json.RawMessage, string, error) {
	url := fmt.Sprintf("%s/cloud/namespaces/%s", a.baseURL(), namespace)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+a.APIKey)

	resp, err := a.httpClient().Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}

	var result getNamespaceResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}
	return result.Namespace.Spec, result.Namespace.ResourceVersion, nil
}

func (a *Activities) updateNamespace(ctx context.Context, namespace string, spec json.RawMessage, resourceVersion string) error {
	url := fmt.Sprintf("%s/cloud/namespaces/%s", a.baseURL(), namespace)

	payload, err := json.Marshal(updateNamespaceRequest{
		Spec:            spec,
		ResourceVersion: resourceVersion,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := a.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
	}
	return nil
}

// AddTRUs increases the provisioned capacity of the target namespace to the requested APS limit.
func (a *Activities) AddTRUs(ctx context.Context, input ProvisionInput) error {
	spec, resourceVersion, err := a.getNamespace(ctx, input.Namespace)
	if err != nil {
		return fmt.Errorf("get namespace %q: %w", input.Namespace, err)
	}

	var specMap map[string]interface{}
	if err := json.Unmarshal(spec, &specMap); err != nil {
		return fmt.Errorf("unmarshal namespace spec: %w", err)
	}
	specMap["capacitySpec"] = map[string]interface{}{
		"provisioned": map[string]interface{}{
			// 1 TRU = 500 APS; convert the requested APS limit to TRUs.
			"value": float64(input.APSLimit) / 500.0,
		},
	}

	updatedSpec, err := json.Marshal(specMap)
	if err != nil {
		return fmt.Errorf("marshal updated spec: %w", err)
	}

	if err := a.updateNamespace(ctx, input.Namespace, updatedSpec, resourceVersion); err != nil {
		return fmt.Errorf("update namespace %q to %d APS: %w", input.Namespace, input.APSLimit, err)
	}
	return nil
}

// RemoveTRUs reverts the target namespace to on-demand capacity mode.
func (a *Activities) RemoveTRUs(ctx context.Context, input ProvisionInput) error {
	spec, resourceVersion, err := a.getNamespace(ctx, input.Namespace)
	if err != nil {
		return fmt.Errorf("get namespace %q: %w", input.Namespace, err)
	}

	var specMap map[string]interface{}
	if err := json.Unmarshal(spec, &specMap); err != nil {
		return fmt.Errorf("unmarshal namespace spec: %w", err)
	}
	specMap["capacitySpec"] = map[string]interface{}{
		"onDemand": map[string]interface{}{},
	}

	updatedSpec, err := json.Marshal(specMap)
	if err != nil {
		return fmt.Errorf("marshal updated spec: %w", err)
	}

	if err := a.updateNamespace(ctx, input.Namespace, updatedSpec, resourceVersion); err != nil {
		return fmt.Errorf("revert namespace %q to on-demand: %w", input.Namespace, err)
	}
	return nil
}
