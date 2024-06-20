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

func New(p *print.Printer) *Spinner {
	return &Spinner{
		printer:   p,
		states:    []string{"|", "/", "-", "\\"},
		startTime: time.Now(),
		delay:     100 * time.Millisecond,
		done:      make(chan bool),
	}
}

func (s *Spinner) Start(message string) {
	s.message = message
	go s.animate()
}

func (s *Spinner) Stop() {
	s.done <- true
	close(s.done)
	s.printer.Info("\r%s ✓ \n", s.message)
}

func (s *Spinner) StopWithError() {
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
