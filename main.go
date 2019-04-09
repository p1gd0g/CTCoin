package main

import (
	"bufio"
	"os"

	"github.com/p1gd0g/CTCoin/node"
)

func main() {

	var node node.Node
	node.StartConn()

	node.Account.New()

	reader := bufio.NewReader(os.Stdin)
	for {
		command, _, _ := reader.ReadLine()
		node.HandleCommand(string(command))
	}
}
