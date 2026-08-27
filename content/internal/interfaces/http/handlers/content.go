package handlers

import (
	"errors"
	"net/http"

	"github.com/IvanDrf/work-hunter/content/internal/domain/models"
	"github.com/IvanDrf/work-hunter/content/internal/domain/ports/service"
	"github.com/gin-gonic/gin"
)

type ContentHandler struct {
	service service.ContentService
}

func NewContentHandler(svc service.ContentService) *ContentHandler {
	return &ContentHandler{service: svc}
}

func (h *ContentHandler) UploadResume(c *gin.Context) {
	h.handleUpload(c, models.TypeResume)
}

func (h *ContentHandler) UploadAvatar(c *gin.Context) {
	h.handleUpload(c, models.TypeAvatar)
}

func (h *ContentHandler) handleUpload(c *gin.Context, cType models.ContentType) {
	userID := c.Param("user_id")
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read file"})
		return
	}
	defer file.Close()

	var svcErr error
	if cType == models.TypeResume {
		svcErr = h.service.UploadResume(c.Request.Context(), userID, file, fileHeader.Size)
	} else {
		svcErr = h.service.UploadAvatar(c.Request.Context(), userID, file, fileHeader.Size)
	}

	if svcErr != nil {
		h.handleError(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": string(cType) + " uploaded"})
}

func (h *ContentHandler) Download(c *gin.Context) {
	userID := c.Param("user_id")
	cType := models.ContentType(c.Param("type"))

	if cType != models.TypeResume && cType != models.TypeAvatar {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content type"})
		return
	}

	file, meta, err := h.service.GetContent(c.Request.Context(), userID, cType)
	if err != nil {
		h.handleError(c, err)
		return
	}
	defer file.Close()

	c.DataFromReader(http.StatusOK, meta.Size, meta.MimeType, file, nil)
}

func (h *ContentHandler) Delete(c *gin.Context) {
	userID := c.Param("user_id")
	cType := models.ContentType(c.Param("type"))

	if err := h.service.DeleteContent(c.Request.Context(), userID, cType); err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *ContentHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, models.ErrContentNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "content not found"})
	case errors.Is(err, models.ErrFileTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
	case errors.Is(err, models.ErrInvalidMimeType), errors.Is(err, models.ErrFileEmpty):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
