package controllers

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
)

type (
	VerificationController struct{}
)

func (c *VerificationController) Form(gctx *gin.Context) {
	submitURL, err := GetFullURL(gctx.Request, "VerificationSubmit", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	Render([]string{"main_layout", path.Join("verification", "form")}, gctx.Writer, struct {
		SubmitURL string
	}{
		SubmitURL: submitURL,
	})
}

func (c *VerificationController) Submit(gctx *gin.Context) {
	email := gctx.PostForm("email")
	fmt.Printf("email = %+v\n", email)

	session, err := models.SessionForEmail(Settings.DB, email)
	if handleError(gctx.Writer, err, http.StatusNotFound) {
		return
	}

	db := Settings.DB.Model(&session).Update("verified", true)
	if db.Error != nil {
		if handleError(gctx.Writer, err, http.StatusInternalServerError) {
			return
		}
	}

	fmt.Printf("Email %s verified\n", email)

	redirectURL, err := GetFullURL(gctx.Request, "VerificationForm", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	gctx.Redirect(http.StatusFound, redirectURL)
}
