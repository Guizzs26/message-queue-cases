package main

import (
	"log"
	"time"
)

func main() {
	log.Println("Notification worker starting...")

	for {
		log.Println("Worker is alive and waiting for messages...")
		time.Sleep(10 * time.Second)
	}
}
