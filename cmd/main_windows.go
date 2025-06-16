package main

import (
	"fmt"
	"os/exec"
)

func LaunchDB() error {
	cmd := exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", "../script/launchDB.ps1")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Erreur :", err)
		return err
	}
	fmt.Println("Sortie :", string(output))
	return nil
}
