package model

// AiRequest :
type AiRequest struct {
	Keyword      string `json:"keyword" binding:"required"`
	FileLocation string `json:"file_location" binding:"required"`
}

// AiOutput :
type AiOutput struct {
	ContentAccracyPercentage string `json:"content_acuracy_percentage"`
	Content                  string `json:"content"`
}

// AiResumeParseOutput :
type AiResumeParseOutput struct {
	Choices []Choices `json:"choices"`
}

// Choices :
type Choices struct {
	Index   int     `json:"index"`
	Message Message `json:"message"`
}

// Message :
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
