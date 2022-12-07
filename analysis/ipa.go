package analysis

import (
	"archive/zip"
	"bytes"
	"errors"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/andrianbdn/iospng"
	"howett.net/plist"
)

var (
	reInfoPlist = regexp.MustCompile(`Payload/[^/]+/Info\.plist`)
	ErrNoIcon   = errors.New("icon not found")
)

type AppInfo struct {
	Name     string
	BundleId string
	Version  string
	Build    string
	Icon     image.Image
	Size     int64
}

type iosPlist struct {
	CFBundleName         string `plist:"CFBundleName"`
	CFBundleDisplayName  string `plist:"CFBundleDisplayName"`
	CFBundleVersion      string `plist:"CFBundleVersion"`
	CFBundleShortVersion string `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string `plist:"CFBundleIdentifier"`
}

func Ipa(name string) (*AppInfo, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	reader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, err
	}
	var plistFile, iosIconFile *zip.File

	for _, f := range reader.File {
		switch {

		case reInfoPlist.MatchString(f.Name):
			plistFile = f
		case strings.Contains(f.Name, "AppIcon60x60"):
			iosIconFile = f
		}
	}
	info, err := parseIpaFile(plistFile)
	icon, err := parseIpaIcon(iosIconFile)
	info.Icon = icon
	info.Size = stat.Size()
	return info, err
}

// func readPackageData(file string) map[string]string {

// }
func parseIpaFile(plistFile *zip.File) (*AppInfo, error) {
	if plistFile == nil {
		return nil, errors.New("info.plist not found")
	}

	rc, err := plistFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	p := new(iosPlist)
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(p); err != nil {
		return nil, err
	}

	info := new(AppInfo)
	if p.CFBundleDisplayName == "" {
		info.Name = p.CFBundleName
	} else {
		info.Name = p.CFBundleDisplayName
	}
	info.BundleId = p.CFBundleIdentifier
	info.Version = p.CFBundleShortVersion
	info.Build = p.CFBundleVersion

	return info, nil
}
func parseIpaIcon(iconFile *zip.File) (image.Image, error) {
	if iconFile == nil {
		return nil, ErrNoIcon
	}

	rc, err := iconFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	var w bytes.Buffer
	iospng.PngRevertOptimization(rc, &w)

	return png.Decode(bytes.NewReader(w.Bytes()))
}
