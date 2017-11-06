package main

type MidiEventType uint8
type MidiChannelId uint8

const (
	NoteOff           MidiEventType = 0x8
	NoteOn            MidiEventType = 0x9
	NoteAftertouch    MidiEventType = 0xA
	Controller        MidiEventType = 0xB
	ProgramChange     MidiEventType = 0xC
	ChannelAftertouch MidiEventType = 0xD
	Pitchbend         MidiEventType = 0xE
	SysExStart        MidiEventType = 0xF0
	SysExEnd          MidiEventType = 0xF7
	MetaEvent         MidiEventType = 0xFF
)

type MidiControllable interface {
	NoteOn(parameter1 uint8, parameter2 uint8)
	NoteOff(parameter1 uint8, parameter2 uint8)
	NoteAftertouch(parameter1 uint8, parameter2 uint8)
	ChannelAftertouch(parameter1 uint8, parameter2 uint8)
}

type MidiEvent struct {
	deltaTime uint32
	EventType MidiEventType
}

type MidiChannelEvent struct {
	MidiEvent
	Channel    MidiChannelId
	Parameter1 uint8
	Parameter2 uint8
}

type MidiMetaEvent struct {
	MidiEvent
	MetaType byte
	Length   uint32
	Data     []byte
}

type MidiSysExEvent struct {
	MidiEvent
	Data []byte
}

func (event *MidiChannelEvent) isOk() bool {
	deltaTimeOk := event.deltaTime <= 0x0FFFFFFF
	channelOk := event.Channel <= 15
	parametersOK := event.Parameter1 <= 127 && event.Parameter2 <= 127
	return deltaTimeOk && channelOk && parametersOK
}
