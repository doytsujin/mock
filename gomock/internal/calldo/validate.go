// Copyright 2020 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package calldo

import (
	"fmt"
	"reflect"

	"github.com/pkg/errors"
)

// ValidateInputAndOutputSig compares the argument and return signatures of the
// function passed to Do against those expected by Call. It returns an error
// unless everything matches.
func ValidateInputAndOutputSig(doFunc, callFunc reflect.Type) error {
	// check number of arguments and type of each argument
	if doFunc.NumIn() != callFunc.NumIn() {
		return fmt.Errorf(
			"Do: expected function to have %d arguments not %d",
			callFunc.NumIn(), doFunc.NumIn())
	}

	lastIdx := callFunc.NumIn()

	// If the function has a variadic argument validate that one first so that
	// we aren't checking for it while we iterate over the other args
	if callFunc.IsVariadic() {
		if ok := validateVariadicArg(lastIdx, doFunc, callFunc); !ok {
			i := lastIdx - 1
			return fmt.Errorf(
				"Do: expected function to have"+
					" arg of type %v at position %d"+
					" not type %v",
				callFunc.In(i), i, doFunc.In(i),
			)
		}

		lastIdx--
	}

	for i := 0; i < lastIdx; i++ {
		callArg := callFunc.In(i)
		doArg := doFunc.In(i)

		if err := validateArg(doArg, callArg); err != nil {
			return fmt.Errorf("input argument at %d: %s", i, err)
		}
	}

	// check number of return vals and type of each val
	if doFunc.NumOut() != callFunc.NumOut() {
		return fmt.Errorf(
			"Do: expected function to have %d return vals not %d",
			callFunc.NumOut(), doFunc.NumOut())
	}

	for i := 0; i < callFunc.NumOut(); i++ {
		callArg := callFunc.Out(i)
		doArg := doFunc.Out(i)

		if err := validateArg(doArg, callArg); err != nil {
			return errors.Wrapf(err, "return argument at %d", i)
		}
	}

	return nil
}

func validateVariadicArg(lastIdx int, doFunc, callFunc reflect.Type) bool {
	if doFunc.In(lastIdx-1) != callFunc.In(lastIdx-1) {
		if doFunc.In(lastIdx-1).Kind() != reflect.Slice {
			return false
		}

		callArgT := callFunc.In(lastIdx - 1)
		callElem := callArgT.Elem()
		if callElem.Kind() != reflect.Interface {
			return false
		}

		doArgT := doFunc.In(lastIdx - 1)
		doElem := doArgT.Elem()

		if ok := doElem.ConvertibleTo(callElem); !ok {
			return false
		}

	}

	return true
}

func validateInterfaceArg(doArg, callArg reflect.Type) error {
	if !doArg.ConvertibleTo(callArg) {
		return fmt.Errorf(
			"expected arg convertible to type %v not type %v",
			callArg, doArg,
		)
	}

	return nil
}

func validateMapArg(doArg, callArg reflect.Type) error {
	callKey := callArg.Key()
	doKey := doArg.Key()

	switch callKey.Kind() {
	case reflect.Interface:
		if err := validateInterfaceArg(doKey, callKey); err != nil {
			return errors.Wrap(err, "map key")
		}
	default:
		if doKey != callKey {
			return fmt.Errorf("expected map key of type %v not type %v",
				callKey, doKey)
		}
	}

	callElem := callArg.Elem()
	doElem := doArg.Elem()

	switch callElem.Kind() {
	case reflect.Interface:
		if err := validateInterfaceArg(doElem, callElem); err != nil {
			return errors.Wrap(err, "map element")
		}
	default:
		if doElem != callElem {
			return fmt.Errorf("expected map element of type %v not type %v",
				callElem, doElem)
		}
	}

	return nil
}

func validateArg(doArg, callArg reflect.Type) error {
	switch callArg.Kind() {
	// If the Call arg is an interface we only care if the Do arg is convertible
	// to that interface
	case reflect.Interface:
		if err := validateInterfaceArg(doArg, callArg); err != nil {
			return err
		}
	default:
		// If the Call arg is not an interface then first check to see if
		// the Do arg is even the same reflect.Kind
		if callArg.Kind() != doArg.Kind() {
			return fmt.Errorf("expected arg of kind %v not %v",
				callArg.Kind(), doArg.Kind())
		}

		switch callArg.Kind() {
		// If the Call arg is a map then we need to handle the case where
		// the map key or element type is an interface
		case reflect.Map:
			if err := validateMapArg(doArg, callArg); err != nil {
				return err
			}
		default:
			if doArg != callArg {
				return fmt.Errorf(
					"Expected arg of type %v not type %v",
					callArg, doArg,
				)
			}
		}
	}

	return nil
}
