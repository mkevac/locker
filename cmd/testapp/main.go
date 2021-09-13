package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkevac/locker"
)

func main() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT)

	locker, err := locker.NewLocker(&locker.Config{
		ConsulAddress: "localhost:8500",
		Key:           "testkey",
		Value:         "testvalue",
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("locking...")

	resultCh, err := locker.Lock(nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("lock acquired")

	select {
	case <-resultCh:
		log.Println("lost lock, exiting")
		return
	case <-signalCh:
		log.Println("unlocking lock...")
		_ = locker.Unlock()
		log.Println("unlocked")
		return
	}
}
