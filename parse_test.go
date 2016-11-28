package unit

import "testing"

func makeUnitMap() map[string]*pUnit {
	m := make(map[string]*pUnit)
	for _, ku := range defaultUnits {
		m[ku.Key] = ku.Unit
	}
	return m
}

func TestParseDimensions(t *testing.T) {
	type testCase struct {
		Unit       string
		Expected   uPoint
		ShouldFail bool
	}

	um := makeUnitMap()

	suite := []testCase{
		testCase{
			Unit:     "m",
			Expected: um["m"].product(),
		},
		testCase{
			Unit:     "   m   ",
			Expected: um["m"].product(),
		},
		testCase{
			Unit:     "k°C",
			Expected: um["°C"].product(),
		},
		testCase{
			Unit:     "k℃",
			Expected: um["°C"].product(),
		},
		testCase{
			Unit:       "°C°C",
			ShouldFail: true,
		},
		testCase{
			Unit:       "(°C)°C",
			ShouldFail: true,
		},
		testCase{
			Unit:       "°C^2°C",
			ShouldFail: true,
		},
		testCase{
			Unit:     "°C^2 °C",
			Expected: um["°C"].Exp(3).product(),
		},
		testCase{
			Unit:     "(m)",
			Expected: um["m"].product(),
		},
		testCase{
			Unit:     "kg",
			Expected: um["g"].product(),
		},
		testCase{
			Unit:     "kmol",
			Expected: um["mol"].product(),
		},
		testCase{
			Unit:     "mm",
			Expected: um["m"].product(),
		},
		testCase{
			Unit:     "kmol / s",
			Expected: um["mol"].Multiply(um["s"].Inverse()).product(),
		},
		testCase{
			Unit:     "g/L",
			Expected: um["g"].Multiply(um["L"].Inverse()).product(),
		},
		testCase{
			Unit:     "ug/uL",
			Expected: um["g"].Multiply(um["L"].Inverse()).product(),
		},
		testCase{
			Unit:     "s^-1",
			Expected: um["s"].Inverse().product(),
		},
		testCase{
			Unit:     "kg·m/s^2",
			Expected: um["N"].product(),
		},
		testCase{
			Unit:     "kg·m/(s^2 s)",
			Expected: um["N"].Multiply(um["Hz"]).product(),
		},
		testCase{
			Unit:     "(kg·m/(s^2 s))^-1",
			Expected: um["N"].Multiply(um["Hz"]).Inverse().product(),
		},
		testCase{
			Unit:     "N/m^2",
			Expected: um["Pa"].product(),
		},
		testCase{
			Unit:       "(m",
			ShouldFail: true,
		},
		testCase{
			Unit:       "molk",
			ShouldFail: true,
		},
	}

	for _, tc := range suite {
		m, err := Parse(1.0, tc.Unit)
		if tc.ShouldFail {
			if err == nil {
				t.Errorf("expecting error but found unit: %q", m.(*measure).unit.product())
			}
			continue
		}

		if err != nil {
			t.Errorf("failed to parse %q: %s", tc.Unit, err)
		} else if e, f := tc.Expected, m.(*measure).unit.product(); e != f {
			t.Errorf("failed to parse %q: expected %q found %q", tc.Unit, e, f)
		} else {
			//t.Logf("%q is %q\n", tc.Unit, f)
		}
	}
}
