package create

import (
	_ "embed"

	"github.com/celsiainternet/elvis/file"
)

//go:embed template/agents.md
var modelAgents string

//go:embed template/settings.json
var modelClaudeSettings string

/**
* MakeAgents: Escribe CLAUDE.md y AGENTS.md con las reglas del proyecto para
* los agentes, y .claude/settings.json, que impide a Claude Code hacer commit
* o push. No toca archivos que ya existan; agentsguide.Install agrega despues
* su bloque con las guias de elvis.
* @param name string
* @return error
**/
func MakeAgents(name string) error {
	_, err := file.MakeFile(".", "CLAUDE.md", "# CLAUDE.md\n\n"+modelAgents, name)
	if err != nil {
		return err
	}

	_, err = file.MakeFile(".", "AGENTS.md", "# AGENTS.md\n\n"+modelAgents, name)
	if err != nil {
		return err
	}

	path, err := file.MakeFolder(".claude")
	if err != nil {
		return err
	}

	_, err = file.MakeFile(path, "settings.json", modelClaudeSettings)
	return err
}
