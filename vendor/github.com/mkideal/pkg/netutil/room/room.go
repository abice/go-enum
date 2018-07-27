package room

import (
	"errors"
)

var (
	ErrOutOfRange = errors.New("out of range")
)

type Role interface {
	Notify(data []byte) error
}

type emptyRole struct{}

func (role emptyRole) Notify(data []byte) error { return nil }

func IsEmptyRole(role Role) bool {
	_, ok := role.(emptyRole)
	return ok
}

type Room interface {
	Size() int
	GetRole(index int) Role
	AddRole(index int, role Role) error
	RemoveRole(index int) Role
	Notify(receivers []int, data []byte)
}

func GetPrevRoleIndex(room Room, index int) int {
	length := room.Size()
	return (index + length - 1) % length
}

func GetNextRoleIndex(room Room, index int) int {
	return (index + 1) % room.Size()
}

func Broadcast(room Room, data []byte) {
	room.Notify(nil, data)
}

// roomBase is base type of FixedSizeRoom and VarSizeRoom
type roomBase struct {
	roles []Role
}

func (room *roomBase) Size() int {
	return len(room.roles)
}

func (room *roomBase) GetRole(index int) Role {
	return room.roles[index]
}

func (room *roomBase) Notify(receivers []int, data []byte) {
	if len(receivers) == 0 {
		for _, role := range room.roles {
			role.Notify(data)
		}
	} else {
		for _, index := range receivers {
			role := room.GetRole(index)
			if role != nil {
				role.Notify(data)
			}
		}
	}
}

// FixedSizeRoom reprensents a fixed role size room
type FixedSizeRoom struct {
	roomBase
}

func NewFixedSizeRoom(size int) Room {
	room := &FixedSizeRoom{
		roomBase: roomBase{
			roles: make([]Role, 0, size),
		},
	}
	for i := 0; i < size; i++ {
		room.roles = append(room.roles, emptyRole{})
	}
	return room
}

func (room *FixedSizeRoom) AddRole(index int, role Role) error {
	if index >= len(room.roles) {
		return ErrOutOfRange
	}
	room.roles[index] = role
	return nil
}

func (room *FixedSizeRoom) RemoveRole(index int) Role {
	if index >= len(room.roles) {
		return nil
	}
	role := room.roles[index]
	room.roles[index] = emptyRole{}
	return role
}

// VarSizeRoom reprensents a variant role size room
type VarSizeRoom struct {
	roomBase
}

func NewVarSizeRoom(roles ...Role) Room {
	room := &VarSizeRoom{
		roomBase: roomBase{
			roles: roles,
		},
	}
	return room
}

func (room *VarSizeRoom) AddRole(index int, role Role) error {
	length := len(room.roles)
	if index > length {
		return ErrOutOfRange
	}
	if index == length {
		room.roles = append(room.roles, role)
	} else {
		room.roles = append(room.roles, nil)
		copy(room.roles[index+1:], room.roles[index:length])
		room.roles[index] = role
	}
	return nil
}

func (room *VarSizeRoom) RemoveRole(index int) Role {
	length := len(room.roles)
	if index >= len(room.roles) {
		return nil
	}
	role := room.roles[index]
	copy(room.roles[index:length-1], room.roles[index+1:])
	room.roles[length-1] = nil
	room.roles = room.roles[:length-1]
	return role
}
