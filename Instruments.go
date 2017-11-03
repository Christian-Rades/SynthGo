package main

import (
	"math"
)

type Note struct {
	on   bool
	adsr *ADSRenvelope
	oszi Oscillator
	LFO  Oscillator
}

func (n *Note) getSample() (out float32) {
	out = n.oszi.getNextValue()
	out *= n.LFO.getNextValue()
	out *= float32(n.adsr.nextMultiplier())
	if (n.adsr.state == R) && (math.Abs(float64(n.adsr.multiplier)) < 0.0001) {
		n.adsr.state = A
		n.on = false
		n.oszi.reset()
		n.LFO.reset()
	}
	return
}

type SineWave struct {
	sampleRate uint
	notes      [128]*Note
}

func newSineWave(sampleRate uint) *SineWave {
	notes := new([128]*Note)

	freq := 8.175799
	ratio := 1.0594631

	for i := range notes {
		notes[i] = &Note{false, newADSRenvelope(1, 1, 500, sampleRate), newSineOszillator(freq, float64(sampleRate)), newSineOszillator(freq/2, float64(sampleRate))}
		freq = freq * ratio
	}

	return &SineWave{sampleRate, *notes}
}

func (sineWave *SineWave) getSample() (out float32) {
	for _, note := range sineWave.notes {
		if note.on {
			out += note.getSample()
		} else {
			out += 0.0
		}
	}
	return
}

func (self *SineWave) NoteOn(parameter1 uint8, parameter2 uint8) {
	self.notes[parameter1].on = true
	self.notes[parameter1].adsr.reset()
}
func (self *SineWave) NoteOff(parameter1 uint8, parameter2 uint8) {
	self.notes[parameter1].adsr.release()
}
func (self *SineWave) NoteAftertouch(parameter1 uint8, parameter2 uint8) {
}
func (self *SineWave) ChannelAftertouch(parameter1 uint8, parameter2 uint8) {
}
func (self *SineWave) generateSound(output *[]float32) {
	for i := range *output {
		(*output)[i] += self.getSample()
	}
}
