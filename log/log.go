package log

import (
	"log"
	"fmt"
	"os"
)

// The const value for log level.
// Use iota.
const (
	VERBOSE = iota + 1
	INFO
	DEBUG
	WARN
	ERROR
	PANIC
	FATAL
)

var _logger *Logger

// todo init log target
func init() {
	_logger = NewLogger()
}

type Logger struct {
	log       *log.Logger // Extent builtin log package.
	callDepth int         // Skip levels.
	maxLevel  uint        // Max log level. Default maxLevel = 0.
}

// New a log target return logger pointer.
func NewLogger() *Logger {
	return &Logger{
		log:       log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile),
		callDepth: 3,
		maxLevel:  0,
	}
}

// Set log output file.
// Default generate a log text file like 'xxxx-xxxx-xxxx.log' in project root path.
// it will be Named with current time.
func (logger *Logger) SetFile() {

}

// Set log output level.
// If log level greater than maxLevel, The call is forbidden.
func (logger *Logger) SetMaxLevel(maxLevel uint) {
	logger.maxLevel = maxLevel
}

// Set std font color for different log level.
// Support color name setting.
func setContentColor(content string, logLevel uint) string {
	if _, ok := colorHandler[logLevel]; ok {
		content = colorHandler[logLevel](content)
	}
	return content
}

// Call builtin logger's Output.
func (logger *Logger) OutPut(prefix string, logLevel uint, v ... interface{}) {
	if prefix != "" {
		logger.log.SetPrefix(prefix + " ")
	}
	// Goroutine safe by builtin log Output
	content := fmt.Sprint(v...)
	content = content[1:len(content)-1]
	content = setContentColor(content, logLevel)
	err := logger.log.Output(logger.callDepth, content)
	if err != nil {
		panic(err)
	}
	return
}

// Declare Arbitrarily log prefix.
// Level = 1.
func (logger *Logger) Verbose(prefix string, v ... interface{}) {
	logger.OutPut("["+prefix+"]", VERBOSE, v)
}

// Level = 2.
func (logger *Logger) Info(v ... interface{}) {
	logger.OutPut("[I]", INFO, v)
}

// Level = 3.
func (logger *Logger) Debug(v ... interface{}) {
	logger.OutPut("[D]", DEBUG, v)
}

// Level = 4.
func (logger *Logger) Warn(v ... interface{}) {
	logger.OutPut("[W]", WARN, v)
}

// Level = 5.
func (logger *Logger) Error(v ... interface{}) {
	logger.OutPut("[E]", ERROR, v)
}

// Level = 6.
func (logger *Logger) Panic(v ... interface{}) {
	logger.OutPut("[P]", PANIC, v)
	panic(v)
}

// Level = 7.
func (logger *Logger) Fatal(v ... interface{}) {
	logger.OutPut("[F]", FATAL, v)
	os.Exit(1)
}

// Output result by call fmt.Print.
func (logger *Logger) Print(v ... interface{}) {
	logger.OutPut("", 0, fmt.Sprint(v...))
}

// Output result format by call fmt.Println.
func (logger *Logger) PrintLn(v ... interface{}) {
	logger.OutPut("", 0, fmt.Sprintln(v...))
}

// Output result format by call fmt.Sprintf.
// Support format.
func (logger *Logger) Printf(format string, v ... interface{}) {
	logger.OutPut("", 0, fmt.Sprintf(format, v...))
}

// Define module rank function.
// Declare Arbitrarily log prefix.
// Level = 1.
func Verbose(prefix string, v ... interface{}) {
	_logger.OutPut("["+prefix+"]", VERBOSE, v)
}

// Level = 2.
func Info(v ... interface{}) {
	_logger.OutPut("[I]", INFO, v)
}

// Level = 3.
func Debug(v ... interface{}) {
	_logger.OutPut("[D]", DEBUG, v)
}

// Level = 4.
func Warn(v ... interface{}) {
	_logger.OutPut("[W]", WARN, v)
}

// Level = 5.
func Error(v ... interface{}) {
	_logger.OutPut("[E]", ERROR, v)
}

// Level = 6.
func Panic(v ... interface{}) {
	_logger.OutPut("[P]", PANIC, v)
	panic(v)
}

// Level = 7.
func Fatal(v ... interface{}) {
	_logger.OutPut("[F]", FATAL, v)
	os.Exit(1)
}

// Output result by call fmt.Print.
func Print(v ... interface{}) {
	_logger.OutPut("", 0, fmt.Sprint(v...))
}

// Output result format by call fmt.Println.
func PrintLn(v ... interface{}) {
	_logger.OutPut("", 0, fmt.Sprintln(v...))
}

// Output result format by call fmt.Sprintf.
// Support format.
func Printf(format string, v ... interface{}) {
	_logger.OutPut("", 0, fmt.Sprintf(format, v...))
}
