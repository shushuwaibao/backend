package model

import (
	"fmt"

	"gorm.io/gorm"
)

// ImageConfig represents the model for an image configuration in the database.
type ImageConfig struct {
	ID           int     `gorm:"primaryKey" json:"id"`
	Nickname     string  `gorm:"size:255" json:"nickname"`
	Name         string  `gorm:"size:255" json:"name"`
	Registry     string  `gorm:"size:255" json:"registry"`
	Version      string  `gorm:"size:255" json:"version"`
	Description  string  `gorm:"type:text" json:"description"`
	Size         float64 `gorm:"type:float" json:"size"`
	BelongsToWho int     `gorm:"type:int" json:"belongsToWho"` // "user" = 0 ,"organization" = 1
	BelongsTo    int     `gorm:"type:int" json:"belongsTo"`    // user id or organization id
	Permission   string  `gorm:"size:255" json:"permission"`   // "0": public, "1": private
}

func GetContainerUrl(ii *ImageConfig) string {
	return fmt.Sprintf("%v/%v:%v", ii.Registry, ii.Name, ii.Version)
}

func GetImageUrlByID(id int) (string, error) {
	var image ImageConfig
	if err := DB.First(&image, id).Error; err != nil {
		return "", err
	}
	return GetContainerUrl(&image), nil
}

// ListAvailableImages lists images based on permission.
// It selects public images and private images that belong to the specified user or organization.
func ListAvailableImages(db *gorm.DB, belongsToWho string, belongsToId int) ([]ImageConfig, error) {
	var images []ImageConfig

	// Build the query to select images that are public or belong to the specified entity.
	query := db.Where("permission = ?", "0") // Public images
	if belongsToWho == "user" || belongsToWho == "organization" {
		// Add condition to include private images belonging to the user or organization
		query = query.Or("(belongs_to_who = ? AND belongs_to = ? AND permission = ?)", belongsToWho, belongsToId, "1")
	}

	// Execute the query
	if err := query.Find(&images).Error; err != nil {
		return nil, err // Handle error, e.g., return it
	}

	return images, nil
}

func InsertImage(image *ImageConfig) error {
	return DB.Create(image).Error
}
