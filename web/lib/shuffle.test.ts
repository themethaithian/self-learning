import { describe, expect, it } from "vitest";
import { shuffleOptions } from "./shuffle";

describe("shuffleOptions", () => {
  it("keeps every source element, only reorders", () => {
    const input = ["a", "b", "c", "d"];
    const result = shuffleOptions(input, () => 0.42);
    expect([...result].sort()).toEqual([...input].sort());
  });

  it("never mutates the source array", () => {
    const input = ["a", "b", "c"];
    shuffleOptions(input, () => 0.9);
    expect(input).toEqual(["a", "b", "c"]);
  });

  it("is deterministic for a given rng sequence", () => {
    // Fisher-Yates with a constant rng()=0 always swaps the current index
    // with index 0 — for [A,B,C] that yields [B,C,A] (i=2 swaps idx2/idx0,
    // then i=1 swaps idx1/idx0). Pinning this exact permutation is what lets
    // component tests assert a known render order without a real RNG.
    expect(shuffleOptions(["A", "B", "C"], () => 0)).toEqual(["B", "C", "A"]);
  });

  it("distributes every option to every position near-uniformly (real Math.random)", () => {
    // A comparator-based "shuffle" (sort(() => Math.random() - 0.5)) or a
    // shuffle that only swaps a couple of elements would show up here as a
    // skewed cell far outside this band — "the order changed" alone can't
    // catch that, only a distribution check can.
    const options = ["A", "B", "C", "D"];
    const n = options.length;
    const runs = 20000;
    const counts = Array.from({ length: n }, () => Array(n).fill(0));

    for (let r = 0; r < runs; r++) {
      const shuffled = shuffleOptions(options);
      shuffled.forEach((opt, pos) => {
        counts[options.indexOf(opt)][pos]++;
      });
    }

    const expected = runs / n;
    const tolerance = expected * 0.15;
    for (const row of counts) {
      for (const cell of row) {
        expect(Math.abs(cell - expected)).toBeLessThan(tolerance);
      }
    }
  });
});
