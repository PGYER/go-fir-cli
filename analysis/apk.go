package analysis

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image"
	"os"

	"github.com/avast/apkparser"
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

func ApkIcon(apkfile string) image.Image {
	pkg, _ := apk.OpenFile(apkfile)
	defer pkg.Close()

	icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image

	return icon
}

func Apk(file string) (appInfo *AppFileInfo, err error) {

	b := bytes.Buffer{}

	enc := xml.NewEncoder(&b)
	enc.Indent("", "\t")
	zipErr, resErr, manErr := apkparser.ParseApk(file, enc)

	if zipErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to open the APK: %s", zipErr.Error())
		os.Exit(1)
		return
	}

	if resErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse resources: %s", resErr.Error())
	}
	if manErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse AndroidManifest.xml: %s", manErr.Error())
		os.Exit(1)
		return
	}
	fmt.Println("======")
	// fmt.Println(b.String())

	type androidLabel struct {
		XMLName xml.Name `xml:"application"`
		Label   string   `xml:"label,attr"`
	}
	var manifestData struct {
		XMLName     xml.Name     `xml:"manifest"`
		Package     string       `xml:"package,attr"`
		VersionName string       `xml:"versionName,attr"`
		VersionCode string       `xml:"versionCode,attr"`
		Application androidLabel `xml:"application"`
	}
	// fmt.Println(b.String())

	err = xml.Unmarshal(b.Bytes(), &manifestData)
	if err != nil {
		panic(err)
	}
	appInfo = &AppFileInfo{}
	appInfo.BundleId = manifestData.Package
	appInfo.Version = manifestData.VersionName
	appInfo.Build = manifestData.VersionCode
	appInfo.Name = manifestData.Application.Label
	appInfo.Type = "android"
	appInfo.ReleaseType = "inhouse"

	icon := ApkIcon(file)
	appInfo.Icon = icon

	return appInfo, nil
}
