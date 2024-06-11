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
	"users/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig1 *oauth2.Config
)

func init() {
	googleOauthConfig1 = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "use your own",
		ClientSecret: "use your own",
		Endpoint:     google.Endpoint,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
	}
}

func (hs *HandlerService) HandleGoogleLogin1(c *gin.Context) {
	random := generateStateOauthCookie1(c)
	url1 := googleOauthConfig1.AuthCodeURL(random)
	c.Redirect(http.StatusTemporaryRedirect, url1)
}
func (hs *HandlerService) HandleCallback(c *gin.Context) {
	cookievalue, err := c.Cookie("oauthstate")
	if err != nil || c.Query("state") != cookievalue {
		logger.Log.Error().Err(err).Msg("value not matched")
		fmt.Println("error value not matched")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	code := c.Query("code")
	token, err := googleOauthConfig1.Exchange(context.Background(), code)
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
	var userInfo1 struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Picture       string `json:"picture"`
		Name          string `json:"name"`
	}
	err = json.NewDecoder(response.Body).Decode(&userInfo1)
	fmt.Println("user deails are", userInfo1)
	if err != nil {
		log.Printf("Failed decoding user info: %s\n", err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}
	// For simplicity, we will just print the user info in the response
	c.JSON(http.StatusOK, gin.H{
		"user": userInfo1,
	})
}
func generateStateOauthCookie1(c *gin.Context) string {
	expirationtime := time.Now().Add(24 * time.Hour)
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, int(expirationtime.Unix()), "/", "localhost", false, true)
	return state
}
