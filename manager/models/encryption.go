package models

type EncryptionKey struct {
	BaseModel
	Key []byte `gorm:"type:binary(32);not null;unique" json:"key"`
}
