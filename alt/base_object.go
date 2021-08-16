package alt

//#ifndef _WIN32
//#include <stdlib.h>
//#endif
// #include "Module.h"
import "C"
import (
	"unsafe"
)

type BaseObjectType = uint8

const (
	PlayerObject BaseObjectType = iota
	VehicleObject
	BlipObject
	WebviewObject
	VoiceChannelObject
	ColshapeObject
	CheckpointObject
	WebSocketClientObject
	HttpClientObject
	AudioObject
	LocalPlayerObject
)

type BaseObject struct {
	Ptr  unsafe.Pointer
	Type BaseObjectType
}

type Base interface {
	HasMetaData(key string) bool
	GetMetaData(key string) interface{}
	SetMetaData(key string, value interface{})
	DeleteMetaData(key string)
}

func (b BaseObject) HasMetaData(key string) bool {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	if b.Type == PlayerObject {
		return int(C.player_has_meta_data(b.Ptr, cKey)) == 1
	} else if b.Type == VoiceChannelObject {
		return int(C.voice_channel_has_meta_data(b.Ptr, cKey)) == 1
	} else if b.Type == CheckpointObject {
		return int(C.checkpoint_has_meta_data(b.Ptr, cKey)) == 1
	} else if b.Type == ColshapeObject {
		return int(C.col_shape_has_meta_data(b.Ptr, cKey)) == 1
	} else if b.Type == VehicleObject {
		return int(C.vehicle_has_meta_data(b.Ptr, cKey)) == 1
	}

	return false
}

func (b BaseObject) GetMetaData(key string) interface{} {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var mValue *MValue

	if b.Type == PlayerObject {
		meta := C.player_get_meta_data(b.Ptr, cKey)
		mValue = &MValue{Ptr: meta.Ptr, Type: uint8(meta.Type), Value: nil}
	} else if b.Type == VoiceChannelObject {
		meta := C.voice_channel_get_meta_data(b.Ptr, cKey)
		mValue = &MValue{Ptr: meta.Ptr, Type: uint8(meta.Type), Value: nil}
	} else if b.Type == CheckpointObject {
		meta := C.checkpoint_get_meta_data(b.Ptr, cKey)
		mValue = &MValue{Ptr: meta.Ptr, Type: uint8(meta.Type), Value: nil}
	} else if b.Type == ColshapeObject {
		meta := C.col_shape_get_meta_data(b.Ptr, cKey)
		mValue = &MValue{Ptr: meta.Ptr, Type: uint8(meta.Type), Value: nil}
	} else if b.Type == VehicleObject {
		meta := C.vehicle_get_meta_data(b.Ptr, cKey)
		mValue = &MValue{Ptr: meta.Ptr, Type: uint8(meta.Type), Value: nil}
	}

	return mValue.GetValue()
}

func (b BaseObject) SetMetaData(key string, value interface{}) {
	mValue := CreateMValue(value)
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	if b.Type == PlayerObject {
		C.player_set_meta_data(b.Ptr, cKey, mValue.Ptr)
	} else if b.Type == VoiceChannelObject {
		C.voice_channel_set_meta_data(b.Ptr, cKey, mValue.Ptr)
	} else if b.Type == CheckpointObject {
		C.checkpoint_set_meta_data(b.Ptr, cKey, mValue.Ptr)
	} else if b.Type == ColshapeObject {
		C.col_shape_set_meta_data(b.Ptr, cKey, mValue.Ptr)
	} else if b.Type == VehicleObject {
		C.vehicle_set_meta_data(b.Ptr, cKey, mValue.Ptr)
	}
}

func (b BaseObject) DeleteMetaData(key string) {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	if b.Type == PlayerObject {
		C.player_delete_meta_data(b.Ptr, cKey)
	} else if b.Type == VoiceChannelObject {
		C.voice_channel_delete_meta_data(b.Ptr, cKey)
	} else if b.Type == CheckpointObject {
		C.checkpoint_delete_meta_data(b.Ptr, cKey)
	} else if b.Type == ColshapeObject {
		C.col_shape_delete_meta_data(b.Ptr, cKey)
	} else if b.Type == VehicleObject {
		C.vehicle_delete_meta_data(b.Ptr, cKey)
	}
}
