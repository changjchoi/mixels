package shared_test

import (
	"../shared"
	"log"
	"testing"
)

func TestConfig(t *testing.T) {
	if shared.Config.Init("../config.json") == false {
		log.Fatal("map config.json loading error")
	}
	shared.Config.Print()
	log.Println("Level2File", shared.Config.Level2File)
}
