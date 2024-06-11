package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"users/ai/chatgpt"
	"users/auth"

	"users/data"
	"users/logger"
	model "users/models"

	"users/parser"
	"users/util"

	"github.com/gin-gonic/gin"
)

// ResumeUpload :
func (hs *HandlerService) ResumeUpload(c *gin.Context) {

	var err error
	var userResumes []model.UserResumes

	userId := auth.JWTClaimUserId

	form, _ := c.MultipartForm()
	if userId == 0 || len(form.File) == 0 {
		errStr := "invalid request parameter"
		c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, errStr, "Something went wrong"))
		return
	}

	files := form.File["files"]
	filePaths := []string{}

	for _, file := range files {
		fileExt := filepath.Ext(file.Filename)
		if !strings.Contains(fileExt, ".pdf") {
			errStr := "file format not valid"
			c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, errStr, "Something went wrong"))
			return
		}
	}

	for _, file := range files {

		var userResume model.UserResumes

		fileExt := filepath.Ext(file.Filename)
		originalFileName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		now := time.Now()
		filename := strings.ReplaceAll(strings.ToLower(originalFileName), " ", "-") + "-" + fmt.Sprintf("%v", now.Unix()) + fileExt
		//filePath := "/etc/ai_resume/user/resume_uploads/"
		//filePath := "/etc/ai_resume/user/resume_uploads/"
		filePath := "F:/upload_resume/"

		filePaths = append(filePaths, filePath)

		created, errStr := util.CheckPathExists(filePath, 1)
		if !created {
			logger.Log.Error().Err(errors.New(errStr)).Msg("while create image, (ResumeUpload())")
			c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, errStr, "Something went wrong"))
			return
		}

		out, err := os.Create(filePath + filename)
		if err != nil {
			logger.Log.Error().Err(err).Msg("while create file, (ResumeUpload())")
			c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Something went wrong"))
			return
		}
		defer out.Close()

		readerFile, _ := file.Open()
		_, err = io.Copy(out, readerFile)
		if err != nil {
			logger.Log.Error().Err(err).Msg("while copy image, (ResumeUpload())")
			c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Something went wrong"))
			return
		}

		userResume.UserId = userId
		userResume.FilePath = filePath + filename
		userResume.CreatedAt = time.Now()
		userResume.ResumeStatusId = 1 // 1 = pending

		userResumes = append(userResumes, userResume)
	}

	var objResumesToUpload []interface{}
	for _, values := range userResumes {
		objResumesToUpload = append(objResumesToUpload, values)
	}

	if len(objResumesToUpload) == 0 {
		c.JSON(http.StatusOK, util.HttpWebResponseSuccess(http.StatusOK, "Record not found"))
		return
	}

	err = data.InsertUserResume(objResumesToUpload)
	if err != nil {
		logger.Log.Error().Err(err).Msg("while insert resume, (ResumeUpload())")
		c.JSON(http.StatusInternalServerError, util.HttpWebResponseError(http.StatusInternalServerError, err.Error(), "Something went wrong"))
		return
	}

	//
	var userResumesUploadParseDatas []model.UserResumesUploadParseData
	for _, userResume := range userResumes {

		//
		var contentReadFile string
		contentReadFile, err = parser.ReadFiles(userResume.FilePath)
		if err != nil {
			fmt.Println("while read file : ", err)
			c.JSON(http.StatusBadRequest, util.HttpWebResponseError(http.StatusBadRequest, err.Error(), "Something went wrong"))
			return
		}

		fmt.Println(contentReadFile)

		// Resume parsing
		var result string
		//chatgpt.ResumeParseRequest
		result, err = chatgpt.ResumeParse(contentReadFile)
		if err != nil {
			fmt.Println("error while genrate resume parse :", err)
			return
		}
		// fmt.Println("success", result)

		var resumeParseOutput model.AiResumeParseOutput
		json.Unmarshal([]byte(result), &resumeParseOutput)

		if len(resumeParseOutput.Choices) > 0 {
			var userResumesUploadParseData model.UserResumesUploadParseData
			// fmt.Println(resumeParseOutput.Choices[0].Message.Content)
			userResumesUploadParseData.Id = userResume.Id
			userResumesUploadParseData.UserId = userResume.UserId
			userResumesUploadParseData.FilePath = userResume.FilePath
			userResumesUploadParseData.CreatedAt = userResume.CreatedAt
			userResumesUploadParseData.ResumeStatusId = userResume.ResumeStatusId
			userResumesUploadParseData.AiResumeScorePercentage = userResume.AiResumeScorePercentage
			userResumesUploadParseData.ResumeParseData = resumeParseOutput.Choices[0].Message.Content

			userResumesUploadParseDatas = append(userResumesUploadParseDatas, userResumesUploadParseData)

			// update in database
			//data.UpdateResumeAiParseDataByPath(userResumesUploadParseData.ResumeParseData, userResumesUploadParseData.FilePath)
		}
	}

	if len(userResumesUploadParseDatas) == 0 {
		c.JSON(http.StatusOK, util.HttpWebResponseSuccess(http.StatusOK, "Record not found!"))
		return
	}

	c.JSON(http.StatusOK, util.HttpWebResponseWithDataSuccess(http.StatusOK, "Resume uploaded successfully", userResumesUploadParseDatas))
}

func (hs *HandlerService) AllResumeOfUser(c *gin.Context) {
	id := auth.JWTClaimUserId
	db := model.DBConn
	var resumes []model.UserResumes
	err := db.Where("user_id=?", id).Find(&resumes).Error
	if err != nil {
		fmt.Println("err", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"resumelist": resumes,
	})

}
