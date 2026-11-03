package agentsguide

import (
	"embed"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

//go:embed MARCO.md ELVIS.md REVISION.md
var files embed.FS

const module = "github.com/celsiainternet/elvis"

/**
* Install: Installs the Elvis agents framework (MARCO.md, ELVIS.md, REVISION.md)
* into the current project: writes the guides to .claude/guides/elvis/ and a
* managed block in CLAUDE.md and AGENTS.md that points to them. Running it again
* replaces the guides and the block with the ones of the installed Elvis version.
* @return error
**/
func Install() error {
	return InstallGuide("elvis", module, files, []string{"MARCO.md", "ELVIS.md", "REVISION.md"})
}

/**
* InstallGuide: Writes the guides of a library into .claude/guides/<name>/ and
* upserts the block delimited by "<!-- <name>:agents-guide:start" and
* "<!-- <name>:agents-guide:end -->" in CLAUDE.md (as @imports) and AGENTS.md
* (as a reading list). Content outside the markers is never touched.
* @param name string, library name used for the folder and the markers
* @param modulePath string, Go module path used to resolve the installed version
* @param guides fs.FS, embedded guide files
* @param names []string, guide file names in reading order
* @return error
**/
func InstallGuide(name, modulePath string, guides fs.FS, names []string) error {
	version := ModuleVersion(modulePath)
	dir := filepath.Join(".claude", "guides", name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	header := fmt.Sprintf("<!-- Generado por %s %s con cmd/install. No lo edites: se sobrescribe al reinstalar. -->\n\n", name, version)
	paths := []string{}
	for _, guideName := range names {
		content, err := fs.ReadFile(guides, guideName)
		if err != nil {
			return err
		}

		path := filepath.ToSlash(filepath.Join(dir, guideName))
		if err := os.WriteFile(path, []byte(header+string(content)), 0644); err != nil {
			return err
		}

		paths = append(paths, path)
	}

	claudeBlock := managedBlock(name, version, paths, "@%s")
	if err := upsertBlock("CLAUDE.md", name, claudeBlock); err != nil {
		return err
	}

	agentsBlock := managedBlock(name, version, paths, "- Lee `%s` antes de escribir o revisar código.")
	return upsertBlock("AGENTS.md", name, agentsBlock)
}

/**
* ModuleVersion: Returns the version of modulePath compiled into the running
* binary, or "devel" when it is a local checkout (workspace or replace).
* @param modulePath string
* @return string
**/
func ModuleVersion(modulePath string) string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "devel"
	}

	modules := append([]*debug.Module{&info.Main}, info.Deps...)
	for _, mod := range modules {
		if mod.Path != modulePath {
			continue
		}

		if mod.Replace != nil {
			mod = mod.Replace
		}

		if mod.Version == "" || mod.Version == "(devel)" {
			return "devel"
		}

		return mod.Version
	}

	return "devel"
}

/**
* managedBlock: Builds the marked block that points to the installed guides.
* @param name string
* @param version string
* @param paths []string
* @param lineFormat string, format applied to each guide path
* @return string
**/
func managedBlock(name, version string, paths []string, lineFormat string) string {
	lines := []string{
		fmt.Sprintf("<!-- %s:agents-guide:start %s -->", name, version),
		fmt.Sprintf("## Marco de trabajo para agentes: %s %s", name, version),
		"",
		fmt.Sprintf("Bloque generado por `go run github.com/celsiainternet/%s/cmd/install`. No lo edites: se reemplaza al volver a ejecutarlo. Las guías están en `.claude/guides/%s/` y tienen prioridad sobre lo que recuerdes de la librería; las reglas propias de este proyecto, fuera de este bloque, tienen prioridad sobre las guías.", name, name),
		"",
	}
	for _, path := range paths {
		lines = append(lines, fmt.Sprintf(lineFormat, path))
	}
	lines = append(lines, fmt.Sprintf("<!-- %s:agents-guide:end -->", name))

	return strings.Join(lines, "\n")
}

/**
* upsertBlock: Replaces the marked block of the library in path, appends it if
* the file has none, or creates the file with it.
* @param path string
* @param name string
* @param block string
* @return error
**/
func upsertBlock(path, name, block string) error {
	startMarker := fmt.Sprintf("<!-- %s:agents-guide:start", name)
	endMarker := fmt.Sprintf("<!-- %s:agents-guide:end -->", name)

	current, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return os.WriteFile(path, []byte(block+"\n"), 0644)
	}
	if err != nil {
		return err
	}

	content := string(current)
	start := strings.Index(content, startMarker)
	if start < 0 {
		content = strings.TrimRight(content, "\n") + "\n\n" + block + "\n"
		return os.WriteFile(path, []byte(content), 0644)
	}

	end := strings.Index(content[start:], endMarker)
	if end < 0 {
		return fmt.Errorf("%s: falta la marca %s", path, endMarker)
	}

	end = start + end + len(endMarker)
	content = content[:start] + block + content[end:]
	return os.WriteFile(path, []byte(content), 0644)
}
