package services

import (
	"users/auth"

	"github.com/gin-gonic/gin"
)

type HandlerService struct{}

func (hs *HandlerService) Bootstrap(r *gin.Engine) {
	public := r.Group("/api/user")
	public.GET("/login", hs.HandleGoogleLogin1)
	r.GET("/callback", hs.HandleCallback)
	// auth
	public.GET("/welcome", hs.HandleMain)
	public.POST("/auth/signup", hs.RegisterUser)
	public.POST("/auth/Login", hs.LoginUser)
	public.GET("/ForgetPassword", hs.ForgetPassword)
	public.POST("/verifyEmail", hs.VerifyEmail)
	public.GET("/ResetPassword", hs.ResetPassword)
	public.GET("/payment", hs.LoadHtmlForm)

	r.Use(auth.Auth())
	groupRoute := r.Group("/api/user")
	groupRoute.POST("/resume-upload", hs.ResumeUpload)
	groupRoute.PATCH("/user-profile/:user_id", hs.UpdateUserProfile)
	groupRoute.GET("/Resume-list", hs.AllResumeOfUser)

}
