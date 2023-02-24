package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	tfBaseURL      = "yourUrl"
	tfContentType  = "application/vnd.api+json"
	tfOrganization = "yourOrganization"
	tfWorkspace    = "yourWorkspace"
)

var (
	tfToken = os.Getenv("TF_TOKEN")
	tfTarGz = os.Args[1]
)

func respCloser(respBody io.Closer) {
	if err := respBody.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getWorkspaceId() string {
	var jsonWorkspaceId any

	client := http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/api/v2/organizations/%s/workspaces/%s", tfBaseURL, tfOrganization, tfWorkspace), nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tfToken))
	req.Header.Add("Content-Type", tfContentType)
	resp, _ := client.Do(req)
	defer respCloser(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &jsonWorkspaceId)
	return fmt.Sprint(jsonWorkspaceId.(map[string]any)["data"].(map[string]any)["id"])
}

func getUploadUrl(workspaceId string) string {
	var jsonUploadUrl map[string]any
	jsonMap := map[string]map[string]string{"data": {"type": "configuration-versions"}}
	jsonBytes, _ := json.Marshal(jsonMap)

	client := http.Client{}
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/api/v2/workspaces/%s/configuration-versions", tfBaseURL, workspaceId), bytes.NewBuffer(jsonBytes))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tfToken))
	req.Header.Add("Content-Type", tfContentType)
	resp, _ := client.Do(req)
	defer respCloser(resp.Body)

	bodyBytes, _ := io.ReadAll(resp.Body)
	json.Unmarshal(bodyBytes, &jsonUploadUrl)
	return fmt.Sprint(jsonUploadUrl["data"].(map[string]any)["attributes"].(map[string]any)["upload-url"])
}

func uploadConfig(uploadUrl string) {
	tfPack, _ := os.Open(tfTarGz)
	defer tfPack.Close()
	client := http.Client{}
	req, _ := http.NewRequest("PUT", uploadUrl, tfPack)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", tfToken))
	req.Header.Add("Content-Type", "application/octet-stream")
	resp, _ := client.Do(req)
	defer respCloser(resp.Body)
	fmt.Println(resp.Status)
}

func main() {
	wsID := getWorkspaceId()
	uploadURL := getUploadUrl(wsID)
	uploadConfig(uploadURL)
}
