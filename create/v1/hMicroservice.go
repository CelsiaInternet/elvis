package create

import "github.com/celsiainternet/elvis/agentsguide"

/**
* MkProject: Crea un proyecto completo. Primero escribe todos los archivos
* (microservicio, README, .env, .gitignore y soporte para agentes) y al final
* ejecuta go mod tidy, que es lo lento: si se interrumpe, no falta ningun archivo.
* @param packageName, name, author, schema string
* @return error
**/
func MkProject(packageName, name, author, schema string) error {
	err := makeMicroservice(packageName, name, schema)
	if err != nil {
		return err
	}

	err = MakeReadme(name)
	if err != nil {
		return err
	}

	err = MakeEnv(name)
	if err != nil {
		return err
	}

	err = MakeGitignore(name)
	if err != nil {
		return err
	}

	err = MakeAgents(name)
	if err != nil {
		return err
	}

	err = agentsguide.Install()
	if err != nil {
		return err
	}

	ProgressNext(20)
	GoModTidy()

	return nil
}

/**
* MkMicroservice: Agrega un microservicio al proyecto actual y ejecuta go mod tidy.
* @param packageName, name, schema string
* @return error
**/
func MkMicroservice(packageName, name, schema string) error {
	err := makeMicroservice(packageName, name, schema)
	if err != nil {
		return err
	}

	ProgressNext(20)
	GoModTidy()

	return nil
}

/**
* makeMicroservice: Escribe los archivos de un microservicio (cmd, deployments,
* internal, pkg, scripts, test, www). Avanza la barra hasta el 80%.
* @param packageName, name, schema string
* @return error
**/
func makeMicroservice(packageName, name, schema string) error {
	ProgressNext(10)
	err := MakeCmd(packageName, name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeDeployments(name)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeInternal(packageName, name, schema)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakePkg(name, schema)
	if err != nil {
		return err
	}

	ProgressNext(10)
	err = MakeScripts(name)
	if err != nil {
		return err
	}

	ProgressNext(15)
	err = MakeTest(name)
	if err != nil {
		return err
	}

	ProgressNext(15)
	err = MakeWWW(name)
	if err != nil {
		return err
	}

	return nil
}

/**
* MkMolue: Agrega el modelo, registra sus rutas y ejecuta go mod tidy.
* @param packageName, modelo, schema string
* @return bool (rutas registradas en router.go), error
**/
func MkMolue(packageName, modelo, schema string) (bool, error) {
	ProgressNext(10)
	routed, err := MakeModel(packageName, modelo, schema)
	if err != nil {
		return false, err
	}

	ProgressNext(80)
	GoModTidy()

	ProgressNext(10)

	return routed, nil
}

/**
* MkRpc
* @param name string
* @return error
**/
func MkRpc(name string) error {
	ProgressNext(10)
	err := MakeRpc(name)
	if err != nil {
		return err
	}

	ProgressNext(90)

	return nil
}

/**
* DeleteMicroservice
* @param name string
* @return error
**/
func DeleteMicroservice(name string) error {
	ProgressNext(50)
	err := DeleteCmd(name)
	if err != nil {
		return err
	}

	ProgressNext(50)

	return nil
}
