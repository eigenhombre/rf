package main

import "os/exec"

// PbCopy sends `data` to the Mac clipboard (e.g. for posting to HN)
func PbCopy(data string) {
	pbcopyCmd := exec.Command("pbcopy")
	pbcopyIn, _ := pbcopyCmd.StdinPipe()
	pbcopyCmd.Start()
	pbcopyIn.Write([]byte(data))
	pbcopyIn.Close()
	pbcopyCmd.Wait()
}
