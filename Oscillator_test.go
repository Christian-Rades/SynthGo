package main

import (
	. "github.com/onsi/gomega"
	"testing"
)

func TestSineWaveOscillator (t *testing.T) {
	RegisterTestingT(t)

	freq := 440.0
	samplerate := float64(SampleRate)
	numOfSamples := uint(100)

	oscillator := newSineOszillator(freq, samplerate)

	Expect(oscillator.tableLength).Should(BeNumerically("==",numOfSamples))
	Expect(oscillator.waveTable[0]).Should(BeNumerically("~",oscillator.waveTable[numOfSamples-1], 0.1)) //checks for popping on loop

	Expect(func() {newSineOszillator(1,0)}).Should(Panic())
	Expect(func() {newSineOszillator(0,1)}).Should(Panic())
	Expect(func() {newSineOszillator(0,0)}).Should(Panic())
}

func TestSawToothOscillator (t *testing.T) {
	RegisterTestingT(t)

	freq := 440.0
	samplerate := float64(SampleRate)
	numOfSamples := uint(100)

	oscillator := newSawToothOscillator(freq, samplerate)

	Expect(oscillator.tableLength).Should(BeNumerically("==",numOfSamples))

	Expect(oscillator.waveTable[0]).Should(BeNumerically("==",0.0))
	Expect(oscillator.waveTable[numOfSamples-1]).Should(BeNumerically("~",1.0, 0.1))

	Expect(func() {newSineOszillator(1,0)}).Should(Panic())
	Expect(func() {newSineOszillator(0,1)}).Should(Panic())
	Expect(func() {newSineOszillator(0,0)}).Should(Panic())
}