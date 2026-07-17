# Constraints → Complexity → Algorithm

The fastest way to narrow down an approach before writing any code: read the
constraints, derive a target complexity, then pick from the small set of
algorithm families that fit that complexity.

## Why this works

Judges (LeetCode, Codeforces, USACO, ...) size `n` so that exactly one
complexity class survives the time limit. The constraint is a direct hint
about the intended algorithm — it is rarely accidental.

Rule of thumb: modern judges execute roughly **10^8 simple operations per
second** (C++-class constant factor; interpreted languages like Python are
often 5-20x slower, which is why LeetCode's limits are generous). Divide the
time limit by that rate to get an operation budget, then find the complexity
that fits `n` into the budget.

## The table (conservative upper bounds)

| n up to | Complexity budget | Typical technique |
|---|---|---|
| 10 | O(n!), O(n^7), O(n^6) | brute-force permutations, deeply nested loops |
| 20 | O(2^n · n), O(n^5) | bitmask DP / subset enumeration, meet-in-the-middle |
| 80 | O(n^4) | 4 nested loops, brute-force over 4 indices |
| 400 | O(n^3) | Floyd-Warshall, interval DP, triple nested loops |
| 7,500 | O(n^2) | nested loops, simple O(n) x O(n) DP, brute-force pairs |
| 7×10^4 | O(n√n) | sqrt decomposition, block-based algorithms |
| 5×10^5 | O(n log n) | sorting, heaps, binary search + scan, segment/BIT builds |
| 5×10^6 | O(n) | linear scan, two pointers, prefix sums, counting sort, hashing |
| 10^9 | O(log n), O(√n) | binary search, math/number theory, fast exponentiation |
| 10^18 | O(log n), O(log^2 n), O(1) | binary search on value (not index), matrix exponentiation, closed-form math |

Source: [USACO Guide — Time Complexity](https://usaco.guide/bronze/time-comp),
cross-checked against the standard [Codeforces 10^8-ops rule of thumb](https://codeforces.com/blog/entry/21772).

## LeetCode gut-check (looser, since limits are generous)

LeetCode's time limits assume a slower reference and tend to leave headroom.
A rougher mental table that's usually good enough mid-interview:

| n up to | Assume | Reach for |
|---|---|---|
| ~10-12 | exponential is fine | full permutations / brute force |
| ~20-24 | O(2^n) | bitmask DP, subsets, backtracking with pruning |
| ~500 | O(n^3) | triple loop DP, Floyd-Warshall |
| ~5,000 | O(n^2) | pairwise brute force, simple DP |
| ~10^5 | O(n log n) | sort + scan, heap, binary search, divide & conquer |
| ~10^6-10^7 | O(n) or O(n log n) | single pass, two pointers, hashing |
| ~10^8+ | O(log n) or O(1) | binary search on answer, math |

If two different constraints appear in the same problem (e.g. `n` and `m` for
a grid, or `n` and `q` for `n` elements / `q` queries), multiply them — the
real work is `O(n · m)` or `O(n · q)`, not `O(n)`.

## Common gotchas

- **Sum-of-constraints across test cases.** "1 ≤ T ≤ 1000 test cases, sum of
  n ≤ 2×10^5" means treat the *sum* as the real n, not the per-test-case
  bound — an O(n) solution run T times against the per-case bound would be
  O(n) against the sum, which is what the setter intended you to target.
- **Two-dimensional input.** An `n × m` grid has `n·m` cells; a "linear scan"
  over it is `O(n·m)`, not `O(n)`. Size your target complexity against the
  total element count, not a single dimension.
- **Memory limit implies a complexity bound too.** 256 MB ÷ 4 bytes/int ≈
  6×10^7 ints. If a DP table would need more cells than that, you need to
  compress a dimension (rolling array) or change approach — independent of
  whether the *time* budget allows it.
- **Amortized vs worst-case.** A dynamic array's `push_back`/`append` is
  amortized O(1) but worst-case O(n) on a single call; Union-Find with path
  compression + union by rank is amortized O(α(n)) ≈ O(1). Don't panic at a
  single expensive call if the amortized bound is what's proven.
- **Expected vs adversarial worst-case.** Hash map lookups are expected O(1)
  but worst-case O(n) under hash collisions; quickselect is expected O(n) but
  worst-case O(n^2) without randomized pivoting. On adversarial judges
  (Codeforces) this matters; on LeetCode it almost never does.
- **Recursion depth is a hidden constraint.** O(n) recursion with n ~ 10^5
  can blow the call stack before it blows the time limit. Either convert to
  iterative, or check the language's stack size (LeetCode's default is often
  too small for `n > ~10^4` deep recursion).
- **Constant factor still matters at the boundary.** `O(n log n)` with a
  segment tree carries a much bigger constant than `O(n log n)` from a single
  `sort()` call. Near the edge of a budget, prefer the simpler structure.

## Reverse lookup: complexity → algorithm family

Once you know the *target* complexity (from the table above), this narrows
which families are even candidates:

| Target | Reach for |
|---|---|
| O(1) | closed-form formula, precomputed/hashed lookup |
| O(log n) | binary search (on index or on the answer), balanced BST ops, fast exponentiation, GCD |
| O(√n) | prime factorization / primality by trial division, sqrt decomposition |
| O(n) | single pass, two pointers, sliding window, prefix sums, Kadane's, counting sort, BFS/DFS on O(n) edges |
| O(n log n) | sorting-based algorithms, heap (top-k / k-way merge), divide & conquer, segment tree / BIT builds, closest-pair |
| O(n^2) | brute-force pairs, simple 1D-state DP with O(n) transition, dense-graph Dijkstra (array, no heap), insertion sort |
| O(n^3) | Floyd-Warshall (all-pairs shortest path), naive matrix multiplication, interval DP (O(n^2) states × O(n) transition), 3D DP |
| O(2^n) | subset/bitmask DP, brute-force subset enumeration, meet-in-the-middle (2^(n/2) halves) |
| O(n · 2^n) | bitmask DP over states with an extra linear factor (e.g. TSP DP) |
| O(n!) | brute-force permutations (e.g. exact TSP without DP) |

See [patterns.md](./patterns.md) for how to map problem *phrasing* (not just
constraints) onto these families, and [glossary.md](./glossary.md) for term
definitions.
