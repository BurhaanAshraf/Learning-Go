package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// ============================================================================
// GO CONTEXT — REVISION + INTERVIEW MASTER FILE
// ============================================================================
//
// WHAT CONTEXT SOLVES
//   Cancellation signals, deadline enforcement, and request-scoped metadata —
//   carried through a call tree without changing every function signature.
//
// MENTAL MODEL: immutable linked-list tree
//   Background → WithTimeout → WithValue → WithValue
//   Every constructor wraps the parent and returns a NEW child. Nothing mutates.
//   Cancelling a node closes its Done() channel and propagates to all descendants.
//
// THREE LAWS (memorise these)
//   1. Cancellation flows DOWNWARD   — cancel a parent, all children die
//   2. Value lookup flows UPWARD     — walks the chain until match or nil at root
//   3. Parents are UNAWARE of child cancellations
//
// CONCURRENCY
//   Context is fully THREAD-SAFE. One context can be passed to thousands of
//   goroutines simultaneously — no mutex needed.
//
// CONSTRUCTOR QUICK REFERENCE
//   Background()         real root — main(), server init, test setup
//   TODO()               placeholder — identical to Background() at runtime,
//                        signals "context design not decided yet"
//   WithCancel()         your code decides when to stop
//   WithCancelCause()    (Go 1.20+) like WithCancel + attach a custom error reason
//   WithTimeout(d)       stops after duration d
//   WithDeadline(t)      stops at absolute time t
//                        WithTimeout(d) == WithDeadline(time.Now().Add(d))
//   WithoutCancel()      (Go 1.21+) detach from parent cancellation, keep values
//   WithValue()          attach request-scoped metadata
//
// ALWAYS defer cancel()
//   For WithCancel / WithTimeout / WithDeadline.
//   Without it, an internal timer goroutine leaks until the deadline fires naturally.
//
// WHAT BELONGS IN WithValue
//   OK      → trace ID, request ID, auth user identity
//   NEVER   → DB connections, config, feature flags  (use dependency injection)
//
// FUNCTION SIGNATURE RULE
//   context.Context is ALWAYS the first parameter. This is Go convention and
//   enforced by the staticcheck linter.
//   ✓  func Save(ctx context.Context, id string) error
//   ✗  func Save(id string, ctx context.Context) error
//
// NEVER STORE CONTEXT IN A STRUCT
//   Context is per-operation, not per-object lifetime.
//   ✗  type Handler struct { ctx context.Context }
//   ✓  Pass ctx as a function argument every time.
// ============================================================================

// ============================================================================
// KEY TYPES
// ============================================================================
//
// RULE: Always use an unexported custom type for context keys.
//   A plain string key (e.g. "requestID") can silently collide with another
//   package that uses the same string. An unexported type is package-private —
//   impossible to guess or collide with from outside this package.
//
//   ✗  ctx = context.WithValue(ctx, "requestID", "abc")
//   ✓  ctx = context.WithValue(ctx, requestIDKey, "abc")

type contextKey string

const (
	requestIDKey contextKey = "requestID"
	osKey        contextKey = "OS"
)

// ============================================================================
// PATTERN 1: Reading values — always comma-ok, never bare assertion
// ============================================================================
//
// ctx.Value() returns `any`. If the key is absent it returns nil.
// A bare type assertion on nil PANICS at runtime.
// Comma-ok costs nothing and never panics.
//
//   ✗  id := ctx.Value(requestIDKey).(string)        // panics if missing
//   ✓  id, ok := ctx.Value(requestIDKey).(string)    // safe always

func logWithContext(ctx context.Context, msg string) {
	requestID, ok := ctx.Value(requestIDKey).(string)
	if !ok {
		requestID = "UNKNOWN"
	}

	osName, ok := ctx.Value(osKey).(string)
	if !ok {
		osName = "UNKNOWN"
	}

	log.Printf("[%s] [%s] %s", requestID, osName, msg)
}

// ============================================================================
// PATTERN 2: Preemptive check — bail before expensive work if already cancelled
// ============================================================================
//
// If a context is already expired when the function is called, there is no
// point starting the work. The select/default idiom is a non-blocking check.

