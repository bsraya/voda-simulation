package main

import (
	"container/heap"
	"math"
)

var assignments []Assignment

type Assignment struct {
	workerID int
	jobID    int
}

type Node struct {
	// stores the parent of a node
	parent *Node

	// the path cost from the root to the node
	pathCost int

	// the cost of the node
	cost     int
	workerID int
	jobID    int
	assigned []bool
}

type Nodes []*Node

func initializeHeap(nodes []*Node) *Nodes {
	h := Nodes(nodes)
	heap.Init(&h)
	return &h
}

func (h Nodes) Len() int { return len(h) }

// order the heap from the lowest cost all the way up to the highest pathcost
func (h Nodes) Less(i, j int) bool { return h[i].pathCost > h[j].pathCost }

func (h Nodes) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *Nodes) Push(x interface{}) {
	*h = append(*h, x.(*Node))
}

func (h *Nodes) Pop() interface{} {
	// pop the node with the smallest path cost
	old := *h
	n := len(old)
	x := old[n-1]
	old[n-1] = nil // avoid memory leak
	*h = old[0 : n-1]
	return x
}

func AssignJobToNode(minimum *Node) {
	if minimum.parent == nil {
		return
	}
	AssignJobToNode(minimum.parent)
	assignments = append(assignments, Assignment{minimum.workerID, minimum.jobID})
}

func CalculateCost(jobs Jobs, x int, assigned []bool) int {
	totalCost := 0

	available := []bool{}
	for i := 0; i < jobs.numberOfJobs; i++ {
		available = append(available, true)
	}

	for i := x + 1; i < jobs.numberOfNodes; i++ {
		max := math.MinInt64
		minIndex := -1
		for j := 0; j < jobs.numberOfJobs; j++ {
			if !assigned[j] && available[j] && jobs.costs[i*jobs.numberOfJobs+j] > max {
				minIndex = j
				max = jobs.costs[i*jobs.numberOfJobs+j]
			}
		}
		totalCost += max
		available[minIndex] = false
	}
	return totalCost
}

// take jobs and heap as parameters
func BranchAndBound(jobs Jobs) (int, []Assignment) {
	h := initializeHeap(Nodes{})
	var assigned []bool
	for i := 0; i < jobs.numberOfJobs; i++ {
		assigned = append(assigned, false)
	}

	heap.Push(h, &Node{
		parent:   nil,
		pathCost: 0,
		cost:     0,
		workerID: -1,
		jobID:    -1,
		assigned: append([]bool{}, assigned...),
	})
	cost := 0

	for h.Len() > 0 {
		// store the node with the smallest cost
		min := heap.Pop(h).(*Node)
		i := min.workerID + 1

		if i == jobs.numberOfNodes {
			AssignJobToNode(min)
			cost = min.cost
			break
		}
		for j := 0; j < jobs.numberOfJobs; j++ {
			if !min.assigned[j] {
				child := &Node{
					parent:   min,
					pathCost: 0,
					cost:     0,
					workerID: i,
					jobID:    j,
					assigned: append([]bool{}, min.assigned...),
				}
				child.assigned[j] = true
				child.pathCost = min.pathCost + jobs.costs[i*jobs.numberOfJobs+j]
				child.cost = child.pathCost + CalculateCost(jobs, i, child.assigned)
				heap.Push(h, child)
			}
		}
	}

	// copy the assignments to a new slice called result
	// and result will be returned
	result := []Assignment{}
	for i := range assignments {
		result = append(result, assignments[i])
	}

	// set global assignment array to empty
	assignments = []Assignment{}

	return cost, result
}
