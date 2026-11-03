package agentsguide

import (
	_ "embed"
	"strings"

	"github.com/celsiainternet/elvis/file"
)

//go:embed ELVIS_AGENTS.md
var content string

const startMarker = "<!-- elvis:agents-guide:start -->"

/**
* Install: Writes the Elvis agents framework guide into the current
* project's CLAUDE.md so Claude Code (and compatible agents) auto-load it.
* Creates CLAUDE.md if missing; if it already exists, appends the guide
* once (idempotent - a second call is a no-op if the marker is present).
* @return error
**/
func Install() error {
	path := "./CLAUDE.md"

	if !file.ExistPath(path) {
		_, err := file.AppendFile(".", "CLAUDE.md", content)
		return err
	}

	current, err := file.ReadFile(path)
	if err != nil {
		return err
	}

	if strings.Contains(current, startMarker) {
		return nil
	}

	_, err = file.AppendFile(".", "CLAUDE.md", "\n\n"+content)
	return err
}
