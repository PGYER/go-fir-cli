package analysis

import (
	"errors"
	"image"
	"image/png"
	"os"
	"strings"
)

type AppFileInfo struct {
	Name        string
	ReleaseType string
	BundleId    string
	Version     string
	Build       string
	Icon        image.Image
	Size        int64
	Type        string
	Udids       []string
	Changelog   string
}

type IUploadApp interface {
	GetIcon() image.Image
	GetAppFileInfo() map[string]string
	FilePath() string
}

type UploadAppService struct {
	AppFileInfo *AppFileInfo
	FilePath    string
	image       image.Image
}

func NewUploadAppService(filePath string) (*UploadAppService, error) {
	var err error
	uploadAppService := &UploadAppService{
		FilePath: filePath,
	}
	if uploadAppService.FilePath != "" {
		uploadAppService.AppFileInfo, err = uploadAppService.GetAppFileInfo()
		uploadAppService.image = uploadAppService.AppFileInfo.Icon
	}

	return uploadAppService, err
}

// 获得App 文件的相关信息
func (a *UploadAppService) GetAppFileInfo() (*AppFileInfo, error) {
	if strings.HasSuffix(a.FilePath, ".ipa") {
		appInfo, err := Ipa(a.FilePath)
		if err != nil {
			return nil, err
		}
		return appInfo, nil

	} else if strings.HasSuffix(a.FilePath, ".apk") {
		appInfo, err := Apk(a.FilePath)
		if err != nil {
			return nil, err
		}
		return appInfo, nil
	} else {
		return nil, errors.New("file type is not ipa or apk")
	}

}

func (a *UploadAppService) SaveImage(imagePath string) error {

	if a.image == nil {
		return errors.New("image is nil")
	}

	if imagePath == "" {
		imagePath = a.FilePath + ".png"
	}
	file, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()
	png.Encode(file, a.image)
	return nil
}
