package proc

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/techniumlabs/cinit/pkg/utils"
)

// RemoveZombies Removes any zombie process
func RemoveZombies(ctx context.Context, wg *sync.WaitGroup) {
	for {
		var status syscall.WaitStatus

		// Wait for orphaned zombie process
		pid, _ := syscall.Wait4(-1, &status, syscall.WNOHANG, nil)

		if pid <= 0 {
			// PID is 0 or -1 if no child waiting
			// so we wait for 1 second for next check
			time.Sleep(1 * time.Second)
		} else {
			// PID is > 0 if a child was reaped
			// we immediately check if another one
			// is waiting
			continue
		}

		// Non-blocking test
		// if context is done
		select {
		case <-ctx.Done():
			// Context is done
			// so we stop goroutine
			wg.Done()
			return
		default:
		}
	}
}

// Run function is to run a given command
func Run(commandStr string, argsSlice []string, envs map[string]string) error {

	// Register chan to receive system signals
	sigs := make(chan os.Signal, 1)
	defer close(sigs)
	signal.Notify(sigs)
	defer signal.Reset()

	// Define command and rebind
	// stdout and stdin
	log.Printf("%s", argsSlice)
	cmd := exec.Command(commandStr, argsSlice...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = utils.ConvertMapToEnv(envs)
	// Create a dedicated pidgroup
	// used to forward signals to
	// main process and all children
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Goroutine for signals forwarding
	go func() {
		for sig := range sigs {
			// Ignore SIGCHLD signals since
			// thez are only usefull for go-init
			if sig != syscall.SIGCHLD {
				// Forward signal to main process and all children
				syscall.Kill(-cmd.Process.Pid, sig.(syscall.Signal))
			}
		}
	}()

	// Start defined command
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for command to exit
	err = cmd.Wait()
	if err != nil {
		return err
	}

	return nil
}

func CleanQuit(cancel context.CancelFunc, wg *sync.WaitGroup, code int) {
	// Signal zombie goroutine to stop
	// and wait for it to release waitgroup
	cancel()
	wg.Wait()

	os.Exit(code)
}
