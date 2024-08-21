package controllers

import (
	"encoding/base64"
	"net/http"
	"path"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/jimmykarily/quizmaker/internal/models"
	"gorm.io/gorm/clause"
)

type (
	SessionController struct{}
)

func (c *SessionController) List(gctx *gin.Context) {
	sessions := []models.Session{}
	err := Settings.DB.Preload(clause.Associations).Find(&sessions).Error
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	var complete, inProgress []models.Session
	for _, s := range sessions {
		oldScore := s.Score
		oldComplete := s.Complete
		s.UpdateCacheColumns()
		if oldScore != s.Score || oldComplete != s.Complete { // if needs update
			err = Settings.DB.Save(s).Error
			if handleError(gctx.Writer, err, http.StatusInternalServerError) {
				return
			}
		}

		if s.Complete {
			complete = append(complete, s)
		} else {
			inProgress = append(inProgress, s)
		}
	}

	sort.Slice(complete, func(i, j int) bool {
		return complete[i].Score > complete[j].Score
	})

	NewQuizURL, err := GetFullURL(gctx.Request, "QuizNew", nil)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	png, err := getQRCodePNG(NewQuizURL)
	if handleError(gctx.Writer, err, http.StatusInternalServerError) {
		return
	}

	viewData := struct {
		QRCodePNG  string
		NewQuizURL string
		Completed  []models.Session
		InProgress []models.Session
	}{
		QRCodePNG:  base64.StdEncoding.EncodeToString(png),
		NewQuizURL: NewQuizURL,
		Completed:  complete,
		InProgress: inProgress,
	}

	Render([]string{"main_layout", path.Join("sessions", "list")}, gctx.Writer, viewData)
}
