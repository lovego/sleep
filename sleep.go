package sleep

import (
	"sync"
	"time"
)

type Sleep struct {
	awakeAt time.Time
	event   interface{}
	asleep  bool
	events  chan interface{}
	sync.RWMutex

	debug bool
}

// Sleep method on a single `Sleep` instance should not be called concurrently,
// otherwise only a random one will be awaken up by `Awake*` methods.
// Sleep method on a single `Sleep` instance can be called many times serially.
// Sleep returns the event that awaken it.
func (s *Sleep) Sleep(d time.Duration, event interface{}) interface{} {
	s.SetAwakeAt(time.Now().Add(d), event)

	return s.Run()
}

// Run method on a single `Sleep` instance should not be called concurrently,
// otherwise only a random one will be awaken up by `Awake*` methods.
// Run method on a single `Sleep` instance can be called many times serially.
// Run returns the event that awaken it.
func (s *Sleep) Run() interface{} {
	s.RLock()
	var d = time.Until(s.awakeAt)
	var event = s.event
	s.RUnlock()

	if d <= 0 {
		return event
	}

	s.Lock()
	if s.events == nil {
		s.events = make(chan interface{})
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
			return event
		case event = <-s.events:
			timer.Stop()
			if d = time.Until(s.GetAwakeAt()); d <= 0 {
				return event
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
func (s *Sleep) SetAwakeAt(at time.Time, event interface{}) {
	s.Lock()
	defer s.Unlock()
	s.awakeAt = at
	s.event = event
}

// set awake at time to zero time.
func (s *Sleep) ClearAwakeAt() {
	s.Lock()
	defer s.Unlock()
	s.awakeAt = time.Time{}
}

// set awake at time to the specified time, and awake at the specified time if asleep.
func (s *Sleep) AwakeAt(at time.Time, event interface{}) {
	s.SetAwakeAt(at, event)
	select {
	case s.events <- event:
	default:
	}
}

// the same as s.AwakeAt(time.Now())
func (s *Sleep) Awake(event interface{}) {
	s.AwakeAt(time.Now(), event)
}

// if current awake at time is zero time or the specified time is ealier than current awake at time,
// call s.AwakeAt(at), otherwise do nothing.
func (s *Sleep) AwakeAtEalier(at time.Time, event interface{}) {
	if awakeAt := s.GetAwakeAt(); awakeAt.IsZero() || at.Before(awakeAt) {
		s.AwakeAt(at, event)
	}
}

// if the specified time is later than current awake at time, call s.AwakeAt(at),
// otherwise do nothing.
func (s *Sleep) AwakeAtLater(at time.Time, event interface{}) {
	if at.After(s.GetAwakeAt()) {
		s.AwakeAt(at, event)
	}
}
