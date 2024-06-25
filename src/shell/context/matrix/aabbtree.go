/*
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package matrix

import (
	"container/list"
	"math"
)

type IAABB interface {
	GetAABB() *AABB
}

type AABB struct {
	minX        float64
	minY        float64
	minZ        float64
	maxX        float64
	maxY        float64
	maxZ        float64
	surfaceArea float64
}

func (a *AABB) calculateSurfaceArea() float64 {
	s := 2.0 * (a.getWidth()*a.getHeight() + a.getWidth()*a.getDepth() + a.getHeight()*a.getDepth())
	return s
}

func NewAABB(minX float64, minY float64, minZ float64, maxX float64, maxY float64, maxZ float64) *AABB {
	a := &AABB{
		minX: minX,
		minY: minY,
		minZ: minZ,
		maxX: maxX,
		maxY: maxY,
		maxZ: maxZ,
	}

	a.surfaceArea = a.calculateSurfaceArea()

	return a
}

func (a *AABB) overlaps(other *AABB) bool {
	// y is deliberately first in the list of checks below as it is seen as more likely than things
	// collide on x,z but not on y than they do on y thus we drop out sooner on a y fail
	return a.maxX > other.minX &&
		a.minX < other.maxX &&
		a.maxY > other.minY &&
		a.minY < other.maxY &&
		a.maxZ > other.minZ &&
		a.minZ < other.maxZ
}

func (a *AABB) contains(other *AABB) bool {
	return other.minX >= a.minX &&
		other.maxX <= a.maxX &&
		other.minY >= a.minY &&
		other.maxY <= a.maxY &&
		other.minZ >= a.minZ &&
		other.maxZ <= a.maxZ
}

func (a *AABB) merge(other *AABB) *AABB {
	b := NewAABB(math.Min(a.minX, other.minX), math.Min(a.minY, other.minY), math.Min(a.minZ, other.minZ),
		math.Max(a.maxX, other.maxX), math.Max(a.maxY, other.maxY), math.Max(a.maxZ, other.maxZ))
	return b
}

func (a *AABB) intersection(other *AABB) *AABB {
	b := NewAABB(math.Max(a.minX, other.minX), math.Max(a.minY, other.minY), math.Max(a.minZ, other.minZ),
		math.Min(a.maxX, other.maxX), math.Min(a.maxY, other.maxY), math.Min(a.maxZ, other.maxZ))
	return b
}

func (a *AABB) getWidth() float64 {
	return a.maxX - a.minX
}

func (a *AABB) getHeight() float64 {
	return a.maxY - a.minY
}

func (a *AABB) getDepth() float64 {
	return a.maxZ - a.minZ
}

const AABBNullNode = 0xffffffff

type AABBNode struct {
	aabb            *AABB
	object          IAABB
	parentNodeIndex uint
	leftNodeIndex   uint
	rightNodeIndex  uint
	nextNodeIndex   uint
}

func NewAABBNode() *AABBNode {
	node := &AABBNode{
		aabb:            &AABB{},
		object:          nil,
		parentNodeIndex: AABBNullNode,
		leftNodeIndex:   AABBNullNode,
		rightNodeIndex:  AABBNullNode,
		nextNodeIndex:   AABBNullNode,
	}
	return node
}

func (a *AABBNode) isLeaf() bool {
	return a.leftNodeIndex == AABBNullNode
}

type AABBTree struct {
	objectNodeIndexMap map[IAABB]uint
	nodes              []*AABBNode
	rootNodeIndex      uint
	allocatedNodeCount uint
	nextFreeNodeIndex  uint
	nodeCapacity       uint
	growthSize         uint
}

func NewAABBTree(initialSize uint) *AABBTree {
	t := &AABBTree{
		rootNodeIndex:      AABBNullNode,
		allocatedNodeCount: 0,
		nextFreeNodeIndex:  0,
		nodeCapacity:       initialSize,
		growthSize:         initialSize,
		nodes:              make([]*AABBNode, initialSize),
		objectNodeIndexMap: make(map[IAABB]uint),
		///nodes.resize(initialSize)
	}
	var nodeIndex uint

	for nodeIndex = 0; nodeIndex < initialSize; nodeIndex++ {
		node := NewAABBNode()
		t.nodes[nodeIndex] = node
		node.nextNodeIndex = nodeIndex + 1
	}
	t.nodes[initialSize-1].nextNodeIndex = AABBNullNode

	return t
}

func (a *AABBTree) allocateNode() (uint, *AABBNode) {
	if a.nextFreeNodeIndex == AABBNullNode {
		//assert(a.allocatedNodeCount == a.nodeCapacity)
		a.nodeCapacity += a.growthSize

		nodes := make([]*AABBNode, a.nodeCapacity)
		copy(nodes, a.nodes)
		a.nodes = nodes

		for nodeIndex := a.allocatedNodeCount; nodeIndex < a.nodeCapacity; nodeIndex++ {
			node := NewAABBNode()
			a.nodes[nodeIndex] = node
			node.nextNodeIndex = nodeIndex + 1
		}
		a.nodes[a.nodeCapacity-1].nextNodeIndex = AABBNullNode
		a.nextFreeNodeIndex = a.allocatedNodeCount
	}

	nodeIndex := a.nextFreeNodeIndex
	allocatedNode := a.nodes[nodeIndex]
	allocatedNode.parentNodeIndex = AABBNullNode
	allocatedNode.leftNodeIndex = AABBNullNode
	allocatedNode.rightNodeIndex = AABBNullNode
	a.nextFreeNodeIndex = allocatedNode.nextNodeIndex
	a.allocatedNodeCount++

	return nodeIndex, allocatedNode
}

func (a *AABBTree) deallocateNode(nodeIndex uint) {
	if len(a.nodes) == 0 {
		return
	}

	if nodeIndex >= 0 && nodeIndex < uint(len(a.nodes)) {
		deallocatedNode := a.nodes[nodeIndex]
		deallocatedNode.nextNodeIndex = a.nextFreeNodeIndex
		a.nextFreeNodeIndex = nodeIndex
		a.allocatedNodeCount--
	}
}

func (a *AABBTree) InsertObject(object IAABB) {
	nodeIndex, node := a.allocateNode()
	node.object = object
	node.aabb = object.GetAABB()

	a.insertLeaf(nodeIndex)
	a.objectNodeIndexMap[object] = nodeIndex
}

func (a *AABBTree) RemoveObject(object IAABB) {
	if nodeIndex, ok := a.objectNodeIndexMap[object]; ok {
		a.removeLeaf(nodeIndex)
		a.deallocateNode(nodeIndex)
		delete(a.objectNodeIndexMap, object)
	}
}

func (a *AABBTree) UpdateObject(object IAABB) {
	if nodeIndex, ok := a.objectNodeIndexMap[object]; ok {
		a.updateLeaf(nodeIndex, object.GetAABB())
	}
}

func (a *AABBTree) QueryOverlaps(object IAABB) []IAABB {
	var overlaps []IAABB

	stack := list.New()
	testAabb := object.GetAABB()

	stack.PushBack(a.rootNodeIndex)

	for stack.Len() > 0 {
		nodeIndex := stack.Back().Value.(uint)

		stack.Remove(stack.Back())

		if nodeIndex == AABBNullNode {
			continue
		}

		node := a.nodes[nodeIndex]
		if node.aabb.overlaps(testAabb) {
			if node.isLeaf() && node.object != object {
				overlaps = append(overlaps, node.object)
			} else {
				stack.PushBack(node.leftNodeIndex)
				stack.PushBack(node.rightNodeIndex)
			}
		}
	}

	return overlaps
}

func (a *AABBTree) insertLeaf(leafNodeIndex uint) {
	// make sure we're inserting a new leaf
	//assert(a.nodes[leafNodeIndex].parentNodeIndex == AABBNullNode)
	//assert(a.nodes[leafNodeIndex].leftNodeIndex == AABBNullNode)
	//assert(a.nodes[leafNodeIndex].rightNodeIndex == AABBNullNode)

	// if the tree is empty then we make the root the leaf
	if a.rootNodeIndex == AABBNullNode {
		a.rootNodeIndex = leafNodeIndex
		return
	}

	// search for the best place to put the new leaf in the tree
	// we use surface area and depth as search heuristics
	treeNodeIndex := a.rootNodeIndex
	leafNode := a.nodes[leafNodeIndex]

	for !a.nodes[treeNodeIndex].isLeaf() {
		//while !a.nodes[treeNodeIndex].isLeaf() {
		// because of the test in the while loop above we know we are never a leaf inside it
		treeNode := a.nodes[treeNodeIndex]
		leftNodeIndex := treeNode.leftNodeIndex
		rightNodeIndex := treeNode.rightNodeIndex
		leftNode := a.nodes[leftNodeIndex]
		rightNode := a.nodes[rightNodeIndex]

		combinedAabb := treeNode.aabb.merge(leafNode.aabb)

		newParentNodeCost := 2.0 * combinedAabb.surfaceArea
		minimumPushDownCost := 2.0 * (combinedAabb.surfaceArea - treeNode.aabb.surfaceArea)

		// use the costs to figure out whether to create a new parent here or descend
		var costLeft float64
		var costRight float64
		if leftNode.isLeaf() {
			costLeft = leafNode.aabb.merge(leftNode.aabb).surfaceArea + minimumPushDownCost
		} else {
			newLeftAabb := leafNode.aabb.merge(leftNode.aabb)
			costLeft = (newLeftAabb.surfaceArea - leftNode.aabb.surfaceArea) + minimumPushDownCost
		}
		if rightNode.isLeaf() {
			costRight = leafNode.aabb.merge(rightNode.aabb).surfaceArea + minimumPushDownCost
		} else {
			newRightAabb := leafNode.aabb.merge(rightNode.aabb)
			costRight = (newRightAabb.surfaceArea - rightNode.aabb.surfaceArea) + minimumPushDownCost
		}

		// if the cost of creating a new parent node here is less than descending in either direction then
		// we know we need to create a new parent node, here and attach the leaf to that
		if newParentNodeCost < costLeft && newParentNodeCost < costRight {
			break
		}

		// otherwise descend in the cheapest direction
		if costLeft < costRight {
			treeNodeIndex = leftNodeIndex
		} else {
			treeNodeIndex = rightNodeIndex
		}
	}

	// the leafs sibling is going to be the node we found above and we are going to create a new
	// parent node and attach the leaf and this item
	leafSiblingIndex := treeNodeIndex
	leafSibling := a.nodes[leafSiblingIndex]
	oldParentIndex := leafSibling.parentNodeIndex

	newParentIndex, newParent := a.allocateNode()
	newParent.parentNodeIndex = oldParentIndex
	newParent.aabb = leafNode.aabb.merge(leafSibling.aabb) // the new parents aabb is the leaf aabb combined with it's siblings aabb
	newParent.leftNodeIndex = leafSiblingIndex
	newParent.rightNodeIndex = leafNodeIndex

	leafNode.parentNodeIndex = newParentIndex
	leafSibling.parentNodeIndex = newParentIndex

	if oldParentIndex == AABBNullNode {
		// the old parent was the root and so this is now the root
		a.rootNodeIndex = newParentIndex
	} else {
		// the old parent was not the root and so we need to patch the left or right index to
		// point to the new node
		oldParent := a.nodes[oldParentIndex]
		if oldParent.leftNodeIndex == leafSiblingIndex {
			oldParent.leftNodeIndex = newParentIndex
		} else {
			oldParent.rightNodeIndex = newParentIndex
		}
	}

	// finally we need to walk back up the tree fixing heights and areas
	treeNodeIndex = leafNode.parentNodeIndex
	a.fixUpwardsTree(treeNodeIndex)
}

func (a *AABBTree) removeLeaf(leafNodeIndex uint) {
	// if the leaf is the root then we can just clear the root pointer and return
	if leafNodeIndex == a.rootNodeIndex {
		a.rootNodeIndex = AABBNullNode
		return
	}

	leafNode := a.nodes[leafNodeIndex]
	parentNodeIndex := leafNode.parentNodeIndex
	parentNode := a.nodes[parentNodeIndex]
	grandParentNodeIndex := parentNode.parentNodeIndex
	var siblingNodeIndex uint
	if parentNode.leftNodeIndex == leafNodeIndex {
		siblingNodeIndex = parentNode.rightNodeIndex
	} else {
		siblingNodeIndex = parentNode.leftNodeIndex
	}
	//parentNode.leftNodeIndex == leafNodeIndex ? parentNode.rightNodeIndex : parentNode.leftNodeIndex
	//assert(siblingNodeIndex != AABBNullNode) // we must have a sibling
	siblingNode := a.nodes[siblingNodeIndex]

	if grandParentNodeIndex != AABBNullNode {
		// if we have a grand parent (i.e. the parent is not the root) then destroy the parent and connect the sibling to the grandparent in its
		// place
		grandParentNode := a.nodes[grandParentNodeIndex]
		if grandParentNode.leftNodeIndex == parentNodeIndex {
			grandParentNode.leftNodeIndex = siblingNodeIndex
		} else {
			grandParentNode.rightNodeIndex = siblingNodeIndex
		}
		siblingNode.parentNodeIndex = grandParentNodeIndex
		a.deallocateNode(parentNodeIndex)

		a.fixUpwardsTree(grandParentNodeIndex)
	} else {
		// if we have no grandparent then the parent is the root and so our sibling becomes the root and has it's parent removed
		a.rootNodeIndex = siblingNodeIndex
		siblingNode.parentNodeIndex = AABBNullNode
		a.deallocateNode(parentNodeIndex)
	}

	leafNode.parentNodeIndex = AABBNullNode
}

func (a *AABBTree) updateLeaf(leafNodeIndex uint, newAaab *AABB) {
	node := a.nodes[leafNodeIndex]

	// if the node contains the new aabb then we just leave things
	// TODO: when we add velocity this check should kick in as often an update will lie within the velocity fattened initial aabb
	// to support this we might need to differentiate between velocity fattened aabb and actual aabb
	if node.aabb.contains(newAaab) {
		return
	}

	a.removeLeaf(leafNodeIndex)
	node.aabb = newAaab
	a.insertLeaf(leafNodeIndex)
}

func (a *AABBTree) fixUpwardsTree(treeNodeIndex uint) {
	for treeNodeIndex != AABBNullNode {
		treeNode := a.nodes[treeNodeIndex]

		// every node should be a parent
		//assert(treeNode.leftNodeIndex != AABBNullNode && treeNode.rightNodeIndex != AABBNullNode)

		// fix height and area
		leftNode := a.nodes[treeNode.leftNodeIndex]
		rightNode := a.nodes[treeNode.rightNodeIndex]
		treeNode.aabb = leftNode.aabb.merge(rightNode.aabb)

		treeNodeIndex = treeNode.parentNodeIndex
	}
}
