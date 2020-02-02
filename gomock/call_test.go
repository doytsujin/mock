package gomock

import (
	"reflect"
	"testing"
)

type mockTestReporter struct {
	errorCalls int
	fatalCalls int
}

func (o *mockTestReporter) Errorf(format string, args ...interface{}) {
	o.errorCalls++
}

func (o *mockTestReporter) Fatalf(format string, args ...interface{}) {
	o.fatalCalls++
}

func (o *mockTestReporter) Helper() {}

func TestCall_After(t *testing.T) {
	t.Run("SelfPrereqCallsFatalf", func(t *testing.T) {
		tr1 := &mockTestReporter{}

		c := &Call{t: tr1}
		c.After(c)

		if tr1.fatalCalls != 1 {
			t.Errorf("number of fatal calls == %v, want 1", tr1.fatalCalls)
		}
	})

	t.Run("LoopInCallOrderCallsFatalf", func(t *testing.T) {
		tr1 := &mockTestReporter{}
		tr2 := &mockTestReporter{}

		c1 := &Call{t: tr1}
		c2 := &Call{t: tr2}
		c1.After(c2)
		c2.After(c1)

		if tr1.errorCalls != 0 || tr1.fatalCalls != 0 {
			t.Error("unexpected errors")
		}

		if tr2.fatalCalls != 1 {
			t.Errorf("number of fatal calls == %v, want 1", tr2.fatalCalls)
		}
	})
}

func TestCall_Do(t *testing.T) {
	t.Run("Do function matches Call function", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function matches Call function and is a interface{}", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x interface{}) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function matches Call function and is a map[interface{}]interface{}", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x map[int]string) bool {
			return true
		}

		callFunc := func(x map[interface{}]interface{}) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function matches Call function and is variadic", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x []int) bool {
			return true
		}

		callFunc := func(x ...int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function matches Call function and is variadic interface{}", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x []int) bool {
			return true
		}

		callFunc := func(x ...interface{}) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("argument to Do is not a function", func(t *testing.T) {
		tr := &mockTestReporter{}

		callFunc := func(x int, y int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do("meow")

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("number of args for Do func don't match Call func", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x int, y int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("arg types for Do func don't match Call func", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x string) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function does not match Call function and is a slice", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x []string) bool {
			return true
		}

		callFunc := func(x []int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function does not match Call function and is a slice interface{}", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x []string) bool {
			return true
		}

		callFunc := func(x []interface{}) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function does not match Call function and is a composite struct", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x b) bool {
			return true
		}

		callFunc := func(x a) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("Do function does not match Call function and is a map", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x map[int]string) bool {
			return true
		}

		callFunc := func(x map[interface{}]int) bool {
			return false
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("number of return vals for Do func don't match Call func", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x int) (bool, error) {
			return false, nil
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})

	t.Run("return types for Do func don't match Call func", func(t *testing.T) {
		tr := &mockTestReporter{}

		doFunc := func(x int) bool {
			if x < 20 {
				return false
			}

			return true
		}

		callFunc := func(x int) error {
			return nil
		}

		c := &Call{
			t:          tr,
			methodType: reflect.TypeOf(callFunc),
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("expected Do to panic")
			}
		}()

		c.Do(doFunc)

		if len(c.actions) != 1 {
			t.Errorf("expected %d actions but got %d", 1, len(c.actions))
		}
	})
}

type a struct {
	name string
}

func (testObj a) Name() string {
	return testObj.name
}

type b struct {
	a
	foo string
}

func (testObj b) Foo() string {
	return testObj.foo
}
