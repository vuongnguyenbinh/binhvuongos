package drive

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Config holds Google Drive OAuth2 credentials
type Config struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	FolderID     string
}

// UploadResult contains the uploaded file info
type UploadResult struct {
	FileID   string `json:"id"`
	FileName string `json:"name"`
	WebViewLink string `json:"webViewLink"`
}

// getAccessToken exchanges refresh token for access token
func getAccessToken(cfg *Config) (string, error) {
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("refresh_token", cfg.RefreshToken)
	data.Set("grant_type", "refresh_token")

	resp, err := http.PostForm("https://oauth2.googleapis.com/token", data)
	if err != nil {
		return "", fmt.Errorf("token request: %w", err)
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("token decode: %w", err)
	}
	if result.Error != "" {
		return "", fmt.Errorf("token error: %s", result.Error)
	}
	return result.AccessToken, nil
}

// UploadFile uploads a file to Google Drive folder
func UploadFile(ctx context.Context, cfg *Config, fileName string, mimeType string, fileData io.Reader) (*UploadResult, error) {
	token, err := getAccessToken(cfg)
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	// Create multipart upload
	// Step 1: Create file metadata
	metadata := fmt.Sprintf(`{"name": "%s", "parents": ["%s"]}`, fileName, cfg.FolderID)

	// Use resumable upload for simplicity with multipart
	// Simple upload with metadata via multipart related
	boundary := "bvos_upload_boundary"
	body := fmt.Sprintf(
		"--%s\r\nContent-Type: application/json; charset=UTF-8\r\n\r\n%s\r\n--%s\r\nContent-Type: %s\r\n\r\n",
		boundary, metadata, boundary, mimeType)

	// Read file content
	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	fullBody := body + string(fileBytes) + "\r\n--" + boundary + "--"

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://www.googleapis.com/upload/drive/v3/files?uploadType=multipart&fields=id,name,webViewLink",
		strings.NewReader(fullBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "multipart/related; boundary="+boundary)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upload request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var result UploadResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result.FileName = fileName
	return &result, nil
}
