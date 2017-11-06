package main

type Sampler interface {
	Sample() float32
}

type Ticker interface {
	Tick()
}

type Mixer struct {
	scheduler Ticker
	sources   []Sampler
}

func (m *Mixer) StartSoundLoop(output chan<- []float32, batchSize uint) {
	go func(m *Mixer) {
		for {
			buffer := make([]float32, batchSize)
			m.getSamples(&buffer)
			output <- buffer
		}
	}(m)
}

func (m *Mixer) getSamples(buffer *[]float32) {
	for i := range *buffer {
		m.scheduler.Tick()
		(*buffer)[i] = 0.0
		for u := range m.sources {
			(*buffer)[i] += m.sources[u].Sample()
		}
	}
}
