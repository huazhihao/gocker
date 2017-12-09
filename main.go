package main

import "os"

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]

	switch cmd {
	case "help":
		help()
	case "rm":
		Rm(args[0])
	case "ps":
		Ps()
	case "run":
		Run(args[0], args[1], args[2:]...)
	//case "exec":
	//Exec(args[0], args[1:]...)
	case "logs":
		Logs(args[0])

	case "images":
		Images()
	case "pull":
		Pull(args[0], args[1])
	//case "commit":
	//TOOD
	case "rmi":
		Rmi(args[0], args[1])
	default:
		help()
		os.Exit(1)
	}

}
