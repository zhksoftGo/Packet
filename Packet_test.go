package Packet

import (
	"fmt"
	"testing"
)

func TestPacket(t *testing.T) {

	str := "This is a test."
	s := make([]byte, 20, 20)
	copy(s, str)
	fmt.Printf("%q %v % v\n", s, len(s), cap(s))

	var pak Packet
	pak.WriteInt32(10000).WriteInt64(10101010).WriteString("This is a test.").WriteBool(true)
	fmt.Println("\npak:", pak)

	val := pak.ReadInt32()
	val2 := pak.ReadInt64()
	vs := pak.ReadString()
	vb := pak.ReadBool()
	fmt.Println(val, val2, vs, vb)
	fmt.Println("\npak after read:", pak)

	var pak2 Packet
	pak.SetReadPos(0)
	pak2.WriteInt32(123456).WritePacket(pak)
	fmt.Println("\npak2:", pak2)

	pak2.ReadInt32()
	pak2Out := pak2.ReadPacket()
	fmt.Println("\npak2Out:", pak2Out)

	val = pak2Out.ReadInt32()
	val2 = pak2Out.ReadInt64()
	vs = pak2Out.ReadString()
	vb = pak2Out.ReadBool()
	fmt.Println(val, val2, vs, vb)

	// Patch
	var pakSend Packet
	pakSend.WriteUint32(0) // length
	pakSend.WriteUint16(0) // type
	pakSend.WriteUint32(0)
	pakSend.WriteString("serverKey")
	l := pakSend.GetWritePos() - 6
	pakSend.PatchInt32(int32(l), 0)
	fmt.Println("\nPatch:", pakSend.ReadUint32(), pakSend.ReadUint16(), pakSend.ReadUint32(), pakSend.ReadString())

	// To base64 string
	pak2Out.SetReadPos(0)
	str = pak2Out.ToBase64String()
	fmt.Println("\nToBase64String:", pak2Out, str)

	// From base64 string
	var pak2FromBase64String Packet
	pak2FromBase64String.FromBase64String(str)
	fmt.Println("\nFromBase64String:", pak2FromBase64String, str)
	val = pak2FromBase64String.ReadInt32()
	val2 = pak2FromBase64String.ReadInt64()
	vs = pak2FromBase64String.ReadString()
	vb = pak2FromBase64String.ReadBool()
	fmt.Println(val, val2, vs, vb)

	strs := []string{
		"twAAAAcAAAB0dHR0NTU1IAAAAGM0Y2E0MjM4YTBiOTIzODIwZGNjNTA5YTZmNzU4NDliAAAAAAoAAAAxMjNAcXEuY29tBgAAADEyMzEyMxEAAAAwMDpFMDozQTo2ODowNDoyRQcAAABQQzo1LjcyIAAAAGRlYTM5Njg2N2Y2MDQ0MWYwOWE0ZjI5MGNjZTUwYmM3AAAAACAAAAAwNmY0NjZmZGQ2MTRlNWNiM2MxZmJlZGRhZjllMzAzOAAAAAA=",
		"twAAAAgAAAB0dHR0MTIzNCAAAABjNGNhNDIzOGEwYjkyMzgyMGRjYzUwOWE2Zjc1ODQ5YgAAAAAKAAAAMTIzQHFxLmNvbQUAAAAxMjMxMxEAAAAwMDpFMDozQTo2ODowNDoyRQcAAABQQzo1LjcyIAAAAGRlYTM5Njg2N2Y2MDQ0MWYwOWE0ZjI5MGNjZTUwYmM3AAAAACAAAAAwNmY0NjZmZGQ2MTRlNWNiM2MxZmJlZGRhZjllMzAzOAAAAAA=",
		"uAAAAAcAAAB0dHQzNDM0IAAAAGM0Y2E0MjM4YTBiOTIzODIwZGNjNTA5YTZmNzU4NDliAAAAAAoAAAAyMzJAcXEuY29tBwAAADEyMzExMTERAAAAMDA6RTA6M0E6Njg6MDQ6MkUHAAAAUEM6NS43MiAAAABkZWEzOTY4NjdmNjA0NDFmMDlhNGYyOTBjY2U1MGJjNwAAAAAgAAAAMDZmNDY2ZmRkNjE0ZTVjYjNjMWZiZWRkYWY5ZTMwMzgAAAAA",
		"uwAAAAYAAAB0dHQ5OTkgAAAAYzRjYTQyMzhhMGI5MjM4MjBkY2M1MDlhNmY3NTg0OWIAAAAADQAAADEyMzEyM0BxcS5jb20IAAAAMTIzMTIzMTIRAAAAMDA6RTA6M0E6Njg6MDQ6MkUHAAAAUEM6NS43MiAAAABkZWEzOTY4NjdmNjA0NDFmMDlhNGYyOTBjY2U1MGJjNwAAAAAgAAAAMDZmNDY2ZmRkNjE0ZTVjYjNjMWZiZWRkYWY5ZTMwMzgAAAAA",
		"ugAAAAgAAAB0dHQ3NTY1NiAAAABjNGNhNDIzOGEwYjkyMzgyMGRjYzUwOWE2Zjc1ODQ5YgAAAAAMAAAAMTMxMjNAcXEuY29tBgAAADEyMzEyMxEAAAAwMDpFMDozQTo2ODowNDoyRQcAAABQQzo1LjcyIAAAAGRlYTM5Njg2N2Y2MDQ0MWYwOWE0ZjI5MGNjZTUwYmM3AAAAACAAAAAwNmY0NjZmZGQ2MTRlNWNiM2MxZmJlZGRhZjllMzAzOAAAAAA=",
		"ugAAAAcAAAB0ZXN0OTk5IAAAAGM0Y2E0MjM4YTBiOTIzODIwZGNjNTA5YTZmNzU4NDliAAAAAA0AAAAxMjMxMjNAcXEuY29tBgAAADEyMzEyMxEAAAAwMDpFMDozQTo2ODowNDoyRQcAAABQQzo1LjcyIAAAAGRlYTM5Njg2N2Y2MDQ0MWYwOWE0ZjI5MGNjZTUwYmM3AAAAACAAAAAwNmY0NjZmZGQ2MTRlNWNiM2MxZmJlZGRhZjllMzAzOAAAAAA=",
	}

	for i := 0; i < len(strs); i++ {
		var pak3 Packet
		pak3.FromBase64String(strs[i])
		fmt.Println(i, pak3.ToBase64String() == strs[i])
	}
}
