"use client";

import Link from "next/link";
import { useState } from "react";
import type { Chapter, Topic, Track } from "@/lib/api";
import { ChevronIcon } from "@/components/icons";

const TRACK_LABELS: Record<string, string> = {
  ddd: "Domain-Driven Design",
  distsys: "Distributed Systems",
  aws: "AWS (SAA-C03)",
  go: "Go",
  dsa: "DSA",
};

function ChapterRow({ topicSlug, chapter }: { topicSlug: string; chapter: Chapter }) {
  const [open, setOpen] = useState(false);
  const conceptsId = `chapter-${chapter.slug}-concepts`;

  return (
    <div className="rounded-xl border border-subtle bg-surface">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        aria-expanded={open}
        aria-controls={conceptsId}
        className="flex w-full flex-wrap items-center justify-between gap-x-3 gap-y-1 rounded-xl px-4 py-3 text-left focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
      >
        <span className="flex min-w-0 items-center gap-2 text-sm font-medium text-body">
          <ChevronIcon open={open} />
          <span className="break-words">{chapter.title}</span>
        </span>
        <span className="shrink-0 text-xs text-faint">{chapter.concepts.length} concepts</span>
      </button>

      {open && (
        <ul id={conceptsId} className="space-y-1 border-t border-subtle px-4 py-3 pl-9">
          {chapter.concepts.map((concept) => (
            <li key={concept.slug}>
              <Link
                href={`/lesson?topic=${encodeURIComponent(topicSlug)}&concept=${encodeURIComponent(concept.slug)}`}
                className="-mx-2 flex flex-wrap items-baseline justify-between gap-x-3 gap-y-0.5 rounded-lg px-2 py-1 text-sm transition-colors duration-150 ease-out hover:bg-accent-tint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
              >
                <span className="min-w-0 break-words text-body">{concept.title}</span>
                <span className="shrink-0 font-mono text-xs text-faint">
                  #{concept.position} · {concept.slug}
                </span>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

function TopicSection({ topic }: { topic: Topic }) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-semibold text-heading">{topic.title}</h3>
      <div className="space-y-2">
        {topic.chapters.map((chapter) => (
          <ChapterRow key={chapter.slug} topicSlug={topic.slug} chapter={chapter} />
        ))}
      </div>
    </div>
  );
}

export function CurriculumTree({ tracks }: { tracks: Track[] }) {
  return (
    <div className="space-y-10">
      {tracks.map((track) => (
        <section key={track.track} className="space-y-4">
          <h2 className="text-xs font-medium uppercase tracking-wide text-faint">
            {TRACK_LABELS[track.track] ?? track.track}
          </h2>
          <div className="space-y-6 pl-1">
            {track.topics.map((topic) => (
              <TopicSection key={topic.slug} topic={topic} />
            ))}
          </div>
        </section>
      ))}
    </div>
  );
}
