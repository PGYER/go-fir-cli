package analysis

import (
	"archive/zip"
	"bytes"
	"errors"
	"image"
	"image/png"
	"io"
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

type iosPlist struct {
	CFBundleName         string `plist:"CFBundleName"`
	CFBundleDisplayName  string `plist:"CFBundleDisplayName"`
	CFBundleVersion      string `plist:"CFBundleVersion"`
	CFBundleShortVersion string `plist:"CFBundleShortVersionString"`
	CFBundleIdentifier   string `plist:"CFBundleIdentifier"`
}

func Ipa(name string) (*AppFileInfo, error) {
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
	var plistFile, iosIconFile, embeddedMobileprovision *zip.File

	for _, f := range reader.File {
		switch {

		case reInfoPlist.MatchString(f.Name):
			plistFile = f
		case strings.Contains(f.Name, "embedded.mobileprovision"):
			embeddedMobileprovision = f
		case strings.Contains(f.Name, "AppIcon60x60"):
			iosIconFile = f
		}
	}
	info, err := parseIpaFile(plistFile)

	info.Udids = readUdids(embeddedMobileprovision)
	// read embedded.mobileprovision file to get udids

	// info.Udids = readUdids(buf)

	icon, _ := parseIpaIcon(iosIconFile)
	info.Icon = icon
	info.Size = stat.Size()
	return info, err
}

func readUdids(embeddedMobileprovision *zip.File) []string {

	rc, err := embeddedMobileprovision.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}

	reg := regexp.MustCompile(`<key>ProvisionedDevices</key>\s*<array>(.*?)</array>`)
	matches := reg.FindStringSubmatch(string(data))
	if len(matches) == 0 {
		return nil
	}

	reg = regexp.MustCompile(`<string>(.*?)</string>`)
	matches1 := reg.FindAllStringSubmatch(matches[1], -1)
	if len(matches) == 0 {
		return nil
	}

	udids := make([]string, 0, len(matches1))
	for _, match := range matches1 {
		udids = append(udids, match[1])
	}
	return udids
}

// func readPackageData(file string) map[string]string {

// }
func parseIpaFile(plistFile *zip.File) (*AppFileInfo, error) {
	if plistFile == nil {
		return nil, errors.New("info.plist not found")
	}

	rc, err := plistFile.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	buf, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	p := new(iosPlist)
	decoder := plist.NewDecoder(bytes.NewReader(buf))
	if err := decoder.Decode(p); err != nil {
		return nil, err
	}

	info := new(AppFileInfo)
	if p.CFBundleDisplayName == "" {
		info.Name = p.CFBundleName
	} else {
		info.Name = p.CFBundleDisplayName
	}
	info.BundleId = p.CFBundleIdentifier
	info.Version = p.CFBundleShortVersion
	info.Build = p.CFBundleVersion
	info.Type = "ios"

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
