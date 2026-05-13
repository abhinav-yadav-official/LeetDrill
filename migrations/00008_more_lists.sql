-- +goose Up
-- +goose StatementBegin

INSERT INTO problem_lists (slug, name, description, source_url, sort_order)
VALUES
  ('grind-169',
   'Grind 169',
   'Extended Grind 75 by Yang Shun — 14 weeks of structured interview prep.',
   'https://www.techinterviewhandbook.org/grind75',
   16),
  ('top-interview-150',
   'Top Interview 150',
   'Official LeetCode Top Interview 150 study plan.',
   'https://leetcode.com/studyplan/top-interview-150/',
   25),
  ('amazon-top-50',
   'Amazon Top 50',
   'Frequently asked problems in Amazon SDE interviews.',
   '',
   50),
  ('google-top-50',
   'Google Top 50',
   'Frequently asked problems in Google SWE interviews.',
   '',
   60)
ON CONFLICT (slug) DO UPDATE SET
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    source_url = EXCLUDED.source_url,
    sort_order = EXCLUDED.sort_order;

-- ─── Grind 169 ────────────────────────────────────────────────────────────────
WITH items(list_slug, position, section, leetcode_slug) AS (
  VALUES
    -- Array (20)
    ('grind-169',   1, 'Array', 'two-sum'),
    ('grind-169',   2, 'Array', 'best-time-to-buy-and-sell-stock'),
    ('grind-169',   3, 'Array', 'majority-element'),
    ('grind-169',   4, 'Array', 'contains-duplicate'),
    ('grind-169',   5, 'Array', 'maximum-subarray'),
    ('grind-169',   6, 'Array', 'insert-interval'),
    ('grind-169',   7, 'Array', 'merge-intervals'),
    ('grind-169',   8, 'Array', 'product-of-array-except-self'),
    ('grind-169',   9, 'Array', 'sort-colors'),
    ('grind-169',  10, 'Array', '3sum'),
    ('grind-169',  11, 'Array', 'move-zeroes'),
    ('grind-169',  12, 'Array', 'squares-of-a-sorted-array'),
    ('grind-169',  13, 'Array', 'remove-duplicates-from-sorted-array'),
    ('grind-169',  14, 'Array', 'merge-sorted-array'),
    ('grind-169',  15, 'Array', 'missing-number'),
    ('grind-169',  16, 'Array', 'maximum-product-subarray'),
    ('grind-169',  17, 'Array', 'trapping-rain-water'),
    ('grind-169',  18, 'Array', 'jump-game'),
    ('grind-169',  19, 'Array', 'jump-game-ii'),
    ('grind-169',  20, 'Array', 'subarray-sum-equals-k'),
    -- String (10)
    ('grind-169',  21, 'String', 'valid-anagram'),
    ('grind-169',  22, 'String', 'valid-palindrome'),
    ('grind-169',  23, 'String', 'longest-substring-without-repeating-characters'),
    ('grind-169',  24, 'String', 'longest-palindromic-substring'),
    ('grind-169',  25, 'String', 'find-all-anagrams-in-a-string'),
    ('grind-169',  26, 'String', 'minimum-window-substring'),
    ('grind-169',  27, 'String', 'roman-to-integer'),
    ('grind-169',  28, 'String', 'reverse-string'),
    ('grind-169',  29, 'String', 'longest-common-prefix'),
    ('grind-169',  30, 'String', 'is-subsequence'),
    -- Stack (7)
    ('grind-169',  31, 'Stack', 'valid-parentheses'),
    ('grind-169',  32, 'Stack', 'implement-queue-using-stacks'),
    ('grind-169',  33, 'Stack', 'evaluate-reverse-polish-notation'),
    ('grind-169',  34, 'Stack', 'min-stack'),
    ('grind-169',  35, 'Stack', 'decode-string'),
    ('grind-169',  36, 'Stack', 'daily-temperatures'),
    ('grind-169',  37, 'Stack', 'largest-rectangle-in-histogram'),
    -- Linked List (11)
    ('grind-169',  38, 'Linked List', 'merge-two-sorted-lists'),
    ('grind-169',  39, 'Linked List', 'linked-list-cycle'),
    ('grind-169',  40, 'Linked List', 'reverse-linked-list'),
    ('grind-169',  41, 'Linked List', 'middle-of-the-linked-list'),
    ('grind-169',  42, 'Linked List', 'palindrome-linked-list'),
    ('grind-169',  43, 'Linked List', 'lru-cache'),
    ('grind-169',  44, 'Linked List', 'remove-nth-node-from-end-of-list'),
    ('grind-169',  45, 'Linked List', 'reorder-list'),
    ('grind-169',  46, 'Linked List', 'add-two-numbers'),
    ('grind-169',  47, 'Linked List', 'copy-list-with-random-pointer'),
    ('grind-169',  48, 'Linked List', 'reverse-nodes-in-k-group'),
    -- Binary Tree (18)
    ('grind-169',  49, 'Binary Tree', 'invert-binary-tree'),
    ('grind-169',  50, 'Binary Tree', 'maximum-depth-of-binary-tree'),
    ('grind-169',  51, 'Binary Tree', 'balanced-binary-tree'),
    ('grind-169',  52, 'Binary Tree', 'validate-binary-search-tree'),
    ('grind-169',  53, 'Binary Tree', 'lowest-common-ancestor-of-a-binary-search-tree'),
    ('grind-169',  54, 'Binary Tree', 'lowest-common-ancestor-of-a-binary-tree'),
    ('grind-169',  55, 'Binary Tree', 'binary-tree-level-order-traversal'),
    ('grind-169',  56, 'Binary Tree', 'binary-tree-right-side-view'),
    ('grind-169',  57, 'Binary Tree', 'binary-tree-maximum-path-sum'),
    ('grind-169',  58, 'Binary Tree', 'construct-binary-tree-from-preorder-and-inorder-traversal'),
    ('grind-169',  59, 'Binary Tree', 'kth-smallest-element-in-a-bst'),
    ('grind-169',  60, 'Binary Tree', 'serialize-and-deserialize-binary-tree'),
    ('grind-169',  61, 'Binary Tree', 'diameter-of-binary-tree'),
    ('grind-169',  62, 'Binary Tree', 'symmetric-tree'),
    ('grind-169',  63, 'Binary Tree', 'subtree-of-another-tree'),
    ('grind-169',  64, 'Binary Tree', 'count-complete-tree-nodes'),
    ('grind-169',  65, 'Binary Tree', 'path-sum'),
    ('grind-169',  66, 'Binary Tree', 'binary-tree-zigzag-level-order-traversal'),
    -- Graph (14)
    ('grind-169',  67, 'Graph', 'flood-fill'),
    ('grind-169',  68, 'Graph', 'number-of-islands'),
    ('grind-169',  69, 'Graph', 'clone-graph'),
    ('grind-169',  70, 'Graph', 'rotting-oranges'),
    ('grind-169',  71, 'Graph', 'course-schedule'),
    ('grind-169',  72, 'Graph', 'pacific-atlantic-water-flow'),
    ('grind-169',  73, 'Graph', 'number-of-connected-components-in-an-undirected-graph'),
    ('grind-169',  74, 'Graph', 'graph-valid-tree'),
    ('grind-169',  75, 'Graph', 'accounts-merge'),
    ('grind-169',  76, 'Graph', 'course-schedule-ii'),
    ('grind-169',  77, 'Graph', 'word-ladder'),
    ('grind-169',  78, 'Graph', 'alien-dictionary'),
    ('grind-169',  79, 'Graph', 'network-delay-time'),
    ('grind-169',  80, 'Graph', 'min-cost-to-connect-all-points'),
    -- Heap (9)
    ('grind-169',  81, 'Heap', 'kth-largest-element-in-an-array'),
    ('grind-169',  82, 'Heap', 'task-scheduler'),
    ('grind-169',  83, 'Heap', 'top-k-frequent-elements'),
    ('grind-169',  84, 'Heap', 'find-median-from-data-stream'),
    ('grind-169',  85, 'Heap', 'merge-k-sorted-lists'),
    ('grind-169',  86, 'Heap', 'k-closest-points-to-origin'),
    ('grind-169',  87, 'Heap', 'last-stone-weight'),
    ('grind-169',  88, 'Heap', 'kth-largest-element-in-a-stream'),
    ('grind-169',  89, 'Heap', 'design-twitter'),
    -- Trie (3)
    ('grind-169',  90, 'Trie', 'implement-trie-prefix-tree'),
    ('grind-169',  91, 'Trie', 'design-add-and-search-words-data-structure'),
    ('grind-169',  92, 'Trie', 'word-search-ii'),
    -- Dynamic Programming (18)
    ('grind-169',  93, 'Dynamic Programming', 'climbing-stairs'),
    ('grind-169',  94, 'Dynamic Programming', 'coin-change'),
    ('grind-169',  95, 'Dynamic Programming', 'house-robber'),
    ('grind-169',  96, 'Dynamic Programming', 'decode-ways'),
    ('grind-169',  97, 'Dynamic Programming', 'longest-increasing-subsequence'),
    ('grind-169',  98, 'Dynamic Programming', 'unique-paths'),
    ('grind-169',  99, 'Dynamic Programming', 'partition-equal-subset-sum'),
    ('grind-169', 100, 'Dynamic Programming', 'word-break'),
    ('grind-169', 101, 'Dynamic Programming', 'house-robber-ii'),
    ('grind-169', 102, 'Dynamic Programming', 'coin-change-ii'),
    ('grind-169', 103, 'Dynamic Programming', 'target-sum'),
    ('grind-169', 104, 'Dynamic Programming', 'interleaving-string'),
    ('grind-169', 105, 'Dynamic Programming', 'edit-distance'),
    ('grind-169', 106, 'Dynamic Programming', 'best-time-to-buy-and-sell-stock-with-cooldown'),
    ('grind-169', 107, 'Dynamic Programming', 'longest-common-subsequence'),
    ('grind-169', 108, 'Dynamic Programming', 'distinct-subsequences'),
    ('grind-169', 109, 'Dynamic Programming', 'burst-balloons'),
    ('grind-169', 110, 'Dynamic Programming', 'regular-expression-matching'),
    -- Binary Search (9)
    ('grind-169', 111, 'Binary Search', 'binary-search'),
    ('grind-169', 112, 'Binary Search', 'first-bad-version'),
    ('grind-169', 113, 'Binary Search', 'search-in-rotated-sorted-array'),
    ('grind-169', 114, 'Binary Search', 'time-based-key-value-store'),
    ('grind-169', 115, 'Binary Search', 'find-minimum-in-rotated-sorted-array'),
    ('grind-169', 116, 'Binary Search', 'koko-eating-bananas'),
    ('grind-169', 117, 'Binary Search', 'search-a-2d-matrix'),
    ('grind-169', 118, 'Binary Search', 'capacity-to-ship-packages-within-d-days'),
    ('grind-169', 119, 'Binary Search', 'median-of-two-sorted-arrays'),
    -- Backtracking (11)
    ('grind-169', 120, 'Backtracking', 'generate-parentheses'),
    ('grind-169', 121, 'Backtracking', 'permutations'),
    ('grind-169', 122, 'Backtracking', 'subsets'),
    ('grind-169', 123, 'Backtracking', 'combination-sum'),
    ('grind-169', 124, 'Backtracking', 'combination-sum-ii'),
    ('grind-169', 125, 'Backtracking', 'subsets-ii'),
    ('grind-169', 126, 'Backtracking', 'letter-combinations-of-a-phone-number'),
    ('grind-169', 127, 'Backtracking', 'word-search'),
    ('grind-169', 128, 'Backtracking', 'palindrome-partitioning'),
    ('grind-169', 129, 'Backtracking', 'n-queens'),
    ('grind-169', 130, 'Backtracking', 'permutations-ii'),
    -- Bit Manipulation (8)
    ('grind-169', 131, 'Bit Manipulation', 'number-of-1-bits'),
    ('grind-169', 132, 'Bit Manipulation', 'missing-number'),
    ('grind-169', 133, 'Bit Manipulation', 'sum-of-two-integers'),
    ('grind-169', 134, 'Bit Manipulation', 'counting-bits'),
    ('grind-169', 135, 'Bit Manipulation', 'reverse-bits'),
    ('grind-169', 136, 'Bit Manipulation', 'single-number'),
    ('grind-169', 137, 'Bit Manipulation', 'reverse-integer'),
    ('grind-169', 138, 'Bit Manipulation', 'add-binary'),
    -- Math & Geometry (7)
    ('grind-169', 139, 'Math & Geometry', 'spiral-matrix'),
    ('grind-169', 140, 'Math & Geometry', 'rotate-image'),
    ('grind-169', 141, 'Math & Geometry', 'set-matrix-zeroes'),
    ('grind-169', 142, 'Math & Geometry', 'happy-number'),
    ('grind-169', 143, 'Math & Geometry', 'plus-one'),
    ('grind-169', 144, 'Math & Geometry', 'powx-n'),
    ('grind-169', 145, 'Math & Geometry', 'palindrome-number'),
    -- Greedy (6)
    ('grind-169', 146, 'Greedy', 'maximum-subarray'),
    ('grind-169', 147, 'Greedy', 'gas-station'),
    ('grind-169', 148, 'Greedy', 'hand-of-straights'),
    ('grind-169', 149, 'Greedy', 'partition-labels'),
    ('grind-169', 150, 'Greedy', 'valid-parenthesis-string'),
    ('grind-169', 151, 'Greedy', 'non-overlapping-intervals')
)
INSERT INTO problem_list_items (list_id, position, section, problem_id)
SELECT pl.id, items.position, items.section, p.id
FROM items
JOIN problem_lists pl ON pl.slug = items.list_slug
JOIN problems p ON p.leetcode_slug = items.leetcode_slug
ON CONFLICT (list_id, problem_id) DO NOTHING;

