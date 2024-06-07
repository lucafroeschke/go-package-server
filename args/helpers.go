package args

import (
	"os"
	"strings"
)

func GetString(arg string) string {
	for _, v := range os.Args {
		if strings.HasPrefix(v, "--"+arg+"=") {
			split := strings.Split(v, "=")
			if len(split) > 1 {
				return split[1]
			}
		}
	}
	return ""
}

func GetBool(arg string) bool {
	for _, v := range os.Args {
		if v == "--"+arg {
			return true
		}
	}
	return false
}
