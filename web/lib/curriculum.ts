import type { Concept } from "./api";

export function countAvailableLessons(concepts: Concept[]): { available: number; total: number } {
  return {
    available: concepts.filter((concept) => concept.has_lesson).length,
    total: concepts.length,
  };
}
