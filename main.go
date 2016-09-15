package chronon

import (
	"fmt"
	"sync"
	"time"
)

type Config struct {
	Size int
	Safe bool // Optionally lock if concurrent access is needed.
}

type timeSlice []time.Duration

type Chronon struct {
	sync.Mutex
	Safe          bool
	Times         timeSlice
	TimesPosition int
	Count         int
}

type Metrics struct {
	Time struct {
		Total    time.Duration
		Avg      time.Duration
		Median   time.Duration
		p95      time.Duration
		Long10p  time.Duration
		Short10p time.Duration
		Max      time.Duration
		Min      time.Duration
	}
	Samples int
	Count   int
}

func New(c *Config) *Chronon {
	return &Chronon{
		Times: make([]time.Duration, c.Size),
		Safe:  c.Safe,
	}
}

// AddTime adds a time.Duration to the Chronon.Times
// slice, then increments the position.
func (m *Chronon) AddTime(t time.Duration) {
	if m.Safe {
		m.Lock()
		defer m.Unlock()
	}

	// If we're at the end, rollover and
	// start overwriting.
	if m.TimesPosition == len(m.Times) {
		m.TimesPosition = 0
	}

	m.Times[m.TimesPosition] = t
	m.TimesPosition++
}

// AddCount simply counts events.
func (m *Chronon) AddCount(i int) {
	if m.Safe {
		m.Lock()
		defer m.Unlock()
	}

	m.Count += i
}

// Dump prints out a generic output of
// all gathered metrics.
func (m *Chronon) Dump() {
	metrics := m.Calc()
	fmt.Printf("%d samples of %d events\n", metrics.Samples, metrics.Count)
	fmt.Printf("Total:\t\t%s\n", metrics.Time.Total)
	fmt.Printf("Avg.:\t\t%s\n", metrics.Time.Avg)
	fmt.Printf("95%%ile:\t\t%s\n", metrics.Time.p95)
	fmt.Printf("Longest 10%%:\t%s\n", metrics.Time.Long10p)
	fmt.Printf("Shortest 10%%:\t%s\n", metrics.Time.Short10p)
	fmt.Printf("Max:\t\t%s\n", metrics.Time.Max)
	fmt.Printf("Min:\t\t%s\n", metrics.Time.Min)
}
