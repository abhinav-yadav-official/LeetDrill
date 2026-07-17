# Patterns: recognizing them, and the generic supersets behind them

Most interview and contest problems are one or two known patterns wearing a
costume. The skill is stripping the costume: read the input shape and the
question, ignore the flavor text, and match against a short list of
canonical techniques.

## Signal → pattern lookup

| If the problem says / implies... | Reach for |
|---|---|
| contiguous subarray/substring with some property | sliding window |
| sorted array, find a pair/triplet matching a sum | two pointers |
| kth largest/smallest, "top k" | heap (priority queue) or quickselect |
| all subsets / combinations / permutations | backtracking |
| shortest path, unweighted graph or grid | BFS |
| shortest path, weighted graph, non-negative edges | Dijkstra |
| shortest path, negative edges allowed | Bellman-Ford |
| all-pairs shortest path, small n (≤ ~400) | Floyd-Warshall |
| connectivity / grouping / cycle detection, undirected | Union-Find (DSU) |
| dependency order / prerequisites / "can this finish" | topological sort |
| minimum spanning tree / connect everything cheaply | Kruskal or Prim |
| range sum/min/max query, static array | prefix sum array |
| range sum/min/max query, array with updates | Fenwick tree (BIT) or segment tree |
| "next greater/smaller element" | monotonic stack |
| sliding window max/min | monotonic deque |
| count ways / min cost / max value, overlapping subproblems | dynamic programming |
| 0/1 choice per item under a capacity constraint | knapsack DP |
| longest/shortest subsequence with an ordering constraint | sequence DP (LIS/LCS family) |
| merge or select among overlapping intervals | sort by start (or end) + linear scan / greedy |
| matching brackets / nesting / "undo" semantics | stack |
| detect a cycle in a linked list, find the middle | fast/slow pointers (Floyd's) |
| running median of a stream | two heaps (max-heap below, min-heap above) |
| merge k sorted lists/arrays | heap (k-way merge) |
| substring search / pattern matching | KMP, Z-function, or rolling hash |
| word search in a grid, count islands, flood fill | DFS/BFS over the grid |
| n ≤ ~20 and the state is "which subset is used" | bitmask DP |
| "minimize the maximum" / "maximize the minimum" and feasibility is monotonic in the answer | binary search on the answer |
| prefix/suffix property over strings (autocomplete, word break) | trie |

This table is a starting point, not a lookup table to memorize verbatim —
the next section is the more durable version of the same idea.

## Generic pattern supersets

Narrow, named patterns nest inside broader families. When you don't
recognize the exact named pattern, recognizing the *family* is often enough
to get unstuck — the family narrows the search space to 2-3 techniques
instead of the full list above.

- **Window/pointer family** — linear scan with pointer(s) maintaining an
  invariant over a contiguous range.
  - Two pointers (opposite ends, converging) ⊂ this family.
  - Sliding window (same-direction pointers, expand/shrink) ⊂ this family.
  - Fast/slow pointers (cycle detection) ⊂ this family.
  - Recognize the family from: "array/string", "contiguous", "O(n) or
    O(n log n) expected", no need to revisit elements once passed.

- **Search family** — exhaustive exploration of a state space, with or
  without pruning.
  - DFS (explore one path fully before backtracking) ⊂ this family.
  - Backtracking (DFS + undo a partial choice on dead end) ⊂ DFS.
  - BFS (explore level by level, shortest-steps property) ⊂ this family,
    but is *not* a subset of DFS — different traversal order, different
    guarantee (shortest path in unweighted graphs).
  - Recognize the family from: "all ways to...", "does there exist a
    path...", tree/graph/grid input, n small enough for exponential blowup
    (see the complexity table) unless pruning or DP collapses it.

- **Graph traversal family** — superset of every graph algorithm.
  - DFS, BFS, Union-Find, topological sort, Dijkstra, Bellman-Ford,
    Floyd-Warshall, Kruskal/Prim (MST) ⊂ this family.
  - Recognize the family from: explicit graph/tree/grid, or an *implicit*
    graph (states are nodes, transitions are edges — e.g. word ladder,
    puzzle-solving BFS).

- **DP family** — optimal substructure + overlapping subproblems, but the
  distinguishing feature between DP *sub-patterns* is the shape of the
  state, not the recurrence style:
  - Linear state (`dp[i]`) — Kadane's, house robber, climbing stairs.
  - Two-sequence state (`dp[i][j]`) — LCS, edit distance.
  - Interval state (`dp[i][j]` = answer over `[i, j]`) — matrix chain
    multiplication, burst balloons, palindrome partitioning.
  - Subset state (`dp[mask]`) — bitmask DP, TSP.
  - Tree state (`dp[node][...]`) — tree DP, computed via post-order DFS.
  - Digit state (`dp[position][tight][...]`) — digit DP for "count numbers
    in [L, R] with property X".
  - Once you see "overlapping subproblems," the real work is picking which
    of these six state shapes fits — that choice is 80% of a DP problem.

- **Range query family** — precompute or maintain an aggregate over ranges.
  - Prefix sum (static array, O(1) query, no updates) ⊂ this family.
  - Difference array (static range *updates*, O(1) per update) ⊂ this
    family — the write-side dual of prefix sums.
  - Fenwick tree / BIT (point update, range query, O(log n) each) ⊂ this
    family.
  - Segment tree (range update, range query, more general aggregates) ⊂
    this family, superset of Fenwick's use cases.
  - Recognize the family from: "range sum/min/max," repeated with updates
    interleaved (if there are no updates, prefix sums alone are enough —
    don't reach for a segment tree).

- **Stack/deque invariant family** — maintain a monotonic (increasing or
  decreasing) sequence to answer "nearest element with property X" in
  amortized O(1) per element.
  - Monotonic stack (next/previous greater or smaller element) ⊂ this
    family.
  - Monotonic deque (sliding window max/min) ⊂ this family — same
    invariant, but elements can also fall off the *front* as the window
    slides.

## How to simplify an unfamiliar problem into a known pattern

1. **Strip the flavor text.** Identify: the input shape (array / string /
   tree / graph / matrix / stream), the constraints (n, value ranges,
   number of queries), and the question type (count / min / max / exists /
   construct / enumerate all).
2. **Read the constraints before designing anything.** They tell you the
   target complexity (see [complexity.md](./complexity.md)), which
   eliminates most candidate families immediately — if n ≤ 20, you're very
   likely looking at bitmask DP or backtracking, not a linear scan.
3. **Classify the input relationship.** Is position/order load-bearing
   (array, string → window/pointer or sequence DP), or is it a relationship
   structure (graph, tree → traversal family)?
4. **Check for monotonicity.** If "can we achieve X?" flips from false to
   true exactly once as some parameter increases, that parameter is
   binary-searchable — binary search on the answer turns an optimization
   problem into a decision problem you can check in isolation.
5. **Check for overlapping subproblems.** If brute force recomputes the
   same `(remaining_input, choices_so_far)` state repeatedly, it's DP — the
   next step is picking the state shape (see the DP family above).
6. **Look for a composition, not a single primitive.** Most "creative"
   problems chain two ordinary patterns: sliding window + hashmap frequency
   count, binary search on the answer + greedy feasibility check, DFS +
   memoization (= tree DP), sort + two pointers. If no single pattern fits
   cleanly, ask which two might combine.
7. **When stuck, measure the gap.** Write down the brute-force complexity
   and the complexity the constraints demand. The ratio between them tells
   you what to add: brute force O(n^2) but n = 10^5 demands O(n log n) →
   you're missing a sort, a hashmap, or a binary search that removes one
   factor of n.

## Further reading

[AlgoMonster's 48-pattern curriculum](https://algo.monster/) is built on the
same idea at interview scale — DFS/BFS, two pointers, and sliding window
alone reportedly cover a large share of FAANG-style interview problems. Use
it as an external cross-check on this list, not a replacement — the
supersets above are meant to generalize past any fixed pattern count.

See [glossary.md](./glossary.md) for term definitions and
[complexity.md](./complexity.md) for the constraint-driven complexity table.
