package sleep_test

import (
	"fmt"
	"time"

	"github.com/lovego/sleep"
)

func ExampleSleep_Sleep() {
	var s sleep.Sleep
	go func() {
		s.Sleep(0)
		fmt.Println("awaken up")
	}()
	time.Sleep(time.Millisecond)
	// Output:
	// awaken up
}

func ExampleSleep_Asleep() {
	var s sleep.Sleep
	go func() {
		time.Sleep(time.Millisecond)
		fmt.Println(s.Asleep())
		time.Sleep(2 * time.Millisecond)
		fmt.Println(s.Asleep())
	}()
	s.Sleep(2 * time.Millisecond)
	fmt.Println(s.Asleep())
	time.Sleep(2 * time.Millisecond)
	// Output:
	// true
	// false
	// false
}

func ExampleSleep_GetAwakeAt() {
	var s sleep.Sleep
	s.Sleep(0)
	fmt.Println(s.GetAwakeAt().IsZero())
	s.ClearAwakeAt()
	fmt.Println(s.GetAwakeAt().IsZero())
	// Output:
	// false
	// true
}

func ExampleSleep_Awake() {
	var s sleep.Sleep
	go func() {
		s.Sleep(time.Minute)
		fmt.Println("awaken up")
	}()
	time.Sleep(time.Millisecond)
	s.Awake()
	time.Sleep(time.Millisecond)

	// Output:
	// awaken up
}

func ExampleSleep_AwakeAtEalier_1() {
	var s sleep.Sleep
	fmt.Println(s.GetAwakeAt().IsZero())
	at := time.Now().Add(time.Minute)
	s.AwakeAtEalier(at)
	fmt.Println(s.GetAwakeAt().Equal(at))
	s.AwakeAtEalier(at.Add(time.Second))
	fmt.Println(s.GetAwakeAt().Equal(at))

	s.AwakeAtEalier(time.Now().Add(time.Millisecond))
	go func() {
		s.Run()
		fmt.Println("awaken up")
	}()
	time.Sleep(3 * time.Millisecond)

	// Output:
	// true
	// true
	// true
	// awaken up
}

func ExampleSleep_AwakeAtEalier_2() {
	var s sleep.Sleep
	s.AwakeAtEalier(time.Now().Add(time.Minute))
	go func() {
		s.Run()
		fmt.Println("awaken up")
	}()
	time.Sleep(time.Millisecond)
	s.AwakeAtEalier(time.Now().Add(time.Millisecond))
	time.Sleep(3 * time.Millisecond)
	// Output:
	// awaken up
}

func ExampleSleep_AwakeAtLater_1() {
	var s sleep.Sleep
	fmt.Println(s.GetAwakeAt().IsZero())
	at := time.Now().Add(-time.Minute)
	s.AwakeAtLater(at)
	fmt.Println(s.GetAwakeAt().Equal(at))
	s.AwakeAtLater(at.Add(-time.Second))
	fmt.Println(s.GetAwakeAt().Equal(at))

	s.AwakeAtLater(time.Now().Add(time.Millisecond))
	go func() {
		s.Run()
		fmt.Println("awaken up")
	}()
	time.Sleep(3 * time.Millisecond)

	// Output:
	// true
	// true
	// true
	// awaken up
}

func ExampleSleep_AwakeAtLater_2() {
	var s sleep.Sleep
	s.AwakeAtLater(time.Now().Add(2 * time.Millisecond))
	go func() {
		s.Run()
		fmt.Println("awaken up")
	}()
	time.Sleep(time.Millisecond)
	s.AwakeAtLater(time.Now().Add(time.Second))
	time.Sleep(3 * time.Millisecond)
	// Output:
}
