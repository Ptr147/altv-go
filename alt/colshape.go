package alt

// #include <stdlib.h>
// #include "Module.h"
import "C"
import (
	"unsafe"

	"github.com/shockdev04/altv-go-pkg/internal/module"
)

type ColShape struct {
	WorldObject
}

func NewColShape(c unsafe.Pointer) *ColShape {
	colShape := &ColShape{}
	colShape.Ptr = c
	colShape.Type = ColshapeObject
	return colShape
}

func CreateColShapeCircle(x float32, y float32, radius float32) *ColShape {
	ptr := C.core_create_col_shape_circle(C.float(x), C.float(y), C.float(radius))
	return NewColShape(ptr)
}
func CreateColShapeCuboid(x1 float32, y1 float32, z1 float32, x2 float32, y2 float32, z2 float32) *ColShape {
	ptr := C.core_create_col_shape_cuboid(C.float(x1), C.float(y1), C.float(z1), C.float(x2), C.float(y2), C.float(z2))
	return NewColShape(ptr)
}
func CreateColShapeCylinder(x float32, y float32, z float32, radius float32, height float32) *ColShape {
	ptr := C.core_create_col_shape_cylinder(C.float(x), C.float(y), C.float(z),C.float(radius), C.float(height))
	return NewColShape(ptr)
}
func CreateColShapeRectangle(x1 float32, y1 float32, x2 float32, y2 float32) *ColShape {
	ptr := C.core_create_col_shape_rectangle(C.float(x1), C.float(y1), C.float(x2), C.float(y2))
	return NewColShape(ptr)
}
func CreateColShapeSphere(x float32, y float32, z float32, radius float32) *ColShape {
	ptr := C.core_create_col_shape_sphere(C.float(x), C.float(y), C.float(z), C.float(radius))
	return NewColShape(ptr)
}

func (c ColShape) IsPlayersOnly() bool {
	if c.Type == ColshapeObject {
		return int(C.col_shape_is_players_only(c.Ptr)) == 1
	} else if c.Type == CheckpointObject {
		return int(C.checkpoint_is_players_only(c.Ptr)) == 1
	}
	return false
}

func (c ColShape) SetPlayersOnly(state bool) {
	if c.Type == ColshapeObject {
		C.col_shape_set_players_only(c.Ptr, C.int(module.Bool2int(state)))
	} else if c.Type == CheckpointObject {
		C.checkpoint_set_players_only(c.Ptr, C.int(module.Bool2int(state)))
	}
}

func (c ColShape) IsPointIn(pos Position) bool {
	if c.Type == ColshapeObject {
		return int(C.col_shape_is_point_in(c.Ptr, C.float(pos.X), C.float(pos.Y), C.float(pos.Z))) == 1
	} else if c.Type == CheckpointObject {
		return int(C.checkpoint_is_point_in(c.Ptr, C.float(pos.X), C.float(pos.Y), C.float(pos.Z))) == 1
	}
	return false
}

func (c ColShape) IsEntityIn(entity *Entity) bool {
	if c.Type == ColshapeObject {
		return int(C.col_shape_is_entity_in(c.Ptr, entity.Ptr)) == 1
	} else if c.Type == CheckpointObject {
		return int(C.checkpoint_is_entity_in(c.Ptr, entity.Ptr)) == 1
	}
	return false
}

func (c ColShape) ColShapeType() int8 {
	if c.Type == ColshapeObject {
		return int8(C.col_shape_get_col_shape_type(c.Ptr))
	} else if c.Type == CheckpointObject {
		return int8(C.checkpoint_get_col_shape_type(c.Ptr))
	}
	return 0
}