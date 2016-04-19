package logging

import (
	"os"
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/fluxio/sync_testing"
)

//  A logger safe for concurrent usage.
type atomicLogger int64

func (l *atomicLogger) Trace(vals ...interface{})              { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Tracef(fmt string, args ...interface{}) { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Debug(vals ...interface{})              { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Debugf(fmt string, args ...interface{}) { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Info(vals ...interface{})               { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Infof(fmt string, args ...interface{})  { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Error(vals ...interface{})              { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) Errorf(fmt string, args ...interface{}) { atomic.AddInt64((*int64)(l), 1) }
func (l *atomicLogger) LogLevel() Level                        { return ErrorLevel }
func (l *atomicLogger) SetLogLevel(lev Level)                  {}

func TestCancellableLogger(t *testing.T) {
	runtime.GOMAXPROCS(10)

	var cl CancellableLogger
	var logger atomicLogger
	cl.Logger = &logger

	sync_testing.MaximizeContention(50,
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Info("x") },
		func() { cl.Cancel() },
	)

	// For some reason, running this test under the race detector works just
	// fine (even if in the regression scenario of a race actually existing in
	// the code) unless the test explicitly fails.  Weird.  So, if you set the
	// "FAIL_FOR_RACE" env var, this test will artificially fail which seems to
	// trigger the race detector's output. Example command line:
	//   FAIL_FOR_RACE=1 go test -race -run CancellableLogger genie/flow
	// Try running that without the write-lock in CancellableLogger.Cancel().
	if os.Getenv("FAIL_FOR_RACE") != "" {
		t.Error("Artificially failing test to poke race detector.")
	}
}
