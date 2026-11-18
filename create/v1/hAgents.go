package create

import (
	_ "embed"
	"strings"

	"github.com/celsiainternet/elvis/file"
	"github.com/celsiainternet/elvis/strs"
)

//go:embed template/agents.md
var modelAgents string

//go:embed template/settings.json
var modelClaudeSettings string

/**
* rulesMarker: Marca de las reglas del proyecto en CLAUDE.md y AGENTS.md, para no agregarlas dos veces.
**/
const rulesMarker = "<!-- elvis:project-rules -->"

/**
* MakeAgents: Escribe las reglas para agentes en CLAUDE.md y AGENTS.md y crea .claude/settings.json.
* @param name string
* @return error
**/
func MakeAgents(name string) error {
	err := upsertRules("CLAUDE.md", name)
	if err != nil {
		return err
	}

	err = upsertRules("AGENTS.md", name)
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

/**
* upsertRules: Crea fileName con las reglas o las antepone si aún no las tiene.
* @param fileName, name string
* @return error
**/
func upsertRules(fileName, name string) error {
	rules := strs.Format("# %s\n\n%s", fileName, fillTemplate(modelAgents, name))
	path := "./" + fileName
	if !file.ExistPath(path) {
		return file.WriteFile(path, rules)
	}

	current, err := file.ReadFile(path)
	if err != nil {
		return err
	}

	if strings.Contains(current, rulesMarker) {
		return nil
	}

	return file.WriteFile(path, rules+"\n"+current)
}
