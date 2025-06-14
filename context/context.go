package context

import (
	"errors"
	"fmt"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key any) any
}

type emptyCtx struct{}

func (e emptyCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (e emptyCtx) Done() <-chan struct{} {
	return nil
}

func (e emptyCtx) Err() error {
	return nil
}

func (e emptyCtx) Value(key any) any {
	return nil
}

type backgroundCtx struct {
	emptyCtx
}

type todoCtx struct {
	emptyCtx
}

func (backgroundCtx) String() string { return "context.Background" }
func (todoCtx) String() string       { return "context.TODO" }

func Background() Context {
	return backgroundCtx{}
}

func TODO() Context {
	return todoCtx{}
}

type CancelFunc func()
type cancelCtx struct {
	parent Context
	done   chan struct{}
	err    error
}

func (c *cancelCtx) Deadline() (time.Time, bool) {
	return c.parent.Deadline()
}

func (c *cancelCtx) Done() <-chan struct{} {
	return c.done
}

func (c *cancelCtx) Err() error {
	return c.err
}

func (c *cancelCtx) String() string {
	if s, ok := c.parent.(fmt.Stringer); ok {
		return s.String() + ".WithCancel"
	}
	return "<unknown>.WithCancel"
}

func (c *cancelCtx) Value(key any) any {
	return c.parent.Value(key)
}

var errCanceled = errors.New("context canceled")

func (c *cancelCtx) cancel() {
	if c.err != nil {
		return
	}
	c.err = errCanceled
	close(c.done)
}

func WithCancel(parent Context) (Context, CancelFunc) {
	if parent == nil {
		panic("no parent set")
	}
	child := &cancelCtx{
		parent: parent,
		done:   make(chan struct{}),
	}

	go func() {
		select {
		case <-parent.Done():
			child.cancel()
		case <-child.done:
			return
		}
	}()

	return child, child.cancel
}

type timerCtx struct {
	*cancelCtx
	deadline time.Time
	timer    *time.Timer
}

func (t *timerCtx) Deadline() (deadline time.Time, ok bool) {
	return t.deadline, true
}

func (t *timerCtx) cancel() {
	t.cancelCtx.cancel()

	if t.timer != nil {
		t.timer.Stop()
		t.timer = nil
	}
}

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	if parent == nil {
		panic("nil parent")
	}

	if pd, ok := parent.Deadline(); ok && pd.Before(d) {
		return WithCancel(parent)
	}

	cctx := &cancelCtx{
		parent: parent,
		done:   make(chan struct{}),
	}
	tc := &timerCtx{cancelCtx: cctx, deadline: d}
	go func() {
		select {
		case <-parent.Done():
			tc.cancel()
		case <-tc.Done():
			return
		}

	}()

	dur := time.Until(d)
	if dur <= 0 {
		tc.cancel()
		return tc, tc.cancel
	}

	tc.timer = time.AfterFunc(dur, tc.cancel)

	return tc, tc.cancel
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}
