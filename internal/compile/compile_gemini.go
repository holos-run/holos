//go:build gemini

package compile

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime" // To get number of CPUs

	"golang.org/x/sync/errgroup"
)

// Placeholder types - replace with your actual definitions
type Component struct {
	ID   int    `json:"id"` // Add an ID to correlate requests/responses if needed later
	Name string `json:"name"`
	Data string `json:"data"`
	// Add other fields as needed
}

type BuildPlan struct {
	ComponentID   int                    `json:"componentId"` // Correlates back to Component.ID
	ComponentUsed string                 `json:"componentUsed"`
	Result        map[string]interface{} `json:"result"`
	Success       bool                   `json:"success"`
	WorkerPID     int                    `json:"workerPid"` // Optional: useful for debugging
	// Add other fields as needed
}

// executeReusableWorkers starts a pool of worker processes, distributes components
// to them, and collects build plans. Uses line-delimited JSON over stdio.
func executeReusableWorkers(ctx context.Context, numWorkers int, components []Component) ([]BuildPlan, error) {
	if len(components) == 0 {
		return nil, nil
	}
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU() // Default to number of CPUs
	}
	if numWorkers > len(components) {
		numWorkers = len(components) // No need for more workers than tasks
	}

	exePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	// Channel for distributing components to worker goroutines
	// Buffer size == len(components) to avoid blocking the sender initially
	tasksChan := make(chan Component, len(components))

	// Channel for collecting results from worker goroutines
	// Buffer size helps prevent worker goroutines blocking if main thread is slow
	resultsChan := make(chan BuildPlan, len(components))

	// Use an error group to manage worker goroutines and capture the first error
	g, childCtx := errgroup.WithContext(ctx)

	// Start worker processes and their managing goroutines
	for i := 0; i < numWorkers; i++ {
		workerIndex := i // Capture loop variable

		// Create the command for the worker process
		// Ensure the child knows it's a worker (e.g., via "--worker" flag)
		cmd := exec.CommandContext(childCtx, exePath, "--worker") // Adapt flag as needed

		stdinPipe, err := cmd.StdinPipe()
		if err != nil {
			return nil, fmt.Errorf("worker %d: failed to get stdin pipe: %w", workerIndex, err)
		}
		stdoutPipe, err := cmd.StdoutPipe()
		if err != nil {
			_ = stdinPipe.Close()
			return nil, fmt.Errorf("worker %d: failed to get stdout pipe: %w", workerIndex, err)
		}
		var stderrBuf bytes.Buffer
		cmd.Stderr = &stderrBuf

		// Start the worker process
		if err := cmd.Start(); err != nil {
			_ = stdinPipe.Close()
			// stdoutPipe closing is less critical before Wait
			return nil, fmt.Errorf("worker %d: failed to start command '%s': %w", workerIndex, exePath, err)
		}
		pid := cmd.Process.Pid // Get PID for logging/debugging
		fmt.Printf("Parent: Started worker %d (PID: %d)\n", workerIndex, pid)

		// Launch a goroutine to manage this specific worker
		g.Go(func() error {
			// Ensure resources are cleaned up for this worker's goroutine
			defer func() {
				fmt.Printf("Parent: Goroutine for worker %d (PID: %d) shutting down stdin.\n", workerIndex, pid)
				// Closing stdin signals the worker to terminate its loop
				stdinPipe.Close()
				// Wait for the process to fully exit after stdin is closed
				waitErr := cmd.Wait()
				stderrContent := stderrBuf.String()
				if waitErr != nil {
					fmt.Fprintf(os.Stderr, "Parent: Worker %d (PID: %d) exited with error (stderr: %q): %v\n", workerIndex, pid, stderrContent, waitErr)
					// Note: Returning an error here might race with errors from I/O below, errgroup handles this.
				} else {
					fmt.Printf("Parent: Worker %d (PID: %d) exited cleanly.\n", workerIndex, pid)
				}
			}()

			// Use buffered I/O for efficiency with line-based protocol
			writer := bufio.NewWriter(stdinPipe)
			// Use bufio.Scanner for reading line-delimited output
			scanner := bufio.NewScanner(stdoutPipe)

			// Process tasks from the channel until it's closed and empty
			for task := range tasksChan {
				// Marshal component to JSON
				jsonData, err := json.Marshal(task)
				if err != nil {
					return fmt.Errorf("worker %d (PID: %d): failed to marshal component %d: %w", workerIndex, pid, task.ID, err)
				}

				// Write JSON line to worker's stdin
				// fmt.Printf("Parent: Sending task %d to worker %d (PID: %d)\n", task.ID, workerIndex, pid) // Debug logging
				if _, err := writer.Write(jsonData); err != nil {
					return fmt.Errorf("worker %d (PID: %d): failed to write component %d to stdin: %w", workerIndex, pid, task.ID, err)
				}
				if err := writer.WriteByte('\n'); err != nil { // Write newline delimiter
					return fmt.Errorf("worker %d (PID: %d): failed to write newline to stdin for component %d: %w", workerIndex, pid, task.ID, err)
				}
				if err := writer.Flush(); err != nil { // Ensure data is sent
					return fmt.Errorf("worker %d (PID: %d): failed to flush stdin for component %d: %w", workerIndex, pid, task.ID, err)
				}

				// Read line (JSON BuildPlan) from worker's stdout
				if !scanner.Scan() {
					// Scanner failed, check for errors or premature EOF
					if err := scanner.Err(); err != nil {
						return fmt.Errorf("worker %d (PID: %d): error scanning stdout after sending component %d: %w", workerIndex, pid, task.ID, err)
					}
					// If no scanner error, it means EOF was reached unexpectedly
					return fmt.Errorf("worker %d (PID: %d): unexpected EOF reading stdout after sending component %d", workerIndex, pid, task.ID)
				}
				line := scanner.Bytes() // Get the line bytes

				// Unmarshal the BuildPlan
				var plan BuildPlan
				if err := json.Unmarshal(line, &plan); err != nil {
					return fmt.Errorf("worker %d (PID: %d): failed to unmarshal build plan (line: %q): %w", workerIndex, pid, string(line), err)
				}
				plan.WorkerPID = pid // Add worker PID for tracking

				// Send the result back to the main goroutine
				select {
				case resultsChan <- plan:
					// fmt.Printf("Parent: Received result for task %d from worker %d (PID: %d)\n", plan.ComponentID, workerIndex, pid) // Debug logging
				case <-childCtx.Done():
					return fmt.Errorf("worker %d (PID: %d): context cancelled while sending result for component %d: %w", workerIndex, pid, task.ID, childCtx.Err())
				}
			}
			// tasksChan was closed and this worker processed all its assigned tasks
			fmt.Printf("Parent: Worker %d (PID: %d) finished processing tasks.\n", workerIndex, pid)
			return nil // Goroutine finished successfully
		})
	}

	// Goroutine to distribute tasks
	// This runs concurrently with the worker goroutines
	go func() {
		fmt.Printf("Parent: Distributing %d tasks...\n", len(components))
		for i, comp := range components {
			comp.ID = i // Assign a unique ID for potential correlation
			select {
			case tasksChan <- comp:
				// Task sent
			case <-childCtx.Done():
				// Context cancelled before all tasks could be sent
				fmt.Fprintf(os.Stderr, "Parent: Task distribution cancelled: %v\n", childCtx.Err())
				// Closing tasksChan here ensures workers eventually stop asking for tasks
				close(tasksChan)
				return
			}
		}
		// After sending all tasks, close the channel to signal workers
		close(tasksChan)
		fmt.Println("Parent: Finished distributing tasks and closed tasks channel.")
	}()

	// Wait for all worker goroutines to complete (or one to error out)
	fmt.Println("Parent: Waiting for workers to finish...")
	err = g.Wait() // Returns the first error encountered by any worker goroutine

	// Close the results channel *after* all worker goroutines have finished
	// This signals the final result collection step
	fmt.Println("Parent: All worker goroutines finished or errored. Closing results channel.")
	close(resultsChan)

	// Collect all results sent by the workers
	// This loop reads until resultsChan is closed
	finalResults := make([]BuildPlan, 0, len(components))
	for result := range resultsChan {
		finalResults = append(finalResults, result)
	}
	fmt.Printf("Parent: Collected %d results.\n", len(finalResults))

	// Return results even if there was an error, they might be partial
	if err != nil {
		// Log the primary error that caused the errgroup to exit
		fmt.Fprintf(os.Stderr, "Parent: executeReusableWorkers returning with error: %v\n", err)
		// The finalResults slice might contain results from before the error occurred
		return finalResults, err
	}

	// Check if we got the expected number of results (only if no error occurred)
	if len(finalResults) != len(components) {
		return finalResults, fmt.Errorf("mismatch: expected %d results, got %d", len(components), len(finalResults))
	}

	return finalResults, nil // Success
}

