package model

import "time"

type UserResumes struct {
	Id                      int       `json:"id" gorm:"primary_key"`
	UserId                  int       `json:"user_id"`
	FilePath                string    `json:"file_path"`
	CreatedAt               time.Time `json:"created_at"`
	ResumeStatusId          int       `json:"resume_status_id"`
	AiResumeScorePercentage string    `json:"ai_resume_score_percentage"`
	AiResumeParseData       string    `json:"ai_resume_parse_data"`
}

// UserResumes :
type UserResumesUploadParseData struct {
	Id                      int       `json:"id" gorm:"primary_key"`
	UserId                  int       `json:"user_id"`
	FilePath                string    `json:"file_path"`
	CreatedAt               time.Time `json:"created_at"`
	ResumeStatusId          int       `json:"resume_status_id"`
	AiResumeScorePercentage string    `json:"ai_resume_score_percentage"`
	ResumeParseData         string    `json:"resume_parse_data"`
}
