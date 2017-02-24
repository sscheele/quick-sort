package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"time"
)

//MyString is just a string with a helpful charAt method that lets us not worry about string length
type MyString string

func main() {
	compareSorts()
}

func compareSorts() {
	str, err := ioutil.ReadFile("/home/sam/Projects/go/src/github.com/sscheele/quick-sort/test.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	test := strings.Split(string(str), "\n")
	test = test[:len(test)-1]
	var inputs = [5][]MyString{make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test))}
	for i := 0; i < len(test); i++ {
		inputs[0][i] = MyString(test[i])
		inputs[1][i] = MyString(test[i])
		inputs[2][i] = MyString(test[i])
		inputs[3][i] = MyString(test[i])
		inputs[4][i] = MyString(test[i])
	}
	fmt.Println("Array length: ", len(test))

	start := time.Now()
	quickSort(inputs[2])
	fmt.Println("Quick sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[2]))

	fmt.Println("Running radix sort on spare array to optimize")
	threeWayRadixSort(inputs[3])
	fmt.Println("Completed")

	start = time.Now()
	threeWayRadixSort(inputs[0])
	fmt.Println("Radix sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[0]))
	//fmt.Println("Sorted array: ", inputs[0])

	start = time.Now()
	combinedSort(inputs[1])
	fmt.Println("Combined sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[1]))
	//fmt.Println("Sorted array: ", inputs[1])

	start = time.Now()
	mySort(inputs[4])
	fmt.Println("Other sort completed in: ", time.Since(start))
	fmt.Println("Verification of sort: ", verifySorted(inputs[4]))
}

func verifyEqual(a, b []MyString) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func combinedSort(arr []MyString) {
	maxDepth := int(math.Log(float64(len(arr))) / 3.26) //note: using 3.26 because it's roughly log(26), where 26 is the size of our character space
	combinedSortHelper(arr, 0, len(arr)-1, 0, maxDepth)
}

func mySort(arr []MyString) {
	myRadixSortHelper(arr, 0, len(arr)-1, 0)
}

func insertionSort(arr []MyString) {
	insertionSortHelper(arr, 0, len(arr)-1)
}

func insertionSortHelper(arr []MyString, start int, stop int) {
	for i := start + 1; i <= stop; i++ {
		key := arr[i]
		j := i - 1
		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
	}
}

func verifySorted(arr []MyString) bool {
	for i := 0; i < len(arr)-1; i++ {
		if arr[i+1] < arr[i] {
			return false
		}
	}
	return true
}

func threeWayRadixSort(arr []MyString) {
	radixSortHelper(arr, 0, len(arr)-1, 0)
}

func radixSortHelper(arr []MyString, start int, stop int, index int) {
	if start >= stop {
		return
	}
	var (
		lt  = start
		gt  = stop
		key = arr[start].charAt(index)
	)
	if stop <= start+3 {
		insertionSortHelper(arr, start, stop)
		return
	}
	for i := start + 1; i <= gt; {
		t := arr[i].charAt(index)
		if t < key {
			exchange(arr, lt, i)
			lt++
		} else if t > key {
			exchange(arr, i, gt)
			gt--
			continue
		}
		i++
	}
	radixSortHelper(arr, start, lt-1, index)
	if key >= 0 {
		radixSortHelper(arr, lt, gt, index+1)
	}
	radixSortHelper(arr, gt+1, stop, index)
}

func combinedSortHelper(arr []MyString, start int, stop int, index int, maxDepth int) {
	if start >= stop {
		return
	}
	var (
		lt  = start
		gt  = stop
		key = arr[start].charAt(index)
	)
	if stop <= start+3 || index >= maxDepth {
		insertionSortHelper(arr, start, stop)
		return
	}
	for i := start + 1; i <= gt; {
		t := arr[i].charAt(index)
		if t < key {
			exchange(arr, lt, i)
			lt++
		} else if t > key {
			exchange(arr, i, gt)
			gt--
			continue
		}
		i++
	}
	combinedSortHelper(arr, start, lt-1, index, maxDepth)
	if key >= 0 {
		combinedSortHelper(arr, lt, gt, index+1, maxDepth)
	}
	combinedSortHelper(arr, gt+1, stop, index, maxDepth)
}

func myRadixSortHelper(arr []MyString, start int, stop int, index int) {
	if start >= stop {
		return
	}
	var (
		lt  = start
		gt  = stop
		key = arr[start].charAt(index)
		eq  []MyString
	)
	if stop <= start+3 {
		insertionSortHelper(arr, start, stop)
		return
	}
	for i := start; i <= gt; {
		item := arr[i]
		//can't fucking implement recursion properly because of anonymous functions
		for len(eq)+lt-gt != 1 {
			t := item.charAt(index)
			if t == key {
				eq = append(eq, item)
				i++
				break
			} else if t > key {
				tmp := arr[gt]
				arr[gt] = item
				item = tmp
				gt--
			} else {
				arr[lt] = item
				lt++
				i++
				break
			}
		}
	}
	for i := range eq {
		arr[lt+i] = eq[i]
	}
	myRadixSortHelper(arr, start, lt-1, index)
	if key >= 0 {
		myRadixSortHelper(arr, lt, gt, index+1)
	}
	myRadixSortHelper(arr, gt+1, stop, index)
}

func exchange(a []MyString, i int, j int) {
	tmp := a[i]
	a[i] = a[j]
	a[j] = tmp
}

func (s MyString) charAt(i int) byte {
	if i > len(s) {
		panic("Index longer than string")
	}
	if i == len(s) {
		return 0
	}
	return s[i]
}

func quickSort(values []MyString) {
	quickSortHelper(values, 0, len(values)-1)
}

func quickSortHelper(values []MyString, l int, r int) {
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
