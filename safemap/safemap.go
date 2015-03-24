package safemap

import (
	"sync"
	"errors"
	"reflect"
)

type Safemap {
	data map[interface{}]interface{}
	mtx sync.RWMutex
	keyType reflect.Type
	elemType reflect.Type
}

func Newsafemap(keyType, elemType reflect.Type) *Safemap {
	return &Safemap {
		data: make(map[interface{}]interface{})
		keyType: keyType
		elemType: elemType
	}
}

func (this *Safemap) Put(key interface{}, value interface{}) (interface{}, error){
	if !this.isAcceptablePair(key, value) {
		return nil, errors.New("Isn't a acceptable key-value pair")
	}
	this.mtx.Lock()
	defer this.mtx.Unlock()
	oldValue := this.data[key]
	this.data[key] = value
	return oldValue, nil
}

func (this *Safemap) Get(key interface{}) interface{} {
	this.mtx.RLock()
	defer this.mtx.RUnlock()
	return this.data[key]
}

func (this *Safemap) Remove(key interface{}) interface{} {
	this.mtx.Lock()
	defer this.mtx.Unlock()
	oldElem := this.data[key]
	delete(this.data, key)
	return oldElem
}


func (this *Safemap) Items() *map[interface{}]interface{} {
	return &(this.data)
}

func (this *Safemap) isAcceptablePair(key, value interface{}) bool {
	if key == nil || reflect.TypeOf(key) != this.keyType {
		return false
	}
	if value == nil || reflect.TypeOf(value) != this.elemType {
		return false
	}
	return true
}