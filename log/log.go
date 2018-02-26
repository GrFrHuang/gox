package log

import (
	"log"
	"fmt"
)

//The const value for log level
//Use iota
const (
	VERBOSE = iota + 1
	INFO
	DEBUG
	WARN
	ERROR
	PANIC
	ASSERT
)

//todo Just support std and file log
type Logger struct {
	log       *log.Logger //Extent builtin log package
	callDepth int         //Skip levels
	maxLevel  int         //Max log level. Default maxLevel = 0
}

//New a log target return logger pointer
func NewLogger() *Logger {
	return &Logger{
		log:       log.New(),
		callDepth: 3,
		maxLevel:  0,
	}
}

//Set log output file
//Default generate a log text file like 'xxxx-xxxx-xxxx.log' in project root path
//it will be Named with current time
func (logger *Logger) SetFile() {

}

//Call builtin logger's Output
func (logger *Logger) OutPut(prefix string, logLevel uint, v ... interface{}) {
	err := logger.log.Output(logger.callDepth, fmt.Sprintf("%v", v))
	if err != nil {
		panic(err)
	}
	return
}

//Declare Arbitrarily log prefix
//Level = 1
func (logger *Logger) Verbose(prefix string, v ... interface{}) {
	logger.OutPut("["+prefix+"]", VERBOSE, v)
}

//Level = 2
func (logger *Logger) Info(v ... interface{}) {
	logger.OutPut("[I]", INFO, v)
}

//Level = 3
func (logger *Logger) Debug(v ... interface{}) {
	logger.OutPut("[D]", DEBUG, v)
}

//Level = 4
func (logger *Logger) Warn(v ... interface{}) {
	logger.OutPut("[W]", WARN, v)
}

//Level = 5
func (logger *Logger) Error(v ... interface{}) {
	logger.OutPut("[E]", ERROR, v)
}

//Level = 6
func (logger *Logger) Panic(v ... interface{}) {
	logger.OutPut("[P]", PANIC, v)
	panic(v)
}

//Level = 7
func (logger *Logger) Assert(v ... interface{}) {
	logger.OutPut("[A]", ASSERT, v)
}

//Output result by call fmt.Print
func (logger *Logger) Print(v ... interface{}) {
	logger.OutPut("", 0, fmt.Print(v...))
}

//Output result format by call fmt.Println
func (logger *Logger) PrintLn(v ... interface{}) {
	logger.OutPut("", 0, fmt.Println(v...))
}

//Output result format by call fmt.Sprintf
func (logger *Logger) Printf(format string, v ... interface{}) {
	logger.OutPut("", 0, fmt.Sprintf(format, v...))
}
