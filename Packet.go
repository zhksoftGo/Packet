package Packet

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
)

// M_Min_BufferSize is an initial allocation minimal capacity.
const M_Min_BufferSize = 64
const M_Packet_Inc_Size = 256

/*
	base_ptr	rd_ptr		   wr_ptr
	|			|			   |
	|___________|___remaining__|____________________
	|_________write_pos________|_______space_______|
	|______________alloc_size(Cap)_________________|
*/

type Packet struct {
	_buf []byte

	/* 写指针的位置, 偏移量, valid data range */
	_write_pos int

	/* 读指针的位置, 偏移量 */
	_read_pos int
}

// 实现Stringer接口
func (b Packet) String() string {
	return fmt.Sprintf("Read pos: %v, Write pos: %v, Cap: %v, Remaining: %v, Space: %v\n_buffer: %X",
		b._read_pos, b._write_pos,
		b.Cap(), b.Remaining(), b.Space(),
		b._buf)
}

func (b Packet) ToHexViewString() string {
	var hex, asc, out string

	for i := 0; i < len(b._buf); i++ {

		// new line
		if i%16 == 0 {
			if i != 0 {
				out += hex + "  " + asc + "\n"
			}

			hex = ""
			asc = ""

			hex += fmt.Sprintf("%08Xh: ", i/16)
		}

		var buf_format string
		if i == b._read_pos && b._read_pos == b._write_pos {
			buf_format = "=%02X"
		} else if i == b._read_pos {
			buf_format = ">%02X"
		} else if i == b._write_pos {
			buf_format = "<%02X"
		} else {
			buf_format = " %02X"
		}

		hex += fmt.Sprintf(buf_format, b._buf[i])

		// 把wr_ptr打印出来
		if i == b._write_pos-1 && len(b._buf) == b._write_pos {
			hex += "<"
		}

		if b._buf[i] >= 32 && b._buf[i] < 127 {
			asc += string(b._buf[i])
		} else {
			asc += "."
		}
	}

	l := len(hex)
	if l > 0 {
		if l < 16*3+11 {
			left := 16*3 + 11 - l%(16*3+11)
			for t := 0; t < left; t++ {
				hex += " "
			}
		}

		out += hex + "  " + asc + "\n"
	}

	return out
}

// 从Buffer构建
func (b *Packet) FromBuff(buf []byte) {
	if b._buf != nil {
		panic(errors.New("Packet: already has a buffer"))
	}

	b._buf = buf
	b._read_pos = 0
	b._write_pos = len(buf)
}

func (b *Packet) GetBuffer() []byte {
	return b._buf
}

func (b *Packet) GetUsedBuffer() []byte {
	return b._buf[:b._write_pos]
}

// 获取剩余可读的数据大小(写的 减 读的)
func (b *Packet) Remaining() int {
	return b._write_pos - b._read_pos
}

// 获取整个分配空间的大小
func (b *Packet) Cap() int {
	if b == nil {
		return 0
	}

	return cap(b._buf)
}

// 获取剩余可用(写)的空间大小(分配的 减 写的)
func (b *Packet) Space() int {
	if b == nil {
		return 0
	}

	return b.Cap() - b._write_pos
}

// Reset resets the _buffer to be empty, but it retains the underlying storage for use by future writes.
func (b *Packet) Reset() {
	b._buf = b._buf[:0]
	b._read_pos = 0
	b._write_pos = 0
}

// 丢弃读指针之前的数据
func (b *Packet) Truncate() {
	b._buf = b._buf[b._read_pos:]
	b._write_pos = b._write_pos - b._read_pos
	b._read_pos = 0
}

// 保证剩余空间大于n
func (b *Packet) GuaranteSpace(n int) bool {

	if b._buf == nil {
		if n <= M_Min_BufferSize {
			b._buf = make([]byte, M_Min_BufferSize)
			return true
		}
	} else if n <= b.Space() {
		return true
	}

	new_size := b.Cap() + n
	new_alloc_size := b.Cap()
	exp := 1
	for new_alloc_size < new_size {
		new_alloc_size += M_Packet_Inc_Size * exp
		exp <<= 1
	}

	// allocate
	_buf := make([]byte, new_alloc_size)
	copy(_buf, b._buf)
	b._buf = _buf

	return true
}

func (b *Packet) GetReadPos() int {
	return b._read_pos
}
func (b *Packet) OffsetReadPos(offset int) {
	if b._read_pos+offset > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	b._read_pos += offset
}
func (b *Packet) SetReadPos(pos int) {
	if pos > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	b._read_pos = pos
}

