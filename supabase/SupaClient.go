package supabase

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type SupaClient struct {
	Client  *http.Client
	BaseUrl string
	AnonKey string
	Config  SupaClientConfig
}

type SupaClientConfig struct {
	Environment string
	AuthToken   string
}

func getExecutablePath() string {
	exePath, err := os.Executable()
	if err != nil {
		// fallback to current dir if somehow can't determine exec path
		return "."
	}
	exeDir := exePath
	// If exePath is not a directory, get its directory
	if fi, err := os.Stat(exePath); err == nil && !fi.IsDir() {
		exeDir = strings.TrimSuffix(exePath, "/"+fi.Name())
	} else if err == nil && fi.IsDir() {
		exeDir = exePath
	} else {
		exeDir = "."
	}
	return exeDir
}

var CONFIG_FILE = getExecutablePath() + string(os.PathSeparator) + ".buildium" + string(os.PathSeparator) + "config.json"

func getDefaultConfig() SupaClientConfig {
	fmt.Println("No config file found, using default config be sure to run `buildium login` to set your credentials")
	return SupaClientConfig{Environment: "PROD", AuthToken: ""}
}

func NewSupaClient(ctx context.Context) *SupaClient {
	config, err := os.ReadFile(CONFIG_FILE)
	var loadedConfig SupaClientConfig
	if err != nil {
		loadedConfig = getDefaultConfig()
	} else {
		err = json.Unmarshal(config, &loadedConfig)
		if err != nil {
			fmt.Printf("Failed to unmarshal config file: %v\n", err)
			loadedConfig = getDefaultConfig()
		}
	}

	switch loadedConfig.Environment {
	case "PROD":
		return &SupaClient{Client: http.DefaultClient, BaseUrl: "https://dpwumtpjesedslulexqz.supabase.co", AnonKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6ImRwd3VtdHBqZXNlZHNsdWxleHF6Iiwicm9sZSI6ImFub24iLCJpYXQiOjE3NjcxNDE4MzksImV4cCI6MjA4MjcxNzgzOX0.JYXW1bzTOmlCtngrlYLAbnGzRXDIcH0mDlwpbg1u8Rs", Config: loadedConfig}
	case "BUILDING":
		return &SupaClient{Client: http.DefaultClient, BaseUrl: "", AnonKey: "", Config: loadedConfig}
	default:
		return &SupaClient{Client: http.DefaultClient, BaseUrl: "http://127.0.0.1:54321", AnonKey: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6ImFub24iLCJleHAiOjE5ODM4MTI5OTZ9.CRXP1A7WOeoJeXxjNni43kdQwgnWNReilDMblYTn_I0", Config: loadedConfig}
	}
}

func (c *SupaClient) VerifyAuthToken(ctx context.Context) error {
	if c.Config.AuthToken == "" {
		return fmt.Errorf("not logged in")
	}

	// TODO: Implement verify auth token logic
	return nil
}

type Project struct {
	ProjectId  string `json:"projectId"`
	Name       string `json:"name"`
	TutorialId string `json:"tutorialId"`
}

func (c *SupaClient) GetProject(ctx context.Context, projectId string) (Project, error) {
	if c.BaseUrl == "" {
		return Project{}, fmt.Errorf("not logged in")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseUrl+"/functions/v1/get-project",
		strings.NewReader(fmt.Sprintf(`{"projectId":"%s"}`, projectId)))
	if err != nil {
		return Project{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AnonKey)
	req.Header.Set("x-buildium-token", c.Config.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return Project{}, err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Project{}, err
		}
		return Project{}, fmt.Errorf("failed to get project: %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Project{}, err
	}
	var project Project
	err = json.Unmarshal(body, &project)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

type Tutorial struct {
	TutorialId  string `json:"tutorialId"`
	Name        string `json:"name"`
	Stages      int    `json:"stages"`
	DockerImage string `json:"dockerImage"`
}

func (c *SupaClient) GetTutorial(ctx context.Context, tutorialId string) (Tutorial, error) {

	if c.BaseUrl == "" {
		return Tutorial{}, fmt.Errorf("not logged in")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseUrl+"/functions/v1/get-tutorial",
		strings.NewReader(fmt.Sprintf(`{"tutorialId":"%s"}`, tutorialId)))
	if err != nil {
		return Tutorial{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.AnonKey)
	req.Header.Set("x-buildium-token", c.Config.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return Tutorial{}, err
	}
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return Tutorial{}, err
		}
		return Tutorial{}, fmt.Errorf("failed to get tutorial: %d: %s", resp.StatusCode, string(body))
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Tutorial{}, err
	}
	var tutorial Tutorial
	err = json.Unmarshal(body, &tutorial)
	if err != nil {
		return Tutorial{}, err
	}
	return tutorial, nil
}

func (c *SupaClient) Login(ctx context.Context, email string, password string) error {
	if c.BaseUrl == "" {
		return nil
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseUrl+"/functions/v1/login",
		strings.NewReader(fmt.Sprintf(`{"email":"%s", "password":"%s"}`, email, password)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.AnonKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to login: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var token struct {
		Token string `json:"token"`
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return err
	}
	c.Config.AuthToken = token.Token
	return c.storeConfig()
}

func (c *SupaClient) storeConfig() error {
	config, err := json.Marshal(c.Config)
	if err != nil {
		return err
	}

	dir := getExecutablePath() + string(os.PathSeparator) + ".buildium"
	if _, statErr := os.Stat(dir); os.IsNotExist(statErr) {
		mkdirErr := os.MkdirAll(dir, 0700)
		if mkdirErr != nil {
			return fmt.Errorf("failed to create config directory: %v", mkdirErr)
		}
	}

	err = os.WriteFile(CONFIG_FILE, config, 0600)
	if err != nil {
		return err
	}
	return nil
}
