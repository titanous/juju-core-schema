package schema_test

import (
	"github.com/titanous/juju-schema"
	. "launchpad.net/gocheck"
	"testing"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct{}

var _ = Suite(&S{})

type Dummy struct{}

func (d *Dummy) Coerce(value interface{}, path []string) (coerced interface{}, err error) {
	return "i-am-dummy", nil
}

var aPath = []string{"<pa", "th>"}

func (s *S) TestConst(c *C) {
	sch := schema.Const("foo")

	out, err := sch.Coerce("foo", aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "foo")

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected "foo", got 42`)

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected "foo", got nothing`)
}

func (s *S) TestAny(c *C) {
	sch := schema.Any()

	out, err := sch.Coerce("foo", aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "foo")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, nil)
}

func (s *S) TestOneOf(c *C) {
	sch := schema.OneOf(schema.Const("foo"), schema.Const(42))

	out, err := sch.Coerce("foo", aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "foo")

	out, err = sch.Coerce(42, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, 42)

	out, err = sch.Coerce("bar", aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: unexpected value "bar"`)
}

func (s *S) TestBool(c *C) {
	sch := schema.Bool()

	out, err := sch.Coerce(true, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, true)

	out, err = sch.Coerce(false, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, false)

	out, err = sch.Coerce(1, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected bool, got 1")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected bool, got nothing")
}

func (s *S) TestInt(c *C) {
	sch := schema.Int()

	out, err := sch.Coerce(42, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, int64(42))

	out, err = sch.Coerce(int8(42), aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, int64(42))

	out, err = sch.Coerce(true, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected int, got true")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected int, got nothing")
}

func (s *S) TestFloat(c *C) {
	sch := schema.Float()

	out, err := sch.Coerce(float32(1.0), aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, float64(1.0))

	out, err = sch.Coerce(float64(1.0), aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, float64(1.0))

	out, err = sch.Coerce(true, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected float, got true")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected float, got nothing")
}

func (s *S) TestString(c *C) {
	sch := schema.String()

	out, err := sch.Coerce("foo", aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "foo")

	out, err = sch.Coerce(true, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected string, got true")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected string, got nothing")
}

func (s *S) TestSimpleRegexp(c *C) {
	sch := schema.SimpleRegexp()
	out, err := sch.Coerce("[0-9]+", aPath)
	c.Assert(err, IsNil)
	c.Assert(out, Equals, "[0-9]+")

	out, err = sch.Coerce(1, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected regexp string, got 1")

	out, err = sch.Coerce("[", aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected valid regexp, got "\["`)

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected regexp string, got nothing`)
}

func (s *S) TestList(c *C) {
	sch := schema.List(schema.Int())
	out, err := sch.Coerce([]int8{1, 2}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, []interface{}{int64(1), int64(2)})

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected list, got 42")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected list, got nothing")

	out, err = sch.Coerce([]interface{}{1, true}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\[1\]: expected int, got true`)
}

func (s *S) TestMap(c *C) {
	sch := schema.Map(schema.String(), schema.Int())
	out, err := sch.Coerce(map[string]interface{}{"a": 1, "b": int8(2)}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[interface{}]interface{}{"a": int64(1), "b": int64(2)})

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got 42")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got nothing")

	out, err = sch.Coerce(map[int]int{1: 1}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected string, got 1")

	out, err = sch.Coerce(map[string]bool{"a": true}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.a: expected int, got true`)

	// First path entry shouldn't have dots in an error message.
	out, err = sch.Coerce(map[string]bool{"a": true}, nil)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `a: expected int, got true`)
}

func (s *S) TestStringMap(c *C) {
	sch := schema.StringMap(schema.Int())
	out, err := sch.Coerce(map[string]interface{}{"a": 1, "b": int8(2)}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"a": int64(1), "b": int64(2)})

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got 42")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got nothing")

	out, err = sch.Coerce(map[int]int{1: 1}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected string, got 1")

	out, err = sch.Coerce(map[string]bool{"a": true}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.a: expected int, got true`)

	// First path entry shouldn't have dots in an error message.
	out, err = sch.Coerce(map[string]bool{"a": true}, nil)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `a: expected int, got true`)
}

