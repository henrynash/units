package units

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	theZero        = &pUnit{}
	symbolNotFound = errors.New("symbol not found")
	prefixNotFound = errors.New("prefix not found")
	unparsedText   = errors.New("unparsed text")
)

// Make a new parser error
func makeParseError(data []byte, pos int, err error) error {
	return errors.New("parse failed at: " +
		strconv.Quote(string(data[:pos])) + " . " +
		strconv.Quote(string(data[pos:])) + ": " +
		err.Error())
}

func makeRuneNotFoundError(r rune) error {
	return errors.New(strconv.QuoteRune(r) + " not found")
}

// Parse a quantity and unit into a measurement.
//
// Unit Grammar:
//   ValidUnit := Unit
//              | ""    # Dimensionless measurement
//   Unit      := Term
//              | ( Unit )        # Grouping
//              | Unit  ^ Integer # Unit exponentiation
//              | Unit  /  Unit   # Unit division
//              | Unit  ·  Unit   # Unit multiplication (· is center dot)
//              | Unit " " Unit   # Unit multiplication (" " is whitespace)
//   Term      := Prefix? Symbol
//   #            1    2   3   6   9  12  15  18  21  24  # Exp
//   Prefix    := da | h | k | M | G | T | P | E | Z | Y  # 10^Exp
//              | d  | c | m | μ | n | p | f | a | z | y  # 10^-Exp
//              |              u
//   Symbol    := m   | g  | s  | A | K  | mol | cd  # Base dimensions
//              | rad | st | Hz | N | Pa | J         # Derived units
//              | W   | C  | V  | F | Ω  | S
//              | Wb  | T  | H  | °C | ℃
//              | lm  | lx | Bq | Gy | Sv | kat
//              | l   | L  | Da                      # Non-SI units
//   Integer   := ..., -2, -1, 0, 1, 2, ...
//
// Examples:
//   - A newton: N, kg m s^-2, kg·m/s^2
//   - A pascal: Pa, N/m^2, kg·m^−1·s−2
//   - A litre: l, L, dm^3
//
// Notes:
//   - The associativity of unit division is unspecified in the International
//   System of Units. For example, a/b/c can mean (a/b)/c or a/(b/c). This
//   library may or may not accept such ambigious units. For portability, users
//   should parenthesize or convert division to exponentiation.
//   - C is Coulomb; °C or ℃ is degree Celsius
func Parse(quantity float64, unitString string) (Measurement, error) {
	data := []byte(unitString)

	if len(data) == 0 {
		return &measure{
			unit: theZero,
		}, nil
	}

	unit, pos, err := parseUnit(data, 0)
	if err != nil {
		return nil, makeParseError(data, pos, err)
	}
	pos, _ = scanToNonSpace(data, pos, false)
	if pos != len(data) {
		return nil, makeParseError(data, pos, unparsedText)
	}

	return &measure{
		Value: quantity,
		Unit:  unitString,
		unit:  unit,
	}, nil
}

func parseUnit(data []byte, pos int) (*pUnit, int, error) {
	var unit *pUnit
	pos, _ = scanToNonSpace(data, pos, false)

	// Unit := ( Unit ) | Term
	pos, err := parseRune(data, pos, '(')
	if err == nil {
		if unit, pos, err = parseUnit(data, pos); err != nil {
			return nil, pos, err
		} else if pos, err = parseRune(data, pos, ')'); err != nil {
			return nil, pos, err
		}
	} else {
		unit, pos, err = parseTerm(data, pos)
		if err != nil {
			return nil, pos, err
		}
	}

	var nextUnit *pUnit
	var exp uComponent

	pos, hadSpace := scanToNonSpace(data, pos, false)

	// Unit := Unit ^ Integer
	if exp, pos, err = parseExponent(data, pos); err == nil {
		unit = unit.Exp(exp)
		pos, hadSpace = scanToNonSpace(data, pos, false)
	}

	// Unit := ...
	if pos, err = parseRune(data, pos, '/'); err == nil {
		// ... | Unit / Unit
		nextUnit, pos, err = parseUnit(data, pos)
		if err != nil {
			return nil, pos, err
		}
		unit = unit.Multiply(nextUnit.Reciprocal())
	} else if pos, err = parseRune(data, pos, '·'); err == nil {
		// ... |  Unit · Unit
		nextUnit, pos, err = parseUnit(data, pos)
		if err != nil {
			return nil, pos, err
		}
		unit = unit.Multiply(nextUnit)
	} else if hadSpace {
		// ... | Unit " " Unit
		nextUnit, pos, err = parseUnit(data, pos)
		if err == nil {
			unit = unit.Multiply(nextUnit)
		} else {
			// Unit parse is done; let caller decide if this is an error
		}
	} else {
		// Unit parse is done; let caller decide if this is an error
	}

	return unit, pos, nil
}

