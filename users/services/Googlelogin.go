package services

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config

// oauthStateString  = "shyamalkaushik"
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}
func (hs *HandlerService) HandleMain(c *gin.Context) {
	fmt.Println("called")

	c.HTML(http.StatusOK, "index.html", nil)

}

func (hs *HandlerService) HandleGoogleLogin(c *gin.Context) {
	oauthStateString := generateStateOauthCookie(c)
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (hs *HandlerService) HandleGoogleCallback(c *gin.Context) {
	oauthState, err := c.Cookie("oauthstate")
	if err != nil || c.Query("state") != oauthState {
		log.Println("Invalid oauth state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Code exchange failed: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Printf("Failed getting user info: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	defer response.Body.Close()

	var userInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
	}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	fmt.Println("user deails are", userInfo)
	if err != nil {
		log.Printf("Failed decoding user info: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// For simplicity, we will just print the user info in the response
	c.JSON(http.StatusOK, gin.H{
		"user": userInfo,
	})
}

// generateStateOauthCookie generates a random string for oauth state and sets it as a cookie
func generateStateOauthCookie(c *gin.Context) string {
	var expiration = time.Now().Add(24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, int(expiration.Unix()), "/", "localhost", false, true)
	return state
}
