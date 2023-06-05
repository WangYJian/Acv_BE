package model

import (
	"time"

	"gorm.io/gorm"
)

type File struct {
	LinkID     string
	UserID     string
	Filename   string
	Path       string
	Result     string
	CreateTime time.Time
}

// Update 函数将根据给定的 ID 更新文件的路径和结果字段。
func (f File) Update(db *gorm.DB) error {
	result := db.Model(&File{}).Where("link_id = ?", f.LinkID).Updates(f)
	return result.Error
}

// CreateFile 函数将创建一个新的文件。
func CreateFile(db *gorm.DB, file File) error {
	result := db.Create(file)
	return result.Error
}

// GetFileByLinkID 函数将返回具有给定 ID 的文件。
func GetFileByLinkID(db *gorm.DB, LinkIDd string) (File, error) {
	var file File
	result := db.First(&file, "link_id = ?", LinkIDd)
	if result.Error != nil {
		return file, result.Error
	}
	return file, nil
}

// GetFilesByUserID 函数将返回具有给定用户 ID 的所有文件。
func GetFilesByUserID(db *gorm.DB, userID string) ([]File, error) {
	var files []File
	result := db.Where("user_id = ?", userID).Order("create_time DESC").Find(&files)
	if result.Error != nil {
		return nil, result.Error
	}
	return files, nil
}
