package encrypt

import (
	"testing"
)

func Test_Encrypt(t *testing.T) {

	uid := "92670036894"
	idCard := "13048119980624067x"
	phone := "19833001203"

	t.Logf("uid: %v\n idCard: %v\nphone: %v\n", uid, Sha256(idCard), Sha256(phone))
}
