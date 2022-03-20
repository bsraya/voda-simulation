package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

type Jobs struct {
	costs         []int
	numberOfNodes int
	numberOfJobs  int
}

const servers = 3

func SimulateBB() {
	var jobs Jobs
	jobs.numberOfNodes = servers

	for i := 0; i < jobs.numberOfNodes*jobs.numberOfNodes; i++ {
		jobs.costs = append(jobs.costs, rand.Intn(100000))
	}

	jobs.numberOfJobs = len(jobs.costs) / jobs.numberOfNodes

	strings := [][]string{}

	csvFile, err := os.Create("bb_3nodes_1000jobs.csv")

	if err != nil {
		log.Fatalf("Failed to create CSV file: %s", err)
	}

	csvWriter := csv.NewWriter(csvFile)

	strings = append(strings, []string{"nodes", "jobs", "duration(ms)", "cost"})
	csvWriter.Write(strings[0])

	for i := 0; i < 1000; i++ {
		start := time.Now()
		cost, _ := BranchAndBound(jobs)
		end := time.Now()
		duration := end.Sub(start).Microseconds()

		csvWriter.Write([]string{strconv.Itoa(jobs.numberOfNodes), strconv.Itoa(len(jobs.costs) / jobs.numberOfNodes), strconv.FormatInt(duration, 10), strconv.Itoa(cost)})

		for j := 0; j < jobs.numberOfNodes; j++ {
			jobs.costs = append(jobs.costs, rand.Intn(100000))
		}
	}

	csvWriter.Flush()
	csvFile.Close()
}

func SimulateKM() {
	var totalScore int64
	nodes := servers
	m := NewMatrix(nodes)
	for j := 0; j < nodes*nodes; j++ {
		m.A[j] = rand.Int63n(10000)
	}

	strings := [][]string{}
	csvFile, err := os.Create("km_3nodes_1000jobs.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %s", err)
	}
	csvWriter := csv.NewWriter(csvFile)
	strings = append(strings, []string{"nodes", "jobs", "duration(ms)", "cost"})
	csvWriter.Write(strings[0])

	for i := 0; i < 1000; i++ {
		start := time.Now()
		result := ComputeMunkresMax(m)
		end := time.Now()
		duration := end.Sub(start).Microseconds()

		for _, rowCol := range result {
			totalScore += m.A[nodes*rowCol.Row+rowCol.Col]
		}

		csvWriter.Write([]string{strconv.Itoa(nodes), strconv.Itoa(len(m.A) / nodes), strconv.FormatInt(duration, 10), strconv.FormatInt(totalScore, 10)})

		for j := 0; j < nodes; j++ {
			m.A = append(m.A, rand.Int63n(100000))
		}

		totalScore = 0
	}
	csvWriter.Flush()
	csvFile.Close()
}

func main() {
	SimulateBB()
	SimulateKM()
}

// run simulations from 1 to 100 nodes
// func main() {
// 	var jobs Jobs
// 	strings := [][]string{}

// 	csvFile, err := os.Create("bb.csv")
// 	csvWriter := csv.NewWriter(csvFile)

// 	if err != nil {
// 		log.Fatalf("Failed to create CSV file: %s", err)
// 	}

// 	strings = append(strings, []string{"nodes", "matrix_size", "duration(ms)", "cost"})
// 	csvWriter.Write(strings[0])

// 	for i := 1; i <= 100; i++ {
// 		jobs.numberOfNodes = i
// 		matrix := []int{}
// 		jobs.costs = []int{}

// 		for j := 0; j < jobs.numberOfNodes*jobs.numberOfNodes; j++ {
// 			matrix = append(matrix, rand.Intn(10000))
// 		}

// 		jobs.costs = append(jobs.costs, matrix...)
// 		jobs.numberOfJobs = len(jobs.costs) / jobs.numberOfNodes

// 		start := time.Now()
// 		cost, _ := BranchAndBound(jobs)
// 		end := time.Now()
// 		duration := end.Sub(start).Microseconds()

// 		csvWriter.Write([]string{strconv.Itoa(jobs.numberOfNodes), strconv.Itoa(len(jobs.costs)), strconv.FormatInt(duration, 10), strconv.Itoa(cost)})
// 	}

// 	csvWriter.Flush()
// 	csvFile.Close()
// }

// func main() {
// 	var totalScore int64
// 	strings := [][]string{}

// 	csvFile, err := os.Create("munkres.csv")
// 	csvWriter := csv.NewWriter(csvFile)

// 	if err != nil {
// 		log.Fatalf("Failed to create CSV file: %s", err)
// 	}

// 	strings = append(strings, []string{"nodes", "matrix_size", "duration(ms)", "cost"})
// 	csvWriter.Write(strings[0])

// 	for i := 1; i <= 100; i++ {
// 		numberOfNodes := i
// 		m := NewMatrix(numberOfNodes)
// 		for j := 0; j < numberOfNodes*numberOfNodes; j++ {
// 			m.A[j] = rand.Int63n(10000)
// 		}

// 		start := time.Now()
// 		result := ComputeMunkresMax(m)
// 		end := time.Now()
// 		duration := end.Sub(start).Microseconds()

// 		for _, rowCol := range result {
// 			totalScore += m.A[numberOfNodes*rowCol.Row+rowCol.Col]
// 		}

// 		csvWriter.Write([]string{strconv.Itoa(numberOfNodes), strconv.Itoa(len(m.A)), strconv.FormatInt(duration, 10), strconv.FormatInt(totalScore, 10)})
// 		totalScore = 0
// 	}

// 	csvWriter.Flush()
// 	csvFile.Close()
// }
