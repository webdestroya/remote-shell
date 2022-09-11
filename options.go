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

type RemoteShellOptions struct {
	username string
	port     int

	timeLimit       time.Duration
	idleTimeout     time.Duration
	connectionGrace time.Duration

	homeDir      string
	shellCommand string

	insecureMode bool

	// TODO: should we allow multiple sessions?
	allowMultipleSessions bool
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
	var graceFlag time.Duration
	var idleTimeFlag time.Duration
	var timeLimitFlag time.Duration
	var insecureModeFlag bool

	fallbackHomeDir, err := os.UserHomeDir()
	if err != nil {
		fallbackHomeDir = fetchEnvValue("HOME", fetchEnvValue("PWD", "/"))
	}
	fallbackUsername := fetchEnvValue("C87RS_USER", "")
	fallbackShell := fetchEnvValue("C87RS_SHELL", "automatic")

	fallBackPort := fetchEnvValueInt("C87RS_PORT", 8722)
	fallBackGrace := fetchEnvValueDuration("C87RS_GRACE", 30*time.Minute)
	fallBackIdle := fetchEnvValueDuration("C87RS_IDLETIME", 0*time.Second)
	fallBackTimeLimit := fetchEnvValueDuration("C87RS_MAXTIME", 12*time.Hour)

	flag.StringVar(&usernameFlag, "user", fallbackUsername, "GitHub username")
	flag.StringVar(&userHomeFlag, "home", fallbackHomeDir, "Home Directory")
	flag.StringVar(&shellCommandFlag, "shell", fallbackShell, "Shell Command")

	flag.IntVar(&portFlag, "port", fallBackPort, "SSH port to use")

	flag.DurationVar(&graceFlag, "grace", fallBackGrace, "How many seconds to wait for a connection before we assume abandonment")
	flag.DurationVar(&idleTimeFlag, "idletime", fallBackIdle, "Idle timeout")
	flag.DurationVar(&timeLimitFlag, "maxtime", fallBackTimeLimit, "Max session duration")

	flag.BoolVar(&insecureModeFlag, "insecure", false, "Skip SSL verification")

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *version {
		fmt.Printf("%s@%s\n", buildVersion, buildSha)
		os.Exit(0)
	}

	if usernameFlag == "" {
		log.Fatal("You must provide a GitHub username")
	}

	userHomePath, err := filepath.Abs(userHomeFlag)
	check(err)

	if shellCommandFlag == "automatic" {
		shellCommandFlag = determineDefaultShell()
	}

	shellCommand, err := exec.LookPath(shellCommandFlag)
	check(err)

	return RemoteShellOptions{
		username:              usernameFlag,
		homeDir:               userHomePath,
		shellCommand:          shellCommand,
		port:                  portFlag,
		idleTimeout:           idleTimeFlag,
		timeLimit:             timeLimitFlag,
		connectionGrace:       graceFlag,
		insecureMode:          insecureModeFlag,
		allowMultipleSessions: false,
	}
}
