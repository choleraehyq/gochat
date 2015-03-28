package vector

import (
	"github.com/choleraehyq/gochat/kademlia"
)

type Element struct {
	Value interface{}
}

type Vector struct {
	data []*Element
}

func NewVector() *Vector {
	return &Vector {
		data: make([]*Element, 0, 0)
	}
}

func (this *Vector) Len() int {
	return len(this.data)
}

func (this *Vector) Swap(i, j int) {
	this.data[i], this.data[j] = this.data[j], this.data[i]
}

func (this *Vector) Less(i, j int) {
	return this.data[i].(*kademlia.ContactRecord).Less(this.data[j].(*kademlia.ContactRecord))
}

func (this *Vector) Push(elem *kademlia.ContactRecord) {
	append(this.data, &Element{Value: elem})
}

func (this *Vector) Cut(start int) {
	if start < 0 {
		start = 0
	}
	this.data = this.data[0:start]
}

func (this *Vector) Resize(int size) {
	if size <= len(this.data) {
		return 
	}
	t := make([]*Element, len(this.data), size)
	copy(t, this.data)
	this.data = t
}