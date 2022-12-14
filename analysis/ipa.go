package analysis

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
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

	if len(info.Udids) == 0 {
		info.ReleaseType = "inhouse"
	} else {
		info.ReleaseType = "adhoc"
	}

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

	pattern := `<key>ProvisionedDevices</key>[\s\S]+?<array>([\s\S]+?)</array>`

	matched, _ := regexp.MatchString(pattern, string(data))

	answer := make([]string, 0)

	if matched {
		r, _ := regexp.Compile(pattern)
		result := r.FindStringSubmatch(string(data))
		if len(result) >= 2 {
			//拆分匹配结果
			arr := strings.Split(result[1], "<string>")
			for _, v := range arr {
				//去除多余的字符
				v = strings.Replace(strings.TrimSpace(v), "</string>", "", -1)
				//保存结果
				if v != "" {
					fmt.Println(v)
					answer = append(answer, v)
				}
			}
			return answer
		}
	}
	return []string{}

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
