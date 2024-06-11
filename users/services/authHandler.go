package services

import (
	"net/http"
	"users/data"
	model "users/models"
	"users/util"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// VerifyPassword :
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GenerateHashPassword(password string) (string, error) {

	var HashedPassword string
	var err error

	Hashvalue, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	HashedPassword = string(Hashvalue)
	return HashedPassword, nil

}
func (hs *HandlerService) VerifyEmail(c *gin.Context) {
	db := model.DBConn
	var datauser model.Users
	var err error
	token := c.Query("token")
	datauser, err = data.GetUsersByVerificationToken(token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	datauser.Verified = true
	datauser.VerificationToken = ""
	db.Save(datauser)
	c.JSON(http.StatusOK, util.HttpWebResponseSuccess(http.StatusOK, "saved"))
}
