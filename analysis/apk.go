package analysis

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"image/png"
	"os"

	"github.com/avast/apkparser"
	"github.com/shogo82148/androidbinary/apk"
)

func ApkIcon(apkfile string) {
	pkg, _ := apk.OpenFile(apkfile)
	defer pkg.Close()

	icon, _ := pkg.Icon(nil) // returns the icon of APK as image.Image
	file, err := os.Create("icon.png")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	png.Encode(file, icon)
}

func Apk(file string) (result map[string]string, err error) {

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
	fmt.Println(b.String())

	err = xml.Unmarshal(b.Bytes(), &manifestData)
	if err != nil {
		panic(err)
	}

	result = make(map[string]string)
	result["bundle_id"] = manifestData.Package
	result["version"] = manifestData.VersionName
	result["build"] = manifestData.VersionCode
	result["name"] = manifestData.Application.Label
	return result, nil
}