func processItem(ctx context.Context, n int) (string, error) {
	select {
	case <-ctx.Done():
		return "", ctx.Err() // context.Canceled or context.DeadlineExceeded
	default:
	}

	// simulate expensive work
	if n%2 == 0 {
		return fmt.Sprintf("%d is even", n), nil
	}
	return fmt.Sprintf("%d is odd", n), nil
}

// ============================================================================
// PATTERN 3: Long-running goroutine — select loop with ctx.Done()
// ============================================================================
//
// ctx.Done()  → channel that closes when context is cancelled or timed out
// ctx.Err()   → WHY it was cancelled:
//                 context.Canceled         (explicit cancel() call)
//                 context.DeadlineExceeded (timeout or deadline elapsed)
// context.Cause(ctx) → the specific error passed to cancel() via WithCancelCause,
//                      or nil if using plain WithCancel/WithTimeout
//
// Without the ctx.Done() case this goroutine runs forever → goroutine leak.

func DoWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stopped | Err:", ctx.Err(), "| Cause:", context.Cause(ctx))
			return
		default:
			fmt.Println("working...")
		}
		time.Sleep(300 * time.Millisecond)
	}
}

// ============================================================================
// PATTERN 4: Value shadowing
// ============================================================================
//
// Storing the same key in a child does NOT mutate the parent.
// Value lookup walks upward and stops at the first match (nearest child wins).
// Use this to override a value for a sub-operation without affecting the caller.

func demonstrateShadowing() {
	root := context.Background()
	ctx1 := context.WithValue(root, requestIDKey, "old-id")
	ctx2 := context.WithValue(ctx1, requestIDKey, "new-id")

	fmt.Println(ctx1.Value(requestIDKey)) // "old-id"  — parent unchanged
	fmt.Println(ctx2.Value(requestIDKey)) // "new-id"  — child shadows parent
}

// ============================================================================
// PATTERN 5: Value retention after cancellation
// ============================================================================
//
// Cancellation closes Done() and sets Err() — it does NOT delete stored values.
// After a timeout you can still read the request ID for logging/error reporting.

func demonstrateValueRetention() {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ctx = context.WithValue(ctx, requestIDKey, "req-abc-123")

	time.Sleep(400 * time.Millisecond) // let timeout fire

	fmt.Println("ctx expired?", ctx.Err())                  // DeadlineExceeded
	fmt.Println("value still readable:", ctx.Value(requestIDKey)) // "req-abc-123"
}

// ============================================================================
// PATTERN 6: WithTimeout — duration-based deadline
// ============================================================================
//
// Always defer cancel() immediately after creation.
// If the work finishes in 10ms on a 5s timeout, defer frees the timer at 10ms
// instead of letting it sit in memory for 5 seconds.

func demonstrateTimeout() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err := processItem(ctx, 10)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	fmt.Println(result)
}

// ============================================================================
// PATTERN 7: WithCancel — manual cancellation
// ============================================================================
//
// Use when YOUR logic decides when to stop, not a clock.
// Common use-cases: first-result-wins across multiple goroutines, error threshold.

func demonstrateManualCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // guard-rail in case of early return paths

	done := make(chan struct{})

	go func() {
		defer close(done)
		time.Sleep(200 * time.Millisecond)
		cancel() // all consumers of ctx stop now
		fmt.Println("goroutine triggered cancel")
	}()

	<-done
	fmt.Println("ctx.Err():", ctx.Err()) // context.Canceled
}

// ============================================================================
// PATTERN 8: WithCancelCause (Go 1.20+) — why was it cancelled?
// ============================================================================
//
// Problem with plain WithCancel: ctx.Err() only returns context.Canceled —
// you can't tell WHICH of 10 goroutines caused the cancellation or why.
//
// WithCancelCause lets you attach a specific error at the call site.
//   ctx.Err()          → always context.Canceled (generic, unchanged)
//   context.Cause(ctx) → the custom error you passed to cancel()
//
// Interview tip: use this for fan-out patterns where you want to surface
// the first real error that caused the abort.

func demonstrateCancelCause() {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(nil) // nil means "no specific cause" — used as guard-rail

	go func() {
		time.Sleep(200 * time.Millisecond)
		cancel(errors.New("db connection lost")) // attach the real reason
	}()

	<-ctx.Done()
	fmt.Println("Err()  →", ctx.Err())          // context.Canceled
	fmt.Println("Cause() →", context.Cause(ctx)) // db connection lost
}

