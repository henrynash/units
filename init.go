package unit

import (
	"errors"
	"sort"
	"strconv"
)

var (
	defaultScales []keyedScale
	defaultUnits  []keyedUnit
)

type keyedUnit struct {
	Key  string
	Unit *pUnit
}

type keyedUnitSlice []keyedUnit

func (a keyedUnitSlice) Len() int {
	return len(a)
}

func (a keyedUnitSlice) Less(i, j int) bool {
	return longLess(a[i].Key, a[j].Key)
}

func (a keyedUnitSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type keyedScale struct {
	Key   string
	Scale int
}

type keyedScaleSlice []keyedScale

func (a keyedScaleSlice) Len() int {
	return len(a)
}

func (a keyedScaleSlice) Less(i, j int) bool {
	return longLess(a[i].Key, a[j].Key)
}

func (a keyedScaleSlice) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Sort longer strings before shorter strings. This ensures that we find the
// longest match first.
func longLess(a, b string) bool {
	if la, lb := len(a), len(b); la == lb {
		return a < b
	} else {
		return la > lb
	}
}

func makeScales() ([]keyedScale, error) {
	var r keyedScaleSlice

	r = append(r, keyedScale{Key: "da", Scale: 1})
	r = append(r, keyedScale{Key: "h", Scale: 2})
	r = append(r, keyedScale{Key: "k", Scale: 3})
	r = append(r, keyedScale{Key: "M", Scale: 6})
	r = append(r, keyedScale{Key: "G", Scale: 9})
	r = append(r, keyedScale{Key: "T", Scale: 12})
	r = append(r, keyedScale{Key: "P", Scale: 15})
	r = append(r, keyedScale{Key: "E", Scale: 18})
	r = append(r, keyedScale{Key: "Z", Scale: 21})
	r = append(r, keyedScale{Key: "Y", Scale: 24})

	r = append(r, keyedScale{Key: "d", Scale: -1})
	r = append(r, keyedScale{Key: "c", Scale: -2})
	r = append(r, keyedScale{Key: "m", Scale: -3})
	r = append(r, keyedScale{Key: "μ", Scale: -6})
	r = append(r, keyedScale{Key: "u", Scale: -6})
	r = append(r, keyedScale{Key: "n", Scale: -9})
	r = append(r, keyedScale{Key: "p", Scale: -12})
	r = append(r, keyedScale{Key: "f", Scale: -15})
	r = append(r, keyedScale{Key: "a", Scale: -18})
	r = append(r, keyedScale{Key: "z", Scale: -21})
	r = append(r, keyedScale{Key: "y", Scale: -24})

	seen := make(map[string]bool)
	for _, v := range r {
		if seen[v.Key] {
			return nil, errors.New("duplicate key " + strconv.Quote(v.Key))
		}
		seen[v.Key] = true
	}

	sort.Sort(r)

	return r, nil
}

func makeUnits() ([]keyedUnit, error) {
	type de map[int]uComponent

	mkpoint := func(m de) (a uPoint) {
		for d, e := range m {
			a[d] = e
		}
		return
	}

	var r keyedUnitSlice

	r = append(r, keyedUnit{
		Key: "m",
		Unit: &pUnit{
			Dim: mkpoint(de{
				lengthDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "g",
		Unit: &pUnit{
			Dim: mkpoint(de{
				massDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "s",
		Unit: &pUnit{
			Dim: mkpoint(de{
				timeDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "A",
		Unit: &pUnit{
			Dim: mkpoint(de{
				currentDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "K",
		Unit: &pUnit{
			Dim: mkpoint(de{
				temperatureDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "mol",
		Unit: &pUnit{
			Dim: mkpoint(de{
				amountDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "cd",
		Unit: &pUnit{
			Dim: mkpoint(de{
				intensityDim: 1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "rad",
		Unit: &pUnit{
			DimLess: []uPoint{
				mkpoint(de{
					lengthDim: 1,
				}),
				mkpoint(de{
					lengthDim: -1,
				}),
			},
		}})
	r = append(r, keyedUnit{
		Key: "sr",
		Unit: &pUnit{
			DimLess: []uPoint{
				mkpoint(de{
					lengthDim: 2,
				}),
				mkpoint(de{
					lengthDim: -2,
				}),
			},
		}})
	r = append(r, keyedUnit{
		Key: "Hz",
		Unit: &pUnit{
			Dim: mkpoint(de{
				timeDim: -1,
			}),
		}})
	r = append(r, keyedUnit{
		Key: "N",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:   1,
					timeDim:   -2,
					lengthDim: 1,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "Pa",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:   1,
					lengthDim: -1,
					timeDim:   -2,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "J",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:   1,
					lengthDim: 2,
					timeDim:   -2,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "W",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:   1,
					lengthDim: 2,
					timeDim:   -3,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "C",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					timeDim:    1,
					currentDim: 1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "V",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    1,
					lengthDim:  2,
					timeDim:    -3,
					currentDim: -1,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "F",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    -1,
					lengthDim:  -2,
					timeDim:    4,
					currentDim: 2,
				},
			),
			Scale: -3,
		}})
	r = append(r, keyedUnit{
		Key: "Ω",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    1,
					lengthDim:  2,
					timeDim:    -3,
					currentDim: -2,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "S",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    -1,
					lengthDim:  -2,
					timeDim:    3,
					currentDim: 2,
				},
			),
			Scale: -3,
		}})
	r = append(r, keyedUnit{
		Key: "Wb",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    1,
					lengthDim:  2,
					timeDim:    -2,
					currentDim: -1,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "T",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    1,
					timeDim:    -2,
					currentDim: -1,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "H",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:    1,
					lengthDim:  2,
					timeDim:    -2,
					currentDim: -2,
				},
			),
			Scale: 3,
		}})
	r = append(r, keyedUnit{
		Key: "°C",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					temperatureCDim: 1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "℃",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					temperatureCDim: 1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "lm",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					intensityDim: 1,
				},
			),
			DimLess: []uPoint{
				mkpoint(de{
					lengthDim: 2,
				}),
				mkpoint(de{
					lengthDim: -2,
				}),
			},
		}})
	r = append(r, keyedUnit{
		Key: "lx",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					lengthDim:    -2,
					intensityDim: 1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "Bq",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					timeDim: -1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "Gy",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					lengthDim: 2,
					timeDim:   -2,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "Sv",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					lengthDim: 2,
					timeDim:   -2,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "kat",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					amountDim: 1,
					timeDim:   -1,
				},
			),
		}})
	r = append(r, keyedUnit{
		Key: "l",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					lengthDim: 3,
				},
			),
			Scale: -3,
		}})
	r = append(r, keyedUnit{
		Key: "L",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					lengthDim: 3,
				},
			),
			Scale: -3,
		}})
	r = append(r, keyedUnit{
		Key: "Da",
		Unit: &pUnit{
			Dim: mkpoint(
				de{
					massDim:   1,
					amountDim: -1,
				},
			),
		}})

	seen := make(map[string]bool)
	for _, v := range r {
		if seen[v.Key] {
			return nil, errors.New("duplicate key " + strconv.Quote(v.Key))
		}
		seen[v.Key] = true
	}

	sort.Sort(r)

	return r, nil
}

func init() {
	var err error
	defaultScales, err = makeScales()
	if err != nil {
		panic(err)
	}

	defaultUnits, err = makeUnits()
	if err != nil {
		panic(err)
	}
}
