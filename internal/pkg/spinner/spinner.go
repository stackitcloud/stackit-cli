package spinner

import (
	"time"

	"github.com/spf13/cobra"
)

type Spinner struct {
	cmd       *cobra.Command
	message   string
	states    []string
	startTime time.Time
	delay     time.Duration
	done      chan bool
}

func New(cmd *cobra.Command) *Spinner {
	return &Spinner{
		cmd:       cmd,
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
	s.cmd.Printf("\r%s ✓ \n", s.message)
}

func (s *Spinner) animate() {
	i := 0
	for {
		select {
		case <-s.done:
			return
		default:
			s.cmd.Printf("\r%s %s ", s.message, s.states[i%len(s.states)])
			i++
			time.Sleep(s.delay)
		}
	}
}
