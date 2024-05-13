package analysis

import (
	"fmt"
	"image"
	"strconv"

	"github.com/shogo82148/androidbinary/apk"
)

type ApkApp struct {
	UploadAppService
}

func (a *ApkApp) GetIcon() image.Image {
	path := a.FilePath
	icon := ApkIcon(path)

	return icon
}

func ApkInfo(file string) (appInfo *AppFileInfo, err error) {
	appInfo = &AppFileInfo{}
	pkg, _ := apk.OpenFile(file)
	defer pkg.Close()
	icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image

	manifest := pkg.Manifest()

	appInfo.Icon = icon
	appInfo.Name, _ = pkg.Label(nil)
	appInfo.BundleId = manifest.Package.MustString()
	appInfo.Version = manifest.VersionName.MustString()
	appInfo.Build = strconv.Itoa(int(manifest.VersionCode.MustInt32()))

	fmt.Print("appInfo: ", appInfo)

	appInfo.Type = "android"
	appInfo.ReleaseType = "inhouse"

	return appInfo, nil
}

func ApkIcon(apkfile string) image.Image {
	pkg, _ := apk.OpenFile(apkfile)
	defer pkg.Close()

	icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image

	return icon
}

func Apk(file string) (appInfo *AppFileInfo, err error) {

	appInfo = &AppFileInfo{}

	appInfo, _ = ApkInfo(file)

	return appInfo, nil
}
