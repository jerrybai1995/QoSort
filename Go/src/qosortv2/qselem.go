package qosortv2

type qselem interface {
    Less(interface{}) bool
    Copy() interface{}
}




