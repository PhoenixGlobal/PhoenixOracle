package test

import (
	"os"
	"syscall"
	"testing"

	. "github.com/onsi/gomega"
)

func TestGracefulShutdown(t *testing.T) {
	t.Parallel()
	RegisterTestingT(t)
	store := Store()
	defer store.Close()

	var completed bool
	store.Exiter = func(code int) {
		completed = true
	}

	store.Start()
	pro, _ := os.FindProcess(syscall.Getpid())

	pro.Kill()

	Eventually(func() bool {
		return completed
	}).Should(BeTrue())
}
