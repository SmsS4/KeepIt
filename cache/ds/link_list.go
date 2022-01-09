package ds

import (
	"fmt"
	"log"
	"sync"
)

type Node struct {
	value string
	next  *Node
	prev  *Node
}

func NewNode(value string) Node {
	return Node{value, nil, nil}
}

func getValue(node *Node) string {
	if node == nil {
		return "nil"
	}
	return node.value
}

func (node *Node) PrintNode() {
	log.Printf(
		"(%s) -> (%s) -> (%s)",
		getValue(node.prev),
		getValue(node),
		getValue(node.next),
	)
}

type LinkList struct {
	head *Node
	tail *Node
	Size int
	lock *sync.Mutex
}

func NewLinkList() LinkList {
	return LinkList{nil, nil, 0, new(sync.Mutex)}
}

func (ll *LinkList) remove(node *Node) {
	log.Printf("\tremove %s from tail", node.value)
	ll.Size -= 1
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}
	if ll.tail == node {
		ll.tail = node.prev
	}
	if ll.head == node {
		ll.head = node.next
	}
	node.prev = nil
	node.next = nil
}

func (ll *LinkList) MoveToTail(node *Node) {
	log.Printf("Move %s to tail", node.value)
	ll.lock.Lock()
	ll.remove(node)
	ll.append(node)
	ll.Print()
	ll.lock.Unlock()
}

func (ll *LinkList) PopHead() {
	ll.lock.Lock()
	log.Print("Remove head from link list\n")
	ll.remove(ll.head)
	ll.Print()
	ll.lock.Unlock()
}

func (ll *LinkList) append(node *Node) {
	log.Printf("\tappend %s to tail", node.value)
	ll.Size += 1
	if ll.head == nil {
		ll.head = node
		ll.tail = node
	} else {
		ll.Print()
		ll.tail.next = node
		node.prev = ll.tail
		ll.tail = node
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
	if ll.head == nil {
		log.Print("LL is empty\n")
	} else {
		log.Printf("LL size is %d\n", ll.Size)
		values := ""
		cur := ll.head
		for cur != nil {
			values += fmt.Sprintf("(%s)->", cur.value)
			cur = cur.next
		}
		log.Print(values)
	}
}