-- ─── LeetCode Top Interview 150 ───────────────────────────────────────────────
WITH items(list_slug, position, section, leetcode_slug) AS (
  VALUES
    -- Array / String (24)
    ('top-interview-150',   1, 'Array / String', 'merge-sorted-array'),
    ('top-interview-150',   2, 'Array / String', 'remove-element'),
    ('top-interview-150',   3, 'Array / String', 'remove-duplicates-from-sorted-array'),
    ('top-interview-150',   4, 'Array / String', 'remove-duplicates-from-sorted-array-ii'),
    ('top-interview-150',   5, 'Array / String', 'majority-element'),
    ('top-interview-150',   6, 'Array / String', 'rotate-array'),
    ('top-interview-150',   7, 'Array / String', 'best-time-to-buy-and-sell-stock'),
    ('top-interview-150',   8, 'Array / String', 'best-time-to-buy-and-sell-stock-ii'),
    ('top-interview-150',   9, 'Array / String', 'jump-game'),
    ('top-interview-150',  10, 'Array / String', 'jump-game-ii'),
    ('top-interview-150',  11, 'Array / String', 'h-index'),
    ('top-interview-150',  12, 'Array / String', 'product-of-array-except-self'),
    ('top-interview-150',  13, 'Array / String', 'gas-station'),
    ('top-interview-150',  14, 'Array / String', 'trapping-rain-water'),
    ('top-interview-150',  15, 'Array / String', 'roman-to-integer'),
    ('top-interview-150',  16, 'Array / String', 'integer-to-roman'),
    ('top-interview-150',  17, 'Array / String', 'length-of-last-word'),
    ('top-interview-150',  18, 'Array / String', 'longest-common-prefix'),
    ('top-interview-150',  19, 'Array / String', 'reverse-words-in-a-string'),
    ('top-interview-150',  20, 'Array / String', 'find-the-index-of-the-first-occurrence-in-a-string'),
    ('top-interview-150',  21, 'Array / String', 'zigzag-conversion'),
    -- Two Pointers (5)
    ('top-interview-150',  22, 'Two Pointers', 'valid-palindrome'),
    ('top-interview-150',  23, 'Two Pointers', 'is-subsequence'),
    ('top-interview-150',  24, 'Two Pointers', 'two-sum-ii-input-array-is-sorted'),
    ('top-interview-150',  25, 'Two Pointers', 'container-with-most-water'),
    ('top-interview-150',  26, 'Two Pointers', '3sum'),
    -- Sliding Window (4)
    ('top-interview-150',  27, 'Sliding Window', 'minimum-size-subarray-sum'),
    ('top-interview-150',  28, 'Sliding Window', 'longest-substring-without-repeating-characters'),
    ('top-interview-150',  29, 'Sliding Window', 'minimum-window-substring'),
    ('top-interview-150',  30, 'Sliding Window', 'substring-with-concatenation-of-all-words'),
    -- Matrix (5)
    ('top-interview-150',  31, 'Matrix', 'valid-sudoku'),
    ('top-interview-150',  32, 'Matrix', 'spiral-matrix'),
    ('top-interview-150',  33, 'Matrix', 'rotate-image'),
    ('top-interview-150',  34, 'Matrix', 'set-matrix-zeroes'),
    ('top-interview-150',  35, 'Matrix', 'game-of-life'),
    -- Hashmap (9)
    ('top-interview-150',  36, 'Hashmap', 'ransom-note'),
    ('top-interview-150',  37, 'Hashmap', 'isomorphic-strings'),
    ('top-interview-150',  38, 'Hashmap', 'word-pattern'),
    ('top-interview-150',  39, 'Hashmap', 'valid-anagram'),
    ('top-interview-150',  40, 'Hashmap', 'group-anagrams'),
    ('top-interview-150',  41, 'Hashmap', 'two-sum'),
    ('top-interview-150',  42, 'Hashmap', 'happy-number'),
    ('top-interview-150',  43, 'Hashmap', 'contains-duplicate-ii'),
    ('top-interview-150',  44, 'Hashmap', 'longest-consecutive-sequence'),
    -- Intervals (4)
    ('top-interview-150',  45, 'Intervals', 'summary-ranges'),
    ('top-interview-150',  46, 'Intervals', 'merge-intervals'),
    ('top-interview-150',  47, 'Intervals', 'insert-interval'),
    ('top-interview-150',  48, 'Intervals', 'minimum-number-of-arrows-to-burst-balloons'),
    -- Stack (5)
    ('top-interview-150',  49, 'Stack', 'valid-parentheses'),
    ('top-interview-150',  50, 'Stack', 'simplify-path'),
    ('top-interview-150',  51, 'Stack', 'min-stack'),
    ('top-interview-150',  52, 'Stack', 'evaluate-reverse-polish-notation'),
    ('top-interview-150',  53, 'Stack', 'basic-calculator'),
    -- Linked List (11)
    ('top-interview-150',  54, 'Linked List', 'linked-list-cycle'),
    ('top-interview-150',  55, 'Linked List', 'add-two-numbers'),
    ('top-interview-150',  56, 'Linked List', 'merge-two-sorted-lists'),
    ('top-interview-150',  57, 'Linked List', 'copy-list-with-random-pointer'),
    ('top-interview-150',  58, 'Linked List', 'reverse-linked-list-ii'),
    ('top-interview-150',  59, 'Linked List', 'reverse-nodes-in-k-group'),
    ('top-interview-150',  60, 'Linked List', 'remove-nth-node-from-end-of-list'),
    ('top-interview-150',  61, 'Linked List', 'remove-duplicates-from-sorted-list-ii'),
    ('top-interview-150',  62, 'Linked List', 'rotate-list'),
    ('top-interview-150',  63, 'Linked List', 'partition-list'),
    ('top-interview-150',  64, 'Linked List', 'lru-cache'),
    -- Binary Tree General (14)
    ('top-interview-150',  65, 'Binary Tree', 'maximum-depth-of-binary-tree'),
    ('top-interview-150',  66, 'Binary Tree', 'same-tree'),
    ('top-interview-150',  67, 'Binary Tree', 'invert-binary-tree'),
    ('top-interview-150',  68, 'Binary Tree', 'symmetric-tree'),
    ('top-interview-150',  69, 'Binary Tree', 'construct-binary-tree-from-preorder-and-inorder-traversal'),
    ('top-interview-150',  70, 'Binary Tree', 'construct-binary-tree-from-inorder-and-postorder-traversal'),
    ('top-interview-150',  71, 'Binary Tree', 'flatten-binary-tree-to-linked-list'),
    ('top-interview-150',  72, 'Binary Tree', 'path-sum'),
    ('top-interview-150',  73, 'Binary Tree', 'sum-root-to-leaf-numbers'),
    ('top-interview-150',  74, 'Binary Tree', 'binary-tree-maximum-path-sum'),
    ('top-interview-150',  75, 'Binary Tree', 'binary-search-tree-iterator'),
    ('top-interview-150',  76, 'Binary Tree', 'count-complete-tree-nodes'),
    ('top-interview-150',  77, 'Binary Tree', 'lowest-common-ancestor-of-a-binary-tree'),
    ('top-interview-150',  78, 'Binary Tree', 'diameter-of-binary-tree'),
    -- Binary Tree BFS (4)
    ('top-interview-150',  79, 'BFS', 'binary-tree-right-side-view'),
    ('top-interview-150',  80, 'BFS', 'average-of-levels-in-binary-tree'),
    ('top-interview-150',  81, 'BFS', 'binary-tree-level-order-traversal'),
    ('top-interview-150',  82, 'BFS', 'binary-tree-zigzag-level-order-traversal'),
    -- Binary Search Tree (3)
    ('top-interview-150',  83, 'Binary Search Tree', 'validate-binary-search-tree'),
    ('top-interview-150',  84, 'Binary Search Tree', 'kth-smallest-element-in-a-bst'),
    ('top-interview-150',  85, 'Binary Search Tree', 'recover-binary-search-tree'),
    -- Graph General (6)
    ('top-interview-150',  86, 'Graph', 'number-of-islands'),
    ('top-interview-150',  87, 'Graph', 'surrounded-regions'),
    ('top-interview-150',  88, 'Graph', 'clone-graph'),
    ('top-interview-150',  89, 'Graph', 'course-schedule'),
    ('top-interview-150',  90, 'Graph', 'course-schedule-ii'),
    ('top-interview-150',  91, 'Graph', 'pacific-atlantic-water-flow'),
    -- Graph BFS (3)
    ('top-interview-150',  92, 'Graph BFS', 'snakes-and-ladders'),
    ('top-interview-150',  93, 'Graph BFS', 'minimum-genetic-mutation'),
    ('top-interview-150',  94, 'Graph BFS', 'word-ladder'),
    -- Trie (3)
    ('top-interview-150',  95, 'Trie', 'implement-trie-prefix-tree'),
    ('top-interview-150',  96, 'Trie', 'design-add-and-search-words-data-structure'),
    ('top-interview-150',  97, 'Trie', 'word-search-ii'),
    -- Backtracking (7)
    ('top-interview-150',  98, 'Backtracking', 'letter-combinations-of-a-phone-number'),
    ('top-interview-150',  99, 'Backtracking', 'combinations'),
    ('top-interview-150', 100, 'Backtracking', 'permutations'),
    ('top-interview-150', 101, 'Backtracking', 'combination-sum'),
    ('top-interview-150', 102, 'Backtracking', 'n-queens-ii'),
    ('top-interview-150', 103, 'Backtracking', 'generate-parentheses'),
    ('top-interview-150', 104, 'Backtracking', 'word-search'),
    -- Divide & Conquer (2)
    ('top-interview-150', 105, 'Divide & Conquer', 'convert-sorted-array-to-binary-search-tree'),
    ('top-interview-150', 106, 'Divide & Conquer', 'sort-list'),
    -- Kadane's Algorithm (2)
    ('top-interview-150', 107, 'Kadane''s Algorithm', 'maximum-subarray'),
    ('top-interview-150', 108, 'Kadane''s Algorithm', 'maximum-sum-circular-subarray'),
    -- Binary Search (7)
    ('top-interview-150', 109, 'Binary Search', 'search-insert-position'),
    ('top-interview-150', 110, 'Binary Search', 'search-a-2d-matrix'),
    ('top-interview-150', 111, 'Binary Search', 'find-peak-element'),
    ('top-interview-150', 112, 'Binary Search', 'search-in-rotated-sorted-array'),
    ('top-interview-150', 113, 'Binary Search', 'find-first-and-last-position-of-element-in-sorted-array'),
    ('top-interview-150', 114, 'Binary Search', 'find-minimum-in-rotated-sorted-array'),
    ('top-interview-150', 115, 'Binary Search', 'median-of-two-sorted-arrays'),
    -- Heap (3)
    ('top-interview-150', 116, 'Heap', 'kth-largest-element-in-an-array'),
    ('top-interview-150', 117, 'Heap', 'ipo'),
    ('top-interview-150', 118, 'Heap', 'find-k-pairs-with-smallest-sums'),
    ('top-interview-150', 119, 'Heap', 'find-median-from-data-stream'),
    -- Bit Manipulation (6)
    ('top-interview-150', 120, 'Bit Manipulation', 'add-binary'),
    ('top-interview-150', 121, 'Bit Manipulation', 'reverse-bits'),
    ('top-interview-150', 122, 'Bit Manipulation', 'number-of-1-bits'),
    ('top-interview-150', 123, 'Bit Manipulation', 'single-number'),
    ('top-interview-150', 124, 'Bit Manipulation', 'bitwise-and-of-numbers-range'),
    -- Math (5)
    ('top-interview-150', 125, 'Math', 'palindrome-number'),
    ('top-interview-150', 126, 'Math', 'plus-one'),
    ('top-interview-150', 127, 'Math', 'powx-n'),
    ('top-interview-150', 128, 'Math', 'sqrtx'),
    ('top-interview-150', 129, 'Math', 'reverse-integer'),
    ('top-interview-150', 130, 'Math', 'max-points-on-a-line'),
    -- 1D DP (5)
    ('top-interview-150', 131, '1D DP', 'climbing-stairs'),
    ('top-interview-150', 132, '1D DP', 'house-robber'),
    ('top-interview-150', 133, '1D DP', 'word-break'),
    ('top-interview-150', 134, '1D DP', 'coin-change'),
    ('top-interview-150', 135, '1D DP', 'longest-increasing-subsequence'),
    -- Multidimensional DP (13)
    ('top-interview-150', 136, 'Multidimensional DP', 'triangle'),
    ('top-interview-150', 137, 'Multidimensional DP', 'minimum-path-sum'),
    ('top-interview-150', 138, 'Multidimensional DP', 'unique-paths-ii'),
    ('top-interview-150', 139, 'Multidimensional DP', 'longest-palindromic-substring'),
    ('top-interview-150', 140, 'Multidimensional DP', 'interleaving-string'),
    ('top-interview-150', 141, 'Multidimensional DP', 'edit-distance'),
    ('top-interview-150', 142, 'Multidimensional DP', 'best-time-to-buy-and-sell-stock-with-cooldown'),
    ('top-interview-150', 143, 'Multidimensional DP', 'coin-change-ii'),
    ('top-interview-150', 144, 'Multidimensional DP', 'target-sum'),
    ('top-interview-150', 145, 'Multidimensional DP', 'longest-common-subsequence'),
    ('top-interview-150', 146, 'Multidimensional DP', 'best-time-to-buy-and-sell-stock-iii'),
    ('top-interview-150', 147, 'Multidimensional DP', 'best-time-to-buy-and-sell-stock-iv'),
    ('top-interview-150', 148, 'Multidimensional DP', 'regular-expression-matching')
)
INSERT INTO problem_list_items (list_id, position, section, problem_id)
SELECT pl.id, items.position, items.section, p.id
FROM items
JOIN problem_lists pl ON pl.slug = items.list_slug
JOIN problems p ON p.leetcode_slug = items.leetcode_slug
ON CONFLICT (list_id, problem_id) DO NOTHING;

