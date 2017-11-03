package main

import (
	"math"
)

type Oscillator interface {
	getNextValue() float32
	reset()
}

type SineOscillator struct {
	counter     uint
	tableLength uint
	waveTable   []float32
}

type SawToothOscillator struct {
	counter     uint
	tableLength uint
	waveTable   []float32
}

func newSineOszillator(frequency float64, sampleRate float64) (oszi *SineOscillator) {
	if frequency <= 0.0 || sampleRate <= 0.0 {
		panic("frequency and samplerate must be greater than 0.0")
	}

	numOfSamples := uint(Round(sampleRate / frequency))
	step := (frequency / sampleRate)
	phase := 0.0
	oszi = &SineOscillator{0, numOfSamples, make([]float32, numOfSamples)}
	for i := range oszi.waveTable {
		oszi.waveTable[i] = float32(math.Sin(phase * math.Pi * 2))
		phase += step
	}
	return
}

func Round(f float64) float64 {
	return math.Floor(f + .5)
}

func (sineOszi *SineOscillator) getNextValue() (output float32) {
	output = sineOszi.waveTable[sineOszi.counter]
	sineOszi.counter++
	if sineOszi.counter >= sineOszi.tableLength {
		sineOszi.counter = 0
	}
	return
}

func (sineOszi *SineOscillator) reset() {
	sineOszi.counter = 0
}

func newSawToothOscillator(frequency float64, sampleRate float64) (oszi *SawToothOscillator) {
	if frequency <= 0.0 || sampleRate <= 0.0 {
		panic("frequency and samplerate must be greater than 0.0")
	}

	numOfSamples := uint(Round(sampleRate / frequency))
	step := 1.0 / float64(numOfSamples)
	phase := 0.0
	oszi = &SawToothOscillator{0, numOfSamples, make([]float32, numOfSamples)}
	for i := range oszi.waveTable {
		oszi.waveTable[i] = float32(phase)
		phase += step
	}
	return
}

func (sawOsci *SawToothOscillator) getNextValue() (output float32) {
	output = sawOsci.waveTable[sawOsci.counter]
	sawOsci.counter++
	if sawOsci.counter >= sawOsci.tableLength {
		sawOsci.counter = 0
	}
	return
}

func (sawOsci *SawToothOscillator) reset() {
	sawOsci.counter = 0
}
