package main

import (
	"fmt"
	"os/exec"
)

func launchDB() error {
	cmd := exec.Command("bash", "../script/launchDB.sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erreur :", err)
		return err
	}
	fmt.Println("Sortie :", string(output))
	return nil
}
