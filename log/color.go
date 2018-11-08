// Set the log output content color of terminal.
// The last settings will cover previous declare.
// The color scope between 30~37.
// The background color default black.
// Current support color name.

package log

import "fmt"

var colorHandler map[uint]func(string) string

func init() {
	colorHandler = make(map[uint]func(string string) string)
}

// Play up content color to be black.
func SetBlack(logLevel uint) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[30m%s\x1b[0m", message)
	}
}

// Play up content color to be red.
func SetRed(logLevel uint, color string) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[31m%s\x1b[0m", message)
	}
}

// Play up content color to be green.
func SetGreen(logLevel uint) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[32m%s\x1b[0m", message)
	}
}

// Play up content color to be yellow.
func SetYellow(logLevel uint) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[33m%s\x1b[0m", message)
	}
}

// Play up content color to be blue.
func SetBlue(logLevel uint) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[34m%s\x1b[0m", message)
	}
}

// Play up content color to be magenta.
func SetMagenta(logLevel uint, color string) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[35m%s\x1b[0m", message)
	}

}

// Play up content color to be cyan.
func SetCyan(logLevel uint, color string) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[36m%s\x1b[0m", message)
	}

}

// Play up content color to be white.
func SetWhite(logLevel uint) {
	colorHandler[logLevel] = func(message string) string {
		return fmt.Sprintf("\x1b[37m%s\x1b[0m", message)
	}
}
