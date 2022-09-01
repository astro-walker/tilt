package cmdimage

import (
	"fmt"

	"github.com/astro-walker/tilt/pkg/apis"
	"github.com/astro-walker/tilt/pkg/model"
)

// Generate the name for the CmdImage API object from an ImageTarget and ManifestName.
func GetName(mn model.ManifestName, id model.TargetID) string {
	return apis.SanitizeName(fmt.Sprintf("%s:%s", mn.String(), id.Name))
}
