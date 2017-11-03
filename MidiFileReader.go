package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
)

type MidiHeader struct {
	FileId       [4]byte
	Size         uint32
	FormatType   uint16
	TrackCount   uint16
	TimeDivision uint16
}

type MidiTrackHeader struct {
	TrackId [4]byte
	Size    uint32
}

type MidiTrack struct {
	MidiTrackHeader
	ChannelEvents []MidiChannelEvent
	MetaEvents    []MidiMetaEvent
	SysExEvents   []MidiSysExEvent
}

type MidiFile struct {
	*os.File
}

type MidiData struct {
	MidiHeader
	Tracks []MidiTrack
}

func concatVariableValue(n uint32, b byte) uint32 {
	b = b & 0x7F // filter out the MSb
	n = n << 7
	return n | uint32(b)
}

func (header *MidiHeader) isOk() error {
	if header.FileId != [4]byte{'M', 'T', 'h', 'd'} {
		return errors.New("wrong file format id")
	}
	if header.Size != 6 {
		return errors.New("size of header must be 6")
	}
	if (header.FormatType != 0) || (header.TrackCount > 1) {
		return errors.New("cant handle files with more than 1 track")
	}
	if header.TrackCount == 0 {
		return errors.New("file does not contain any track")
	}
	if header.TimeDivision == 0 {
		return errors.New("could not determine time division")
	}
	return nil
}

func (header *MidiTrackHeader) isOk() error {
	if header.TrackId != [4]byte{'M', 'T', 'r', 'k'} {
		return errors.New("malformed track")
	}
	return nil
}

func (f *MidiFile) readSingleByte() (output byte, err error) {
	currentByte := make([]byte, 1)
	_, err = f.Read(currentByte)
	output = currentByte[0]
	return
}

func (f MidiFile) readVariableValue() (value uint32, err error) {
	b := byte(0)
	for i := uint(0); i < 4; i++ {
		b, err = f.readSingleByte()
		value = concatVariableValue(value, b)
		if int8(b) >= 0 {
			fmt.Printf("%d", i)
			break
		}
	}
	return
}

func (f *MidiFile) readMidiChannelEvent(parentEvent MidiEvent) (event MidiChannelEvent, err error) {
	event.MidiEvent = parentEvent
	event.EventType = MidiEventType((parentEvent.EventType & 0xF0) >> 4)
	event.Channel = MidiChannelId(parentEvent.EventType & 0x0F)

	event.Parameter1, err = f.readSingleByte()
	event.Parameter2, err = f.readSingleByte()

	return
}

func (f *MidiFile) readMidiMetaEvent(parentEvent MidiEvent) (event MidiMetaEvent, err error) {
	event.MidiEvent = parentEvent

	event.MetaType, err = f.readSingleByte()

	event.Length, err = f.readVariableValue()

	event.Data = make([]byte, event.Length)
	n := int(0)
	n, err = f.Read(event.Data)

	if uint32(n) != event.Length {
		err = errors.New("bad meta event length")
	}

	return
}

func (f *MidiFile) readMidiSysExEvent(parentEvent MidiEvent) (event MidiSysExEvent, err error) {
	event.MidiEvent = parentEvent
	for err == nil {
		b := byte(0)
		b, err = f.readSingleByte()
		if MidiEventType(b) == SysExEnd {
			return
		}
		event.Data = append(event.Data, b)
	}
	return
}

func (f *MidiFile) readMidiEvents() (channelEvents []MidiChannelEvent, metaEvents []MidiMetaEvent, sysExEvents []MidiSysExEvent, err error) {
	deltaTime := uint32(0)
	typeByte := byte(0)
	channelEvent := MidiChannelEvent{MidiEvent{0, 0}, 0, 0, 0}
	metaEvent := MidiMetaEvent{MidiEvent{0, 0}, 0, 0, nil}
	sysExEvent := MidiSysExEvent{MidiEvent{0, 0}, nil}
	for err == nil {
		deltaTime, err = f.readVariableValue()
		typeByte, err = f.readSingleByte()
		event := MidiEvent{deltaTime, MidiEventType(typeByte)}
		switch MidiEventType(typeByte) {
		case MetaEvent:
			metaEvent, err = f.readMidiMetaEvent(event)
			metaEvents = append(metaEvents, metaEvent)
			if metaEvent.MetaType == 0x2F {
				return
			}
		case SysExStart, SysExEnd:
			sysExEvent, err = f.readMidiSysExEvent(event)
			sysExEvents = append(sysExEvents, sysExEvent)
		default:
			channelEvent, err = f.readMidiChannelEvent(event)
			channelEvents = append(channelEvents, channelEvent)
		}
	}
	return
}

func (f *MidiFile) parseFile() (data MidiData, err error) {
	header := MidiHeader{}
	binary.Read(f, binary.BigEndian, &header)
	err = header.isOk()
	if err != nil {
		return
	} else {
		data.MidiHeader = header
	}

	trackHeader := MidiTrackHeader{}
	binary.Read(f, binary.BigEndian, &trackHeader)
	err = trackHeader.isOk()
	if err != nil {
		return
	}

	channelEvents, metaEvents, sysExEvents, eventErr := f.readMidiEvents()
	err = eventErr
	if err != nil {
		return
	}

	data.Tracks = make([]MidiTrack, 1)
	track := MidiTrack{trackHeader, channelEvents, metaEvents, sysExEvents}
	data.Tracks[0] = track
	return
}

func (data *MidiData) getSamplesPerTick(sampleRate uint) uint {
	if int16(data.TimeDivision) <= 0 {
		panic("!!!implement SMPTE frame rates! ! !")
	}

	bpm := uint16(120) * data.TimeDivision
	bps := uint(bpm) / 60
	return sampleRate / bps
}
