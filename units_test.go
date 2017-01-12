package units

import (
	"math"
	"testing"
)

func TestNew(t *testing.T) {
	m, err := New("mm", Must(Parse(3.0, "m")))
	if err != nil {
		t.Error(err)
	} else if e, f := "mm", m.MeasurementUnit(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 3000.0, m.Quantity(); e != f {
		t.Errorf("expecting %v found %v", e, f)
	}

	m, err = New("mg/L", Must(Parse(3.0, "g")), Must(Reciprocal(Must(Parse(3.0, "ml")))))
	if err != nil {
		t.Error(err)
	} else if e, f := "mg/L", m.MeasurementUnit(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 1e6, m.Quantity(); e != f {
		t.Errorf("expecting %v found %v", e, f)
	}

	m, err = New("mg/(cm)^3", Must(Parse(1.0, "g")), Must(Reciprocal(Must(Parse(2.0, "ml")))))
	if err != nil {
		t.Error(err)
	} else if e, f := "mg/(cm)^3", m.MeasurementUnit(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 0.5, m.Quantity(); e != f {
		t.Errorf("expecting %v found %v", e, f)
	}

	m, err = New("g", Must(Parse(1.0, "g")), Must(Reciprocal(Must(Parse(1.0, "l")))))
	if err == nil {
		t.Errorf("expecting error got %v", m)
	}
}

func TestDimensionLess(t *testing.T) {
	m, err := New("", Must(Parse(1.0, "g")), Must(Reciprocal(Must(Parse(2.0, "g")))))
	if err != nil {
		t.Error(err)
	} else if e, f := "", m.MeasurementUnit(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 0.5, m.Quantity(); e != f {
		t.Errorf("expecting %v found %v", e, f)
	}

	m, err = New("g / g", Must(Parse(1.0, "g")), Must(Reciprocal(Must(Parse(2.0, "g")))))
	if err != nil {
		t.Error(err)
	} else if e, f := "g / g", m.MeasurementUnit(); e != f {
		t.Errorf("expecting %q found %q", e, f)
	} else if e, f := 0.5, m.Quantity(); e != f {
		t.Errorf("expecting %v found %v", e, f)
	}
}

func TestScaleOverflow(t *testing.T) {
	// MaxFloat ~ 1e308
	// Yg -> yg ~ 1e48
	// 1e308 / 1e48 = 6.4...
	big := Must(Parse(1.0, "Yg"))
	var rest []Measurement
	for i := 0; i < 6; i++ {
		rest = append(rest, big)
	}
	m, err := New("yg^7", big, rest...)
	if err == nil {
		t.Errorf("expecting error got %v", m)
	}

	m, err = New("yg^6", big, rest[:len(rest)-1]...)
	if err != nil {
		t.Error(err)
	}
}

func TestValueOverflow(t *testing.T) {
	// MaxFloat ~ 1e308
	// Yg -> yg ~ 1e48
	// 1e308 / 1e48 = 6.4...
	m, err := New("g^2", Must(Parse(math.MaxFloat64, "g")), Must(Parse(2.0, "g")))
	if err == nil {
		t.Errorf("expecting error got %v", m)
	}

	m, err = New("g", Must(Parse(math.MaxFloat64, "g^2")), Must(Reciprocal(Must(Parse(2.0, "g")))))
	if err != nil {
		t.Error(err)
	}
}

func TestValueUnderflow(t *testing.T) {
	m, err := New("g", Must(Parse(math.SmallestNonzeroFloat64, "g^2")), Must(Reciprocal(Must(Parse(2.0, "g")))))
	if err == nil {
		t.Errorf("expecting error got %v", m)
	}

	// da = 10^1
	m, err = New("dag", Must(Parse(math.SmallestNonzeroFloat64, "g")))
	if err == nil {
		t.Errorf("expecting error got %v", m)
	}

	// d = 10^-1
	m, err = New("dg", Must(Parse(math.SmallestNonzeroFloat64, "g")))
	if err != nil {
		t.Error(err)
	}
}
