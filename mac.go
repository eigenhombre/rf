package main

import (
	"log"
	"os/exec"
)

// PbCopy sends `data` to the Mac clipboard (e.g. for posting to HN)
func PbCopy(data string) {
	pbcopyCmd := exec.Command("pbcopy")
	pbcopyIn, _ := pbcopyCmd.StdinPipe()
	pbcopyCmd.Start()
	pbcopyIn.Write([]byte(data))
	pbcopyIn.Close()
	pbcopyCmd.Wait()
}

// macOpen calls the open command on a URL or a file
func macOpen(target string) {
	cmd := exec.Command("open", target)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
