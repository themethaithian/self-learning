"use client";

import Link from "next/link";
import { useState } from "react";
import type { Chapter, Topic } from "@/lib/api";
import { ChevronIcon } from "@/components/icons";

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

function ChapterList({ topicSlug, chapters }: { topicSlug: string; chapters: Chapter[] }) {
  return (
    <div className="space-y-2">
      {chapters.map((chapter) => (
        <ChapterRow key={chapter.slug} topicSlug={topicSlug} chapter={chapter} />
      ))}
    </div>
  );
}

function TopicSection({ topic }: { topic: Topic }) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-semibold text-heading">{topic.title}</h3>
      <ChapterList topicSlug={topic.slug} chapters={topic.chapters} />
    </div>
  );
}

// A single-topic track's topic title duplicates the page's own heading (the
// curriculum currently models most tracks as one topic per book/subject), so
// that layer is skipped and its chapters render directly.
export function TrackTopics({ topics }: { topics: Topic[] }) {
  if (topics.length === 1) {
    return <ChapterList topicSlug={topics[0].slug} chapters={topics[0].chapters} />;
  }
  return (
    <div className="space-y-6">
      {topics.map((topic) => (
        <TopicSection key={topic.slug} topic={topic} />
      ))}
    </div>
  );
}
