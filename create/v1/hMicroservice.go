package create

import "github.com/celsiainternet/elvis/agentsguide"

/**
* MkProject
* @param packageName, name, author, schema string
* @return error
**/
func MkProject(packageName, name, author, schema string) error {
	ProgressNext(16)
	err := MkMicroservice(packageName, name, schema)
	if err != nil {
		return err
	}

	ProgressNext(16)
	err = MakeReadme(name)
	if err != nil {
		return err
	}

	ProgressNext(16)
	err = MakeEnv(name)
	if err != nil {
		return err
	}

	ProgressNext(16)
	err = MakeGitignore(name)
	if err != nil {
		return err
	}

	ProgressNext(16)
	err = agentsguide.Install()
	if err != nil {
		return err
	}

	ProgressNext(20)

	return nil
}

/**
* MkMicroservice
* @param packageName, name, schema string
* @return error
**/
func MkMicroservice(packageName, name, schema string) error {
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
* MkMolue
* @param packageName, modelo, schema string
* @return error
**/
func MkMolue(packageName, modelo, schema string) error {
	ProgressNext(10)
	err := MakeModel(packageName, modelo, schema)
	if err != nil {
		return err
	}

	ProgressNext(90)

	return nil
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
