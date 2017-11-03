package main

import(
	AD "github.com/gordonklaus/portaudio"
)

const(
	buf1 uint = iota
	buf2 uint = iota
)

const (
	AudioStart uint = iota
	AudioStop uint = iota
)

type AudioContext  struct{
	stream 	*AD.Stream
	sampleRate uint
	currentBuffer uint
	bufferChannel chan *[]float32
	buffer1 []float32
	buffer2 []float32
}


func newAudioContext(sampleRate uint, bufferSize uint) (context *AudioContext, err error) {
	context = &AudioContext{nil, sampleRate, buf1, make(chan *[]float32),make([]float32, bufferSize),make([]float32, bufferSize)}

	callback := func(out []float32) {
		switch context.currentBuffer {
		case buf1:
			context.bufferChannel <- &context.buffer2
			for i := range out{
				out[i] = context.buffer1[i] * 0.1
			}
			context.currentBuffer = buf2
			return
		case buf2:
			context.bufferChannel <- &context.buffer1
			for i := range out{
				out[i] = context.buffer2[i] * 0.1
			}
			context.currentBuffer = buf1
			return
		}
	}

	context.stream, err = AD.OpenDefaultStream(0, 1, float64(sampleRate), int(bufferSize), callback)

	return
}