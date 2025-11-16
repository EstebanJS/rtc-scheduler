// pkg/logger/logger.go
package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level representa el nivel de logging
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// Logger define la interfaz para logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// SimpleLogger es una implementación simple de Logger
type SimpleLogger struct {
	level  Level
	logger *log.Logger
}

// New crea un nuevo logger con nivel INFO por defecto
func New() Logger {
	return &SimpleLogger{
		level:  InfoLevel,
		logger: log.New(os.Stdout, "", 0),
	}
}

// NewWithLevel crea un nuevo logger con un nivel específico
func NewWithLevel(level Level) Logger {
	return &SimpleLogger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// Debug registra un mensaje de nivel DEBUG
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.log(DebugLevel, msg, args...)
	}
}

// Info registra un mensaje de nivel INFO
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.log(InfoLevel, msg, args...)
	}
}

// Warn registra un mensaje de nivel WARN
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.log(WarnLevel, msg, args...)
	}
}

// Error registra un mensaje de nivel ERROR
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.log(ErrorLevel, msg, args...)
	}
}

// Fatal registra un mensaje de nivel FATAL y termina el programa
func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	l.log(FatalLevel, msg, args...)
	os.Exit(1)
}

// log formatea y registra un mensaje
func (l *SimpleLogger) log(level Level, msg string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	levelName := levelNames[level]
	
	// Formatear mensaje con colores
	coloredLevel := l.colorize(level, levelName)
	
	// Construir mensaje
	fullMsg := fmt.Sprintf("[%s] %s %s", timestamp, coloredLevel, msg)
	
	// Agregar argumentos adicionales si existen
	if len(args) > 0 {
		fullMsg += " " + l.formatArgs(args...)
	}
	
	l.logger.Println(fullMsg)
}

// formatArgs formatea argumentos adicionales como pares clave=valor
func (l *SimpleLogger) formatArgs(args ...interface{}) string {
	if len(args) == 0 {
		return ""
	}
	
	var result string
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			if i > 0 {
				result += " "
			}
			result += fmt.Sprintf("%v=%v", args[i], args[i+1])
		}
	}
	return result
}

// colorize agrega colores ANSI al nivel de log
func (l *SimpleLogger) colorize(level Level, text string) string {
	colors := map[Level]string{
		DebugLevel: "\033[36m", // Cyan
		InfoLevel:  "\033[32m", // Green
		WarnLevel:  "\033[33m", // Yellow
		ErrorLevel: "\033[31m", // Red
		FatalLevel: "\033[35m", // Magenta
	}
	
	reset := "\033[0m"
	
	if color, ok := colors[level]; ok {
		return color + text + reset
	}
	
	return text
}

// SetLevel cambia el nivel de logging
func (l *SimpleLogger) SetLevel(level Level) {
	l.level = level
}

// NoopLogger es un logger que no hace nada (útil para tests)
type NoopLogger struct{}

// NewNoop crea un logger que no registra nada
func NewNoop() Logger {
	return &NoopLogger{}
}

func (n *NoopLogger) Debug(msg string, args ...interface{}) {}
func (n *NoopLogger) Info(msg string, args ...interface{})  {}
func (n *NoopLogger) Warn(msg string, args ...interface{})  {}
func (n *NoopLogger) Error(msg string, args ...interface{}) {}
func (n *NoopLogger) Fatal(msg string, args ...interface{}) {
	os.Exit(1)
}