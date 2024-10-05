package config

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"

	"github.com/mviner000/eyymi/eyygo/constants"
)

var debugLogger *log.Logger

func init() {
	// Initialize debugLogger immediately with default output
	debugLogger = log.New(os.Stdout, constants.ColorCyan+"DEBUG: "+constants.ColorReset, log.Ldate|log.Ltime|log.Lshortfile)
}

func InitLogger(debug bool) {
	if !debug {
		debugLogger.SetOutput(io.Discard)
		log.SetOutput(io.Discard)
	} else {
		// Ensure that the default logger also outputs to stdout if in debug mode
		log.SetOutput(os.Stdout)
	}
}

// DebugLogf is a wrapper that formats the message and uses Println to ensure a newline
func DebugLogf(format string, v ...interface{}) {
	if debugLogger != nil {
		message := fmt.Sprintf(format, v...)
		debugLogger.Println(constants.ColorGreen + message + constants.ColorReset) // Use Println to ensure a newline
	}
}

// DebugLog is a simple wrapper for Println
func DebugLog(v ...interface{}) {
	if debugLogger != nil {
		message := fmt.Sprint(v...)
		debugLogger.Println(constants.ColorGreen + message + constants.ColorReset) // Use Println to ensure a newline
	}
}

// LogStruct logs the fields of a struct dynamically, handling nested structs and maps
func LogStruct(name string, s interface{}) {
	debugLogger.Println(constants.ColorYellow + name + constants.ColorReset)
	logStructFields(s, "  ")
}

func logStructFields(s interface{}, indent string) {
	value := reflect.ValueOf(s)
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldName := value.Type().Field(i).Name

		switch field.Kind() {
		case reflect.Struct:
			debugLogger.Printf("%s%s%s:%s", indent, constants.ColorBlue, fieldName, constants.ColorReset)
			logStructFields(field.Interface(), indent+"  ")
		case reflect.Map:
			if field.Len() == 0 {
				continue // Skip logging empty maps
			}
			debugLogger.Printf("%s%s%s: {", indent, constants.ColorBlue, fieldName)
			for _, key := range field.MapKeys() {
				val := field.MapIndex(key)
				debugLogger.Printf("%s  %v: %v", indent, key, val)
			}
			debugLogger.Printf("%s}", indent)
		default:
			debugLogger.Printf("%s%s%s: %v%s", indent, constants.ColorBlue, fieldName, field.Interface(), constants.ColorReset)
		}
	}
}

// CheckAllowedOrigins checks for wildcard in AllowedOrigins and logs a warning if in production
func CheckAllowedOrigins(settings SettingsStruct) {
	if !settings.IsDevelopment && strings.Contains(settings.AllowedOrigins, "*") {
		debugLogger.Println(constants.ColorRed + "WARNING: Using wildcard '*' in AllowedOrigins in production is not recommended!" + constants.ColorReset)
	}
}
