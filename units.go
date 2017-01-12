// Package units provides dimensional analysis and conversions
package units

import (
	"errors"
	"math"
	"strconv"
	"strings"
)

const (
	ZeroCelsiusInKelvin = 273.15 // 0 °C in K
)

var (
	tbd              = errors.New("tbd")
	divideByZero     = errors.New("divide by zero")
	wrongDimension   = errors.New("wrong dimension")
	underflow        = errors.New("underflow")
	overflow         = errors.New("overflow")
	nothingToConvert = errors.New("nothing to convert")
)

var (
	// Zero
	zeroValue = &measure{
		unit: &pUnit{},
	}
)

// Base Dimensions
const (
	currentDim      = iota // I: Electric current
	intensityDim           // J: Luminous intensity
	lengthDim              // L
	massDim                // M
	amountDim              // N
	timeDim                // T
	temperatureDim         // Θ: Absolute temperature
	temperatureCDim        // ΘC: Celsius temperature
	numDim
)

// TODO(ddn): Consider canonicalizing uPoints with *uPoint when the number of
// dimensions grows. Currently sizeof(uPoint) ~ sizeof(float64) so it's not a
// big deal now.

// Component value in dimensional unit space
type uComponent int8

// Point in dimensional unit space
type uPoint [numDim]uComponent

func (a uPoint) String() string {
	labels := []string{
		"I", "J", "L", "M", "N", "T", "Θ", "ΘC",
	}
	var terms []string
	for idx, v := range a {
		if v == 0 {
			continue
		}
		label := "X"
		if idx < len(labels) {
			label = labels[idx]
		}
		terms = append(terms, label+"^"+strconv.FormatInt(int64(v), 10))
	}
	return strings.Join(terms, " ")
}

type Measurement interface {
	Quantity() float64
	MeasurementUnit() string
}

// Parsed unit
type pUnit struct {
	// For normal units, the unit dimensions
	Dim uPoint
	// For dimensionless units, the unnormalized factors (whose product should
	// be one).
	DimLess []uPoint
	Scale   int
}

// TODO: Currently pUnit operations are strongly normalizing. Need to revisit
// for better support for dimensionless values

// Return product of all dimension factors
func (a *pUnit) product() uPoint {
	var newDim uPoint
	for _, dim := range append(a.DimLess, a.Dim) {
		for idx, v := range dim {
			newDim[idx] += v
		}
	}
	return newDim
}

func (a *pUnit) Multiply(b *pUnit) *pUnit {
	r := a.product()
	for idx, v := range b.product() {
		r[idx] += v
	}

	return &pUnit{
		Dim:   r,
		Scale: a.Scale + b.Scale,
	}
}

func (a *pUnit) Reciprocal() *pUnit {
	r := a.product()
	for idx, v := range r {
		r[idx] = -v
	}

	return &pUnit{
		Dim:   r,
		Scale: -a.Scale,
	}
}

func (a *pUnit) Exp(e uComponent) *pUnit {
	r := a.product()
	for idx, v := range r {
		r[idx] = e * v
	}

	return &pUnit{
		Dim:   r,
		Scale: int(e) * a.Scale,
	}
}

// Parsed measurement
type measure struct {
	Value float64 // Quantity value
	Unit  string  // Given unit
	unit  *pUnit  // Parsed unit
}

func (a *measure) Quantity() float64 {
	return a.Value
}

func (a *measure) MeasurementUnit() string {
	return a.Unit
}

// Returns the reciprocal of a measurement. E.g., Reciprocal(2 m/s) = 1/2 s/m.
// The unit of the reciprocal is implementation dependent; use New to convert
// it to a specific unit of measure.
func Reciprocal(mm Measurement) (Measurement, error) {
	m, err := parse(mm)
	if err != nil {
		return zeroValue, err
	}
	if m.Value == 0.0 {
		return zeroValue, divideByZero
	}
	m.Value = 1.0 / m.Value
	m.unit = m.unit.Reciprocal()
	if len(m.Unit) != 0 {
		// TODO: canonicalize?
		m.Unit = "(" + m.Unit + ")^-1"
	}
	return m, nil
}

// Convert one measurement to another dimension or scale by applying conversion
// factors.
//
// The units for intermediate terms is unspecified and may change. If either an
// overflow or underflow occurs, an error will be returned.
func New(unitString string, m0 Measurement, ms ...Measurement) (Measurement, error) {
	m, err := parse(m0)
	if err != nil {
		return zeroValue, err
	}

	unit := m.unit
	value := m.Value
	for _, mm := range ms {
		m, err := parse(mm)
		if err != nil {
			return zeroValue, err
		}

		value = value * m.Value
		unit = unit.Multiply(m.unit)
	}

	targetM, err := Parse(0.0, unitString)
	if err != nil {
		return zeroValue, err
	}
	target := targetM.(*measure)

	if target.unit.product() != unit.product() {
		return zeroValue, wrongDimension
	}

	scaleDiff := unit.Scale - target.unit.Scale
	value *= math.Pow10(scaleDiff)
	if value == 0.0 {
		return zeroValue, underflow
	}
	if math.IsInf(value, 0) {
		return zeroValue, overflow
	}

	return &measure{
		Value: value,
		Unit:  unitString,
		unit:  unit,
	}, nil
}

// Convenience function that panics if measurement operation fails.
func Must(m Measurement, err error) Measurement {
	if err != nil {
		panic(err)
	}
	return m
}
