package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Command int

const (
	Help Command = iota
	Exit
	Request
)

type Handler func(argv []string, stop chan bool, requests chan<- string)

var commandNames = map[string]Command{
	"help": Help,
	"exit": Exit,
	"req":  Request,
}

var commandHandlers = map[Command]Handler{
	Help:    helpHandler,
	Exit:    exitHandler,
	Request: requestHandler,
}

func Run(stop chan bool, wg *sync.WaitGroup, requests chan<- string) {
	defer wg.Done()
	printWelcome()
	for {
		select {
		case <-stop:
			fmt.Println("Stopping CLI...")
			return
		default:
			input, err := readInput()
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			handleInput(*input, stop, requests)
		}
	}
}

func handleInput(input []string, stop chan bool, requests chan<- string) {
	if len(input) == 0 {
		return
	}
	commandName, ok := commandNames[input[0]]
	if !ok {
		fmt.Println("Unknown command " + input[0])
		return
	}
	handler, ok := commandHandlers[commandName]
	if !ok {
		fmt.Println("Failed to handle " + input[0])
		return
	}
	handler(input, stop, requests)
}

func readInput() (*[]string, error) {
	buf := bufio.NewReader(os.Stdin)
	fmt.Print(">>> ")
	input, err := buf.ReadString('\n')
	if err != nil {
		return nil, err
	}
	inputFields := strings.Fields(input)
	return &inputFields, nil
}

func helpHandler(argv []string, stop chan bool, requests chan<- string) {
	fmt.Println("Help!")
}

func exitHandler(argv []string, stop chan bool, requests chan<- string) {
	fmt.Println("Stopping server...")
	close(stop)
}

func requestHandler(argv []string, stop chan bool, requests chan<- string) {
	id := uuid()
	fmt.Println(fmt.Sprintf("Making request [%s]", id))
	requests <- fmt.Sprintf("[Request %s]", id)
}

func printWelcome() {
	fmt.Println("Kerboros 0.0.1")
	fmt.Println("[Dev Branch " + time.Now().Format("Mon Jan 2006, 3:04pm") + "]")
	fmt.Println("\033[1;31mWarning: Kerboros is still in development, and could be unstable.\033[0m")
}

func uuid() string {
	f, _ := os.Open("/dev/urandom")
	b := make([]byte, 16)
	f.Read(b)
	f.Close()
	return fmt.Sprintf("%x", b[0:4])
}
