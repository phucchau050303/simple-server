package credential

import (
	"testing"
)

// Test hasing function with phrase "isk4P2ost@%"
func TestHash(t *testing.T) {
	password := "isk4P2ost@%"
	hashed, err := Hash(password)
	t.Logf("hashed password is: %q", hashed)
	if hashed == "" || err != nil {
		t.Errorf(`Hash("isk4P2ost%%") can't be processed, error: %v`, err)
	}
}

// Test verify function with password "isk4P2ost@%"
// and hash pE7m65TYdV9Fw$atZkkSApihNocY6aQX3zXkP8wy06s22UM+2lQdxTIEY
func TestVerify(t *testing.T) {
	password := "isk4P2ost@%"
	hashed := "$argon2id$v=19$m=65536,t=1,p=8$IyrBdcszwuvo3jOLD/Lj8Q$rlRlQHBZ5g/x1IPbAkM+kurZz2mrBsRDGeDl2d7zGkw"
	match, err := Verify(password, hashed)
	if !match || err != nil {
		t.Errorf(`Verify("isk4P2ost@%%", "$argon2id$v=19$m=65536,t=1,p=8$IyrBdcszwuvo3jOLD/Lj8Q$rlRlQHBZ5g/x1IPbAkM+kurZz2mrBsRDGeDl2d7zGkw") = %t, %v, want true, error`, match, err)
	}
}
