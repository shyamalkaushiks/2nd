package services

import (
	"fmt"
	"html"
	"net/http"
	"strconv"
	"strings"
	"time"
	"users/auth"
	"users/data"
	"users/logger"
	model "users/models"
	"users/util"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (hs *HandlerService) RegisterUser(c *gin.Context) {
	var err error
	var request model.Users
	request.UserUuid = GenerateId()
	if err = c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Invalid request parameters"))
		return
	}
	request.Password, err = GenerateHashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "something went wrong"))
		return
	}
	request.VerificationToken, err = data.GenerateVerificationToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.HttpWebResponseError(http.StatusInternalServerError, err.Error(), "Something went wrong"))
		return
	}
	request.Verified = false

	//email space
	request.Email = html.EscapeString(strings.TrimSpace(request.Email))
	err = data.CreateUSer(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.HttpWebResponseError(http.StatusInternalServerError, err.Error(), "Something went wrong"))
		return
	}

	var userlogin model.UserLoginRequest

	userlogin.Email = request.Email
	var response model.Users
	response, err = data.GetUSerByEmail(userlogin)
	fmt.Println("email id is", response.Email)
	if err != nil {
		logger.Log.Error().Err(err).Msg("while getting user records, " + userlogin.Email + " (/UserRegister)")
		c.JSON(http.StatusInternalServerError, util.HttpWebResponseError(http.StatusInternalServerError, err.Error(), "Something went wrong"))
		return
	}
	if response.Id == 0 { // when user not found in database
		c.JSON(http.StatusUnauthorized, util.HttpWebResponseError(http.StatusUnauthorized, err.Error(), "Your email and password do not match. Please try again."))
		return
	}
	token, err2 := auth.GenerateJwt(int(response.Id), response.Email)
	if err2 != nil {
		logger.Log.Error().Err(err).Msg("while verify password email or password is incorect," + response.Email + "(/UserRegister)")
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "invalid request paramaters"))
		return
	}
	response.Password = ""
	response.Token = token
	c.JSON(http.StatusOK, util.HttpWebResponseWithDataSuccess(http.StatusOK, "User created successfully", response))

}

func (hs *HandlerService) LoginUser(c *gin.Context) {
	var err error
	var logindetails model.UserLoginRequest
	//binding json
	if err = c.ShouldBindBodyWithJSON(&logindetails); err != nil {
		logger.Log.Error().Err(err).Msg("while binding json (/UserLogin)")
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Invalid request parameters"))
		return
	}
	var response model.Users
	response, err = data.GetUSerByEmail(logindetails)

	if err != nil && err.Error() != "record not found" {
		logger.Log.Error().Err(err).Msg("while getting user records, " + response.Email + " (/UserLogin)")
		c.JSON(http.StatusInternalServerError, util.HttpWebResponseError(http.StatusInternalServerError, err.Error(), "Something went wrong"))
		return
	}
	if response.Id == 0 { // when user not found in database
		c.JSON(http.StatusUnauthorized, util.HttpWebResponseError(http.StatusUnauthorized, err.Error(), "Your email and password do not match. Please try again."))
		return
	}

	//
	err = VerifyPassword(logindetails.Password, response.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		logger.Log.Error().Err(err).Msg("while verify password, email or password is incorrect, " + response.Email + " (/UserLogin)")
		c.JSON(http.StatusUnauthorized, util.HttpWebResponseError(http.StatusUnauthorized, err.Error(), "Your email and password do not match. Please try again."))
		return
	}

	// token generate
	token, err := auth.GenerateJwt(int(response.Id), response.Email)
	if err != nil {
		logger.Log.Error().Err(err).Msg("while verify password, email or password is incorrect, " + response.Email + " (/UserLogin)")
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Invalid request parameters"))
		return
	}

	response.Token = token

	c.JSON(http.StatusOK, util.HttpWebResponseWithDataSuccess(http.StatusOK, "Logged in successfully", response))

}

// GenerateId :
func GenerateId() string {
	currentTime := time.Now()
	refIdStr := strconv.Itoa(int(currentTime.Unix()))

	return refIdStr
}
func (hs *HandlerService) UpdateUserProfile(c *gin.Context) {
	user_id := c.Param("user_id")
	if user_id == "" {
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, "Invalid user ID", "User ID is empty"))
		return
	}
	fmt.Println(user_id)

	//user send data
	var userUpdate model.Users
	//db data
	userIDInt, err := strconv.Atoi(user_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Invalid user ID"))
		return
	}
	users, err := data.GetUSerByUserid(userIDInt)
	if err != nil {
		fmt.Println("users not found", err)
	}
	if err := c.ShouldBindJSON(&userUpdate); err != nil {
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Invalid request parameters"))
		return
	}
	if userUpdate.FirstName != "" {
		users.FirstName = userUpdate.FirstName
	}
	if userUpdate.LastName != "" {
		users.LastName = userUpdate.LastName
	}
	if userUpdate.Phone != "" {
		users.Phone = userUpdate.Phone
	}
	if userUpdate.Photo != "" {
		users.Photo = userUpdate.Photo
	}
	err1 := data.UpdateUser(users)
	if err1 != nil {
		fmt.Println("error in updation", err1)
	} else {
		c.JSON(http.StatusOK, "updated")
	}

}

var jwtKey = []byte("conscious")

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func (hs *HandlerService) ForgetPassword(c *gin.Context) {
	//searchindbifaccountexists
	//if yes
	email := c.Query("emailid")
	fmt.Println("email id is", email)
	data, err := data.FindUserByEmailInDb(email)
	fmt.Println(data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": `no user found in db ${err} `,
		})

	}
	expirationTime := time.Now().Add(15 * time.Minute)
	jwtclaims := &Claims{
		Email: data.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaims)
	tokenstring, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "token error",
		})
	}
	//resetString = "localhost:8080/reset-password?token=%s", tokenString
	fmt.Println(tokenstring)

}

type Reset struct {
	Token    string `json:"token"`
	Password string `json:"password"`
}

func (hs *HandlerService) ResetPassword(c *gin.Context) {
	var resetdata Reset
	if err := c.BindJSON(&resetdata); err != nil {
		c.JSON(http.StatusOK, gin.H{"err": err.Error()})
		return
	}
	Claims := &Claims{}
	token, err := jwt.ParseWithClaims(resetdata.Token, Claims,
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return

	}
	hashedpassword, err := bcrypt.GenerateFromPassword([]byte(resetdata.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	var user model.Users
	db := model.DBConn
	err = db.Where("users.email=?", Claims.Email).First(&user).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"err": err})
		return
	}
	password := string(hashedpassword)
	user.Password = password
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password successfully updated"})

}
