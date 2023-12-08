package execscripts

import (
	"os"
	"os/exec"
)

func CopyFileLocally(export_directory, file_name string) error {
	cmd := exec.Command("/bin/bash", "./scripts/copier.sh", export_directory, file_name)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
