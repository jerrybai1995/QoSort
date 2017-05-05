package qosortv2

// implements qselem interface

type doublepair struct {
    x, y float64
}

func (p doublepair) Less(i interface{}) bool {
    p2 := i.(doublepair)
    return p.x < p2.x
}

func (p doublepair) Copy() interface{} {
    return doublepair{p.x, p.y}
}