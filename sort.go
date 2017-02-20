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
	//trainGE()
	compareSorts()
}

func trainGE() {
	str, err := ioutil.ReadFile("/home/sam/Projects/go/src/github.com/sscheele/quick-sort/test.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	test := strings.Split(string(str), "\n")
	test = test[:len(test)-1]

	for j := 0; j < 3; j++ {
		throwAway := make([]MyString, len(test))
		for i := 0; i < len(test); i++ {
			throwAway[i] = MyString(test[i])
		}
		combinedSort(throwAway)
	}
	alpha := .1
	divVal := 4.0
	for {
		if time.Now().UnixNano()%100 == 0 {
			alpha *= 4 //get out of sticking points
		}
		fmt.Println("DivVal: ", divVal)
		var inputs = [3][]MyString{make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test))}
		for i := 0; i < len(test); i++ {
			inputs[0][i] = MyString(test[i])
			inputs[1][i] = MyString(test[i])
			inputs[2][i] = MyString(test[i])
		}
		start := time.Now()
		cSortTrain(inputs[0], divVal-alpha)
		sinceLow := time.Since(start)

		start = time.Now()
		cSortTrain(inputs[1], divVal+alpha)
		sinceHigh := time.Since(start)

		start = time.Now()
		cSortTrain(inputs[2], divVal)
		sinceBase := time.Since(start)

		if sinceHigh > sinceLow && sinceHigh > sinceBase {
			divVal += alpha
			continue
		}
		if sinceLow > sinceHigh && sinceLow > sinceBase {
			divVal -= alpha
			continue
		}
		//base is the greatest, lower alpha
		alpha = alpha / 2.0
	}
}

func compareSorts() {
	str, err := ioutil.ReadFile("/home/sam/Projects/go/src/github.com/sscheele/quick-sort/test.dat")
	if err != nil {
		fmt.Println(err)
		return
	}
	test := strings.Split(string(str), "\n")
	test = test[:len(test)-1]
	var inputs = [4][]MyString{make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test)), make([]MyString, len(test))}
	for i := 0; i < len(test); i++ {
		inputs[0][i] = MyString(test[i])
		inputs[1][i] = MyString(test[i])
		inputs[2][i] = MyString(test[i])
		inputs[3][i] = MyString(test[i])
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
	//fmt.Println("Sorted array: ", tmp)

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
	radixSortHelper(arr, 0, len(arr)-1, 0, maxDepth)
}

func cSortTrain(arr []MyString, x float64) {
	maxDepth := int(math.Log(float64(len(arr))) / x)
	radixSortHelper(arr, 0, len(arr)-1, 0, maxDepth)
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
	radixSortHelper(arr, 0, len(arr)-1, 0, -1)
}

func radixSortHelper(arr []MyString, start int, stop int, index int, maxChar int) {
	if start >= stop {
		return
	}
	var (
		lt  = start
		gt  = stop
		key = arr[start].charAt(index)
	)
	if stop <= start+3 || (maxChar != -1 && index >= maxChar) {
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
	radixSortHelper(arr, start, lt-1, index, maxChar)
	if key >= 0 {
		radixSortHelper(arr, lt, gt, index+1, maxChar)
	}
	radixSortHelper(arr, gt+1, stop, index, maxChar)
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