-- ─── Amazon Top 50 ────────────────────────────────────────────────────────────
WITH items(list_slug, position, section, leetcode_slug) AS (
  VALUES
    -- Array & Hashing (10)
    ('amazon-top-50',  1, 'Array & Hashing', 'two-sum'),
    ('amazon-top-50',  2, 'Array & Hashing', 'contains-duplicate'),
    ('amazon-top-50',  3, 'Array & Hashing', 'maximum-subarray'),
    ('amazon-top-50',  4, 'Array & Hashing', 'product-of-array-except-self'),
    ('amazon-top-50',  5, 'Array & Hashing', 'trapping-rain-water'),
    ('amazon-top-50',  6, 'Array & Hashing', 'group-anagrams'),
    ('amazon-top-50',  7, 'Array & Hashing', 'top-k-frequent-elements'),
    ('amazon-top-50',  8, 'Array & Hashing', 'longest-consecutive-sequence'),
    ('amazon-top-50',  9, 'Array & Hashing', 'subarray-sum-equals-k'),
    ('amazon-top-50', 10, 'Array & Hashing', 'sliding-window-maximum'),
    -- String (5)
    ('amazon-top-50', 11, 'String', 'longest-substring-without-repeating-characters'),
    ('amazon-top-50', 12, 'String', 'minimum-window-substring'),
    ('amazon-top-50', 13, 'String', 'valid-parentheses'),
    ('amazon-top-50', 14, 'String', 'longest-palindromic-substring'),
    ('amazon-top-50', 15, 'String', 'encode-and-decode-strings'),
    -- Linked List (5)
    ('amazon-top-50', 16, 'Linked List', 'merge-two-sorted-lists'),
    ('amazon-top-50', 17, 'Linked List', 'reorder-list'),
    ('amazon-top-50', 18, 'Linked List', 'add-two-numbers'),
    ('amazon-top-50', 19, 'Linked List', 'lru-cache'),
    ('amazon-top-50', 20, 'Linked List', 'merge-k-sorted-lists'),
    -- Trees (8)
    ('amazon-top-50', 21, 'Trees', 'binary-tree-level-order-traversal'),
    ('amazon-top-50', 22, 'Trees', 'maximum-depth-of-binary-tree'),
    ('amazon-top-50', 23, 'Trees', 'serialize-and-deserialize-binary-tree'),
    ('amazon-top-50', 24, 'Trees', 'lowest-common-ancestor-of-a-binary-tree'),
    ('amazon-top-50', 25, 'Trees', 'binary-tree-maximum-path-sum'),
    ('amazon-top-50', 26, 'Trees', 'binary-tree-right-side-view'),
    ('amazon-top-50', 27, 'Trees', 'validate-binary-search-tree'),
    ('amazon-top-50', 28, 'Trees', 'word-search-ii'),
    -- Graph (7)
    ('amazon-top-50', 29, 'Graph', 'number-of-islands'),
    ('amazon-top-50', 30, 'Graph', 'course-schedule'),
    ('amazon-top-50', 31, 'Graph', 'clone-graph'),
    ('amazon-top-50', 32, 'Graph', 'pacific-atlantic-water-flow'),
    ('amazon-top-50', 33, 'Graph', 'number-of-connected-components-in-an-undirected-graph'),
    ('amazon-top-50', 34, 'Graph', 'accounts-merge'),
    ('amazon-top-50', 35, 'Graph', 'min-cost-to-connect-all-points'),
    -- Heap (5)
    ('amazon-top-50', 36, 'Heap', 'kth-largest-element-in-an-array'),
    ('amazon-top-50', 37, 'Heap', 'top-k-frequent-elements'),
    ('amazon-top-50', 38, 'Heap', 'find-median-from-data-stream'),
    ('amazon-top-50', 39, 'Heap', 'k-closest-points-to-origin'),
    ('amazon-top-50', 40, 'Heap', 'task-scheduler'),
    -- Dynamic Programming (8)
    ('amazon-top-50', 41, 'Dynamic Programming', 'coin-change'),
    ('amazon-top-50', 42, 'Dynamic Programming', 'word-break'),
    ('amazon-top-50', 43, 'Dynamic Programming', 'unique-paths'),
    ('amazon-top-50', 44, 'Dynamic Programming', 'longest-increasing-subsequence'),
    ('amazon-top-50', 45, 'Dynamic Programming', 'jump-game'),
    ('amazon-top-50', 46, 'Dynamic Programming', 'edit-distance'),
    ('amazon-top-50', 47, 'Dynamic Programming', 'partition-equal-subset-sum'),
    ('amazon-top-50', 48, 'Dynamic Programming', 'house-robber'),
    -- Binary Search & Misc (4)
    ('amazon-top-50', 49, 'Binary Search', 'search-in-rotated-sorted-array'),
    ('amazon-top-50', 50, 'Binary Search', 'median-of-two-sorted-arrays')
)
INSERT INTO problem_list_items (list_id, position, section, problem_id)
SELECT pl.id, items.position, items.section, p.id
FROM items
JOIN problem_lists pl ON pl.slug = items.list_slug
JOIN problems p ON p.leetcode_slug = items.leetcode_slug
ON CONFLICT (list_id, problem_id) DO NOTHING;

