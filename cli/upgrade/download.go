package upgrade

import (
	"fmt"
	"runtime"
)

const downloadURL = "https://mercury.com/dist/%s-amd64/%s/%s"

func DownloadURL(ver, file string) string {
	return fmt.Sprintf(downloadURL, runtime.GOOS, ver, file)
}
