package gofun

import (
	"fmt"
	"go/types"
	"reflect"

	"git.woa.com/vasd_masc_ba/YitihuaOteam/base/gofun/internal/importer"
	"git.woa.com/vasd_masc_ba/YitihuaOteam/base/gofun/internal/value"
)

// builtinTypes golang的内置类型
var builtinTypes = map[types.BasicKind]reflect.Type{
	types.Bool:       reflect.TypeOf(true),
	types.Int:        reflect.TypeOf(int(0)),
	types.Int8:       reflect.TypeOf(int8(0)),
	types.Int16:      reflect.TypeOf(int16(0)),
	types.Int32:      reflect.TypeOf(int32(0)),
	types.Int64:      reflect.TypeOf(int64(0)),
	types.Uint:       reflect.TypeOf(uint(0)),
	types.Uint8:      reflect.TypeOf(uint8(0)),
	types.Uint16:     reflect.TypeOf(uint16(0)),
	types.Uint32:     reflect.TypeOf(uint32(0)),
	types.Uint64:     reflect.TypeOf(uint64(0)),
	types.Uintptr:    reflect.TypeOf(uintptr(0)),
	types.Float32:    reflect.TypeOf(float32(0)),
	types.Float64:    reflect.TypeOf(float64(0)),
	types.Complex64:  reflect.TypeOf(complex64(0)),
	types.Complex128: reflect.TypeOf(complex128(0)),
	types.String:     reflect.TypeOf(""),

	types.UntypedBool:    reflect.TypeOf(true),
	types.UntypedInt:     reflect.TypeOf(int(0)),
	types.UntypedRune:    reflect.TypeOf(rune(0)),
	types.UntypedFloat:   reflect.TypeOf(float64(0)),
	types.UntypedComplex: reflect.TypeOf(complex128(0)),
	types.UntypedString:  reflect.TypeOf(""),
}

// typeChange 将types.Type转换为对应的reflect.Type类型
func typeChange(typ types.Type) reflect.Type {
	rType := importer.GetExternalType(typ)
	if rType != nil {
		return rType
	}
	switch t := typ.Underlying().(type) {
	case *types.Array:
		rType = reflect.ArrayOf(int(t.Len()), typeChange(t.Elem()))
	case *types.Basic:
		rtype := builtinTypes[t.Kind()]
		if rtype == nil {
			panic(t.Kind())
		}
		rType = rtype
	case *types.Chan:
		var dir reflect.ChanDir
		switch t.Dir() {
		case types.RecvOnly:
			dir = reflect.RecvDir
		case types.SendOnly:
			dir = reflect.SendDir
		case types.SendRecv:
			dir = reflect.BothDir
		}
		rType = reflect.ChanOf(dir, typeChange(t.Elem()))
	case *types.Interface:
		rType = reflect.TypeOf(func(interface{}) {}).In(0)
	case *types.Map:
		rType = reflect.MapOf(typeChange(t.Key()), typeChange(t.Elem()))
	case *types.Pointer:
		rType = reflect.PtrTo(typeChange(t.Elem()))
	case *types.Slice:
		rType = reflect.SliceOf(typeChange(t.Elem()))
	case *types.Struct:
		fields := make([]reflect.StructField, t.NumFields())
		for i := range fields {
			field := t.Field(i)
			fmt.Println("PkgPath: ", field)
			fields[i] = reflect.StructField{
				Name:      field.Name(),
				Type:      typeChange(t.Field(i).Type()),
				Tag:       reflect.StructTag(t.Tag(i)),
				Offset:    0,
				Index:     []int{i},
				Anonymous: field.Anonymous(),
			}
		}
		rType = reflect.StructOf(fields)
	default:
		rType = reflect.TypeOf(func(interface{}) {}).In(0)
	}
	return rType
}

// conv 将变量v的类型转换为typ
func conv(v interface{}, typ types.Type) value.Value {
	var reflectValue reflect.Value
	rtype := typeChange(typ)
	if v == nil {
		return value.RValue{Value: reflect.Zero(rtype)}
	}
	reflectValue = reflect.ValueOf(v).Convert(rtype)
	return value.RValue{Value: reflectValue}
}
