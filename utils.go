package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
)

func spawnProcess(cmd string, args []string, app *fiber.App) error {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		cmd := exec.Command(cmd, args...)
		cmd.Stderr = os.Stderr

		// use a preocess group since otherwise the node processes spawned by bun wont die
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		err := cmd.Start()
		if err != nil {
			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
				fmt.Println("Error sending SIGTERM to sub process group:", err)
			}

			fmt.Println("Error starting sub process:", err)
			return
		}

		go func() {
			<-shutdown

			if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM); err != nil {
				fmt.Println("Error sending SIGTERM to sub process group:", err)
			}
		}()

		if err := cmd.Wait(); err != nil {
			fmt.Println("Error waiting for sub process to exit:", err)
		}

		fmt.Println("sub process server stopped")

		if err := app.Shutdown(); err != nil {
			fmt.Println("Error shutting down HTTP server:", err)
		}
	}()

	return nil
}
