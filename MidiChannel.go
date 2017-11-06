package main

type MidiChannel struct {
	id         MidiChannelId
	instrument MidiControllable
}

func (channel *MidiChannel) handleEvent(event MidiChannelEvent) {
	if channel.id != event.Channel {
		panic("Event has wrong ID") // maybe switch to fail silently?
	}
	param1, param2 := event.Parameter1, event.Parameter2
	switch event.EventType {
	case NoteOff:
		channel.instrument.NoteOff(param1, param2)
	case NoteOn:
		channel.instrument.NoteOn(param1, param2)
	case NoteAftertouch:
		channel.instrument.NoteAftertouch(param1, param2)
	case ChannelAftertouch:
		channel.instrument.ChannelAftertouch(param1, param2)
	}
}