// ============================================================================
// PATTERN 9: WithoutCancel (Go 1.21+) — detach from parent cancellation
// ============================================================================
//
// Problem: a request times out, but you still need to write an audit log or
// send a billing event. You can't use the expired request context for that.
//
// WithoutCancel returns a child that:
//   - inherits all key-values from the parent chain (still readable)
//   - is NOT cancelled when the parent times out or is cancelled
//   - has no deadline of its own
//
// Use sparingly — the detached context needs its own timeout for the cleanup work.

func demonstrateWithoutCancel() {
	parent, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	parent = context.WithValue(parent, requestIDKey, "REQ-99")

	// detach: inherits values, immune to parent timeout
	cleanup := context.WithoutCancel(parent)

	time.Sleep(100 * time.Millisecond) // parent times out

	fmt.Println("parent.Err()  →", parent.Err())              // DeadlineExceeded
	fmt.Println("cleanup.Err() →", cleanup.Err())             // nil (still alive)
	fmt.Println("cleanup value →", cleanup.Value(requestIDKey)) // "REQ-99"

	// best practice: give the cleanup work its own tight deadline
	cleanupCtx, cleanupCancel := context.WithTimeout(cleanup, 2*time.Second)
	defer cleanupCancel()
	_ = cleanupCtx // use for audit log, billing event, etc.
}

// ============================================================================
// PATTERN 10: Real-world middleware chain
// ============================================================================
//
// This is the most common production pattern.
// A request enters at the HTTP handler. The request ID is injected into context
// and flows through every layer — service, repo, DB query — without any layer
// needing to know about request IDs explicitly.
//
// Every layer only receives (ctx, businessArgs) and passes ctx down.
// Cancellation and timeouts propagate automatically.

func httpHandler(w http.ResponseWriter, r *http.Request) {
	// inject request-scoped metadata at the edge
	ctx := context.WithValue(r.Context(), requestIDKey, r.Header.Get("X-Request-ID"))

	if err := serviceLayer(ctx, "order-42"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serviceLayer(ctx context.Context, orderID string) error {
	// add a tight deadline for this specific operation
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	logWithContext(ctx, "processing "+orderID)
	return repoLayer(ctx, orderID)
}

func repoLayer(ctx context.Context, orderID string) error {
	// context carries both the request ID and the 3s deadline into the DB call
	_ = orderID
	_ = ctx
	return nil // in real code: db.QueryContext(ctx, ...)
}

// ============================================================================
// PATTERN 11: stdlib I/O — always use the Context variant
// ============================================================================
//
// HTTP and SQL both have Context-aware versions.
// Without them, a cancelled context does NOT abort the in-flight network call —
// the goroutine blocks until the OS-level TCP timeout fires (could be minutes).

func fetchWithContext(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func queryWithContext(ctx context.Context, db *sql.DB) error {
	rows, err := db.QueryContext(ctx, "SELECT id FROM orders WHERE status = $1", "pending")
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

// ============================================================================
// EXECUTION
// ============================================================================

func Context() {
	fmt.Println("=== 1. comma-ok value read ===")
	ctx := context.WithValue(context.Background(), requestIDKey, "req-001")
	ctx = context.WithValue(ctx, osKey, "Ubuntu")
	logWithContext(ctx, "boot")

	fmt.Println("\n=== 2. value shadowing ===")
	demonstrateShadowing()

	fmt.Println("\n=== 3. timeout ===")
	demonstrateTimeout()

	fmt.Println("\n=== 4. manual cancel ===")
	demonstrateManualCancel()

	fmt.Println("\n=== 5. value retention after timeout ===")
	demonstrateValueRetention()

	fmt.Println("\n=== 6. WithCancelCause (Go 1.20+) ===")
	demonstrateCancelCause()

	fmt.Println("\n=== 7. WithoutCancel (Go 1.21+) ===")
	demonstrateWithoutCancel()

	fmt.Println("\n=== 8. worker loop with cancellation ===")
	workerCtx, workerCancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer workerCancel()
	DoWork(workerCtx)
}