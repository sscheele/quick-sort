package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

func main() {
	str, err := ioutil.ReadFile("/home/sam/Projects/go/src/github.com/sscheele/quick-sort/test.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	test := strings.Split(string(str), "\n")
	start := time.Now()
	threeWayRadixSort(test)
	fmt.Println("Radix sort completed in: ", time.Since(start))
	start = time.Now()
	combinedSort(test)
	fmt.Println("Combined sort completed in: ", time.Since(start))
}

func combinedSort(arr []string) []string {
	maxDepth := int(math.Log(float64(len(arr))))
	almostSorted := radixSortHelper(arr, 0, len(arr), 0, maxDepth)
	return insertionSort(almostSorted, maxDepth)
}

func insertionSort(arr []string, maxDepth int) []string {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && strCmp(arr[j], key, maxDepth) == 1 {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
	return arr
}

//return 1 iff a > b, 0 iff a == b, -1 iff a < b
func strCmp(a string, b string, i int) int {
	for ; i < len(a) && i < len(b); i++ {
		if a[i] > b[i] {
			return 1
		}
		if b[i] > a[i] {
			return -1
		}
	}
	return 0
}

func threeWayRadixSort(arr []string) []string {
	return radixSortHelper(arr, 0, len(arr), 0, -1)
}

func radixSortHelper(arr []string, start int, stop int, index int, maxDepth int) []string {
	if start == stop || (maxDepth != -1 && maxDepth >= index) {
		return arr
	}
	retVal, newStart, newStop := triage(arr, start, stop, index)
	if newStart == -1 && newStop == -1 {
		return arr
	}
	retVal = radixSortHelper(retVal, start, newStart, index, maxDepth)     //sort the lesser array
	retVal = radixSortHelper(retVal, newStart, newStop, index+1, maxDepth) //sort the equal array on the basis of the next
	retVal = radixSortHelper(retVal, newStop, stop, index, maxDepth)       //sort the greater array
	return retVal
}

//triage performs a radix sort based only on a letter index (stop is non-inclusive)
//it returns the sorted array with, first element with an equal letter, and the first element with a larger letter
func triage(arr []string, start int, stop int, index int) ([]string, int, int) {
	sorted := make([]string, len(arr))
	//sorted begins as a clone of arr
	for i := 0; i < len(arr); i++ {
		sorted[i] = arr[i]
	}
	var indices = [3]int{start, start, start}
	var firstIndex int

	//set firstIndex to the first index such that index < len(arr[i])
	for firstIndex = start; firstIndex < stop; firstIndex++ {
		if index < len(arr[firstIndex]) {
			break
		}
	}
	if firstIndex == stop {
		//there are no words with letters at index, simply return arr
		return arr, -1, -1
	}

	//initialize an index array for the radix sort
	initialLetter := arr[firstIndex][index]
	for i := start; i < stop; i++ {
		if index >= len(arr[i]) || arr[i][index] < initialLetter {
			//shorter strings come first alphabetically
			indices[1]++
			indices[2]++
		} else if arr[i][index] == initialLetter {
			indices[2]++
		}
	}
	firstEqual := indices[1]
	firstMax := indices[2]

	for i := start; i < stop; i++ {
		if index >= len(arr[i]) || arr[i][index] < initialLetter {
			sorted[indices[0]] = arr[i]
			indices[0]++
		} else if arr[i][index] == initialLetter {
			sorted[indices[1]] = arr[i]
			indices[1]++
		} else {
			sorted[indices[2]] = arr[i]
			indices[2]++
		}
	}
	return sorted, firstEqual, firstMax
}