-- ─── Google Top 50 ────────────────────────────────────────────────────────────
WITH items(list_slug, position, section, leetcode_slug) AS (
  VALUES
    -- Array & String (10)
    ('google-top-50',  1, 'Array & String', 'maximum-subarray'),
    ('google-top-50',  2, 'Array & String', 'trapping-rain-water'),
    ('google-top-50',  3, 'Array & String', 'minimum-window-substring'),
    ('google-top-50',  4, 'Array & String', 'longest-substring-without-repeating-characters'),
    ('google-top-50',  5, 'Array & String', 'product-of-array-except-self'),
    ('google-top-50',  6, 'Array & String', 'next-permutation'),
    ('google-top-50',  7, 'Array & String', 'merge-intervals'),
    ('google-top-50',  8, 'Array & String', 'container-with-most-water'),
    ('google-top-50',  9, 'Array & String', 'meeting-rooms-ii'),
    ('google-top-50', 10, 'Array & String', 'largest-rectangle-in-histogram'),
    -- String & DP (8)
    ('google-top-50', 11, 'String & DP', 'word-break'),
    ('google-top-50', 12, 'String & DP', 'longest-palindromic-substring'),
    ('google-top-50', 13, 'String & DP', 'decode-ways'),
    ('google-top-50', 14, 'String & DP', 'regular-expression-matching'),
    ('google-top-50', 15, 'String & DP', 'edit-distance'),
    ('google-top-50', 16, 'String & DP', 'word-break'),
    ('google-top-50', 17, 'String & DP', 'encode-and-decode-strings'),
    ('google-top-50', 18, 'String & DP', 'alien-dictionary'),
    -- Trees & Graphs (14)
    ('google-top-50', 19, 'Trees & Graphs', 'number-of-islands'),
    ('google-top-50', 20, 'Trees & Graphs', 'word-ladder'),
    ('google-top-50', 21, 'Trees & Graphs', 'course-schedule'),
    ('google-top-50', 22, 'Trees & Graphs', 'clone-graph'),
    ('google-top-50', 23, 'Trees & Graphs', 'serialize-and-deserialize-binary-tree'),
    ('google-top-50', 24, 'Trees & Graphs', 'binary-tree-maximum-path-sum'),
    ('google-top-50', 25, 'Trees & Graphs', 'lowest-common-ancestor-of-a-binary-tree'),
    ('google-top-50', 26, 'Trees & Graphs', 'pacific-atlantic-water-flow'),
    ('google-top-50', 27, 'Trees & Graphs', 'word-search-ii'),
    ('google-top-50', 28, 'Trees & Graphs', 'graph-valid-tree'),
    ('google-top-50', 29, 'Trees & Graphs', 'min-cost-to-connect-all-points'),
    ('google-top-50', 30, 'Trees & Graphs', 'network-delay-time'),
    ('google-top-50', 31, 'Trees & Graphs', 'accounts-merge'),
    ('google-top-50', 32, 'Trees & Graphs', 'swim-in-rising-water'),
    -- Dynamic Programming (10)
    ('google-top-50', 33, 'Dynamic Programming', 'longest-increasing-subsequence'),
    ('google-top-50', 34, 'Dynamic Programming', 'coin-change'),
    ('google-top-50', 35, 'Dynamic Programming', 'unique-paths'),
    ('google-top-50', 36, 'Dynamic Programming', 'maximum-product-subarray'),
    ('google-top-50', 37, 'Dynamic Programming', 'jump-game'),
    ('google-top-50', 38, 'Dynamic Programming', 'burst-balloons'),
    ('google-top-50', 39, 'Dynamic Programming', 'longest-common-subsequence'),
    ('google-top-50', 40, 'Dynamic Programming', 'partition-equal-subset-sum'),
    ('google-top-50', 41, 'Dynamic Programming', 'interleaving-string'),
    ('google-top-50', 42, 'Dynamic Programming', 'distinct-subsequences'),
    -- Heap & Misc (8)
    ('google-top-50', 43, 'Heap', 'find-median-from-data-stream'),
    ('google-top-50', 44, 'Heap', 'merge-k-sorted-lists'),
    ('google-top-50', 45, 'Heap', 'task-scheduler'),
    ('google-top-50', 46, 'Heap', 'top-k-frequent-elements'),
    ('google-top-50', 47, 'Binary Search', 'median-of-two-sorted-arrays'),
    ('google-top-50', 48, 'Binary Search', 'search-in-rotated-sorted-array'),
    ('google-top-50', 49, 'Backtracking', 'n-queens'),
    ('google-top-50', 50, 'Backtracking', 'word-search')
)
INSERT INTO problem_list_items (list_id, position, section, problem_id)
SELECT pl.id, items.position, items.section, p.id
FROM items
JOIN problem_lists pl ON pl.slug = items.list_slug
JOIN problems p ON p.leetcode_slug = items.leetcode_slug
ON CONFLICT (list_id, problem_id) DO NOTHING;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM problem_lists
WHERE slug IN ('grind-169', 'top-interview-150', 'amazon-top-50', 'google-top-50');
-- +goose StatementEnd
