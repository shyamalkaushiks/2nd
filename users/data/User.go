package data

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	model "users/models"
)

func CreateUSer(detail model.Users) error {
	db := model.DBConn
	db.AutoMigrate(model.Users{})
	err := db.Create(&detail).Error
	if err != nil {
		return err
	}
	return nil
}
func GetUSerByEmail(request model.UserLoginRequest) (model.Users, error) {
	var udetail model.Users
	db := model.DBConn
	err := db.Where("users.email=?", request.Email).First(&udetail).Error

	// err := db.Where("users.email=?", request.Email).Select("users.id, users.first_name, users.last_name, users.email, users.password, users.phone, users.account_id").
	// 	First(&udetail).Error

	if err != nil {
		return udetail, err
	}
	fmt.Println(udetail)
	return udetail, err

}
func GetUSerByUserid(no int) (model.Users, error) {
	var users model.Users
	db := model.DBConn
	err := db.Where("users.id=?", no).First(&users).Error
	if err != nil {
		return users, err
	}
	return users, nil

}
func UpdateUser(user model.Users) error {
	db := model.DBConn
	result := db.Save(&user)
	return result.Error
}
func FindUserByEmailInDb(email string) (model.Users, error) {
	var user model.Users
	db := model.DBConn
	err := db.Where("users.email=?", email).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, err
}

// GenerateVerificationToken generates a random token for email verification
func GenerateVerificationToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
func GetUsersByVerificationToken(token string) (model.Users, error) {

	db := model.DBConn
	var Users model.Users
	err := db.Where("users.verification_token=?", token).First(&Users).Error
	if err != nil {
		return Users, err
	}
	return Users, err
}
