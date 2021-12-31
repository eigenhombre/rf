package main

import (
	"fmt"
	"os/exec"
)

// pbCopy sends `data` to the Mac clipboard (e.g. for posting to HN)
func pbCopy(data string) {
	pbcopyCmd := exec.Command("pbcopy")
	pbcopyIn, _ := pbcopyCmd.StdinPipe()
	pbcopyCmd.Start()
	pbcopyIn.Write([]byte(data))
	pbcopyIn.Close()
	pbcopyCmd.Wait()
}

// macOpen calls the open command on a URL or a file
func macOpen(target string) error {
	cmd := exec.Command("open", target)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("unable to open %s", target)
	}
	return nil
}
