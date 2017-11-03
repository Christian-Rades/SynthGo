package main

import (
	//"github.com/gordonklaus/portaudio"
	//"math"
	//"time"
	"fmt"
	//"github.com/gordonklaus/portaudio"
	"os"
	//"time"
	"github.com/gordonklaus/portaudio"
	"time"
)

const SampleRate = uint(44100)

//func main() {
//
//	portaudio.Initialize()
//	defer portaudio.Terminate()
//	adsr := newADSRenvelope(1, 200, 1, 44100)
//	step := float64(440 / 44100.0)
//	phase := float64(0.0)
//	count := uint(0)
//	stream, err := portaudio.OpenDefaultStream(0, 1, 44100.0, 0, func(out []float32) {
//		for i := range out {
//			out[i] = float32(math.Sin(phase*math.Pi*2)*adsr.nextMultiplier()) / 5.0
//			_, phase = math.Modf(phase + step)
//			count++
//			if count > 3000 {
//				adsr.release()
//			}
//		}
//	})
//	chk(err)
//	defer stream.Close()
//	chk(stream.Start())
//	time.Sleep(3 * time.Second)
//	chk(stream.Stop())
//}
func main() {
	f, err := os.Open("MidiGen/chords.mid")
	chk(err)
	defer f.Close()

	mf := MidiFile{f}

	data, err := mf.parseFile()
	chk(err)

	fmt.Printf("\n")
	fmt.Printf("%+v\n", data.MidiHeader)
	fmt.Printf("%+v\n", data.Tracks[0].MidiTrackHeader)

	fmt.Printf("MidiChanEvents: %d \n", len(data.Tracks[0].ChannelEvents))

	scheduler := newScheduler(uint32(data.getSamplesPerTick(SampleRate)))

	for _, e := range data.Tracks[0].ChannelEvents {
		fmt.Printf("%+v\n", e)
		scheduler.addEvent(e)
	}

	//output := make([]float32, 44100)
	fmt.Printf("samplespertick: %d", data.getSamplesPerTick(SampleRate))

	//oszi := newSineOszillator(880.0, float64(SampleRate))

	//for i := range output {
	//	output[i] = oszi.getNextValue()
	//}

	portaudio.Initialize()
	defer portaudio.Terminate()
	context, contextErr := newAudioContext(SampleRate, data.getSamplesPerTick(SampleRate))
	stream := context.stream

	chk(contextErr)
	go mainloop(scheduler, uint32(data.getSamplesPerTick(SampleRate)), context.bufferChannel)

	defer stream.Close()
	chk(stream.Start())
	time.Sleep(10 * time.Second)
	chk(stream.Stop())

}

func chk(err error) {
	if err != nil {
		panic(err)
	}
}

func mainloop(scheduler *Scheduler,numOfSamples uint32, bufferChan chan *[]float32) {
	for {
		buffer := <- bufferChan
		scheduler.getNextSamples(numOfSamples, buffer)
	}
}