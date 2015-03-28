package kademlia

import (
	"github.com/"
	"encoding/hex"
	"container/list"
	"math/rand"
	"sort"
	"log"
	"os"
)

const (
	IdLength = 20
	BucketSize = 20
)

type NodeID [IdLength]byte

func NewNodeID(data string) (ret NodeID) {
	decoded, err := hex.DecodeString(data)
	checkError(err)
	for i:=0; i < IdLength; i++ {
		ret[i] = decoded[i]
	}
	return
}

func NewRandomNodeID() (ret NodeID) {
	for i:=0; i < IdLength; i++ {
		ret[i] = uint8(rand.Intn(256))
	}
	return
}

func (node NodeID) String() string {
	return hex.EncodeToString(node[0:IdLength])
}

func (node NodeID) Equals(rhs NodeID) bool {
	for i:=0; i < IdLength; i++ {
		if node[i] != rhs[i] {
			return false
		}
	}
	return true
}

func (node NodeID) Less(rhs interface{}) bool {
	for i:=0; i < IdLength; i++ {
		if node[i] != rhs.(NodeID)[i] {
			return node[i] < rhs.(NodeID)[i]
		}
	}
	return false
}

func (node NodeID) Xor(rhs NodeID) (ret NodeID) {
	for i:=0; i < IdLength; i++ {
		ret[i] = node[i] ^ rhs[i]
	}
	return
}

func (node NodeID) PrefixLen() (ret int) {
	for i:=0; i < IdLength; i++ {
		for j:=0; j < 8; j++ {
			if (node[i] >> uint8(7-j)) & 0x1 != 0 {
				return i*8+j
			}
		}
	}
	return IdLength*8-1
}

type Contact struct {
	id NodeID
}

type RoutingTable struct {
	node NodeID
	buckets [IdLength*8]*list.List
}

func NewRoutingTable(node NodeID)(ret RoutingTable) {
	for i:=0; i < IdLength*8; i++ {
		ret.buckets[i] = list.New()
	}
	ret.node = node
	return 
}

func (table *RoutingTable) Update(contact *Contact) {
	prefix_length := contact.id.Xor(table.node.id).PrefixLen()
	bucket := table.buckets[prefix_length]
	element := bucket.Front()
	for element != nil {
		if element.Value.(*Contact).id.Equals(contact.id) {
			break
		}
		element = element.Next()
	}
	if element == nil {
		if bucket.Len() < BucketSize {
			bucket.PushFront(contact)
		}
		else {
			//TODO: Handle insertion when the list is full by evicting old elements if
			//they don't respond to a ping.
		}
	}
	else {
		bucket.MoveToFront(element)
	}
}

type ContactRecord struct {
	node *Contact
	sortKey NodeID
}

func (rec *ContactRecord) Less(other interface{}) bool {
	return rec.sortKey.Less(other.(*ContactRecord).sortKey)
}

func copyToVector(start, end *list.Element, vec vector.Vector, target NodeID) {
	for elt:=start; elt != end; elt = elt.Next() {
		contact := elt.Value.(*Contact)
		vec.Push(&ContactRecord{contact, contact.id.Xor(target)})	
	}
}

func (table *RoutingTable) FindClosest(target NodeID, count int) {
	ret = vector.NewVector()
	ret.Resize(count)
	bucket_num := target.Xor(table.node).PrefixLen()
	bucket := table.buckets[bucket_num]
	copyToVector(bucket.Front(), nil, ret, target)
	
	for i:=1; (bucket_num-i >= 0 || bucket_num+i < IdLength*8) && ret.Len() < count; i++ {
		if bucket_num-i >= 0 {
			bucket = table.buckets[bucket_num-i]
			copyToVector(bucket.Front(), nil, ret, target)
		}
		if bucket_num+i < count {
			bucket = table.buckets[bucket_num+i]
			copyToVector(bucket.Front(), nil, ret, target)
		}
	}
	sort.Sort(ret)
	if ret.Len() > count {
		ret.Cut(count)
	}
	return
}

type Kademlia struct {
	routes *RoutingTable
	NetworkId string
}

func NewKademlia(this *Contact, networkId string) *Kademlia {
	return &Kademlia {
		routes: NewRoutingTable(this),
		NetworkId: networkId
	}
}

type RPCHeader struct {
	Sender *Contact
	NetworkId string
}
func checkError(err error) {
	if err != nil {
		log.Printf("Error is: %v\n", err)
		os.Exit(1)
	}
}