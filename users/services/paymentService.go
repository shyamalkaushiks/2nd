package services

import (
	"fmt"
	"net/http"
	model "users/models"

	"github.com/gin-gonic/gin"
)

func (hs *HandlerService) LoadHtmlForm(c *gin.Context) {

	c.HTML(http.StatusOK, "payment.html", nil)

}
func (hs *HandlerService) AcceptData(c *gin.Context) {
	var payment model.Paymentdatails
	err := c.ShouldBindJSON(&payment)
	if err != nil {
		fmt.Println("binding errors", err)
	}
}
