// Package unit provides dimensional analysis and conversions
package unit

import (
	"errors"
	"strconv"
	"strings"
)

const (
	ZeroCelsiusInKelvin = 273.15 // 0 °C in K
)

var (
	tbd = errors.New("tbd")
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

func (a *pUnit) Inverse() *pUnit {
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

// Returns the inverse of a measurement. E.g., Inverse(2 m/s) = 1/2 s/m
func Inverse(mm Measurement) (Measurement, error) {
	m, err := parse(mm)
	if err != nil {
		return nil, err
	}
	return m, tbd
}

// Rescale measurement to target unit scale.
//
// If the quantity would be too large to represent in the target scale, return
// an overflow error. Ignore underflows when a quantity is too small to
// represent in the target scale. It is the responsibility of the caller to
// provide a unit that preserves the desired precision.
func ScaleTo(unit string, m Measurement) (Measurement, error) {
	return nil, tbd
}

// Convert measurement by applying factor. Rescale measurement m and factor to
// target unit scale.
//
// If either quantity would be too large to represent in the target scale,
// return an overflow error. Ignore underflows when a quantity is too small to
// represent in the target scale. It is the responsibility of the caller to
// provide a unit that preserves the desired precision.
func ConvertTo(unitString string, m, factor Measurement) (Measurement, error) {
	return nil, tbd
}

// Convert one measurement to another dimension or scale by applying conversion
// factors.
//
// The units for intermediate terms is unspecified and may change. If either an
// overflow or underflow occurs, an error will be returned. For more control
// over conversions consider using ConvertTo() and ScaleTo().
func New(unitString string, ms ...Measurement) (Measurement, error) {
	return nil, tbd
}

// Convenience function that panics if measurement operation fails.
func Must(m Measurement, err error) Measurement {
	if err != nil {
		panic(err)
	}
	return m
}
