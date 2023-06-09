package utils

import (
	"errors"
	"mime/multipart"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadImage(destination string, file *multipart.FileHeader, context *gin.Context) (string, error) {
	imagesName := GenerateUUID() + file.Filename

	os.Chmod(destination, 0755)

	err := context.SaveUploadedFile(file, destination+imagesName)

	if err != nil {
		return "", errors.New(Translation("fail_upload", nil, nil))
	}

	return imagesName, nil
}

func RemoveImage(filePath string) error {
	err := os.Remove(filePath)
	if err != nil {
		return errors.New(Translation("remove_file_fail", nil, nil))
	}

	return nil
}

func IsImageFile(file *multipart.FileHeader) bool {
	contentType := file.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		return false
	}

	return true
}
