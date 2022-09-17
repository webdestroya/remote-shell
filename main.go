package main

import (
	"log"
	"os"
)

var (
	buildVersion = "development"
	buildSha     = "devel"
)

func main() {

	defer func() {
		exitCode := 0
		if r := recover(); r != nil {
			exitCode = 1
			log.Println("FATAL ERROR! ", r)
		}
		os.Exit(exitCode)
	}()

	log.Println("Starting Cloud87 Remote Shell Service")
	log.Println("")

	log.Println("")
	log.Printf("Version: %s@%s\n", buildVersion, buildSha)
	log.Println("")

	if globalOptions.username != "" {
		log.Println("GitHubUser:", globalOptions.username)
	}
	log.Println("HomeDir:   ", globalOptions.homeDir)
	log.Println("ShellCmd:  ", globalOptions.shellCommand)
	log.Println("Port:      ", globalOptions.port)
	log.Println("IdleTime:  ", globalOptions.idleTimeout.String())
	log.Println("MaxTime:   ", globalOptions.timeLimit.String())
	log.Println("GraceTime: ", globalOptions.connectionGrace.String())
	if globalOptions.insecureMode {
		log.Println("Insecure:  YES!! HTTPS ENDPOINTS WILL NOT BE VERIFIED")
	}

	startSSHService()
}
