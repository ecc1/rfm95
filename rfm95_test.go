package rfm95

import (
	"testing"
	"unsafe"
)

func TestRFConfiguration(t *testing.T) {
	have := int(unsafe.Sizeof(RFConfiguration{}))
	want := RegVersion - RegOpMode + 1
	if have != want {
		t.Errorf("Sizeof(RFConfiguration) == %d, want %d", have, want)
	}
}
