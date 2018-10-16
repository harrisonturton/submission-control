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

type Handler func(argv []string, stop chan bool)

var commandHandlers = map[string]Handler{
	"exit": exitHandler,
}

// Run is the main loop of the CLI, it manages
// fetching and parsing input, and delegating
// functionality.
func Run(stop chan bool, wg *sync.WaitGroup) {
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
			handleInput(*input, stop)
		}
	}
}

// ReadInput simply reads from the terminal,
// and returns a space-seperated array of
// the input words.
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

// Handle input takes the CLI input and tries
// to find the appropriate handler based on the
// first word (seperated by spaces)
func handleInput(input []string, stop chan bool) {
	if len(input) == 0 {
		return
	}
	handler, ok := commandHandlers[input[0]]
	if !ok {
		fmt.Println(input[0] + ": unknown command")
		return
	}
	handler(input, stop)
}

// Pretty-print the welcome message.
func printWelcome() {
	fmt.Println("Hydra 0.0.1")
	fmt.Println("[Dev Branch " + time.Now().Format("Mon Jan 2006, 3:04pm") + "]")
	fmt.Println("\033[1;31mWarning: Hydra is still in development, and could be unstable.\033[0m")
}
