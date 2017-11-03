package main

type Scheduler struct {
	timeCounter    uint32
	sampleCounter  uint32
	samplesPerTick uint32
	channels       []MidiChannel
	quque          []MidiChannelEvent
}

func (scheduler *Scheduler) addEvent(event MidiChannelEvent) {
	if !event.isOk() {
		return
	}
	scheduler.quque = append(scheduler.quque, event)
}

func (scheduler *Scheduler) readSoundBite(wavetable *[]float32) {
	for i := range *wavetable {
		(*wavetable)[i] = 0.0
	}
	for _, channel := range scheduler.channels {
		channel.getSound(wavetable)
	}
	return
}

func (scheduler *Scheduler) getNextSamples(sampleCount uint32, output *[]float32) {
	if len(scheduler.quque) == 0 {
		for i := range *output {
			(*output)[i] = 0.0
		}
		return
	}

	scheduler.tick()

	scheduler.sampleCounter += sampleCount
	if scheduler.sampleCounter >= scheduler.samplesPerTick {
		scheduler.timeCounter++
		scheduler.sampleCounter = sampleCount % scheduler.samplesPerTick
	}
	scheduler.readSoundBite(output)
	return
}

func (scheduler *Scheduler) tick() {
	for len(scheduler.quque) > 0 {
		currentEvent := scheduler.quque[0]
		if currentEvent.deltaTime > scheduler.timeCounter {
			break
		}
		scheduler.quque = scheduler.quque[1:]
		scheduler.channels[currentEvent.Channel].handleEvent(currentEvent)
		scheduler.timeCounter = 0
	}
}

func (scheduler *Scheduler) nextTick() (out []float32) {
	if len(scheduler.quque) == 0 {
		panic("queue empty")
	}
	for len(scheduler.quque) > 0 {
		currentEvent := scheduler.quque[0]
		if currentEvent.deltaTime > scheduler.timeCounter {
			break
		}
		scheduler.quque = scheduler.quque[1:]
		scheduler.channels[currentEvent.Channel].handleEvent(currentEvent)
		scheduler.timeCounter = 0
	}
	scheduler.timeCounter++
	out = make([]float32, scheduler.samplesPerTick)
	scheduler.readSoundBite(&out)
	return
}

func newScheduler(samplesPerTick uint32) (scheduler *Scheduler) {
	channels := make([]MidiChannel, 1)
	for i := 0; i < 1; i++ {
		channels[i] = MidiChannel{MidiChannelId(i), newSineWave(44100)}
	}

	scheduler = &Scheduler{0, 0, samplesPerTick, channels, make([]MidiChannelEvent, 0)}
	return
}
