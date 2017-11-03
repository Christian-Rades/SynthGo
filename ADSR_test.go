package main

import "testing"

func TestAdsrCreation(t *testing.T) {
	envelope := newADSRenvelope(100, 300, 50, 44100)
	if envelope.state != A {
		t.Error("Envelope was expected to be created in state Attack instead got", envelope.state)
	}
	if envelope.multiplier != 0.0 {
		t.Error("Envelope was expected to be created with Multiplier 0.0 instead got", envelope.multiplier)
	}
}

func TestAdsrErrorOnCreationWithNoSampleRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("NewADSR doesn't panic when Samplerate is 0 ")
		}
	}()
	newADSRenvelope(100, 300, 50, 0)
}

func TestAdsrWithZeroTimes(t *testing.T) {
	envelope := newADSRenvelope(0, 0, 0, 123)
	if envelope.attackRate != 1.0 {
		t.Error("ADSR doesn't handle 0ms attack time correctly")
	}
	if envelope.decayRate != -0.3 {
		t.Error("ADSR doesn't handle 0ms decay time correctly")
	}
	if envelope.releaseRate != -0.7 {
		t.Error("ADSR doesn't handle 0ms release time correctly")
	}
}

func TestAdsrEnvelope(t *testing.T) {
	envelope := newADSRenvelope(1, 2, 3, 1000) //1000hz sample rate so each tick is 1ms
	multipliers := make([]float64, 9)
	for i := range multipliers {
		multipliers[i] = envelope.nextMultiplier()
		if i == 3 {
			envelope.release()
		}
	}

	envelope.reset()

	if envelope.multiplier != 0.0 {
		t.Error("Envelope was expected to have Multiplier 0.0 after reset instead got", envelope.multiplier)
	}

	for i := range multipliers {
		switch i {
		case 0:
			if multipliers[i] != 0.0 {
				t.Error("value expected 0.0 got", multipliers[i])
			}
		case 1:
			if multipliers[i] != 1.0 {
				t.Error("value expected 1.0 got", multipliers[i])
			}
		case 2:
			if (multipliers[i] >= 1.0) && (multipliers[i] <= 0.7) {
				t.Error("value expected less than 1.0 more than 0.7 got", multipliers[i])
			}
		case 3:
			if multipliers[i] != 0.7 {
				t.Error("value expected 0.7 got", multipliers[i])
			}
		case 4:
			if multipliers[i] != 0.7 {
				t.Error("value expected 0.7 got", multipliers[i])
			}
		case 5:
			if (multipliers[i] >= 0.7) && (multipliers[i] <= 0.3) {
				t.Error("value expected less than 0.7 more than 0.3 got", multipliers[i])
			}
		case 6:
			if (multipliers[i] >= 0.3) && (multipliers[i] <= 0.1) {
				t.Error("value expected less than 0.3 more than 0.1 got", multipliers[i])
			}
		case 7:
			if multipliers[i] >= 0.0001 {
				t.Error("value expected below 0.0001 got", multipliers[i])
			}
		case 8:
			if (multipliers[i] >= 0.0001) && (multipliers[i] < 0.0) {
				t.Error("non negative value expected below 0.0001 got", multipliers[i])
			}
		}
	}
}
