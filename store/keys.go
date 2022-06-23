package store

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"github.com/sharop/bone/common"
	"math"
	"strings"
)

func RKey(attribute string, resource Resource) string {
	// :resource:source:ventas
	// :
	//TODO: Parse path
	//Level Resources
	return NSAttr(Resources, resource, attribute)
}

func NSAttr(ns Level, resource Resource, attribute string) string {
	return common.UintToStr(uint64(ns)) + Separator + common.UintToStr(uint64(resource)) + Separator + attribute
}
func ParseNSBytes(attribute string) ([]byte, string) {
	splits := strings.SplitN(attribute, Separator, 3)
	ns := make([]byte, 8*2)

	for idx, x := range splits[:2] {

		binary.BigEndian.PutUint64(ns[((1<<idx)-1)*8:], common.StrToUint(x))
	}
	//0:7
	//8:15
	return ns, splits[2]

}

func ParseAttr(attr string) string {
	return strings.SplitN(attr, Separator, 3)[2]
}

// writeAttr validate attr and write in binary format
func writeAttr(buf []byte, attr string) []byte {
	AssertTrue(len(attr) < math.MaxUint16)
	binary.BigEndian.PutUint16(buf[:2], uint16(len(attr)))

	rest := buf[2:]
	AssertTrue(len(attr) == copy(rest, attr))

	return rest[len(attr):]
}

func generateKey(keyType byte, attribute string, reserv int) ([]byte, int) {

	namespace, attribute := ParseNSBytes(attribute)
	pLen := 1 + 16 + 2 + len(attribute) // (Level + Resource) + LenAttr  + len(attribute)
	buf := make([]byte, pLen+reserv)
	buf[0] = keyType
	AssertTrue(copy(buf[1:], namespace) == 16)
	rest := buf[17:]

	writeAttr(rest, attribute)
	return buf, pLen
}

// DataKey  generates a key used in data records
// --------------------------------------------------------------------------------------------------------------------
// |	BYTE 0		|	BYTE 1		|	BYTE 2-3	| 	BYTE[4:METALEN]	|	1 Byte		|	8 Bytes		|	8 Bytes   |
// |	LEVEL		|	RESOURCE	|	KNameLen	| 	KeyName			|	KeyType		|	 UID		|    Reserved |
// --------------------------------------------------------------------------------------------------------------------
// : Resource: Source:0-ventas:01
func DataKey(attr string, uid uint64) []byte {
	xRes := 1 + 8                                     // RecordType + UID
	buf, pLen := generateKey(DefaultKind, attr, xRes) // DefaultKind for KData, KIndex and KPredicate
	//buf[1] = byte(ns[0])
	rest := buf[pLen:]
	rest[0] = KData // StoreData

	rest = rest[1:]
	binary.BigEndian.PutUint64(rest, uid)
	return buf
}

// Parse would parse the key. ParsedKey does not reuse the key slice, so the key slice can change
// without affecting the contents of ParsedKey.
func Parse(key []byte) (Key, error) {
	var p Key

	//fmt.Printf("key:\n%s\n", hex.Dump(key))
	if len(key) < 9 {
		return p, errors.New("Key length less than 9")
	}
	p.bytePrefix = key[0]
	namespace := key[1:17]
	//fmt.Printf("namespace:\n%s\n", hex.Dump(namespace))

	key = key[17:]
	//fmt.Printf("key:\n%s\n", hex.Dump(key))

	if p.bytePrefix == ByteUnused {
		return p, nil
	}

	if len(key) < 3 {
		return p, errors.Errorf("Invalid format for key %v", key)
	}
	sz := int(binary.BigEndian.Uint16(key[:2]))
	k := key[2:]
	//fmt.Printf("k:\n%s\n", hex.Dump(k))

	if len(k) < sz {
		return p, errors.Errorf("Invalid size %v for key %v", sz, key)
	}

	//TODO: Change 0 to correct attribute
	p.Attr = NSAttr(Level(binary.BigEndian.Uint64(namespace)), 0, string(k[:sz]))
	k = k[sz:]
	//fmt.Printf("k:\n%s\n", hex.Dump(k))

	switch p.bytePrefix {
	case KIndex, KReverse:
		return p, nil
	default:
	}

	p.KeyType = k[0]
	k = k[1:]
	//fmt.Printf("k:\n%s\n", hex.Dump(k))

	switch p.KeyType {
	case KData:
		if len(k) < 8 {
			return p, errors.Errorf("uid length < 8 for key: %q, parsed key: %+v", key, p)
		}
		//TODO: Correct size of Key actual 11bytes
		p.UId = binary.BigEndian.Uint64(k)
		if p.UId == 0 {
			return p, errors.Errorf("Invalid UID with value 0 for key: %v", key)
		}

		k = k[8:]
		//fmt.Printf("k:\n%s\n", hex.Dump(k))

	default:
		// Some other data type.
		return p, errors.Errorf("Invalid data type")
	}
	return p, nil
}
