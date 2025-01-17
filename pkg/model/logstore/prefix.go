package logstore

import (
	"fmt"
	"strings"

	"github.com/astro-walker/tilt/pkg/model"
)

func SourcePrefix(n model.ManifestName) string {
	if n == "" || n == model.MainTiltfileManifestName {
		return ""
	}
	max := 13
	spaces := ""
	if len(n) > max {
		n = n[:max-1] + "…"
	} else {
		spaces = strings.Repeat(" ", max-len(n))
	}
	return fmt.Sprintf("%s%s │ ", spaces, n)
}
