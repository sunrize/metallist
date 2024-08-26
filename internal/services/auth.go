package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"metallist/internal/urlhelper"
)

type AuthService struct {
	Name     string
	Config   *oauth2.Config
	Verifier string
}

func (s AuthService) LoginPathURL() string {
	return fmt.Sprintf("/auth/%v/login", s.Name)
}

func (s AuthService) CallbackPathURL() string {
	return fmt.Sprintf("/auth/%v/callback", s.Name)
}

func AuthServices() []AuthService {
	var authServices = []AuthService{
		{
			Name: "anilist",
			Config: &oauth2.Config{
				ClientID:     "",
				ClientSecret: "",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "https://anilist.co/api/v2/oauth/authorize",
					TokenURL:  "https://anilist.co/api/v2/oauth/token",
					AuthStyle: oauth2.AuthStyleInParams,
				},
				RedirectURL: "http://localhost:1212/auth/anilist/callback",
			},
			Verifier: oauth2.GenerateVerifier(),
		},
		{
			Name: "myanimelist",
			Config: &oauth2.Config{
				ClientID:     "",
				ClientSecret: "",
				Endpoint: oauth2.Endpoint{
					AuthURL:   "https://myanimelist.net/v1/oauth2/authorize",
					TokenURL:  "https://myanimelist.net/v1/oauth2/token",
					AuthStyle: oauth2.AuthStyleInParams,
				},
				RedirectURL: "http://localhost:1212/auth/myanimelist/callback",
			},
			Verifier: oauth2.GenerateVerifier(),
		},
	}

	// Reading the configuration file
	viper.SetConfigFile(fmt.Sprintf("/config/%v", "secrets.yaml"))
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("Fatal error config file: %s \n", err))
	}

	// Filling in the configuration for each service
	for i, service := range authServices {
		service.Config.ClientID = viper.GetString(service.Name + ".client_id")
		service.Config.ClientSecret = viper.GetString(service.Name + ".client_secret")
		authServices[i].Config = service.Config
	}

	return authServices
}

func LoginHandler(service AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Setting RedirectUrl to the current server URL
		service.Config.RedirectURL = urlhelper.GetFullURLOverridePath(r, service.CallbackPathURL())

		url := service.Config.AuthCodeURL("state", oauth2.SetAuthURLParam("code_challenge", service.Verifier))
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}
}

func CallbackHandler(service AuthService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		state := r.URL.Query().Get("state")
		if state != "state" {
			http.Error(w, "Invalid state parameter", http.StatusBadRequest)
			return
		}

		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code parameter", http.StatusBadRequest)
			return
		}

		ctx := context.Background()

		token, err := service.Config.Exchange(ctx, code, oauth2.VerifierOption(service.Verifier))
		if err != nil {
			log.Println("Error exchanging code for token:", err)
			http.Error(w, "Error exchanging code for token", http.StatusInternalServerError)
			return
		}

		// Saving the token to the cache
		err = saveTokensToCache(service, token)
		if err != nil {
			log.Println("Error saving tokens to cache:", err)
			http.Error(w, "Error saving tokens to cache", http.StatusInternalServerError)
			return
		}

		// After successfully saving the token, we send the user back to the ~~original~~ page
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}
}

// saveTokensToCache saves tokens to a cache file
func saveTokensToCache(service AuthService, tokens *oauth2.Token) error {
	cacheFilePath := filepath.Join("/config", "user", "auth", service.Name+".yaml")

	// Checking if there is a path to the cache file
	if _, err := os.Stat(filepath.Dir(cacheFilePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(cacheFilePath), 0644); err != nil {
			fmt.Println("Error creating directories:", err)
			return err
		}
	}

	// Saving tokens in JSON
	data, err := json.Marshal(tokens)
	if err != nil {
		return err
	}

	// Opening the file for writing
	file, err := os.OpenFile(cacheFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing data to a file
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

// loadTokensFromCache loads tokens from the cache file
func loadTokensFromCache(service AuthService) (*oauth2.Token, error) {
	cacheFilePath := filepath.Join("/config", "user", "auth", service.Name+".yaml")
	// Checking if the cache file exists
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		return nil, errors.New("The cache file could not be found.")
	}

	// Opening the file for reading
	file, err := os.Open(cacheFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Reading data from a file
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Deserializing tokens
	var tokens oauth2.Token
	if err := json.Unmarshal(data, &tokens); err != nil {
		return nil, err
	}

	return &tokens, nil
}

// refreshTokens updates tokens if they are outdated
func refreshTokens(service AuthService) (*oauth2.Token, error) {
	// Loading tokens from the cache
	tokens, err := loadTokensFromCache(service)
	if err != nil {
		return nil, err
	}

	// Checking whether the tokens need to be updated
	if time.Now().After(tokens.Expiry) {
		// Updating tokens
		newTokens, err := service.Config.TokenSource(context.Background(), tokens).Token()
		if err != nil {
			return nil, err
		}

		// Saving new tokens to the cache
		if err := saveTokensToCache(service, newTokens); err != nil {
			return nil, err
		}

		return newTokens, nil
	}

	return tokens, nil
}

// getAuthenticatedClient returns the http client with authentication using tokens from the cache
func GetAuthenticatedClient(service AuthService) (*http.Client, error) {
	// Loading tokens from the cache
	tokens, err := refreshTokens(service)
	if err != nil {
		panic(fmt.Sprintf("Fatal error cache file: %s \n", err))
	}

	// We return the http client with authentication
	return service.Config.Client(context.Background(), tokens), nil
}

func TestSaveTokens() error {
	token := new(oauth2.Token)
	service := AuthService{
		Name: "test_service",
		Config: &oauth2.Config{
			ClientID:     "",
			ClientSecret: "",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://example.com/oauth/authorize",
				TokenURL:  "https://example.com/oauth/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
			// RedirectURL: "http://localhost:1212/auth/anilist/callback",
		},
		Verifier: "",
	}

	err := saveTokensToCache(service, token)
	if err != nil {
		return err
	}

	return nil
}

func TestLoadTokens() error {
	service := AuthService{
		Name: "test_service",
		Config: &oauth2.Config{
			ClientID:     "",
			ClientSecret: "",
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://example.com/oauth/authorize",
				TokenURL:  "https://example.com/oauth/token",
				AuthStyle: oauth2.AuthStyleInParams,
			},
			// RedirectURL: "http://localhost:1212/auth/anilist/callback",
		},
		Verifier: "",
	}

	_, err := loadTokensFromCache(service)
	if err != nil {
		return err
	}

	return nil
}