func (b *Packet) GetWritePos() int {
	return b._write_pos
}
func (b *Packet) OffWritePos(offset int) {
	if b._write_pos+offset >= b.Cap() {
		panic(errors.New("Packet: exceed capacity"))
	}

	b._write_pos += offset
}
func (b *Packet) SetWritePos(pos int) {
	if pos > b.Cap() {
		panic(errors.New("Packet: exceed capacity"))
	}

	b._write_pos = pos
}

func (b *Packet) PatchInto(p []byte, pos int) {
	l := len(p)
	if pos+l > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	tb := b._buf[pos:]
	copy(tb, p)
}

func (b *Packet) PatchInt32(v int32, pos int) {
	if pos < 0 || pos+4 > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	back_w_pos := b._write_pos
	b._write_pos = pos
	b.WriteInt32(v)
	b.SetWritePos(back_w_pos)
}

func (b *Packet) PeekOut(pos int, len int) []byte {
	if pos+len > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	tb := b._buf[pos : pos+len]
	return tb
}

func (b *Packet) PeekInt32(pos int) int32 {
	if pos < 0 || pos+4 > b._write_pos {
		panic(errors.New("Packet: exceed valid data range"))
	}

	back_r_pos := b._read_pos
	b._read_pos = pos
	v := b.ReadInt32()
	b.SetReadPos(back_r_pos)

	return v
}

// For io.Writer
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
func (b *Packet) Write(p []byte) (n int, err error) {

	l := len(p)
	if !b.GuaranteSpace(l) {
		panic(errors.New("Packet: not enough space"))
	}

	s := b._buf[b._write_pos:]
	copy(s, p)

	b._write_pos += l

	return l, nil
}

func (b *Packet) WriteInt8(v int8) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}
func (b *Packet) WriteUint8(v uint8) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteInt16(v int16) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}
func (b *Packet) WriteUint16(v uint16) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteInt32(v int32) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}
func (b *Packet) WriteUint32(v uint32) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteInt64(v int64) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}
func (b *Packet) WriteUint64(v uint64) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteBool(v bool) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteFloat32(v float32) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}
func (b *Packet) WriteFloat64(v float64) *Packet {
	binary.Write(b, binary.LittleEndian, v)
	return b
}

func (b *Packet) WriteString(v string) *Packet {
	b.WriteInt32(int32(len(v)))
	b.Write([]byte(v))

	return b
}

func (b *Packet) WritePacket(v Packet) *Packet {
	dataLen := v._write_pos
	b.WriteInt32(int32(dataLen))

	if dataLen != 0 {
		tb := v._buf[:v._write_pos]
		b.Write(tb)

		b.WriteInt32(int32(v._read_pos))
	}

	return b
}

func (b *Packet) ToBase64String() string {
	var pak Packet
	pak.WritePacket(*b)
	return base64.StdEncoding.EncodeToString(pak._buf[:pak._write_pos])
}

// For io.Reader
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
func (b *Packet) Read(p []byte) (n int, err error) {

	l := len(p)
	if l == 0 {
		return 0, nil
	}

	r := b.Remaining()
	if r < l {
		return 0, errors.New("Packet: not enough remaining data")
	}

	s := b._buf[b._read_pos:]
	n = copy(p, s)

	b._read_pos += n

	return n, nil
}

func (b *Packet) ReadInt8() int8 {
	var v int8
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}
func (b *Packet) ReadUint8() uint8 {
	var v uint8
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadInt16() int16 {
	var v int16
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}
func (b *Packet) ReadUint16() uint16 {
	var v uint16
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadInt32() int32 {
	var v int32
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}
func (b *Packet) ReadUint32() uint32 {
	var v uint32
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadInt64() int64 {
	var v int64
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}
func (b *Packet) ReadUint64() uint64 {
	var v uint64
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadBool() bool {
	var v bool
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadFloat32() float32 {
	var v float32
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}
func (b *Packet) ReadFloat64() float64 {
	var v float64
	err := binary.Read(b, binary.LittleEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

func (b *Packet) ReadString() string {

	l := b.ReadInt32()

	tb := make([]byte, l)
	sb := b._buf[b._read_pos:]
	copy(tb, sb)

	b._read_pos += int(l)

	return string(tb)
}

func (b *Packet) ReadPacket() Packet {

	var v Packet
	v._write_pos = int(b.ReadInt32())
	v._buf = make([]byte, v._write_pos)

	if v._write_pos != 0 {
		b.Read(v._buf)
		v._read_pos = int(b.ReadInt32())
	}

	return v
}

func (b *Packet) FromBase64String(s string) {

	buf, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	} else {
		var pak Packet
		pak.Write(buf)
		*b = pak.ReadPacket()
	}
}
