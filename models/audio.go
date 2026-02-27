package models

package models

type Audio struct {
    ID                 uint   `json:"id" gorm:"primaryKey"`
    Title              string `json:"title"`
    AudioFilePath      string `json:"video_file_path"`
    CoverImageFilePath string `json:"cover_image_file_path"`
}

type AudioInput struct {
    Title string `form:"title" binding:"required"`
}
