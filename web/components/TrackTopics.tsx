"use client";

import Link from "next/link";
import { useState } from "react";
import type { Chapter, Concept, Topic } from "@/lib/api";
import { countAvailableLessons } from "@/lib/curriculum";
import { ChevronIcon } from "@/components/icons";

function ConceptRow({ topicSlug, concept, step }: { topicSlug: string; concept: Concept; step: number }) {
  const trailing = concept.has_lesson
    ? typeof concept.est_minutes === "number"
      ? `~${concept.est_minutes} min`
      : null
    : "No lesson yet";

  const rowClassName =
    "-mx-2 flex flex-wrap items-baseline justify-between gap-x-3 gap-y-0.5 rounded-lg px-2 py-1 text-sm";
  const content = (
    <>
      <span className="min-w-0 break-words text-body">
        <span className="mr-2 text-xs text-faint">{step}.</span>
        {concept.title}
      </span>
      {trailing && <span className="shrink-0 text-xs text-faint">{trailing}</span>}
    </>
  );

  if (!concept.has_lesson) {
    return <div className={rowClassName}>{content}</div>;
  }

  return (
    <Link
      href={`/lesson?topic=${encodeURIComponent(topicSlug)}&concept=${encodeURIComponent(concept.slug)}`}
      className={`${rowClassName} transition-colors duration-150 ease-out hover:bg-accent-tint focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2`}
    >
      {content}
    </Link>
  );
}

function ChapterRow({ topicSlug, chapter }: { topicSlug: string; chapter: Chapter }) {
  const [open, setOpen] = useState(false);
  const conceptsId = `chapter-${chapter.slug}-concepts`;
  const { available, total } = countAvailableLessons(chapter.concepts);

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
        <span className="shrink-0 text-xs text-faint">
          {available}/{total} lessons
        </span>
      </button>

      {open && (
        <ul id={conceptsId} className="space-y-1 border-t border-subtle px-4 py-3 pl-9">
          {chapter.concepts.map((concept, index) => (
            <li key={concept.slug}>
              <ConceptRow topicSlug={topicSlug} concept={concept} step={index + 1} />
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
