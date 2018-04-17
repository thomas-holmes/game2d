package main

import (
	"fmt"

	"github.com/thomas-holmes/game2d/pkg/game"
)

func main() {
	fmt.Println("Hello")

	game := game.NewGame()

	err := game.Init()

	if err != nil {
		panic(err)
	}

	game.Run()
}
