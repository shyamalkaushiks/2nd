package auth

import (
	"errors"
	"time"
	"users/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	JWTClaimUserId int
	JWTClaimEmail  string
)

var jwtKey = []byte(config.Config.API_SECRET)

type JwtClaims struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

func GenerateJwt(UserId int, email string) (tokenString string, err error) {
	expirationTime := time.Now().Add(time.Hour * 24)
	jwtclaims := &JwtClaims{
		Id:    UserId,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaims)
	// fmt.Println(token)
	tokenString, err = token.SignedString(jwtKey)
	// fmt.Println(tokenString)
	return
}
func Auth() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")
		if tokenString == "" {
			context.JSON(401, gin.H{
				"error": "request does not contain an access token",
			})
			context.Abort()
			return
		}
		err := ValidateToken(tokenString)
		if err != nil {
			context.JSON(401, gin.H{"error": err.Error()})
			context.Abort()
			return
		}
		context.Next()
	}

}

func ValidateToken(token string) (err error) {
	tokenv, err := jwt.ParseWithClaims(
		token,
		&JwtClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	Claims, ok := tokenv.Claims.(*JwtClaims)
	if !ok {
		err = errors.New("couldn't parse claims")
		return
	}
	if Claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token expired")
		return
	}
	JWTClaimUserId = Claims.Id
	JWTClaimEmail = Claims.Email
	if JWTClaimUserId == 0 || JWTClaimEmail == "" {
		err = errors.New("token not valid")
		return
	}
	return

}
