package alt

import (
	"fmt"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/timo972/altv-go-pkg/pb"
	"golang.org/x/exp/slices"
	"google.golang.org/protobuf/proto"
)

type Decoder struct {
	Buffer    []byte
	RootValue reflect.Value
	RootType  reflect.Type
	MValue    *pb.MValue
}

func NewDecoder(data []byte) *Decoder {
	return &Decoder{
		Buffer: data,
	}
}

func parsePointer(ptrStr string) (unsafe.Pointer, error) {
	ptrUint, err := strconv.ParseUint(ptrStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return unsafe.Pointer(uintptr(ptrUint)), nil
}

func (d *Decoder) Decode(v interface{}) error {
	if d.MValue == nil && len(d.Buffer) > 0 {
		err := proto.Unmarshal(d.Buffer, d.MValue)
		if err != nil {
			return err
		}
	} else if d.MValue == nil {
		return fmt.Errorf("no data to decode")
	}

	d.RootValue = reflect.ValueOf(v)
	d.RootType = d.RootValue.Type()

	return d.decode()
}

func (d *Decoder) decode() error {
	switch d.RootType.Kind() {
	case reflect.Ptr:
		rt := d.RootType.Elem()
		rv := d.RootValue.Elem()
		kind := rt.Kind()

		if kind == reflect.Struct {
			// struct pointer
			d.decodeStruct(rt, rv)
		} else if kind == reflect.Slice || kind == reflect.Array {
			// slice / array pointer
			d.decodeSlice(rt, rv)
		} else if kind == reflect.Map {
			// map pointer
			d.decodeMap(rt, rv)
		}
	case reflect.String:
		d.RootValue.SetString(d.MValue.GetStringValue())
	case reflect.Bool:
		d.RootValue.SetBool(d.MValue.GetBoolValue())
	case reflect.Float32, reflect.Float64:
		d.RootValue.SetFloat(d.MValue.GetDoubleValue())
	case reflect.Func:
		// function
		f := d.MValue.GetExternFunctionValue()
		ptr, err := parsePointer(f.GetPtr())
		if err != nil {
			return err
		}

		d.RootValue.Set(reflect.ValueOf(ExternFunction{
			Ptr: ptr,
		}))
	// integers
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d.RootValue.SetInt(d.MValue.GetIntValue())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d.RootValue.SetUint(d.MValue.GetUintValue())
	case reflect.Array, reflect.Slice:
		d.decodeSlice(d.RootType, d.RootValue)
	case reflect.Struct:
		// vector3, rgba, vector2
		d.decodeStruct(d.RootType, d.RootValue)
	case reflect.Map:
		// map
		d.decodeMap(d.RootType, d.RootValue)
	default:
		return fmt.Errorf("unsupported type: %s", d.RootType.Kind())
	}

	return nil
}

func (d *Decoder) decodeStruct(rt reflect.Type, rv reflect.Value) error {
	structName := rt.Name()

	if structName == "RGBA" {
		rgba := d.MValue.GetRgbaValue()

		rv.Set(reflect.ValueOf(RGBA{
			R: uint8(rgba.GetR()),
			G: uint8(rgba.GetG()),
			B: uint8(rgba.GetB()),
			A: uint8(rgba.GetA()),
		}))
	} else if structName == "Vector2" {
		v2 := d.MValue.GetVector2Value()

		rv.Set(reflect.ValueOf(Vector2{
			X: v2.GetX(),
			Y: v2.GetY(),
		}))
	} else if structName == "Vector3" {
		v3 := d.MValue.GetVector3Value()

		rv.Set(reflect.ValueOf(Vector3{
			X: v3.GetX(),
			Y: v3.GetY(),
			Z: v3.GetZ(),
		}))
	} else if structName == "Player" || structName == "Entity" || structName == "Vehicle" || structName == "ColShape" || structName == "Checkpoint" || structName == "VoiceChannel" || structName == "Blip" {
		base := d.MValue.GetBaseObjectValue()
		t := BaseObjectType(base.GetType())
		ptr, err := parsePointer(base.GetPtr())
		if err != nil {
			return err
		}

		var v reflect.Value

		switch t {
		case PlayerObject:
			if structName == "Entity" {
				e := &Entity{}
				e.Ptr = ptr
				e.Type = t
				v = reflect.ValueOf(e)
			} else {
				v = reflect.ValueOf(newPlayer(ptr))
			}
		case VehicleObject:
			if structName == "Entity" {
				e := &Entity{}
				e.Ptr = ptr
				e.Type = t
				v = reflect.ValueOf(e)
			} else {
				v = reflect.ValueOf(newVehicle(ptr))
			}
		case ColshapeObject:
			v = reflect.ValueOf(newColShape(ptr))
		case CheckpointObject:
			v = reflect.ValueOf(newCheckpoint(ptr))
		case VoiceChannelObject:
			v = reflect.ValueOf(newVoiceChannel(ptr))
		case BlipObject:
			v = reflect.ValueOf(newBlip(ptr))
		}

		rv.Set(v)
	} else {
		keys := d.MValue.GetDict()
		values := d.MValue.GetList()

		fieldCount := rv.NumField()

		for i := 0; i < fieldCount; i++ {
			field := rv.Field(i)
			fieldType := rt.Field(i)
			name := getFieldName(fieldType)

			i := slices.Index(keys, name)
			if i < 0 {
				// field is not set
				continue
			}

			valueDecoder := &Decoder{
				MValue:    values[i],
				RootValue: field,
				RootType:  field.Type(),
			}

			err := valueDecoder.decode()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (d *Decoder) decodeMap(rt reflect.Type, rv reflect.Value) error {
	m := reflect.MakeMap(rt)

	keys := d.MValue.GetDict()
	values := d.MValue.GetList()

	for i, key := range keys {
		valueDecoder := &Decoder{
			MValue: values[i],
			// FIXME: ???
			RootValue: reflect.New(rt.Elem()),
			RootType:  rt.Elem(),
		}

		err := valueDecoder.decode()
		if err != nil {
			return err
		}

		m.SetMapIndex(reflect.ValueOf(key), valueDecoder.RootValue)
	}

	rv.Set(m)

	return nil
}

func (d *Decoder) decodeSlice(rt reflect.Type, rv reflect.Value) error {
	values := d.MValue.GetList()
	size := len(values)
	l := reflect.MakeSlice(rt, size, size)
	elemType := rt.Elem()

	for i, pbVal := range values {
		valueDecoder := &Decoder{
			MValue:    pbVal,
			RootValue: l.Index(i),
			RootType:  elemType,
		}

		err := valueDecoder.decode()
		if err != nil {
			return err
		}
	}

	return nil
}
