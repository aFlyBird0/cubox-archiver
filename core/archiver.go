package core

// Archiver 归档器，能够提供原有的数据的Keys，并且能够操作数据
type Archiver interface {
	Operator
	KeysInitiator
}
