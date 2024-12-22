/*
 * Copyright 2024 Ardan Labs
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * This file is part of the [Service] project by Ardan Labs.
 * Repository URL: https://github.com/ardanlabs/service
 *
 * Changes Made:
 * - Stephen O'Dwyer Combined logger.go handler.go and model.go foundation/logger package
 * For more information, see the repository's changelog or commit history.
 */
// Package logger provides support for initializing the log system.

package secretsengine

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"runtime"
	"time"

	"log/slog"
)

// -------------------------------------------------------------------------
// logger

// TraceIDFn represents a function that can return the trace id from
// the specified context.
type TraceIDFn func(ctx context.Context) string

// Logger represents a logger for logging information.
type Logger struct {
	handler   slog.Handler
	traceIDFn TraceIDFn
}

// New constructs a new log for application use.
func NewLogger(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn) *Logger {
	return newLogger(w, minLevel, serviceName, traceIDFn, Events{})
}

// NewWithEvents constructs a new log for application use with events.
func NewWithEvents(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {
	return newLogger(w, minLevel, serviceName, traceIDFn, events)
}

// NewWithHandler returns a new log for application use with the underlying
// handler.
func NewWithHandler(h slog.Handler) *Logger {
	return &Logger{handler: h}
}

// NewStdLogger returns a standard library Logger that wraps the slog Logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(logger.handler, slog.Level(level))
}

// Debug logs at LevelDebug with the given context.
func (log *Logger) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelDebug, 3, msg, args...)
}

// Debugc logs the information at the specified call stack position.
func (log *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelDebug, caller, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (log *Logger) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack position.
func (log *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelInfo, caller, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (log *Logger) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelWarn, 3, msg, args...)
}

// Warnc logs the information at the specified call stack position.
func (log *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelWarn, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (log *Logger) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelError, 3, msg, args...)
}

// Errorc logs the information at the specified call stack position.
func (log *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelError, caller, msg, args...)
}

func (log *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	slogLevel := slog.Level(level)

	if !log.handler.Enabled(ctx, slogLevel) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), slogLevel, msg, pcs[0])

	if log.traceIDFn != nil {
		args = append(args, "trace_id", log.traceIDFn(ctx))
	}
	r.Add(args...)

	log.handler.Handle(ctx, r)
}

func newLogger(w io.Writer, minLevel Level, serviceName string, traceIDFn TraceIDFn, events Events) *Logger {

	// Convert the file name to just the name.ext when this key/value will
	// be logged.
	f := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			if source, ok := a.Value.Any().(*slog.Source); ok {
				v := fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line)
				return slog.Attr{Key: "file", Value: slog.StringValue(v)}
			}
		}

		return a
	}

	// Construct the slog JSON handler for use.
	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.Level(minLevel), ReplaceAttr: f}))

	// If events are to be processed, wrap the JSON handler around the custom
	// log handler.
	if events.Debug != nil || events.Info != nil || events.Warn != nil || events.Error != nil {
		handler = newLogHandler(handler, events)
	}

	// Attributes to add to every log.
	attrs := []slog.Attr{
		{Key: "service", Value: slog.StringValue(serviceName)},
	}

	// Add those attributes and capture the final handler.
	handler = handler.WithAttrs(attrs)

	return &Logger{
		handler:   handler,
		traceIDFn: traceIDFn,
	}
}

//---------------------------------------------------------------------
// handler

// logHandler provides a wrapper around the slog handler to capture which
// log level is being logged for event handling.
type logHandler struct {
	handler slog.Handler
	events  Events
}

func newLogHandler(handler slog.Handler, events Events) *logHandler {
	return &logHandler{
		handler: handler,
		events:  events,
	}
}

// Enabled reports whether the handler handles records at the given level.
// The handler ignores records whose level is lower.
func (h *logHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// WithAttrs returns a new JSONHandler whose attributes consists
// of h's attributes followed by attrs.
func (h *logHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &logHandler{handler: h.handler.WithAttrs(attrs), events: h.events}
}

// WithGroup returns a new Handler with the given group appended to the receiver's
// existing groups. The keys of all subsequent attributes, whether added by With
// or in a Record, should be qualified by the sequence of group names.
func (h *logHandler) WithGroup(name string) slog.Handler {
	return &logHandler{handler: h.handler.WithGroup(name), events: h.events}
}

// Handle looks to see if an event function needs to be executed for a given
// log level and then formats its argument Record.
func (h *logHandler) Handle(ctx context.Context, r slog.Record) error {
	switch r.Level {
	case slog.LevelDebug:
		if h.events.Debug != nil {
			h.events.Debug(ctx, toRecord(r))
		}

	case slog.LevelError:
		if h.events.Error != nil {
			h.events.Error(ctx, toRecord(r))
		}

	case slog.LevelWarn:
		if h.events.Warn != nil {
			h.events.Warn(ctx, toRecord(r))
		}

	case slog.LevelInfo:
		if h.events.Info != nil {
			h.events.Info(ctx, toRecord(r))
		}
	}

	return h.handler.Handle(ctx, r)
}

//---------------------------------------------------------------------
// model

// Level represents different logging levels.
type Level slog.Level

// A set of possible logging levels.
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// Record represents the data that is being logged.
type Record struct {
	Time       time.Time
	Message    string
	Level      Level
	Attributes map[string]any
}

func toRecord(r slog.Record) Record {
	atts := make(map[string]any, r.NumAttrs())

	f := func(attr slog.Attr) bool {
		atts[attr.Key] = attr.Value.Any()
		return true
	}
	r.Attrs(f)

	return Record{
		Time:       r.Time,
		Message:    r.Message,
		Level:      Level(r.Level),
		Attributes: atts,
	}
}

// EventFn is a function to be executed when configured against a log level.
type EventFn func(ctx context.Context, r Record)

// Events contains an assignment of an event function to a log level.
type Events struct {
	Debug EventFn
	Info  EventFn
	Warn  EventFn
	Error EventFn
}
