package internal

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	BaseURL    = "https://api.github.com/repos/cat-milk/Anime-Girls-Holding-Programming-Books/contents"
	RawBaseURL = "https://raw.githubusercontent.com/cat-milk/Anime-Girls-Holding-Programming-Books/master"
)

type GitHubContent struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	Type        string `json:"type"`
	DownloadURL string `json:"download_url"`
	URL         string `json:"url"`
}

type Fetcher struct {
	client *http.Client
}

type authTransport struct {
	Token     string
	Transport http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	// Set the Authorization header with the Bearer token
	req.Header.Set("Authorization", "Bearer "+t.Token)
	return t.Transport.RoundTrip(req)
}

func NewFetcher() *Fetcher {

	// Create a new HTTP client with a timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Add authentication if token is available
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		client.Transport = &authTransport{
			Token:     token,
			Transport: http.DefaultTransport,
		}
	}

	return &Fetcher{client: client}
}

func (f *Fetcher) GetLanguageFolders() ([]string, error) {

	// Fetch the contents of the repository
	resp, err := f.client.Get(BaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repository contents: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	// Decode the JSON response
	var contents []GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Extract language folders
	var languages []string
	for _, content := range contents {
		if content.Type == "dir" && !strings.HasPrefix(content.Name, ".") {
			languages = append(languages, content.Name)
		}
	}

	if len(languages) == 0 {
		return nil, fmt.Errorf("no language folders found in repository")
	}

	return languages, nil
}

func (f *Fetcher) GetImagesInFolder(language string) ([]string, error) {

	// Validate language input
	url := fmt.Sprintf("%s/%s", BaseURL, language)

	// Validate language folder
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch folder contents for %s: %w", language, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("language folder '%s' not found", language)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d for folder %s", resp.StatusCode, language)
	}

	// Decode the JSON response
	var contents []GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&contents); err != nil {
		return nil, fmt.Errorf("failed to decode response for folder %s: %w", language, err)
	}

	// Extract image files
	var images []string
	supportedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	// Check each content item for supported image extensions
	for _, content := range contents {
		if content.Type == "file" {
			ext := strings.ToLower(filepath.Ext(content.Name))
			if supportedExtensions[ext] {
				images = append(images, content.Name)
			}
		}
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("no images found in folder '%s'", language)
	}

	return images, nil
}

func (f *Fetcher) DownloadImage(language, filename string) (string, error) {

	// Validate filename
	url := fmt.Sprintf("%s/%s/%s", RawBaseURL, language, filename)

	// Make HTTP GET request to download the image
	resp, err := f.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download image, status: %d", resp.StatusCode)
	}

	// Create temp file
	tempDir := os.TempDir()
	ext := filepath.Ext(filename)
	tempFile, err := os.CreateTemp(tempDir, "waifetch-*"+ext)
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tempFile.Close()

	// Copy image data to temp file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write image data: %w", err)
	}

	return tempFile.Name(), nil
}

func (f *Fetcher) DisplayImage(imagePath string) error {

	// Check if chafa is installed
	_, err := exec.LookPath("chafa")
	if err != nil {
		return fmt.Errorf("chafa is not installed or not in PATH. Please install chafa to display images")
	}

	// Get terminal size
	width, height, err := term.GetSize(int(os.Stdin.Fd()))
	if err != nil {
		width, height, err = term.GetSize(int(os.Stdout.Fd()))
		if err != nil {
			// Fallback to default size
			width, height = 80, 40
		}
	}

	// Calculate display size with padding
	displayWidth := width - 4
	displayHeight := height - 6

	// Ensure reasonable minimum and maximum sizes
	if displayWidth < 20 {
		displayWidth = 20
	} else if displayWidth > 300 {
		displayWidth = 300
	}

	if displayHeight < 10 {
		displayHeight = 10
	} else if displayHeight > 200 {
		displayHeight = 200
	}

	sizeStr := fmt.Sprintf("%dx%d", displayWidth, displayHeight)

	// chafa command with options
	cmd := exec.Command("chafa",
		"--size", sizeStr,
		"--symbols", "block",
		"--colors", "256",
		"--dither", "ordered",
		imagePath)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to display image with chafa: %w", err)
	}

	return nil
}

func (f *Fetcher) FetchRandomImage(language string) error {

	var languages []string
	var err error

	// Seed the random number generator
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// If a language is specified, validate it
	if language != "" {
		languages, err = f.GetLanguageFolders()
		if err != nil {
			return fmt.Errorf("failed to get language folders: %w", err)
		}

		found := false
		for _, lang := range languages {
			if strings.EqualFold(lang, language) {
				language = lang
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("language '%s' not found. Available languages: %s",
				language, strings.Join(languages, ", "))
		}

	} else { // If no language is specified, select a random one
		languages, err = f.GetLanguageFolders()
		if err != nil {
			return fmt.Errorf("failed to get language folders: %w", err)
		}

		language = languages[r.Intn(len(languages))]
	}

	// Get images in the selected language folder
	images, err := f.GetImagesInFolder(language)
	if err != nil {
		return err
	}

	selectedImage := images[r.Intn(len(images))]

	// Download the selected image
	tempPath, err := f.DownloadImage(language, selectedImage)
	if err != nil {
		return err
	}

	// Ensure the temporary file is cleaned up after use
	defer func() {
		if err := os.Remove(tempPath); err != nil {
			fmt.Printf("Warning: failed to clean up temporary file: %v\n", err)
		}
	}()

	fmt.Println()
	return f.DisplayImage(tempPath)
}
