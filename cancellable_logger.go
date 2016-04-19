package logging

import (
	"sync"
)

// CancellableLogger is a logger that can be cancelled.  Once cancelled, all
// subsequent log calls are dropped.  The cancellation is thread-safe.
type CancellableLogger struct {
	Logger
	m sync.RWMutex
}

func ul(m *sync.RWMutex) func() { m.RLock(); return m.RUnlock }

var _ Logger = &CancellableLogger{}

func (c *CancellableLogger) Trace(v ...interface{}) { defer ul(&c.m)(); c.Logger.Trace(v...) }
func (c *CancellableLogger) Tracef(f string, a ...interface{}) {
	defer ul(&c.m)()
	c.Logger.Tracef(f, a...)
}
func (c *CancellableLogger) Debug(v ...interface{}) { defer ul(&c.m)(); c.Logger.Debug(v...) }
func (c *CancellableLogger) Debugf(f string, a ...interface{}) {
	defer ul(&c.m)()
	c.Logger.Debugf(f, a...)
}
func (c *CancellableLogger) Info(v ...interface{}) { defer ul(&c.m)(); c.Logger.Info(v...) }
func (c *CancellableLogger) Infof(f string, a ...interface{}) {
	defer ul(&c.m)()
	c.Logger.Infof(f, a...)
}
func (c *CancellableLogger) Error(v ...interface{}) { defer ul(&c.m)(); c.Logger.Error(v...) }
func (c *CancellableLogger) Errorf(f string, a ...interface{}) {
	defer ul(&c.m)()
	c.Logger.Errorf(f, a...)
}
func (c *CancellableLogger) LogLevel() Level { defer ul(&c.m)(); return c.Logger.LogLevel() }
func (c *CancellableLogger) SetLogLevel(newLev Level) {
	defer ul(&c.m)()
	c.Logger.SetLogLevel(newLev)
}

func (c *CancellableLogger) Cancel() { c.m.Lock(); c.Logger = DiscardLogger{}; c.m.Unlock() }

// DiscardLogger is a Logger that drops all logging calls.
type DiscardLogger struct{}

func (l DiscardLogger) Trace(vals ...interface{})              {}
func (l DiscardLogger) Tracef(fmt string, args ...interface{}) {}
func (l DiscardLogger) Debug(vals ...interface{})              {}
func (l DiscardLogger) Debugf(fmt string, args ...interface{}) {}
func (l DiscardLogger) Info(vals ...interface{})               {}
func (l DiscardLogger) Infof(fmt string, args ...interface{})  {}
func (l DiscardLogger) Error(vals ...interface{})              {}
func (l DiscardLogger) Errorf(fmt string, args ...interface{}) {}
func (l DiscardLogger) LogLevel() Level                        { return ErrorLevel }
func (l DiscardLogger) SetLogLevel(newLevel Level)             {}