// --- Child Worker Logic (Must be added to main.go) ---

// // Example main function incorporating the worker logic
// func main() {
// 	// Check if running as a worker
// 	if len(os.Args) > 1 && os.Args[1] == "--worker" {
// 		// Set GOMAXPROCS? Maybe not necessary if CUE isn't concurrent anyway.
// 		runWorker() // Run the dedicated worker function
// 		return      // Important: worker exits via runWorker
// 	}
//
// 	// --- Parent Process Logic ---
// 	fmt.Println("Running as parent...")
// 	componentsToProcess := []Component{
// 		{Name: "CompA", Data: "data1"},
// 		{Name: "CompB", Data: "data2"},
// 		{Name: "CompC", Data: "data3"},
// 		{Name: "CompD", Data: "data4"},
// 		{Name: "CompE", Data: "data5"},
// 	}
//
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second) // Example timeout
// 	defer cancel()
//
// 	// Use, for example, 2 worker processes
// 	buildPlans, err := executeReusableWorkers(ctx, 2, componentsToProcess)
// 	if err != nil {
// 		// Note: buildPlans might contain partial results even if err is non-nil
// 		fmt.Fprintf(os.Stderr, "Error executing reusable workers: %v\n", err)
// 		// Decide if partial results are useful or should be discarded
// 		// os.Exit(1) // Optionally exit
// 	}
//
// 	fmt.Printf("Parent: Received %d build plans:\n", len(buildPlans))
// 	// Note: Order is not guaranteed relative to input order.
// 	// Sort or process based on ComponentID if needed.
// 	sort.Slice(buildPlans, func(i, j int) bool {
// 		return buildPlans[i].ComponentID < buildPlans[j].ComponentID
// 	})
// 	for _, plan := range buildPlans {
// 		fmt.Printf("  Plan for Component %d (from worker %d): Success=%t, Result=%v\n",
// 			plan.ComponentID, plan.WorkerPID, plan.Success, plan.Result)
// 	}
// }

