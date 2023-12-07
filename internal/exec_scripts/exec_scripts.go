package execscripts

import (
	"os"
	"os/exec"
)

func RunShellCommand() error {
	export_directory := "/"
	file_name := "produ:ct_export-2023-12-07 12:50:19.00291 +0000 UTC.csv"
	cmd := exec.Command("/bin/bash", "./scripts/runner.sh", export_directory, file_name)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
