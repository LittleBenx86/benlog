package core

import (
	"github.com/LittleBenx86/Benlog/internal/global/consts"
	"github.com/LittleBenx86/Benlog/internal/global/variables"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGTERM) // listen all exit signals
		receivedSignal := <-ch
		variables.Logger.Warn(consts.SIGNAL_PROCESS_KILL, zap.String("signal value", receivedSignal.String()))
		close(ch)
		os.Exit(1)
	}()
}
