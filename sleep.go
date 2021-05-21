package sleep

import (
	"sync"
	"time"
)

type Sleep struct {
	awakeAt time.Time
	asleep  bool
	refresh chan struct{}
	sync.RWMutex

	debug bool
}

// `Sleep` method on a single `Sleep` instance should not be called concurrently,
// otherwise only a random one will be awaken up by `Awake*` methods.
// `Sleep` method on a single `Sleep` instance can be called many times serially.
func (s *Sleep) Sleep(d time.Duration) {
	s.SetAwakeAt(time.Now().Add(d))

	s.Run()
}

// `Run` method on a single `Sleep` instance should not be called concurrently,
// otherwise only a random one will be awaken up by `Awake*` methods.
// `Run` method on a single `Sleep` instance can be called many times serially.
func (s *Sleep) Run() {
	d := time.Until(s.GetAwakeAt())
	if d <= 0 {
		return
	}

	s.Lock()
	if s.refresh == nil {
		s.refresh = make(chan struct{})
	}
	s.asleep = true
	s.Unlock()

	defer func() {
		s.Lock()
		s.asleep = false
		s.Unlock()
	}()

	for {
		var timer = time.NewTimer(d)
		select {
		case <-timer.C:
			timer.Stop()
			return
		case <-s.refresh:
			timer.Stop()
			if d = time.Until(s.GetAwakeAt()); d <= 0 {
				return
			}
		}
	}
}

func (s *Sleep) Asleep() bool {
	s.RLock()
	defer s.RUnlock()
	return s.asleep
}

// get awake at time.
func (s *Sleep) GetAwakeAt() time.Time {
	s.RLock()
	defer s.RUnlock()
	return s.awakeAt
}

// set awake at time.
func (s *Sleep) SetAwakeAt(at time.Time) {
	s.Lock()
	defer s.Unlock()
	s.awakeAt = at
}

// set awake at time to zero time.
func (s *Sleep) ClearAwakeAt() {
	s.Lock()
	defer s.Unlock()
	s.awakeAt = time.Time{}
}

// set awake at time to the specified time, and awake at the specified time if asleep.
func (s *Sleep) AwakeAt(at time.Time) {
	s.SetAwakeAt(at)
	select {
	case s.refresh <- struct{}{}:
	default:
	}
}

// the same as s.AwakeAt(time.Now())
func (s *Sleep) Awake() {
	s.AwakeAt(time.Now())
}

// if current awake at time is zero time or the specified time is ealier than current awake at time,
// call s.AwakeAt(at), otherwise do nothing.
func (s *Sleep) AwakeAtEalier(at time.Time) {
	if awakeAt := s.GetAwakeAt(); awakeAt.IsZero() || at.Before(awakeAt) {
		s.AwakeAt(at)
	}
}

// if the specified time is later than current awake at time, call s.AwakeAt(at),
// otherwise do nothing.
func (s *Sleep) AwakeAtLater(at time.Time) {
	if at.After(s.GetAwakeAt()) {
		s.AwakeAt(at)
	}
}
