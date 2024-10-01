package main

import (
	"buttplugosu/internal/gameplay"
)

func main() {
	go gameplay.HandlePlug()
	gameplay.Init()
}
