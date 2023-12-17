// // main.go
// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"sort"
// 	"strings"
// 	"sync"
// 	"time"
// )

// func HomeHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprint(w, "Hello, Go Server!")
// }

// func SortHandler(w http.ResponseWriter, r *http.Request) {
// 	// Sample array to sort
// 	arr := []string{"banana", "apple", "orange", "grape"}

// 	// Sort the array
// 	sort.Strings(arr)

// 	// Display the sorted array
// 	fmt.Fprint(w, "Sorted Array: "+strings.Join(arr, ", "))
// }

// // func main() {
// // 	http.HandleFunc("/", HomeHandler)
// // 	http.HandleFunc("/sort", SortHandler) // New route for sorting array
// // 	fmt.Println("Server is running on http://localhost:8000")
// // 	// creating a serverwhihc will run on 8000
// // 	http.ListenAndServe(":8000", nil)
// // }
// // main.go
// // package main

// // import (
// // 	"fmt"
// // 	"log"
// // 	"net/http"
// // 	"time"

// // 	"github.com/gorilla/mux"
// // 	"github.com/myusername/GOsortArrayProject/handlers"
// // )

// func main() {
// 	r := mux.NewRouter()

// 	r.HandleFunc("/process-single", ProcessSingle).Methods("POST")
// 	r.HandleFunc("/process-concurrent", ProcessConcurrent).Methods("POST")

// 	port := ":8000"
// 	fmt.Printf("Server is running on http://localhost%s\n", port)

// 	log.Fatal(http.ListenAndServe(port, r))
// }

// // building logic  for sorting the array in single and concurrent process

// type RequestPayload struct {
// 	ToSort [][]int `json:"to_sort"`
// }

// type ResponsePayload struct {
// 	SortedArrays [][]int `json:"sorted_arrays"`
// 	TimeNS       int64   `json:"time_ns"`
// }

// func ProcessSingle(w http.ResponseWriter, r *http.Request) {
// 	var request RequestPayload
// 	decodeJSON(w, r, &request)

// 	startTime := time.Now()
// 	sortedArrays := make([][]int, len(request.ToSort))
// 	for i, subArray := range request.ToSort {
// 		sortedArrays[i] = sortSubArray(subArray)
// 	}
// 	endTime := time.Now()

// 	sendResponse(w, ResponsePayload{SortedArrays: sortedArrays, TimeNS: endTime.Sub(startTime).Nanoseconds()})
// }

// func ProcessConcurrent(w http.ResponseWriter, r *http.Request) {
// 	var request RequestPayload
// 	decodeJSON(w, r, &request)

// 	startTime := time.Now()
// 	var wg sync.WaitGroup
// 	sortedArrays := make([][]int, len(request.ToSort))
// 	for i, subArray := range request.ToSort {
// 		wg.Add(1)
// 		go func(i int, subArray []int) {
// 			defer wg.Done()
// 			sortedArrays[i] = sortSubArray(subArray)
// 		}(i, subArray)
// 	}
// 	wg.Wait()
// 	endTime := time.Now()

// 	sendResponse(w, ResponsePayload{SortedArrays: sortedArrays, TimeNS: endTime.Sub(startTime).Nanoseconds()})
// }

// func sortSubArray(subArray []int) []int {
// 	sortedArray := make([]int, len(subArray))
// 	copy(sortedArray, subArray)
// 	sort.Ints(sortedArray)
// 	return sortedArray
// }

// func decodeJSON(w http.ResponseWriter, r *http.Request, target interface{}) {
// 	decoder := json.NewDecoder(r.Body)
// 	if err := decoder.Decode(target); err != nil {
// 		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
// 		return
// 	}
// }

// func sendResponse(w http.ResponseWriter, response ResponsePayload) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)

//		encoder := json.NewEncoder(w)
//		if err := encoder.Encode(response); err != nil {
//			http.Error(w, "Error encoding response", http.StatusInternalServerError)
//			return
//		}
//	}
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

const port = 8000

func main() {
	http.HandleFunc("/process-single", processSingle)
	http.HandleFunc("/process-concurrent", processConcurrent)
	fmt.Printf("Server listening on port %d...\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

// input and resoponse format
type request struct {
	ToSort [][]int `json:"to_sort"`
}

type response struct {
	SortedArrays [][]int `json:"sorted_arrays"`
	TimeNs       int64   `json:"time_ns"`
}

// sequential processing
func processSingle(w http.ResponseWriter, r *http.Request) {
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	for _, subArr := range req.ToSort {
		sort.Ints(subArr)
	}
	endTime := time.Now()

	resp := response{
		SortedArrays: req.ToSort,
		TimeNs:       endTime.Sub(startTime).Nanoseconds(),
	}

	json.NewEncoder(w).Encode(resp)
}

// for concurrent processing
func processConcurrent(w http.ResponseWriter, r *http.Request) {
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startTime := time.Now()

	// Channel to receive sorted sub-arrays
	sortedChannel := make(chan []int)

	// Launch goroutines for each sub-array
	for _, subArr := range req.ToSort {
		go func(arr []int) {
			sort.Ints(arr)
			sortedChannel <- arr
		}(subArr)
	}

	// Collect sorted sub-arrays from channel
	var sortedArrays [][]int
	// for i := range req.ToSort {
	// 	sortedArrays = append(sortedArrays, <-sortedChannel)
	// }

	var channelCount = len(req.ToSort)
	for range req.ToSort {
		copy(sortedArrays, sortedArrays[:channelCount])
		sortedArrays = append(sortedArrays, <-sortedChannel)
	}

	endTime := time.Now()

	resp := response{
		SortedArrays: sortedArrays,
		TimeNs:       endTime.Sub(startTime).Nanoseconds(),
	}

	json.NewEncoder(w).Encode(resp)
}
