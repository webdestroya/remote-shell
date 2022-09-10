package main

import (
	"log"
)

var Version = "development"

func main() {

	globalOptions := parseCommandFlags()

	log.Println("Starting Cloud87 Remote Shell Service")
	log.Println("")

	log.Println("")
	log.Println("Version:   ", Version)
	log.Println("")

	log.Println("GitHubUser:", globalOptions.username)
	log.Println("HomeDir:   ", globalOptions.homeDir)
	log.Println("ShellCmd:  ", globalOptions.shellCommand)
	log.Println("Port:      ", globalOptions.port)
	log.Println("IdleTime:  ", globalOptions.idleTimeout.String())
	log.Println("MaxTime:   ", globalOptions.timeLimit.String())
	log.Println("GraceTime: ", globalOptions.connectionGrace.String())
	if globalOptions.insecureMode {
		log.Println("Insecure:  YES!! HTTPS ENDPOINTS WILL NOT BE VERIFIED")
	}

	startSSHService(globalOptions)
}
