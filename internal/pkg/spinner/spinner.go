package spinner

import (
	"time"

	"github.com/stackitcloud/stackit-cli/internal/pkg/print"
)

type Spinner struct {
	printer   *print.Printer
	message   string
	states    []string
	startTime time.Time
	delay     time.Duration
	done      chan bool
}

// Run starts a spinner and stops it when f completes
func Run(p *print.Printer, message string, f func() error) error {
	_, err := Run2(p, message, func() (struct{}, error) {
		err := f()
		return struct{}{}, err
	})
	return err
}

// Run2 starts a spinner and stops it when f (result arity 2) completes.
func Run2[T any](p *print.Printer, message string, f func() (T, error)) (T, error) {
	var r T
	spinner := newSpinner(p)
	spinner.start(message)
	r, err := f()
	if err != nil {
		spinner.stopWithError()
		return r, err
	}
	spinner.stop()
	return r, nil
}

func newSpinner(p *print.Printer) *Spinner {
	return &Spinner{
		printer:   p,
		states:    []string{"|", "/", "-", "\\"},
		startTime: time.Now(),
		delay:     100 * time.Millisecond,
		done:      make(chan bool),
	}
}

func (s *Spinner) start(message string) {
	s.message = message
	go s.animate()
}

func (s *Spinner) stop() {
	s.done <- true
	close(s.done)
	s.printer.Info("\r%s ✓ \n", s.message)
}

func (s *Spinner) stopWithError() {
	s.done <- true
	close(s.done)
	s.printer.Info("\r%s ✗ \n", s.message)
}

func (s *Spinner) animate() {
	i := 0
	for {
		select {
		case <-s.done:
			return
		default:
			s.printer.Info("\r%s %s ", s.message, s.states[i%len(s.states)])
			i++
			time.Sleep(s.delay)
		}
	}
}
