package main

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
	"os"
	"time"
)

const SampleRate = uint(44100)

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

	sinewaves := make([]*SineWave, 1)
	sinewaves[0] = newSineWave(SampleRate)

	instruments := make([]MidiControllable, 1)
	instruments[0] = sinewaves[0]

	tickers := make([]Sampler, 1)
	tickers[0] = sinewaves[0]

	scheduler := newScheduler(instruments, data.getSamplesPerTick(SampleRate))
	mixer := Mixer{scheduler, tickers}

	go func() {
		for i := range data.Tracks[0].ChannelEvents {
			scheduler.eventQueue <- &data.Tracks[0].ChannelEvents[i]
			fmt.Printf("%d\n", i)
		}
	}()

	fmt.Printf("samplespertick: %d", data.getSamplesPerTick(SampleRate))

	portaudio.Initialize()
	defer portaudio.Terminate()
	context, contextErr := newAudioContext(SampleRate, SampleRate / 90)
	stream := context.stream

	chk(contextErr)

	go mixer.StartSoundLoop(context.bufferChannel, SampleRate / 90)

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
