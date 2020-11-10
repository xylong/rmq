package model

import (
	"fmt"
)

// Trans 转账
type Trans struct {
	ID    int    `json:"id"`
	From  string `json:"from" form:"from"`
	To    string `json:"to" form:"to"`
	Money int    `json:"money" form:"money"`
}

func NewTrans() *Trans {
	return &Trans{}
}

func (t *Trans) String() string {
	return fmt.Sprintf("%s转给%s：%d\n", t.From, t.To, t.Money)
}
