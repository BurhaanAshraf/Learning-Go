package main

import (
	"fmt"
	"sync"
)

/*
================================================================
CORE CONCEPT — sync.Pool
================================================================

WHY IT EXISTS:
  Repeatedly allocating and discarding short-lived objects (e.g.
  buffers on every HTTP request) creates GC pressure — the
  garbage collector has to keep tracking and cleaning up objects
  that are almost immediately dead. sync.Pool lets you RECYCLE
  objects instead of allocating fresh ones each time.

  Without Pool:  New -> Use -> Garbage Collected  (every time)
  With Pool:     Get -> Use -> Put back -> reused by someone else

WHAT IT ACTUALLY IS:
  A temporary object cache — but NOT a general-purpose cache like
  Redis or an LRU map. The Go runtime is free to silently evict
  ANY object from the pool during GC, with zero notice. So: it's
  cache-like in behavior (reuse when available), but you must
  never rely on something still being there — treat every Get()
  as "might be new, might be recycled," never "will definitely
  still hold what I put in."

⚠ THE #1 REAL-WORLD GOTCHA — stale data:
  Put() does NOT clear/reset the object's fields. If you don't
  manually reset before Put() (or right after Get()), the next
  caller can receive an object still carrying YOUR old data.
  This is a genuine bug class in production code, not just a
  theoretical concern — always reset before returning to the pool.

WHEN TO USE / AVOID:
  ✓ bytes.Buffer, JSON/XML encoders, HTTP req/resp buffers,
    any short-lived, frequently-allocated temporary object.
  ✗ DB connections, file handles, or anything with a lifecycle
    that needs explicit cleanup (Close(), etc.) — Pool gives no
    guarantee about WHEN or IF Put() objects get reused, so you
    can't rely on it for resource lifecycle management.
================================================================
*/

type person struct {
	name string
	age  int
}

// reset clears a person's fields before it goes back into the
// pool — without this, the next Get() caller could silently
// inherit stale data from whoever used this object last.
func (p *person) reset() {
	p.name = ""
	p.age = 0
}

func syncPool() {
	// New is optional — if Get() finds no reusable object, New()
	// is called automatically to create one.
	pool := sync.Pool{
		New: func() any {
			fmt.Println("Creating a new Person.")
			return &person{}
		},
	}

	// --- Get(): pool is empty -> New() runs ---
	person1 := pool.Get().(*person)
	person1.name = "John"
	person1.age = 18
	fmt.Printf("Person 1 -> Name: %s Age: %d\n", person1.name, person1.age)

	// --- Put(): return to pool for reuse. Object is NOT ---
	// destroyed or cleared automatically — reset manually first.
	person1.reset()
	pool.Put(person1)
	fmt.Println("Returned Person 1 to Pool.")

	// --- Get(): pool has person1 available -> reused, not new ---
	// Because we reset() before Put(), person2 starts clean
	// instead of inheriting "John"/18.
	person2 := pool.Get().(*person)
	fmt.Println("Got Person 2 (reused, reset):", person2)

	// --- Pool empty again -> New() runs ---
	person3 := pool.Get().(*person)
	person3.name = "Jane"
	person3.age = 21
	fmt.Println("Got Person 3 (new):", person3)

	person2.reset()
	person3.reset()
	pool.Put(person2)
	pool.Put(person3)
	fmt.Println("Returned Person 2 & Person 3.")

	person4 := pool.Get().(*person)
	person5 := pool.Get().(*person)
	fmt.Println("Got Person 4:", person4)
	fmt.Println("Got Person 5:", person5)

	person4.reset()
	person5.reset()
	pool.Put(person4)
	pool.Put(person5)
	fmt.Println("Returned Person 4 & Person 5.")
}

/*
================================================================
HOW sync.Pool WORKS
================================================================
Get() called
     │
Pool has an object available?
     ├── No  -> call New() -> return new object
     └── Yes -> return existing (possibly stale!) object

Put(obj) called
     │
Object stored for potential reuse
(NOT copied — the SAME object may be handed back by a future Get())
(NOT guaranteed to survive — GC may drop it at any time)
================================================================
QUICK REVISION
================================================================
sync.Pool -> temporary reusable-object cache, reduces GC pressure.

Get()  -> returns an available object, or calls New() if empty.
Put()  -> returns an object for potential reuse (does NOT reset it).
New()  -> factory function, called automatically when pool is empty.

⚠ Always reset an object's state before Put() — Pool never does
  this for you, and forgetting it is the most common real bug.

Not guaranteed: objects can vanish from the pool after any GC
cycle — never rely on something still being there.

Good fits:    bytes.Buffer, encoders, temporary/short-lived structs.
Bad fits:     DB connections, file handles, anything needing
              explicit lifecycle management (Close, etc.)
================================================================
INTERVIEW Q&A
================================================================
Q: What problem does sync.Pool solve?
A: Reduces allocations and GC pressure by reusing short-lived
   temporary objects instead of allocating fresh ones each time.

Q: What happens if Get() finds no object in the pool?
A: It calls New() if provided; if New() is nil, Get() returns nil.

Q: Does sync.Pool guarantee an object stays available for reuse?
A: No — the GC may evict pooled objects at any time, with no
   notification. Never rely on persistence.

Q: Is sync.Pool a "cache" in the general sense (like Redis, LRU)?
A: It's cache-like (objects may be reused instead of recreated),
   but it's NOT persistent or guaranteed storage — think of it as
   a GC-pressure optimization, not a data store.

Q: Does Put() reset an object's fields automatically?
A: No — this is the most common real-world sync.Pool bug. You
   must manually reset state before Put() (or right after Get()),
   or callers can silently receive stale data from a previous user.

Q: Why is sync.Pool useful?
A: It improves performance under high allocation churn by
   recycling objects, reducing both allocation cost and GC work.
================================================================
*/
