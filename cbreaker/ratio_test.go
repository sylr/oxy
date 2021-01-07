package cbreaker

import (
	"math"
	"testing"
	"time"

	"abstraction.fr/oxy/v2/utils"
	"github.com/mailgun/holster/v3/clock"
	"github.com/stretchr/testify/assert"
)

func TestRampUp(t *testing.T) {
	defer clock.Freeze(time.Now()).Unfreeze()

	duration := 10 * time.Second
	rc := newRatioController(duration, &utils.DefaultLogger{})

	allowed, denied := 0, 0
	for i := 0; i < int(duration/time.Millisecond); i++ {
		ratio := sendRequest(&allowed, &denied, rc)
		expected := rc.targetRatio()
		diff := math.Abs(expected - ratio)
		assert.EqualValues(t, 0, round(diff, 0.5, 1))
		clock.Advance(time.Millisecond)
	}
}

func sendRequest(allowed, denied *int, rc *ratioController) float64 {
	if rc.allowRequest() {
		*allowed++
	} else {
		*denied++
	}
	if *allowed+*denied == 0 {
		return 0
	}
	return float64(*allowed) / float64(*allowed+*denied)
}

func round(val float64, roundOn float64, places int) float64 {
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	var round float64
	if div >= roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	return round / pow
}
