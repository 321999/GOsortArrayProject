// handlers/handler.go
package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, Go Server!")
}

func SortHandler(w http.ResponseWriter, r *http.Request) {
	// Sample array to sort
	arr := []string{"banana", "apple", "orange", "grape"}

	// Sort the array
	sort.Strings(arr)

	// Display the sorted array
	fmt.Fprint(w, "Sorted Array: "+strings.Join(arr, ", "))
}
