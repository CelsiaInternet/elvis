package create

import (
	"encoding/json"
	"fmt"
	"go/format"
	"strings"
	"testing"
)

/**
* render: Sustituye $1, $2, ... igual que file.MakeFile.
* @param model string, args ...string
* @return string
**/
func render(model string, args ...string) string {
	result := model
	for i, v := range args {
		result = strings.ReplaceAll(result, fmt.Sprintf("$%d", i+1), v)
	}

	return result
}

/**
* TestGoTemplatesAreGofmtClean: Cada plantilla Go, con los argumentos que le
* pasa el generador, debe ser codigo valido y ya formateado con gofmt.
**/
func TestGoTemplatesAreGofmtClean(t *testing.T) {
	cases := []struct {
		name  string
		model string
		args  []string
	}{
		{"main.go", modelMain, []string{"github.com/acme/demo", "demo"}},
		{"service.go", modelService, []string{"github.com/acme/demo", "demo"}},
		{"api.go (schema)", modelDbApi, []string{"github.com/acme/demo", "demo"}},
		{"api.go", modelApi, []string{"github.com/acme/demo", "demo"}},
		{"event.go", modelEvent, []string{"demo"}},
		{"msg.go", modelMsg, []string{"demo"}},
		{"config.go", modelConfig, []string{"demo"}},
		{"rpc.go", modelhRpc, []string{"demo"}},
		{"controller.go (schema)", modelDbController, []string{"demo"}},
		{"controller.go", modelController, []string{"demo"}},
		{"schema.go", modelSchema, []string{"demo", "schemaVentas", "ventas"}},
		{"model.go", modelModel, []string{"demo", "Demo"}},
		{"hDemo.go (schema)", modelDbHandler, []string{"demo", "Demo", "schemaVentas", "DEMO", "demo"}},
		{"hDemo.go", modelHandler, []string{"demo", "Demo", "demo"}},
		{"router.go (schema)", modelDbRouter, []string{"demo", "Demo"}},
		{"router.go", modelRouter, []string{"demo", "demo"}},
	}

	for _, c := range cases {
		src := render(c.model, c.args...)
		got, err := format.Source([]byte(src))
		if err != nil {
			t.Errorf("%s: no es Go valido: %v", c.name, err)
			continue
		}

		if string(got) != src {
			t.Errorf("%s: no esta formateado con gofmt", c.name)
		}
	}
}

/**
* TestTemplatesHaveNoTokens: Las plantillas no deben traer JWT reales.
**/
func TestTemplatesHaveNoTokens(t *testing.T) {
	for name, model := range map[string]string{
		"restHttp":            restHttp,
		"modelDeploy":         modelDeploy,
		"modelOkeDeployment":  modelOkeDeployment,
		"modelOkeStatefulSet": modelOkeStatefulSet,
	} {
		if strings.Contains(model, "eyJhbGci") {
			t.Errorf("%s contiene un JWT", name)
		}
	}
}

/**
* TestEnvTemplate: El .env usa el nombre del servicio y trae las variables
* que el servicio generado lee al arrancar.
**/
func TestEnvTemplate(t *testing.T) {
	env := render(modelEnvar, "demo")
	for _, want := range []string{
		"APP=demo\n", "PORT=", "PATH_URL=/api/services/demo\n", "RPC_PORT=", "RPC_HOST=",
		"DB_DRIVER=", "DB_APPLICATION_NAME=demo\n", "REDIS_HOST=", "NATS_HOST=", "SECRET=",
	} {
		if !strings.Contains(env, want) {
			t.Errorf(".env no contiene %q", want)
		}
	}

	if strings.Contains(env, "$") {
		t.Errorf(".env tiene un placeholder sin reemplazar")
	}
}

/**
* TestAgentsTemplates: Las reglas para agentes estan en la plantilla y el
* settings.json de Claude Code es JSON valido que bloquea commit y push.
**/
func TestAgentsTemplates(t *testing.T) {
	rules := render(modelAgents, "demo")
	for _, want := range []string{
		"entendible, mantenible y escalable por un humano",
		"Los commits los hace el usuario",
		"go run ./cmd/demo",
	} {
		if !strings.Contains(rules, want) {
			t.Errorf("agents.md no contiene %q", want)
		}
	}

	var settings struct {
		Permissions struct {
			Deny []string `json:"deny"`
		} `json:"permissions"`
	}
	if err := json.Unmarshal([]byte(modelClaudeSettings), &settings); err != nil {
		t.Fatalf("settings.json no es JSON valido: %v", err)
	}

	deny := strings.Join(settings.Permissions.Deny, " ")
	for _, want := range []string{"Bash(git commit:*)", "Bash(git push:*)"} {
		if !strings.Contains(deny, want) {
			t.Errorf("settings.json no bloquea %s", want)
		}
	}
}
