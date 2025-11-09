package util

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"io"
	"strings"

	"github.com/alecthomas/chroma/quick"
)

func PrettyDumpJSON(w io.Writer, v any) error {
	var s strings.Builder

	err := json.MarshalWrite(&s, v, jsontext.WithIndent("  "))
	if err != nil {
		return err
	}

	return quick.Highlight(w, s.String(), "json", "terminal256", "monokai")
}
