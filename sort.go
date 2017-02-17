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
	var inputs = [4][]string{make([]string, len(test)), make([]string, len(test)), make([]string, len(test))}
	for i := 0; i < len(test); i++ {
		inputs[0][i] = test[i]
		inputs[1][i] = test[i]
		inputs[2][i] = test[i]
	}
	fmt.Println("Array length: ", len(test))

	start := time.Now()
	quickSort(inputs[2])
	fmt.Println("Quick sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[2]))

	start = time.Now()
	threeWayRadixSort(inputs[0])
	fmt.Println("Radix sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[0]))
	//fmt.Println("Sorted array: ", tmp)

	start = time.Now()
	combinedSort(inputs[1])
	fmt.Println("Combined sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[1]))
	//fmt.Println("Sorted array: ", tmp)
	/*
		start = time.Now()
		tmp = insertionSort(inputs[2])
		fmt.Println("Insertion sort completed in: ", time.Since(start))
		fmt.Println("Verification of sort: ", verifySorted(tmp))
		//fmt.Println("Sorted array: ", tmp)
	*/
}

func verifyEqual(a, b []string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func combinedSort(arr []string) {
	maxDepth := int(math.Log(float64(len(arr))) / 3.26) //note: using 3.26 because it's roughly log(26), where 26 is the size of our character space
	radixSortHelper(arr, 0, len(arr), 0, maxDepth)
	insertionSort(arr)
}

func insertionSort(arr []string) {
	insertionSortHelper(arr, 0, len(arr))
}

func insertionSortHelper(arr []string, start int, stop int) {
	for i := start + 1; i < stop; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func verifySorted(arr []string) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i+1] < arr[i] {
			return false
		}
	}
	return true
}

func threeWayRadixSort(arr []string) {
	radixSortHelper(arr, 0, len(arr), 0, -1)
}

func radixSortHelper(arr []string, start int, stop int, index int, maxDepth int) {
	if start == stop || (maxDepth != -1 && index >= maxDepth) {
		return
	}
	newStart, newStop := triage(arr, start, stop, index)
	if newStart == -1 && newStop == -1 {
		return
	}
	radixSortHelper(arr, start, newStart, index, maxDepth)     //sort the lesser array
	radixSortHelper(arr, newStart, newStop, index+1, maxDepth) //sort the equal array on the basis of the next
	radixSortHelper(arr, newStop, stop, index, maxDepth)       //sort the greater array
}

//triage performs an in-place radix sort based only on a letter index (stop is non-inclusive)
//it returns the sorted array with, first element with an equal letter, and the first element with a larger letter
func triage(arr []string, start int, stop int, index int) (int, int) {
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
		return -1, -1
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
	return firstEqual, firstMax
}

func quickSort(values []string) {
	quickSortHelper(values, 0, len(values)-1)
}

func quickSortHelper(values []string, l int, r int) {
	if l >= r {
		return
	}
	pivot := values[l]
	i := l + 1
	for j := l; j <= r; j++ {
		if pivot > values[j] {
			values[i], values[j] = values[j], values[i]
			i++
		}
	}
	values[l], values[i-1] = values[i-1], pivot
	quickSortHelper(values, l, i-2)
	quickSortHelper(values, i, r)
}
