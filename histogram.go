package circonusgometrics

import (
	"sync"

	"github.com/circonus-labs/circonusllhist"
)

// A Histogram measures the distribution of a stream of values.
type Histogram struct {
	name string
	hist *circonusllhist.Histogram
	rw   sync.RWMutex
}

func (m *CirconusMetrics) Timing(metric string, val float64) {
	m.SetHistogramValue(metric, val)
}

func (m *CirconusMetrics) RecordValue(metric string, val float64) {
	m.SetHistogramValue(metric, val)
}

func (m *CirconusMetrics) SetHistogramValue(metric string, val float64) {
	m.hm.Lock()
	if _, ok := m.histograms[metric]; !ok {
		m.histograms[metric] = &Histogram{
			name: metric,
			hist: circonusllhist.New(),
		}
	}
	m.hm.Unlock()

	m.histograms[metric].rw.Lock()
	defer m.histograms[metric].rw.Unlock()

	m.histograms[metric].hist.RecordValue(val)
}

func (m *CirconusMetrics) RemoveHistogram(metric string) {
	m.hm.Lock()
	defer m.hm.Unlock()
	delete(m.histograms, metric)
}

// // NewHistogram returns a new Circonus histogram that accumulates until reported on.
// func (m *CirconusMetrics) NewHistogram(name string) *Histogram {
// 	hm.Lock()
// 	defer hm.Unlock()
//
// 	if hist, ok := histograms[name]; ok {
// 		return hist
// 	}
//
// 	hist := &Histogram{
// 		name: name,
// 		hist: circonusllhist.New(),
// 	}
// 	histograms[name] = hist
// 	return hist
// }
//
// // Remove removes the given histogram.
// func (h *Histogram) Remove() {
// 	hm.Lock()
// 	defer hm.Unlock()
// 	delete(histograms, h.name)
// }
//
// type hname string // unexported to prevent collisions
//
// // A Histogram measures the distribution of a stream of values.
// type Histogram struct {
// 	name string
// 	hist *circonusllhist.Histogram
// 	rw   sync.RWMutex
// }
//
// // Name returns the name of the histogram
// func (h *Histogram) Name() string {
// 	return h.name
// }
//
// // RecordValue records the given value
// func (h *Histogram) RecordValue(v float64) {
// 	h.rw.Lock()
// 	defer h.rw.Unlock()
//
// 	h.hist.RecordValue(v)
// }
