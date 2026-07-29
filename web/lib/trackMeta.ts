export const TRACK_DISPLAY_ORDER = ["ddd", "distsys", "aws", "go", "dsa", "ddia", "ai-systems"] as const;

type KnownTrack = (typeof TRACK_DISPLAY_ORDER)[number];

const TRACK_LABELS: Record<KnownTrack, string> = {
  ddd: "Domain-Driven Design",
  distsys: "Distributed Systems",
  aws: "AWS (SAA-C03)",
  go: "Go",
  dsa: "DSA",
  ddia: "Designing Data-Intensive Applications",
  "ai-systems": "AI & LLM Systems",
};

export function trackLabel(track: string): string {
  return TRACK_LABELS[track as KnownTrack] ?? track;
}

function trackDisplayIndex(track: string): number {
  const idx = (TRACK_DISPLAY_ORDER as readonly string[]).indexOf(track);
  return idx === -1 ? Number.MAX_SAFE_INTEGER : idx;
}

// TRACK_DISPLAY_ORDER is a sort key only, never a filter — a track the
// backend sends that isn't listed here yet is appended, not dropped (UX-2's
// review: iterating a hardcoded allowlist instead of the data made unknown
// tracks vanish silently).
export function sortTracksByDisplayOrder<T extends { track: string }>(items: T[]): T[] {
  return [...items].sort((a, b) => trackDisplayIndex(a.track) - trackDisplayIndex(b.track));
}
