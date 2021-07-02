package images

import "testing"

type compareVersionTest struct {
	v1  string
	v2  string
	out bool
}

var compareVersionsLtTest = []compareVersionTest{
	{"1.0.0", "1.0.0", false},
	{"1.0.0", "1.0.1", true},
	{"0.0.1", "0.1.1", true},
	{"1.0.1", "0.1.1", false},
	{"1.0.1", "1.0.12", true},
	{"1.0.12", "1.0.1", false},
	{"1.0.12", "1.0.2", false},
}

var compareVersionsEqTest = []compareVersionTest{
	{"1.0.0", "1.0.0", true},
	{"1.0.0", "1.0.1", false},
}

var compareVersionsGtTest = []compareVersionTest{
	{"1.0.0", "1.0.0", false},
	{"1.0.0", "1.0.1", false},
	{"1.0.1", "1.0.0", true},
}

func TestCompareVersions(t *testing.T) {
	vc := getGenericVersionComparer()

	for _, test := range compareVersionsLtTest {
		if outLt := vc.Lt(test.v1, test.v2); outLt != test.out {
			t.Errorf("Incorrectly determined if %s < %s", test.v1, test.v2)
		}
	}

	for _, test := range compareVersionsEqTest {
		if outEq := vc.Eq(test.v1, test.v2); outEq && !test.out {
			t.Errorf("Determined versions %s and %s to be equal but they are not", test.v1, test.v2)
		} else if !outEq && test.out {
			t.Errorf("Determined versions %s and %s not to be equal but they are", test.v1, test.v2)
		}
	}

	for _, test := range compareVersionsGtTest {
		if outGt := vc.Gt(test.v1, test.v2); outGt != test.out {
			t.Errorf("Incorrectly determined if %s > %s", test.v1, test.v2)
		}
	}
}