func parseRune(data []byte, pos int, r rune) (int, error) {
	if len(data) <= pos {
		return pos, makeRuneNotFoundError(r)
	}
	dr, width := utf8.DecodeRune(data[pos:])
	if dr != r {
		return pos, makeRuneNotFoundError(r)
	}
	return pos + width, nil
}

func parseTerm(data []byte, startPos int) (*pUnit, int, error) {
	// Term := Prefix?
	scale, pos, _ := parsePrefix(data, startPos)
	//   ... Symbol
	unit, pos, err := parseSymbol(data, pos)
	if err != nil {
		// Some symbols are substrings of prefixes (e.g., m(illi) and m(eter)),
		// so try Term := Symbol as well
		unit, pos, err = parseSymbol(data, startPos)
		if err != nil {
			return nil, pos, err
		}
		scale = 0
	}
	return &pUnit{
		Dim:   unit.Dim,
		Scale: scale,
	}, pos, nil
}

func parsePrefix(data []byte, pos int) (int, int, error) {
	if len(data) <= pos {
		return 0, pos, prefixNotFound
	}

	str := string(data[pos:])

	for _, ks := range defaultScales {
		if strings.HasPrefix(str, ks.Key) {
			return ks.Scale, pos + len(ks.Key), nil
		}
	}

	return 0, pos, prefixNotFound
}

func parseSymbol(data []byte, pos int) (*pUnit, int, error) {
	if len(data) <= pos {
		return nil, pos, prefixNotFound
	}

	str := string(data[pos:])
	for _, ku := range defaultUnits {
		if strings.HasPrefix(str, ku.Key) {
			return ku.Unit, pos + len(ku.Key), nil
		}
	}

	return nil, pos, symbolNotFound
}

func parseExponent(data []byte, pos int) (uComponent, int, error) {
	pos, err := parseRune(data, pos, '^')
	if err != nil {
		return 1, pos, err
	}

	// Scan to end of integer string
	end := scanToNonDigit(data, pos)

	// Parse
	i, err := strconv.ParseInt(string(data[pos:end]), 10, 8)
	if err != nil {
		return 1, pos, err
	}
	return uComponent(i), end, nil
}

// Return position of first non-space or len(data) if none and if whitespace
// was found
func scanToNonSpace(data []byte, pos int, hadSpace bool) (int, bool) {
	if len(data) <= pos {
		return pos, hadSpace
	}
	idx := bytes.IndexFunc(data[pos:], func(r rune) bool {
		return !unicode.IsSpace(r)
	})
	if idx < 0 {
		return len(data), hadSpace || true
	}
	return pos + idx, hadSpace || idx != 0
}

// Return position of first non-digit or len(data) if none
func scanToNonDigit(data []byte, pos int) int {
	idx := pos
	for width := 0; idx < len(data); idx += width {
		v, w := utf8.DecodeRune(data[idx:])
		width = w
		if idx == pos && v == '-' {
			continue
		}
		if !unicode.IsDigit(v) {
			break
		}
	}
	return idx
}

func parse(m Measurement) (*measure, error) {
	if m, ok := m.(*measure); ok {
		return m, nil
	}
	m, err := Parse(m.Quantity(), m.MeasurementUnit())
	if err != nil {
		return nil, err
	}
	return m.(*measure), nil
}
