package cli

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// A progress prints the progress lines of one command on standard error
// (docs/spec/README.md, Commands: Progress): one line when a step starts, and
// one line for each beat while the step runs. A package that runs the steps
// gets its step method as a function, so it imports nothing for this. It is
// for one goroutine; the beats print from their own goroutine, only while a
// step runs.
type progress struct {
	w       io.Writer
	command string        // for example gate
	every   time.Duration // the time between two beats
	ticker  func(time.Duration) (<-chan time.Time, func())

	stop, done chan struct{} // of the beats of the running step
	i, n       int           // the running step
	name       string
	mu         sync.Mutex // over the lines, and noted
	noted      bool       // a line of a wait came since the last beat
	afterBeat  func()     // nil; a test waits on it for each beat to end
}

// newProgress gives the progress of command, with a beat every ten seconds.
func newProgress(w io.Writer, command string) *progress {
	return &progress{w: w, command: command, every: 10 * time.Second, ticker: func(d time.Duration) (<-chan time.Time, func()) {
		t := time.NewTicker(d)
		return t.C, t.Stop
	}}
}

// step ends the beats of the step before, prints the line of step i of n,
// and starts the beats of this step.
func (p *progress) step(i, n int, name string) {
	p.end()
	p.i, p.n, p.name, p.noted = i, n, name, false
	fmt.Fprintf(p.w, "layup %s: [%d/%d] %s\n", p.command, i, n, name)
	p.beats()
}

// note prints a line of a wait of the running step (docs/spec/run.md, The
// command), and the next beat prints nothing, so a wait that prints a line
// every ten seconds shows one line every ten seconds; with no running step it
// prints the line alone.
func (p *progress) note(line string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.stop == nil {
		fmt.Fprintf(p.w, "layup %s: %s\n", p.command, line)
		return
	}
	p.noted = true
	fmt.Fprintf(p.w, "layup %s: [%d/%d] %s: %s\n", p.command, p.i, p.n, p.name, line)
}

// beats starts the beats of the running step; a beat counts from the start of
// the step.
func (p *progress) beats() {
	tick, stopTicker := p.ticker(p.every)
	p.stop, p.done = make(chan struct{}), make(chan struct{})
	go func(stop, done chan struct{}) {
		defer close(done)
		defer stopTicker()
		for beat := 1; ; beat++ {
			select {
			case <-stop:
				return
			case <-tick:
				p.mu.Lock()
				if !p.noted {
					fmt.Fprintf(p.w, "layup %s: [%d/%d] %s: %d s\n", p.command, p.i, p.n, p.name, int64(beat)*int64(p.every/time.Second))
				}
				p.noted = false
				p.mu.Unlock()
				if p.afterBeat != nil {
					p.afterBeat()
				}
			}
		}
	}(p.stop, p.done)
}

// end stops the beats of the running step, and returns when no beat can
// print any more.
func (p *progress) end() { p.stopBeats() }

func (p *progress) stopBeats() {
	if p.stop != nil {
		close(p.stop)
		<-p.done
		p.stop, p.done = nil, nil
	}
}
