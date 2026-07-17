# Glossary

Terse definitions for terms used in [complexity.md](./complexity.md) and
[patterns.md](./patterns.md).

- **Big-O / Big-Theta / Big-Omega** — O is an upper bound on growth rate
  ("at most this bad"), Ω is a lower bound ("at least this bad"), Θ is both
  ("exactly this bad" up to constants). Interview usage almost always means
  Θ even when people say "O."
- **Amortized complexity** — the average cost per operation over a
  worst-case *sequence* of operations, even if individual operations vary.
  Dynamic array `append`: O(1) amortized, O(n) worst case on the single call
  that triggers a resize.
- **Expected / average-case complexity** — the average over random inputs
  or randomized algorithm choices, distinct from worst-case. Hash map
  lookup: O(1) expected, O(n) worst case under adversarial collisions.
- **In-place vs auxiliary space** — in-place uses O(1) extra memory beyond
  the input; auxiliary space is memory used beyond the output, whether or
  not it's O(1).
- **Stable sort** — equal elements keep their relative input order (merge
  sort is stable; typical quicksort/heapsort are not).
- **Memoization vs tabulation** — memoization is top-down recursion with a
  cache (write the brute force, add a cache); tabulation is bottom-up,
  filling a DP table in dependency order. Same complexity, different
  control flow and stack usage.
- **Greedy** — commit to the locally-best choice at each step and never
  revisit it. Correct only when the problem has the exchange-argument /
  matroid-like property that local optimality implies global optimality;
  when unsure whether greedy is provably correct, DP is the safe fallback.
- **Divide and conquer** — split into independent subproblems, solve
  recursively, combine. Distinct from DP in that subproblems don't overlap
  (merge sort splits disjoint halves; DP revisits shared states).
- **Meet in the middle** — split an exponential search space (size 2^n) into
  two halves of size 2^(n/2), solve each independently, then combine —
  turns O(2^n) into O(2^(n/2)), typically for n up to ~40.
- **Two pointers vs sliding window** — two pointers usually means opposite
  ends converging toward each other (sorted-array pair sum); sliding window
  usually means same-direction pointers where one expands and one shrinks
  to maintain an invariant over a contiguous range. Both are instances of
  the window/pointer family in [patterns.md](./patterns.md).
- **Monotonic stack / deque** — a stack or deque kept strictly increasing or
  decreasing by popping elements that violate the order before pushing;
  answers "nearest element satisfying X" queries in amortized O(1) per
  element.
- **Prefix sum / difference array** — prefix sum precomputes cumulative
  sums for O(1) range-sum queries on a static array; a difference array is
  its dual, giving O(1) range *updates* at the cost of needing a prefix-sum
  pass to read any single value back out.
- **Fenwick tree (BIT)** — binary indexed tree; supports point update and
  prefix/range query in O(log n) with a small constant, less flexible but
  simpler than a segment tree.
- **Segment tree** — a tree over array ranges supporting range update and
  range query (sum, min, max, or any associative op) in O(log n); more
  general than a Fenwick tree, larger constant factor.
- **Union-Find (Disjoint Set Union / DSU)** — maintains a partition of
  elements into disjoint sets with near-O(1) `union` and `find`
  (amortized O(α(n)), α = inverse Ackermann, effectively constant) when
  using path compression and union by rank/size.
- **Topological sort** — a linear ordering of a DAG's nodes such that every
  edge points forward; undefined (and detectable) if the graph has a cycle.
- **Backtracking** — DFS over a decision tree that undoes a choice on
  reaching a dead end or invalid state, optionally pruning branches early
  to avoid full exponential exploration.
- **Bitmask** — represent a subset of up to ~20-24 items as an integer, one
  bit per item; enables O(1) subset membership tests and is the state
  representation for bitmask DP.
- **Binary search on the answer** — instead of searching for a value's
  index, binary search over the space of *possible answers*, using a
  feasibility check (usually O(n) or O(n log n)) at each candidate; valid
  whenever feasibility is monotonic in the candidate answer.
- **Heap / priority queue** — supports O(log n) insert and O(log n)
  extract-min/max, O(1) peek; the standard structure for "top k," "k-way
  merge," and Dijkstra's frontier.
- **Trie (prefix tree)** — a tree where each root-to-node path spells a
  prefix; supports O(length) prefix queries, used for autocomplete, word
  search, and prefix-based DP.
- **Kadane's algorithm** — O(n) DP for maximum-subarray-sum: at each index,
  either extend the running subarray or restart at the current element,
  whichever is larger.
- **Floyd's cycle detection (tortoise and hare)** — two pointers moving at
  different speeds through a linked structure; they meet iff there's a
  cycle, and further pointer arithmetic locates the cycle's start in O(n)
  time, O(1) space.
- **Dijkstra vs Bellman-Ford vs Floyd-Warshall** — Dijkstra: single-source
  shortest path, non-negative weights only, O((V+E) log V) with a heap.
  Bellman-Ford: single-source, handles negative weights and detects
  negative cycles, O(V·E). Floyd-Warshall: all-pairs shortest path, any
  weights (no negative cycles), O(V^3) — only viable when V is small (see
  the complexity table).
- **NP-hard (informal)** — no known polynomial-time algorithm; if a problem
  is a disguised version of a known NP-hard problem (TSP, subset sum,
  knapsack decision, graph coloring), expect the intended solution to be
  exponential-but-bounded (small n), a DP with pseudo-polynomial complexity
  (bounded by a value range, not just n), or an approximation — not an
  exact polynomial algorithm.
