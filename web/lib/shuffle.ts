/**
 * Fisher-Yates shuffle. `rng` defaults to `Math.random` but accepts an
 * injected generator so callers (component lazy-state init, tests) can pin
 * the resulting order without special-casing test environments.
 *
 * `Array.prototype.sort(() => Math.random() - 0.5)` is not a fix for this:
 * a non-transitive comparator biases the result under V8's sort, and even a
 * "correct" comparator can't produce a uniform permutation.
 */
export function shuffleOptions<T>(options: readonly T[], rng: () => number = Math.random): T[] {
  const result = options.slice();
  for (let i = result.length - 1; i > 0; i--) {
    const j = Math.floor(rng() * (i + 1));
    [result[i], result[j]] = [result[j], result[i]];
  }
  return result;
}
