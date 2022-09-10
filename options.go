package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	_defaultGithubUser = "_not_provided_"
)

type RemoteShellOptions struct {
	username        string
	timeLimit       time.Duration
	port            int
	idleTimeout     time.Duration
	homeDir         string
	shellCommand    string
	connectionGrace time.Duration
	insecureMode    bool
}

func determineDefaultShell() string {
	defaultShells := [5]string{"bash", "sh", "dash", "zsh", "rbash"}

	for _, shellExec := range defaultShells {
		shellPath, err := exec.LookPath(shellExec)
		if err != nil {
			continue
		}
		return shellPath
	}

	return "/tmp/no-shell-was-found"
}

func parseCommandFlags() RemoteShellOptions {
	var help = flag.Bool("help", false, "Show help")
	var version = flag.Bool("version", false, "Show version")
	var usernameFlag string
	var userHomeFlag string
	var shellCommandFlag string
	var portFlag int
	var graceFlag int
	var idleTimeFlag int
	var timeLimitFlag int
	var insecureModeFlag bool

	fallbackHomeDir, homeErr := os.UserHomeDir()
	if homeErr != nil {
		fallbackHomeDir = fetchEnvValue("HOME", fetchEnvValue("PWD", "/"))
	}
	fallbackUsername := fetchEnvValue("C87_RSHELL_GHUSER", _defaultGithubUser)
	fallbackShell := fetchEnvValue("C87_RSHELL_SHELL", determineDefaultShell())

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
	flag.BoolVar(&insecureModeFlag, "insecure", false, "Skip SSL verification")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Println(Version)
		os.Exit(0)
	}

	if usernameFlag == _defaultGithubUser {
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
		insecureMode:    insecureModeFlag,
	}
}