// // runWorker implements the logic for a child worker process.
// // It reads line-delimited JSON Components from stdin and writes
// // line-delimited JSON BuildPlans to stdout.
// func runWorker() {
// 	workerPID := os.Getpid()
// 	fmt.Fprintf(os.Stderr, "Worker (PID: %d): Starting\n", workerPID)

// 	// Use buffered I/O for stdin and stdout
// 	stdinScanner := bufio.NewScanner(os.Stdin)
// 	stdoutWriter := bufio.NewWriter(os.Stdout)
// 	defer stdoutWriter.Flush() // Ensure buffer is flushed on exit

// 	// Loop reading tasks line by line from stdin
// 	for stdinScanner.Scan() {
// 		line := stdinScanner.Bytes()
// 		if len(line) == 0 { // Skip empty lines if any occur
// 			continue
// 		}

// 		var comp Component
// 		if err := json.Unmarshal(line, &comp); err != nil {
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Failed to decode component (line: %q): %v\n", workerPID, string(line), err)
// 			// Decide strategy: exit? skip? For now, exit.
// 			os.Exit(1)
// 		}

// 		fmt.Fprintf(os.Stderr, "Worker (PID: %d): Processing component %d (%s)\n", workerPID, comp.ID, comp.Name)

// 		// --- Simulate CUE processing ---
// 		time.Sleep(time.Duration(500+rand.Intn(500)) * time.Millisecond) // Simulate variable work
// 		success := true
// 		resultData := map[string]interface{}{
// 			"processedData": fmt.Sprintf("Processed %s by %d", comp.Data, workerPID),
// 			"timestamp":     time.Now().UnixNano(),
// 		}
// 		// --- End Simulation ---

// 		plan := BuildPlan{
// 			ComponentID:   comp.ID, // Echo back the ID
// 			ComponentUsed: comp.Name,
// 			Result:        resultData,
// 			Success:       success,
// 			// WorkerPID added by parent, not needed here
// 		}

// 		// Marshal result to JSON
// 		planJSON, err := json.Marshal(plan)
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Failed to marshal build plan for component %d: %v\n", workerPID, comp.ID, err)
// 			// Decide strategy: exit? skip? For now, exit.
// 			os.Exit(1)
// 		}

// 		// Write JSON result line to stdout
// 		if _, err := stdoutWriter.Write(planJSON); err != nil {
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Failed to write build plan for component %d to stdout: %v\n", workerPID, comp.ID, err)
// 			os.Exit(1)
// 		}
// 		if err := stdoutWriter.WriteByte('\n'); err != nil { // Add newline delimiter
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Failed to write newline for component %d to stdout: %v\n", workerPID, comp.ID, err)
// 			os.Exit(1)
// 		}
// 		// Flush the buffer after each line to ensure parent receives it
// 		if err := stdoutWriter.Flush(); err != nil {
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Failed to flush stdout for component %d: %v\n", workerPID, comp.ID, err)
// 			os.Exit(1)
// 		}
// 		fmt.Fprintf(os.Stderr, "Worker (PID: %d): Finished processing component %d (%s)\n", workerPID, comp.ID, comp.Name)
// 	}

// 	// Check for scanner errors (e.g., read errors) after the loop finishes
// 	if err := stdinScanner.Err(); err != nil {
// 		// Don't report EOF as an error, it's the signal to exit cleanly
// 		if err != io.EOF {
// 			fmt.Fprintf(os.Stderr, "Worker (PID: %d): Error reading stdin: %v\n", workerPID, err)
// 			os.Exit(1)
// 		}
// 	}

// 	// EOF reached on stdin, parent closed the pipe. Exit cleanly.
// 	fmt.Fprintf(os.Stderr, "Worker (PID: %d): Stdin closed, exiting cleanly.\n", workerPID)
// 	os.Exit(0)
// }

// NOTE: Need these imports for the example usage and worker logic:
// import (
// 	"io"
// 	"math/rand"
// 	"sort"
// 	"time"
// )
