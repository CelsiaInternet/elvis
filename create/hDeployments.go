package create

import "github.com/cgalvisleon/elvis/utility"

func MakeDeployments(name string) error {
	_, err := utility.MakeFolder("deployments", "dev")
	if err != nil {
		return err
	}

	_, err = utility.MakeFolder("deployments", "local")
	if err != nil {
		return err
	}

	_, err = utility.MakeFolder("deployments", "prd")
	if err != nil {
		return err
	}

	return nil
}
