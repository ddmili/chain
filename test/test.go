package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
)

func main() {

}

func tt(val int, arr []int) {
	left := search(val, arr)
	if left != -1 {
		fmt.Println(arr[left-1])
	}
	right := search(val+1, arr)
	if left != -1 {
		fmt.Println(arr[right])
	}

}

func search(val int, arr []int) int {
	index := sort.SearchInts(arr, val)
	if len(arr) == index || arr[index] != val {
		return -1
	}
	return index
}

func read(filePath string) {
	m := map[string]int{}
	hot := [10]string{}
	fileHandle, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer fileHandle.Close()
	buffer := make([]byte, 255)
	for {
		n, err := fileHandle.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
		if n == 0 {
			break
		}
		str := string(buffer[:n])
		if v, ok := m[str]; ok {
			v++
			for k, s := range hot {
				if s == "" || v > m[s] {
					hot[k] = str
					break
				}
			}
			m[str] = v
		}
	}

	fmt.Println(hot)
}

func BinarySearch(target int, data []int) int {
	if len(data) < 0 {
		return -1
	}
	right := len(data) - 1
	left := 0
	for left <= right {
		mid := (left + right) / 2
		if data[mid] == target {
			return mid
		}
		if data[mid] > target {
			right = mid - 1
		}
		if data[mid] < target {
			left = mid + 1
		}
	}
	return -1
}

func test() {
	wg := &sync.WaitGroup{}
	t := map[string]int{}
	y := map[string]int{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		t = readGame("new_user")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		y = readGame("active_user")
	}()
	wg.Wait()
	for k, v := range y {
		if val, ok := t[k]; ok {
			fmt.Println(k, val/v)
		} else {
			fmt.Println(k, 0)
		}
	}
}

func readGame(filePath string) map[string]int {
	game := make(map[string]int)
	fileHandle, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return game
	}
	defer fileHandle.Close()
	buffer := make([]byte, 255)
	for {
		n, err := fileHandle.Read(buffer)
		if err != nil && err != io.EOF {
			fmt.Println(err)
		}
		if n == 0 {
			break
		}
		str := string(buffer[:n])
		arr := strings.Split(str, " ")
		if _, ok := game[arr[1]]; ok && arr[1] != "" {
			game[arr[1]]++
		}
	}

	return game

}
