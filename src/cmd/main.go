package main

const (
	DefaultMaxJitter = 2000
	DefaultMinJitter = 100
)

func sleepWithJitter(min int, maax int) {
	if min < 1 {
		min = DefaultMinJitter
	}

	if max < 1 || max < min {
		max = DefaultMaxJitter
	}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd := rand.Intn(max-min) + min
	time.Sleep(time.Duration(rnd) * time.Millisecond)
}
