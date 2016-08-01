// Copyright (c) 2014-2016 Bitmark Inc.
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package avl

import ()

// delete: tree balancer
func balanceLeft(pp **node) bool {
	h := true
	p := *pp
	// h; left branch has shrunk
	if -1 == p.balance {
		p.balance = 0
	} else if 0 == p.balance {
		p.balance = 1
		h = false
	} else { // balance = 1, rebalance
		p1 := p.right
		if p1.balance >= 0 { // single rr rotation
			p.right = p1.left
			p1.left = p
			if 0 == p1.balance {
				p.balance = 1
				p1.balance = -1
				h = false
			} else {
				p.balance = 0
				p1.balance = 0
			}

			p1.up = p.up
			p.up = p1
			if nil != p.right {
				p.right.up = p
			}

			*pp = p1
		} else { // double rl rotation
			p2 := p1.left
			p1.left = p2.right
			p2.right = p1
			p.right = p2.left
			p2.left = p
			if +1 == p2.balance {
				p.balance = -1
			} else {
				p.balance = 0
			}
			if -1 == p2.balance {
				p1.balance = 1
			} else {
				p1.balance = 0
			}
			p2.balance = 0

			p2.up = p.up
			if nil != p.right {
				p.right.up = p
			}
			if nil != p1.left {
				p1.left.up = p1
			}
			p.up = p2
			p1.up = p2

			*pp = p2
		}
	}
	return h
}

// delete: tree balancer
func balanceRight(pp **node) bool {
	h := true
	p := *pp
	// h; right branch has shrunk
	if 1 == p.balance {
		p.balance = 0
	} else if 0 == p.balance {
		p.balance = -1
		h = false
	} else { // balance = -1, rebalance
		p1 := p.left
		if p1.balance <= 0 { // single ll rotation
			p.left = p1.right
			p1.right = p
			if 0 == p1.balance {
				p.balance = -1
				p1.balance = 1
				h = false
			} else {
				p.balance = 0
				p1.balance = 0
			}

			p1.up = p.up
			p.up = p1
			if nil != p.left {
				p.left.up = p
			}

			*pp = p1
		} else { // double lr rotation
			p2 := p1.right
			p1.right = p2.left
			p2.left = p1
			p.left = p2.right
			p2.right = p
			if -1 == p2.balance {
				p.balance = 1
			} else {
				p.balance = 0
			}
			if +1 == p2.balance {
				p1.balance = -1
			} else {
				p1.balance = 0
			}
			p2.balance = 0

			p2.up = p.up
			if nil != p.left {
				p.left.up = p
			}
			if nil != p1.right {
				p1.right.up = p1
			}
			p.up = p2
			p1.up = p2

			*pp = p2
		}
	}
	return h
}

// delete: rearrange deleted node
func del(qq **node, rr **node) bool {
	h := false
	if nil != (*rr).right {
		h = del(qq, &(*rr).right)
		if h {
			h = balanceRight(rr)
		}
	} else {
		q := *qq
		r := *rr
		rl := r.left
		if nil != rl {
			rl.up = r.up
		}

		if r != q.left {
			r.left = q.left
		}
		r.right = q.right
		r.up = q.up
		r.balance = q.balance

		if nil != r.right {
			r.right.up = r
		}
		if nil != r.left {
			r.left.up = r
		}

		*qq = r
		*rr = rl

		h = true
	}
	return h
}

// delete a specific item
func (tree *Tree) Delete(key item) interface{} {
	value, removed, _ := delete(key, &tree.root)
	if removed {
		tree.count -= 1
	}
	return value
}

// internal delete routine
func delete(key item, pp **node) (interface{}, bool, bool) {
	h := false
	if nil == *pp { // key not in tree
		return nil, false, h
	}
	value := interface{}(nil)
	removed := false
	switch (*pp).key.Compare(key) {
	case +1: // (*pp).key > key
		value, removed, h = delete(key, &(*pp).left)
		if h {
			h = balanceLeft(pp)
		}
	case -1: // (*pp).key < key
		value, removed, h = delete(key, &(*pp).right)
		if h {
			h = balanceRight(pp)
		}
	default: // found: delete p
		q := *pp
		value = q.value // preserve the value part
		if nil == q.right {
			if nil != q.left {
				q.left.up = q.up
			}
			*pp = q.left
			h = true
		} else if nil == q.left {
			if nil != q.right {
				q.right.up = q.up
			}
			*pp = q.right
			h = true
		} else {
			h = del(pp, &q.left)
			(*pp).left = q.left // p has changed, but q.left has left link value
			if h {
				h = balanceLeft(pp)
			}
		}
		freeNode(q)    // return deleted node to pool
		removed = true // indicate that an item was removed
	}
	return value, removed, h
}