package create

import utl "github.com/cgalvisleon/elvis/utilities"

func MakeDeployments(name string) error {
	_, err := utl.MakeFolder("deployments", "dev")
	if err != nil {
		return err
	}

	_, err = utl.MakeFolder("deployments", "local")
	if err != nil {
		return err
	}

	_, err = utl.MakeFolder("deployments", "prd")
	if err != nil {
		return err
	}

	return nil
}
