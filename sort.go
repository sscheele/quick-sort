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
	test = test[:len(test)-1]
	var inputs = [3][]string{make([]string, len(test)), make([]string, len(test)), make([]string, len(test))}
	for i := 0; i < len(test); i++ {
		inputs[0][i] = test[i]
		inputs[1][i] = test[i]
		inputs[2][i] = test[i]
	}
	fmt.Println("Array length: ", len(test))

	start := time.Now()
	tmp := threeWayRadixSort(inputs[0])
	fmt.Println("Radix sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(tmp))
	//fmt.Println("Sorted array: ", tmp)

	start = time.Now()
	tmp = combinedSort(inputs[1])
	fmt.Println("Combined sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(tmp))
	//fmt.Println("Sorted array: ", tmp)

	start = time.Now()
	tmp = insertionSort(inputs[2])
	fmt.Println("Insertion sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(tmp))
	//fmt.Println("Sorted array: ", tmp)
}

func verifyEqual(a, b []string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func combinedSort(arr []string) []string {
	maxDepth := int(math.Log(float64(len(arr))) / 3.26) //note: using 3.26 because it's roughly log(26), where 26 is the size of our character space
	almostSorted := radixSortHelper(arr, 0, len(arr), 0, maxDepth)
	return insertionSort(almostSorted)
}

func insertionSort(arr []string) []string {
	for i := 1; i < len(arr); i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
	return arr
}

func verifySorted(arr []string) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i+1] < arr[i] {
			return false
		}
	}
	return true
}

func threeWayRadixSort(arr []string) []string {
	return radixSortHelper(arr, 0, len(arr), 0, -1)
}

func radixSortHelper(arr []string, start int, stop int, index int, maxDepth int) []string {
	if start == stop || (maxDepth != -1 && index >= maxDepth) {
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
	sorted := make([]string, stop-start)
	var indices [3]int

	//set firstIndex to the first index such that index < len(arr[i])
	var firstIndex int
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
		} else if arr[i][index] == initialLetter {
			indices[2]++
		}
	}
	indices[2] += indices[1]

	firstEqual := indices[1] + start
	firstMax := indices[2] + start

	for i := start; i < stop; i++ {
		toChange := 2
		if index >= len(arr[i]) || arr[i][index] < initialLetter {
			toChange = 0
		} else if arr[i][index] == initialLetter {
			toChange = 1
		}
		sorted[indices[toChange]] = arr[i]
		indices[toChange]++
	}

	for i := 0; i < len(sorted); i++ {
		arr[start+i] = sorted[i]
	}
	return arr, firstEqual, firstMax
}
