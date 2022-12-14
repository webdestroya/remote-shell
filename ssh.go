package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
	gossh "golang.org/x/crypto/ssh"
)

// has someone connected yet?
var hasConnection bool = false

func setWinsize(f *os.File, w, h int) {
	// best effort set the window size

	//nolint:errcheck
	syscall.Syscall(
		syscall.SYS_IOCTL,
		f.Fd(),
		uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

//nolint:errcheck
func dismissSessionHandler(s ssh.Session) {
	io.WriteString(s, "This session has already been connected to. You cannot have multiple connections to the same session.\n")
	s.Exit(1)
}

func sessionHandler(options *RemoteShellOptions, notify chan bool, s ssh.Session) {
	log.Printf("Session User (requested): %s\n", s.User())
	log.Printf("Session RemoteIP: %s\n", s.RemoteAddr().String())
	log.Printf("Session ClientVersion: %s\n", s.Context().ClientVersion())
	log.Printf("Session ServerVersion: %s\n", s.Context().ServerVersion())
	log.Printf("Session SessionID: %s\n", s.Context().SessionID())

	hasConnection = true

	cmd := exec.Command(options.shellCommand)

	ptyReq, winCh, isPty := s.Pty()
	if isPty {
		log.Println("Starting PTY Session")

		cmd.Env = append(filteredEnvironmentVars(),
			fmt.Sprintf("TERM=%s", ptyReq.Term),
			fmt.Sprintf("HOME=%s", options.homeDir),
			fmt.Sprintf("USER=%s", options.currentUserName),
			fmt.Sprintf("LOGNAME=%s", options.currentUserName),
			fmt.Sprintf("SHELL=%s", options.shellCommand),
			fmt.Sprintf("C87RS_SESSIONID=%s", s.Context().SessionID()),
		)
		f, err := pty.Start(cmd)
		check(err)

		go func() {
			for win := range winCh {
				setWinsize(f, win.Width, win.Height)
			}
		}()

		// if these error, then it will abort the session.

		go func() {
			io.Copy(f, s) // stdin
		}()

		io.Copy(s, f) // stdout

		log.Println("Shell command output has ended. Waiting for command to end.")

		err = cmd.Wait()
		if err != nil {
			log.Println("The requested shell errored out. Are you sure it was correct?")
			s.Exit(1)
		} else {
			log.Println("Shell command terminated successfully")
			s.Exit(0)
		}
	} else {
		log.Println("NoPTY requested")
		io.WriteString(s, "No PTY requested.\n")
		s.Exit(1)
	}
	log.Println("Session ended")
	notify <- true
}

func openSSHSocket(options *RemoteShellOptions) net.Listener {
	connAddr := fmt.Sprintf(":%d", options.port)

	// lc := net.ListenConfig{
	// 	C
	// }

	sock, err := net.Listen("tcp", connAddr)
	check(err)
	return sock
}

func startSSHService() {
	publicKeys := exportAuthorizedKeys()

	if len(publicKeys) == 0 {
		panic("No SSH keys were provided")
	}

	log.Println("Starting SSH Service")

	notificationChannel := make(chan bool)

	ssh.Handle(func(s ssh.Session) {
		// log.Println("Handler invoked")
		if hasConnection {
			dismissSessionHandler(s)
		} else {
			sessionHandler(&globalOptions, notificationChannel, s)
		}
	})

	// Auth Request
	pubKeyAuthHandle := func(ctx ssh.Context, key ssh.PublicKey) bool {
		for _, pubKey := range publicKeys {
			if ssh.KeysEqual(key, pubKey) {
				return true
			}
		}
		return false
	}

	serverConfCallback := func(ctx ssh.Context) *gossh.ServerConfig {
		return &gossh.ServerConfig{
			NoClientAuth: false,
			BannerCallback: func(conn gossh.ConnMetadata) string {
				return "You are now connected to Cloud87 Remote Shell instance.\n\n"
			},
			AuthLogCallback: func(conn gossh.ConnMetadata, method string, err error) {
				log.Printf("Auth Attempt: %s method=%s err=%s\n", conn.RemoteAddr().String(), method, err)
			},
		}
	}

	sessionRequestCallback := func(sess ssh.Session, requestType string) bool {
		log.Printf("Session Requested by %s for type=%s\n", sess.Context().SessionID(), requestType)
		return true
	}

	server := &ssh.Server{
		Version:                       fmt.Sprintf("Cloud87-%s", buildVersion),
		PublicKeyHandler:              pubKeyAuthHandle,
		PasswordHandler:               nil,
		KeyboardInteractiveHandler:    nil,
		LocalPortForwardingCallback:   nil,
		ReversePortForwardingCallback: nil,
		SessionRequestCallback:        sessionRequestCallback,
		ServerConfigCallback:          serverConfCallback,
		MaxTimeout:                    globalOptions.timeLimit,
		IdleTimeout:                   globalOptions.idleTimeout,
	}

	log.Println("Creating socket")
	sock := openSSHSocket(&globalOptions)

	// boot the server socket
	go func() {
		log.Println("Booting SSH Listener")
		err := server.Serve(sock)
		// if err != nil && err != ssh.ErrServerClosed {
		if err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			panic(err)
		}
		log.Println("Serve goroutine ending")
		notificationChannel <- true
	}()

	go func() {
		log.Printf("Starting connection grace timer: %s\n", globalOptions.connectionGrace)
		time.Sleep(globalOptions.connectionGrace)

		if hasConnection {
			log.Println("Server has received a connection within grace period. Not killing self.")
			return
		}

		log.Println("Triggering socket shutdown. No longer accepting connections")

		err := server.Shutdown(context.Background())
		check(err)
		notificationChannel <- true
	}()

	select {
	case <-time.After(globalOptions.timeLimit):
		log.Println("Reached deadline for service. Dying")
		checkNoPanic(server.Close())
	case <-notificationChannel:
		log.Println("Death has been requested!")
	}

	// <-notificationChannel

	log.Println("End of startSSHService")

}

// https://gist.github.com/protosam/53cf7970e17e06135f1622fa9955415f
