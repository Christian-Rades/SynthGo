package main

type State int

const (
	A State = iota
	D
	S
	R
)

type ADSRenvelope struct {
	attackTime  uint
	attackRate  float64
	decayTime   uint
	decayRate   float64
	releaseTime uint
	releaseRate float64
	sampleRate  uint
	state       State
	multiplier  float64
}

func newADSRenvelope(attack uint, decay uint, release uint, sampleRate uint) *ADSRenvelope {
	if sampleRate == 0 {
		panic("samplerate cannot be 0")
	}
	samplesPerMillisecond := sampleRate / 1000

	attackTime := attack * samplesPerMillisecond
	decayTime := decay * samplesPerMillisecond
	releaseTime := release * samplesPerMillisecond

	attackRate := 1.0
	decayRate := -0.3
	releaseRate := -0.7
	if attackTime != 0 {
		attackRate /= float64(attackTime)
	}
	if decayTime != 0 {
		decayRate /= float64(decayTime)
	}
	if releaseTime != 0 {
		releaseRate /= float64(releaseTime)
	}

	envelope := ADSRenvelope{attackTime, attackRate, decayTime,
		decayRate, releaseTime,
		releaseRate, sampleRate, A, 0.0}
	return &envelope
}

func (p *ADSRenvelope) nextMultiplier() (m float64) {
	m = p.multiplier

	switch p.state {
	case A:
		p.multiplier += p.attackRate
		if p.multiplier >= 1.0 {
			p.state = D
			return
		}
	case D:
		p.multiplier += p.decayRate
		if p.multiplier <= 0.7 {
			p.state = S
			return
		}
	case R:
		if p.multiplier <= 0.0001 {
			return
		}
		p.multiplier += p.releaseRate
	}

	return
}

func (p *ADSRenvelope) release() {
	p.state = R
}

func (p *ADSRenvelope) reset() {
	p.state = A
	p.multiplier = 0.0
}
