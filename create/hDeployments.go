package create

import "github.com/cgalvisleon/elvis/utilities"

func MakeDeployments(name string) error {
	_, err := utilities.MakeFolder("deployments", "dev")
	if err != nil {
		return err
	}

	_, err = utilities.MakeFolder("deployments", "local")
	if err != nil {
		return err
	}

	_, err = utilities.MakeFolder("deployments", "prd")
	if err != nil {
		return err
	}

	return nil
}
