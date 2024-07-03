//Task: Concurrent Sorting with Goroutines and Channels
//
//Description:
//Write a program that sorts a large array of integers concurrently.
//The main goroutine should divide the array into several subArrays and send these subArrays to a channel.
//Multiple worker goroutines should receive these subArrays, sort them,
//and send the sorted subArrays back to another channel.
//The main goroutine should collect these sorted subArrays and merge them into a single sorted array.

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"
)

// Function to generate a large array of random integers
func generateArray(size int) []int {
	rand.Seed(time.Now().UnixNano())
	array := make([]int, size)
	for i := range array {
		array[i] = rand.Intn(100)
	}
	return array
}

// Function to divide an array into subArrays
func divideArray(array []int, numSubArrays int) [][]int {
	subArrays := make([][]int, numSubArrays)
	subarraySize := len(array) / numSubArrays
	for i := 0; i < numSubArrays; i++ {
		if i == numSubArrays-1 {
			subArrays[i] = array[i*subarraySize:]
		} else {
			subArrays[i] = array[i*subarraySize : (i+1)*subarraySize]
		}
	}
	return subArrays
}

// Worker function to sort subArrays
func worker(subArrayChan, sortedSubArrayChan chan []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for subArray := range subArrayChan {
		sort.Ints(subArray)
		sortedSubArrayChan <- subArray
	}
}

// Function to merge sorted subArrays into a single sorted array
func mergeSortedSubArrays(sortedSubArrays [][]int) []int {
	if len(sortedSubArrays) == 0 {
		return []int{}
	}

	result := sortedSubArrays[0]
	for i := 1; i < len(sortedSubArrays); i++ {
		result = mergeTwoSortedArrays(result, sortedSubArrays[i])
	}
	return result
}

// Function to merge two sorted arrays
func mergeTwoSortedArrays(arr1, arr2 []int) []int {
	result := make([]int, 0, len(arr1)+len(arr2))
	i, j := 0, 0
	for i < len(arr1) && j < len(arr2) {
		if arr1[i] < arr2[j] {
			result = append(result, arr1[i])
			i++
		} else {
			result = append(result, arr2[j])
			j++
		}
	}
	result = append(result, arr1[i:]...)
	result = append(result, arr2[j:]...)
	return result
}

func main() {
	size := 1000
	numSubArrays := 10 // Number of subArrays
	numWorkers := 5    // Number of worker goroutines

	var wg sync.WaitGroup
	sortedSubArrayChan := make(chan []int, numSubArrays)
	subArrayChan := make(chan []int, numSubArrays)

	// Generate a large array of random integers
	array := generateArray(size)

	// Divide the array into subArrays
	subArrays := divideArray(array, numSubArrays)

	go func() {
		for _, subArray := range subArrays {
			sort.Ints(subArray)
			subArrayChan <- subArray
		}
		close(subArrayChan)
	}()

	go func() {
		wg.Wait()
		close(sortedSubArrayChan)
	}()

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(subArrayChan, sortedSubArrayChan, &wg)
	}

	sortedSubArrays := make([][]int, numSubArrays)
	for sortedSubArray := range sortedSubArrayChan {
		sortedSubArrays = append(sortedSubArrays, sortedSubArray)
	}

	finalSortedArray := mergeSortedSubArrays(sortedSubArrays)

	fmt.Println("Final Sorted Array:", finalSortedArray)
}
