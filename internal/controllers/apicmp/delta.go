package apicmp

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/astro-walker/tilt/internal/timecmp"
)

func Comparators() []interface{} {
	return []interface{}{
		func(a, b resource.Quantity) bool {
			return a.Cmp(b) == 0
		},
		func(a, b metav1.MicroTime) bool {
			return timecmp.Equal(a, b)
		},
		func(a, b metav1.Time) bool {
			return timecmp.Equal(a, b)
		},
		func(a, b labels.Selector) bool {
			return a.String() == b.String()
		},
		func(a, b fields.Selector) bool {
			return a.String() == b.String()
		},
	}
}

// A deep equality check to see if a client object and
// a server object are different, such that the server object
// needs to be updated.
var delta = conversion.EqualitiesOrDie(Comparators()...)

func DeepEqual(a, b interface{}) bool {
	typeA := reflect.TypeOf(a)
	typeB := reflect.TypeOf(b)
	if typeA != typeB {
		panic(fmt.Sprintf("internal error: comparing incommensurable objects: %T, %T", a, b))
	}
	return delta.DeepEqual(a, b)
}
