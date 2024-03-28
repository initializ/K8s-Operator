package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"net/http"
)

const USER_AGENT_TYPE = "k8-operator"

// this function is used to get the key of the organisation
func CallGetEncryptedWorkspaceKey(httpClient *resty.Client, serviceToken string, request GetEncryptedWorkspaceKeyRequest) (GetEncryptedWorkspaceKeyResponse, error) {
	endpoint := API_HOST_URL + "getkey"
	var result GetEncryptedWorkspaceKeyResponse
	response, err := httpClient.
		R().
		SetResult(&result).
		SetHeader("org_id", request.OrgID). // Set the content type header
		SetHeader("User-Agent", USER_AGENT_TYPE).
		SetHeader("App-User", "testsecert").
		SetHeader("Authorization", "Bearer "+serviceToken). // Set the Authorization header with the bearer token
		Get(endpoint)
	if err != nil {
		return GetEncryptedWorkspaceKeyResponse{}, fmt.Errorf("CallGetEncryptedWorkspaceKey: Unable to complete API request [err=%s]", err)
	}
	if response.IsError() {
		return GetEncryptedWorkspaceKeyResponse{}, fmt.Errorf("CallGetEncryptedWorkspaceKey: Unsuccessful response: [response=%s]", response)
	}
	return result, nil
}

// CallGetEncryptedSecrets makes the API call to fetch encrypted secrets.
// CallGetEncryptedSecrets makes the API call to fetch encrypted secrets.
func CallGetEncryptedSecrets(request GetEncryptedSecretsRequest, serviceToken string) (GetEncryptedSecretResponse, error) {
	endpoint := API_HOST_URL + "getsecrets"
	var result GetEncryptedSecretResponse

	requestBody, err := json.Marshal(request)
	if err != nil {
		return GetEncryptedSecretResponse{}, fmt.Errorf("CallGetEncryptedSecrets: Unable to marshal request: %v", err)
	}
	req, _ := http.NewRequest("GET", endpoint, bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+serviceToken)
	req.Header.Set("User-Agent", USER_AGENT_TYPE)
	req.Header.Set("App-User", "testsecert")
	req.Header.Set("org_id", request.OrgId)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return GetEncryptedSecretResponse{}, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return GetEncryptedSecretResponse{}, err
	}

	// // Print the response status code and body
	// fmt.Println("Response Status:", resp.Status)
	// fmt.Println("Response Body:", string(body))

	// Unmarshal response body into result
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error unmarshaling response body:", err)
		return GetEncryptedSecretResponse{}, err
	}

	return result, nil
}
