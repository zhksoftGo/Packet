/*=============================================================================
*	Copyright (C) 2006-2021, Zhang Kun(zhk.tiger@gmail.com). All Rights Reserved.
*	Generated by: ProtocolGen 1.93 2021-6-26
=============================================================================*/
package Cactus

//打包解包类
import "github.com/zhksoftGo/Packet"

//整数数组
type VectorInt []int32

func (v *VectorInt) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteInt32((*v)[item])
	}
}

func (v *VectorInt) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t int32
		t = pak.ReadInt32()
		*v = append(*v, t)
	}
}

//Short数组
type VectorShort []int16

func (v *VectorShort) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteInt16((*v)[item])
	}
}

func (v *VectorShort) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t int16
		t = pak.ReadInt16()
		*v = append(*v, t)
	}
}

//int 64数组
type VectorInt64 []int64

func (v *VectorInt64) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteInt64((*v)[item])
	}
}

func (v *VectorInt64) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t int64
		t = pak.ReadInt64()
		*v = append(*v, t)
	}
}

//unsigned int 64数组
type VectorUint64 []uint64

func (v *VectorUint64) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteUint64((*v)[item])
	}
}

func (v *VectorUint64) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t uint64
		t = pak.ReadUint64()
		*v = append(*v, t)
	}
}

//Cactus::String数组
type VectorString []string

func (v *VectorString) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteString((*v)[item])
	}
}

func (v *VectorString) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t string
		t = pak.ReadString()
		*v = append(*v, t)
	}
}

//浮点数组
type VectorFloat []float32

func (v *VectorFloat) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		pak.WriteFloat32((*v)[item])
	}
}

func (v *VectorFloat) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t float32
		t = pak.ReadFloat32()
		*v = append(*v, t)
	}
}

//int-int map
type MapIntInt map[int32]int32

func (v *MapIntInt) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for key, item := range *v {
		pak.WriteInt32(key)
		pak.WriteInt32(item)
	}
}

func (v *MapIntInt) Read(pak *Packet.Packet) {
	*v = make(MapIntInt)
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var key int32
		key = pak.ReadInt32()
		var item int32
		item = pak.ReadInt32()
		(*v)[key] = item
	}
}

//int-bool map
type MapIntBool map[int32]bool

func (v *MapIntBool) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for key, item := range *v {
		pak.WriteInt32(key)
		pak.WriteBool(item)
	}
}

func (v *MapIntBool) Read(pak *Packet.Packet) {
	*v = make(MapIntBool)
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var key int32
		key = pak.ReadInt32()
		var item bool
		item = pak.ReadBool()
		(*v)[key] = item
	}
}

//String-String map
type MapStringString map[string]string

func (v *MapStringString) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for key, item := range *v {
		pak.WriteString(key)
		pak.WriteString(item)
	}
}

func (v *MapStringString) Read(pak *Packet.Packet) {
	*v = make(MapStringString)
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var key string
		key = pak.ReadString()
		var item string
		item = pak.ReadString()
		(*v)[key] = item
	}
}

//单条消息记录
type SMsgRecordItem struct {
	time int64 `xml:"-"    json:"-"`				//发生时间
	msg Packet.Packet `xml:"-"    json:"-"`				//去掉消息头(长度+类型)的消息数据
}

func (v *SMsgRecordItem) Write(pak *Packet.Packet) {
	pak.WriteInt64(v.time)
	pak.WritePacket(v.msg)
}

func (v *SMsgRecordItem) Read(pak *Packet.Packet) {
	v.time = pak.ReadInt64()
	v.msg = pak.ReadPacket()
}

//消息记录数组
type VectorMsgRecord []SMsgRecordItem

func (v *VectorMsgRecord) Write(pak *Packet.Packet) {
	l := uint32(len(*v))
	pak.WriteUint32(l)
	for item := range *v {
		(*v)[item].Write(pak)
	}
}

func (v *VectorMsgRecord) Read(pak *Packet.Packet) {
	l := pak.ReadUint32()
	var i uint32
	for i = 0; i < l; i++ {
		var t SMsgRecordItem
		t.Read(pak)
		*v = append(*v, t)
	}
}

