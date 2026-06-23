package tools

import (
	"log"
	"math/rand"
	"time"
)

func SendPasta(who string) {
	log.Println("Sending Pasta to: ", who)
	time.Sleep(time.Duration(rand.Intn(31)) * time.Second)
	log.Println("Pasta Sent Successfully")
	log.Printf("%s is eating the pasta", who)
	time.Sleep(time.Duration(rand.Intn(5)) * time.Second)
	log.Println("Pasta has been eaten, they thank you!")
}
