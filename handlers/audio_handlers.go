package handlers

import (
	"fmt"
	"net/http"
	"os"
	"time"
	"path/filepath"
	"strings"
	"audio-server/models"
	"audio-server/utils"
	"audio-server/database"
	"github.com/gin-gonic/gin"
)

func CreateAudio(c *gin.Context) {
	var input models.AudioInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	audioFile, audioHeader, err := c.Request.FormFile("audio")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "audio file is required"})
		return
	}
	defer audioFile.Close()

	ext := strings.ToLower(filepath.Ext(audioHeader.Filename))
	allowed := map[string]bool{".mp3": true, ".wav": true, ".aac": true, ".m4a": true}
	if !allowed[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported audio format"})
		return
	}

	storage := utils.LocalStorage{BasePath: "./audios"}
	audioPath, err := storage.UploadFile(audioFile, audioHeader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload audio"})
		return
	}

	coverPath := ""
	coverFile, coverHeader, err := c.Request.FormFile("cover_image")
	if err == nil {
		defer coverFile.Close()
		coverPath, err = storage.UploadFile(coverFile, coverHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload cover"})
			return
		}
	}

	hlsDir := fmt.Sprintf("./audios/hls/%d",  time.Now().UnixNano())
	_, err = utils.GenerateHLS(audioPath, hlsDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate HLS"})
		return
	}

	audio := models.Audio{
		Title:              input.Title,
		AudioFilePath:      audioPath,
		CoverImageFilePath: coverPath,
		HLSPlaylistPath:    hlsDir, // store folder only
	}

	if err := database.DB.Create(&audio).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create audio"})
		return
	}

	c.JSON(http.StatusCreated, audio)
}

func GetAudios(c *gin.Context) {
	var audios []models.Audio
	database.DB.Find(&audios)
	c.JSON(http.StatusOK, audios)
}

func GetAudio(c *gin.Context) {
	id := c.Param("id")
	var audio models.Audio
	if err := database.DB.First(&audio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio not found"})
		return
	}

	if _, err := os.Stat(audio.AudioFilePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio file missing"})
		return
	}

	ext := strings.ToLower(filepath.Ext(audio.AudioFilePath))
	switch ext {
	case ".mp3":
		c.Header("Content-Type", "audio/mpeg")
	case ".wav":
		c.Header("Content-Type", "audio/wav")
	case ".aac":
		c.Header("Content-Type", "audio/aac")
	case ".m4a":
		c.Header("Content-Type", "audio/mp4")
	default:
		c.Header("Content-Type", "application/octet-stream")
	}

	c.Header("Accept-Ranges", "bytes")
	http.ServeFile(c.Writer, c.Request, audio.AudioFilePath)
}

func UpdateAudio(c *gin.Context) {
	id := c.Param("id")
	var audio models.Audio
	if err := database.DB.First(&audio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio not found"})
		return
	}

	if title := c.PostForm("title"); title != "" {
		audio.Title = title
	}

	storage := utils.LocalStorage{BasePath: "./audios"}

	newAudioFile, newAudioHeader, err := c.Request.FormFile("audio")
	if err == nil {
		defer newAudioFile.Close()
		if audio.AudioFilePath != "" {
			os.Remove(audio.AudioFilePath)
		}
		newPath, err := storage.UploadFile(newAudioFile, newAudioHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload new audio"})
			return
		}
		audio.AudioFilePath = newPath

		// Regenerate HLS
		hlsDir := fmt.Sprintf("./audios/hls/%d",  time.Now().UnixNano())
		_, err = utils.GenerateHLS(audio.AudioFilePath, hlsDir)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to regenerate HLS"})
			return
		}
		audio.HLSPlaylistPath = hlsDir
	}

	newCoverFile, newCoverHeader, err := c.Request.FormFile("cover_image")
	if err == nil {
		defer newCoverFile.Close()
		if audio.CoverImageFilePath != "" {
			os.Remove(audio.CoverImageFilePath)
		}
		newCoverPath, err := storage.UploadFile(newCoverFile, newCoverHeader)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload new cover"})
			return
		}
		audio.CoverImageFilePath = newCoverPath
	}

	database.DB.Save(&audio)
	c.JSON(http.StatusOK, audio)
}

func DeleteAudio(c *gin.Context) {
	id := c.Param("id")
	var audio models.Audio
	if err := database.DB.First(&audio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio not found"})
		return
	}

	if audio.AudioFilePath != "" {
		os.Remove(audio.AudioFilePath)
	}
	if audio.CoverImageFilePath != "" {
		os.Remove(audio.CoverImageFilePath)
	}
	if audio.HLSPlaylistPath != "" {
		os.RemoveAll(audio.HLSPlaylistPath)
	}

	database.DB.Delete(&audio)
	c.JSON(http.StatusOK, gin.H{"message": "audio deleted successfully"})
}

func StreamHLS(c *gin.Context) {
	id := c.Param("id")
	var audio models.Audio
	if err := database.DB.First(&audio, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "audio not found"})
		return
	}

	requestedFile := c.Param("filepath")
	if requestedFile == "" || requestedFile == "/" {
		requestedFile = "index.m3u8"
	} else {
		requestedFile = strings.TrimPrefix(requestedFile, "/")
	}

	path := filepath.Join(audio.HLSPlaylistPath, requestedFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}

	c.File(path)
}