package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const defaultGithubUser = "_not_provided_"

type RemoteShellOptions struct {
	username        string
	timeLimit       time.Duration
	port            int
	idleTimeout     time.Duration
	homeDir         string
	shellCommand    string
	connectionGrace time.Duration
}

var globalOptions RemoteShellOptions

func parseCommandFlags() RemoteShellOptions {
	var help = flag.Bool("help", false, "Show help")
	var usernameFlag string
	var userHomeFlag string
	var shellCommandFlag string
	var portFlag int
	var graceFlag int
	var idleTimeFlag int
	var timeLimitFlag int

	fallbackHomeDir := fetchEnvValue("HOME", fetchEnvValue("PWD", "/"))
	fallbackUsername := fetchEnvValue("C87_RSHELL_GHUSER", defaultGithubUser)
	fallbackShell := fetchEnvValue("C87_RSHELL_SHELL", "/bin/sh")

	fallBackPort := fetchEnvValueInt("C87_RSHELL_PORT", 8722)
	fallBackGrace := fetchEnvValueInt("C87_RSHELL_GRACE", 600)
	fallBackIdle := fetchEnvValueInt("C87_RSHELL_IDLE_TIMEOUT", 0)
	fallBackTimeLimit := fetchEnvValueInt("C87_RSHELL_MAX_RUNTIME", 43200)

	flag.StringVar(&usernameFlag, "user", fallbackUsername, "GitHub username")
	flag.StringVar(&userHomeFlag, "home", fallbackHomeDir, "Home Directory")
	flag.StringVar(&shellCommandFlag, "shell", fallbackShell, "Shell Command")
	flag.IntVar(&portFlag, "port", fallBackPort, "SSH port to use")
	flag.IntVar(&graceFlag, "grace", fallBackGrace, "How many seconds to wait for a connection before we assume abandonment")
	flag.IntVar(&idleTimeFlag, "idletime", fallBackIdle, "Idle timeout")
	flag.IntVar(&timeLimitFlag, "maxtime", fallBackTimeLimit, "Max session duration")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if usernameFlag == defaultGithubUser {
		log.Fatal("You must provide a GitHub username")
	}

	userHomePath, uHomeErr := filepath.Abs(userHomeFlag)
	check(uHomeErr)

	shellCommand, shellErr := exec.LookPath(shellCommandFlag)
	check(shellErr)

	return RemoteShellOptions{
		username:        usernameFlag,
		homeDir:         userHomePath,
		shellCommand:    shellCommand,
		port:            portFlag,
		idleTimeout:     (time.Duration(idleTimeFlag) * time.Second),
		timeLimit:       (time.Duration(timeLimitFlag) * time.Second),
		connectionGrace: (time.Duration(graceFlag) * time.Second),
	}
}

// https://github.com/gliderlabs/ssh/blob/master/_examples/ssh-publickey/public_key.go

func main() {

	globalOptions = parseCommandFlags()

	log.Println("GitHubUser:", globalOptions.username)
	log.Println("HomeDir:   ", globalOptions.homeDir)
	log.Println("ShellCmd:  ", globalOptions.shellCommand)
	log.Println("Port:      ", globalOptions.port)
	log.Println("IdleTime:  ", globalOptions.idleTimeout.String())
	log.Println("MaxTime:   ", globalOptions.timeLimit.String())
	log.Println("GraceTime: ", globalOptions.connectionGrace.String())

	startSSHService(globalOptions)
}
