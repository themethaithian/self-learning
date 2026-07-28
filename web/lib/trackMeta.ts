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
