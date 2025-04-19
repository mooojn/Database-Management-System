package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)
var tempResults [][]string // Temporary store for rollback

// CORS Middleware
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Global Variables
var (
	transactionActive = false
	transactionMutex  sync.Mutex
	results           [][]string // Store structured data (Operand1, Operator, Operand2, Result)
)

// JSON Response Helper
func jsonResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

// Start Transaction
func startTransaction(w http.ResponseWriter, r *http.Request) {
	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	if transactionActive {
		http.Error(w, "Transaction already active", http.StatusConflict)
		return
	}

	transactionActive = true
	tempResults = append([][]string{}, results...) // Save previous state
	results = [][]string{} // Start fresh transaction
	fmt.Println("Transaction started")
	fmt.Fprintln(w, "Transaction started")
}


// Perform Arithmetic Operation
func performOperation(w http.ResponseWriter, r *http.Request) {
	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	if !transactionActive {
		jsonResponse(w, http.StatusBadRequest, "No active transaction")
		return
	}

	op := r.URL.Query().Get("operation") // add, sub, mul, div
	a, err1 := strconv.Atoi(r.URL.Query().Get("a"))
	b, err2 := strconv.Atoi(r.URL.Query().Get("b"))

	if err1 != nil || err2 != nil {
		jsonResponse(w, http.StatusBadRequest, "Invalid numbers")
		return
	}

	var result float64
	switch op {
	case "add":
		result = float64(a + b)
	case "sub":
		result = float64(a - b)
	case "mul":
		result = float64(a * b)
	case "div":
		if b == 0 {
			jsonResponse(w, http.StatusBadRequest, "Cannot divide by zero")
			return
		}
		result = float64(a) / float64(b)
	default:
		jsonResponse(w, http.StatusBadRequest, "Invalid operation")
		return
	}

	// Store structured result
	results = append(results, []string{strconv.Itoa(a), op, strconv.Itoa(b), fmt.Sprintf("%.2f", result)})

	fmt.Printf("Operation performed: %d %s %d = %.2f\n", a, op, b, result)
	jsonResponse(w, http.StatusOK, fmt.Sprintf("Result: %.2f", result))
}

// Commit Transaction (Save to CSV)
func commitTransaction(w http.ResponseWriter, r *http.Request) {
	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	if !transactionActive {
		jsonResponse(w, http.StatusBadRequest, "No active transaction")
		return
	}

	fmt.Println("Committing transaction...")

	// Open CSV file (append mode)
	file, err := os.OpenFile("transactions.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// If file is empty, write headers
	fileInfo, _ := file.Stat()
	if fileInfo.Size() == 0 {
		if err := writer.Write([]string{"Operand1", "Operator", "Operand2", "Result"}); err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to write CSV header")
			return
		}
	}

	// Write transactions to CSV
	for _, entry := range results {
		if err := writer.Write(entry); err != nil {
			jsonResponse(w, http.StatusInternalServerError, "Failed to write to CSV")
			return
		}
	}

	fmt.Println("Transaction committed and saved to CSV")

	// Reset transaction
	transactionActive = false
	results = [][]string{}

	jsonResponse(w, http.StatusOK, "Transaction committed and saved to CSV")
}
// Rollback Transaction (Discard Operations)
func rollbackTransaction(w http.ResponseWriter, r *http.Request) {
	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	if !transactionActive {
		http.Error(w, "No active transaction", http.StatusBadRequest)
		return
	}

	// Discard changes by resetting tempResults
	results = tempResults

	transactionActive = false
	fmt.Println("Transaction rolled back. Changes discarded.")
	fmt.Fprintln(w, "Transaction rolled back. No changes saved.")
}

// Cancel Transaction (Discard changes)
func cancelTransaction(w http.ResponseWriter, r *http.Request) {
	transactionMutex.Lock()
	defer transactionMutex.Unlock()

	if !transactionActive {
		http.Error(w, "No active transaction", http.StatusBadRequest)
		return
	}

	transactionActive = false
	results = [][]string{} // Clear stored results

	fmt.Println("Transaction cancelled")
	fmt.Fprintln(w, "Transaction cancelled")
}


func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/start", startTransaction)
	mux.HandleFunc("/operate", performOperation)
	mux.HandleFunc("/commit", commitTransaction)
    mux.HandleFunc("/cancel", cancelTransaction)
	mux.HandleFunc("/rollback", rollbackTransaction)


	// Wrap handlers with CORS middleware
	handler := enableCORS(mux)

	log.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
