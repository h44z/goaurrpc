package rpc

import (
	"log/slog"
)

// Log writes a log message either to a file or stdout
func (s *Server) Log(msg string, args ...any) {
	slog.Info(msg, args...)
}

// LogVerbose writes log messages if the verbose flag is set
func (s *Server) LogVerbose(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// LogVeryVerbose writes log messages if the very verbose flag is set
func (s *Server) LogVeryVerbose(msg string, args ...any) {
	if s.veryVerbose {
		slog.Debug(msg, args...)
	}
}
