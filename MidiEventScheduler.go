package main

type Scheduler struct {
	eventQueue     chan *MidiChannelEvent
	nextEvent      *MidiChannelEvent
	instruments    []MidiChannel
	sampleCounter  uint
	deltaCounter   uint32
	samplesPerTick uint
}

func (s *Scheduler) getNextEvent() (out *MidiChannelEvent, isOk bool) {
	if s.nextEvent != nil {
		out = s.nextEvent
		isOk = true
	} else {
		select {
		case s.nextEvent = <-s.eventQueue:
			out = s.nextEvent
		default:
			s.nextEvent = nil
			isOk = false
		}
	}
	return
}

func (s *Scheduler) midiTick() {
	ev,ok := s.getNextEvent()
	if !ok {
		return
	}
	for ev.deltaTime <= s.deltaCounter {
		s.instruments[ev.Channel].handleEvent(*ev)

		s.deltaCounter = 0
		s.nextEvent = nil
		ev, ok = s.getNextEvent()
		if !ok {
			return
		}
	}
}

func (s *Scheduler) Tick() {
	s.midiTick()
	s.sampleCounter++
	if s.sampleCounter >= s.samplesPerTick {
		s.deltaCounter++
		s.sampleCounter = 0
	}
}

func newScheduler(instruments []MidiControllable, samplesPerTick uint) *Scheduler {
	channels := make([]MidiChannel, len(instruments))
	for i := 0; i < 1; i++ {
		channels[i] = MidiChannel{MidiChannelId(i), instruments[i]}
	}
	return &Scheduler{make(chan *MidiChannelEvent, 100), nil, channels, 0, 0, samplesPerTick}
}
