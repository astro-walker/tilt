package encoding

import "github.com/astro-walker/tilt/internal/tiltfile/starkit"

type Plugin struct{}

func NewPlugin() Plugin {
	return Plugin{}
}

const (
	readYAMLN         = "read_yaml"
	readYAMLStreamN   = "read_yaml_stream"
	decodeYAMLN       = "decode_yaml"
	decodeYAMLStreamN = "decode_yaml_stream"
	encodeYAMLN       = "encode_yaml"
	encodeYAMLStreamN = "encode_yaml_stream"

	readJSONN   = "read_json"
	decodeJSONN = "decode_json"
	encodeJSONN = "encode_json"
)

func (Plugin) OnStart(env *starkit.Environment) error {
	for _, b := range []struct {
		name string
		f    starkit.Function
	}{
		{readYAMLN, readYAML},
		{readYAMLStreamN, readYAMLStream},
		{decodeYAMLN, decodeYAML},
		{decodeYAMLStreamN, decodeYAMLStream},
		{encodeYAMLN, encodeYAML},
		{encodeYAMLStreamN, encodeYAMLStream},

		{readJSONN, readJSON},
		{decodeJSONN, decodeJSON},
		{encodeJSONN, encodeJSON},
	} {
		err := env.AddBuiltin(b.name, b.f)
		if err != nil {
			return err
		}
	}

	return nil
}
