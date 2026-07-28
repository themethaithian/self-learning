import type { Concept } from "./api";

export function countAvailableLessons(concepts: Concept[]): { available: number; total: number } {
  return concepts.reduce(
    (acc, concept) => ({
      available: acc.available + (concept.has_lesson ? 1 : 0),
      total: acc.total + 1,
    }),
    { available: 0, total: 0 },
  );
}
