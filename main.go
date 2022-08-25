package main

import (
	"log"
)

func init() {
	Init()
}

func main() {
	log.Println("start...")
	Serve()
	log.Println("end...")
}
