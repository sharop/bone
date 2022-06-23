package store

import (
	"github.com/google/uuid"
	"time"
)

type Level byte

const (
	Core Level = iota
	Engine
	Services
	Resources
)

type Resource byte

const (
	Source Resource = iota + 1
	Space
	Asset
	Team
	Project
)

const (
	DefaultKind = byte(0x0A)
	KData       = byte(0x0)
	KIndex      = byte(0x1)
	KReverse    = byte(0x2)
	KPredicate  = byte(0x3)
	ByteUnused  = byte(0xF) //Discard type

)

const (
	NSSeparator = "|"
	Separator   = "-"
)

type source int

// Kind of sources
const (
	DataBase source = iota
	File
	RestService
	Custom
)

type Key struct {
	KeyType      byte
	Attr         string
	UId          uint64
	Root         Level
	ResourceType Resource
	bytePrefix   byte
}

type Value struct {
	Value     *[]byte
	Version   uint16
	Created   uint64
	Modified  uint64
	Meta      []byte
	IsLast    bool
	Available bool
}

type Predicate struct {
	Value    []byte
	SUId     uint64
	DUId     uint64
	Created  uint64
	Modified uint64
}

func (k Key) IsData() bool {
	return k.KeyType == KData

}

type Record interface {
	SetID(uuid.UUID)
	GetID() uuid.UUID
	GetVersion() int64
	GetCreationDate() int64
	GetRelations() []uuid.UUID
}
type Tag struct {
	UId      uint64
	Name     string
	Created  uint64
	Modified uint64
}

type Item struct {
	ID       uuid.UUID
	Type     Resource
	Path     string `badgerholdIndex:"Path"`
	Key      string `badgerholdIndex:"Key"`
	Value    *string
	Created  time.Time
	Modified time.Time
	Meta     []byte
	Version  int64
}

/*

(ITEM) --[Relation]-->   (ITEM)


/System/services/
	efiscal
	DataWallet
/Core

/tenant/sources/Balance
/tenant/sources/Balance/tables
/tenant/sources/Balance/meta
/tenant/sources/Balance/Info



/Tenant/Sources/[... Balance]/Metada
/Tenant/Sources/[... Balance]/Doc
/Tenant/Sources/[... Balance]/API Info


*/
