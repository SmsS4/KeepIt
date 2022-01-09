package ds

import (
	"fmt"
	"log"
	"sync"
)

type Node struct {
	Value string
	Next  *Node
	Prev  *Node
}

func NewNode(value string) Node {
	return Node{value, nil, nil}
}

func getValue(node *Node) string {
	if node == nil {
		return "nil"
	}
	return node.Value
}

func (node *Node) PrintNode() {
	log.Printf(
		"(%s) -> (%s) -> (%s)",
		getValue(node.Prev),
		getValue(node),
		getValue(node.Next),
	)
}

type LinkList struct {
	Head *Node
	Tail *Node
	Size int
	lock *sync.Mutex
}

func NewLinkList() LinkList {
	return LinkList{nil, nil, 0, new(sync.Mutex)}
}

func (ll *LinkList) remove(node *Node) {
	log.Printf("\tremove %s from tail", node.Value)
	ll.Size -= 1
	if node.Prev != nil {
		node.Prev.Next = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	}
	if ll.Tail == node {
		ll.Tail = node.Prev
	}
	if ll.Head == node {
		ll.Head = node.Next
	}
	node.Prev = nil
	node.Next = nil
}

func (ll *LinkList) MoveToTail(node *Node) {
	log.Printf("Move %s to tail", node.Value)
	ll.lock.Lock()
	ll.remove(node)
	ll.append(node)
	ll.Print()
	ll.lock.Unlock()
}

func (ll *LinkList) PopHead() {
	ll.lock.Lock()
	log.Print("Remove head from link list\n")
	ll.remove(ll.Head)
	ll.Print()
	ll.lock.Unlock()
}

func (ll *LinkList) append(node *Node) {
	log.Printf("\tappend %s to tail", node.Value)
	ll.Size += 1
	if ll.Head == nil {
		ll.Head = node
		ll.Tail = node
	} else {
		ll.Print()
		ll.Tail.Next = node
		node.Prev = ll.Tail
		ll.Tail = node
	}
}

func (ll *LinkList) AppendValue(value string) *Node {
	log.Printf("Append %s to link list", value)
	node := NewNode(value)
	ll.lock.Lock()
	ll.append(&node)
	ll.Print()
	ll.lock.Unlock()
	return &node
}

func (ll *LinkList) Print() {
	if ll.Head == nil {
		log.Print("LL is empty\n")
	} else {
		log.Printf("LL size is %d\n", ll.Size)
		values := ""
		cur := ll.Head
		for cur != nil {
			values += fmt.Sprintf("(%s)->", cur.Value)
			cur = cur.Next
		}
		log.Print(values)
	}
}
