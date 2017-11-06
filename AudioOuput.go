package main

import (
	AD "github.com/gordonklaus/portaudio"
)

type AudioContext struct {
	stream        *AD.Stream
	sampleRate    uint
	bufferChannel chan []float32
}

func newAudioContext(sampleRate uint, bufferSize uint) (context *AudioContext, err error) {
	context = &AudioContext{nil, sampleRate, make(chan []float32, 1)}

	callback := func(out []float32) {
		buffer := <-context.bufferChannel
		for i := range out {
			out[i] = buffer[i] * 0.1
		}
	}

	context.stream, err = AD.OpenDefaultStream(0, 1, float64(sampleRate), int(bufferSize), callback)

	return
}