func assertFieldMap(c *C, sch schema.Checker) {
	out, err := sch.Coerce(map[string]interface{}{"a": "A", "b": "B"}, aPath)

	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"a": "A", "b": "B", "c": "C"})

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got 42")

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, "<path>: expected map, got nothing")

	out, err = sch.Coerce(map[string]interface{}{"a": "A", "b": "C"}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.b: expected "B", got "C"`)

	out, err = sch.Coerce(map[string]interface{}{"b": "B"}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.a: expected "A", got nothing`)

	// b is optional
	out, err = sch.Coerce(map[string]interface{}{"a": "A"}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"a": "A", "c": "C"})

	// First path entry shouldn't have dots in an error message.
	out, err = sch.Coerce(map[string]bool{"a": true}, nil)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `a: expected "A", got true`)
}

func (s *S) TestFieldMap(c *C) {
	fields := schema.Fields{
		"a": schema.Const("A"),
		"b": schema.Const("B"),
		"c": schema.Const("C"),
	}
	defaults := schema.Defaults{
		"b": schema.Omit,
		"c": "C",
	}
	sch := schema.FieldMap(fields, defaults)
	assertFieldMap(c, sch)

	out, err := sch.Coerce(map[string]interface{}{"a": "A", "b": "B", "d": "D"}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"a": "A", "b": "B", "c": "C"})

	out, err = sch.Coerce(map[string]interface{}{"a": "A", "d": "D"}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"a": "A", "c": "C"})
}

func (s *S) TestFieldMapDefaultInvalid(c *C) {
	fields := schema.Fields{
		"a": schema.Const("A"),
	}
	defaults := schema.Defaults{
		"a": "B",
	}
	sch := schema.FieldMap(fields, defaults)
	_, err := sch.Coerce(map[string]interface{}{}, aPath)
	c.Assert(err, ErrorMatches, `<path>.a: expected "A", got "B"`)
}

func (s *S) TestStrictFieldMap(c *C) {
	fields := schema.Fields{
		"a": schema.Const("A"),
		"b": schema.Const("B"),
		"c": schema.Const("C"),
	}
	defaults := schema.Defaults{
		"b": schema.Omit,
		"c": "C",
	}
	sch := schema.StrictFieldMap(fields, defaults)
	assertFieldMap(c, sch)

	out, err := sch.Coerce(map[string]interface{}{"a": "A", "b": "B", "d": "D"}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>.d: expected nothing, got "D"`)
}

func (s *S) TestSchemaMap(c *C) {
	fields1 := schema.FieldMap(schema.Fields{
		"type": schema.Const(1),
		"a":    schema.Const(2),
	}, nil)
	fields2 := schema.FieldMap(schema.Fields{
		"type": schema.Const(3),
		"b":    schema.Const(4),
	}, nil)
	sch := schema.FieldMapSet("type", []schema.Checker{fields1, fields2})

	out, err := sch.Coerce(map[string]int{"type": 1, "a": 2}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"type": 1, "a": 2})

	out, err = sch.Coerce(map[string]int{"type": 3, "b": 4}, aPath)
	c.Assert(err, IsNil)
	c.Assert(out, DeepEquals, map[string]interface{}{"type": 3, "b": 4})

	out, err = sch.Coerce(map[string]int{}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.type: expected supported selector, got nothing`)

	out, err = sch.Coerce(map[string]int{"type": 2}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.type: expected supported selector, got 2`)

	out, err = sch.Coerce(map[string]int{"type": 3, "b": 5}, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>\.b: expected 4, got 5`)

	out, err = sch.Coerce(42, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected map, got 42`)

	out, err = sch.Coerce(nil, aPath)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `<path>: expected map, got nothing`)

	// First path entry shouldn't have dots in an error message.
	out, err = sch.Coerce(map[string]int{"a": 1}, nil)
	c.Assert(out, IsNil)
	c.Assert(err, ErrorMatches, `type: expected supported selector, got nothing`)
}
