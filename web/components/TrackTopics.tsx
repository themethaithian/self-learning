"use client";

import Link from "next/link";
import { useState } from "react";
import type { Chapter, Concept, Topic } from "@/lib/api";
import { chapterReadStats, conceptReadMarker, type ProgressByKey } from "@/lib/curriculum";
import { ChevronIcon, CheckIcon, InProgressIcon } from "@/components/icons";

interface ProgressProp {
  progress: ProgressByKey | null;
}

// Every row reserves this slot whether or not it renders a marker, so
// a track that mixes marked and unmarked concepts still keeps every title
// starting at the same x position — otherwise "not_started" rows (no
// marker) sit flush left while "passed"/"in_progress" rows indent one icon
// further, and the list goes ragged at narrow widths.
function MarkerSlot({ marker }: { marker: "passed" | "in_progress" | "none" }) {
  if (marker === "none") return <span data-testid="concept-marker-slot" className="inline-block h-4 w-4 shrink-0" />;
  if (marker === "passed") {
    return (
      <span
        data-testid="concept-marker-slot"
        role="img"
        aria-label="Read"
        className="inline-flex h-4 w-4 shrink-0 items-center text-success-strong"
      >
        <CheckIcon />
      </span>
    );
  }
  return (
    <span
      data-testid="concept-marker-slot"
      role="img"
      aria-label="In progress"
      className="inline-flex h-4 w-4 shrink-0 items-center text-warning-strong"
    >
      <InProgressIcon />
    </span>
  );
}

function ConceptRow({
  topicSlug,
  concept,
  step,
  progress,
}: { topicSlug: string; concept: Concept; step: number } & ProgressProp) {
  const trailing = concept.has_lesson
    ? typeof concept.est_minutes === "number" && concept.est_minutes > 0
      ? `~${concept.est_minutes} min`
      : null
    : "No lesson yet";

  const marker = conceptReadMarker(concept.has_lesson, progress, topicSlug, concept.slug);

  const rowClassName =
    "-mx-2 flex flex-wrap items-baseline justify-between gap-x-3 gap-y-0.5 rounded-lg px-2 py-1 text-sm";
  const titleClassName = concept.has_lesson ? "text-body" : "text-muted";
  const content = (
    <>
      <span className={`flex min-w-0 items-baseline gap-2 break-words ${titleClassName}`}>
        <span className="text-xs text-faint">{step}.</span>
        <MarkerSlot marker={marker} />
        {concept.title}
      </span>
      {trailing && <span className="shrink-0 text-xs text-faint">{trailing}</span>}
    </>
  );

  if (!concept.has_lesson) {
    return <div className={`${rowClassName} cursor-default`}>{content}</div>;
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

function ChapterRow({ topicSlug, chapter, progress }: { topicSlug: string; chapter: Chapter } & ProgressProp) {
  const [open, setOpen] = useState(false);
  const conceptsId = `chapter-${chapter.slug}-concepts`;
  const { read, withLesson, planned } = chapterReadStats(topicSlug, chapter.concepts, progress);

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
        {read !== null && withLesson > 0 ? (
          <span className="flex shrink-0 flex-col items-end text-xs leading-tight text-faint">
            <span>
              {read}/{withLesson} read
            </span>
            {withLesson < planned && (
              <span>
                {withLesson}/{planned} ready
              </span>
            )}
          </span>
        ) : (
          <span className="shrink-0 text-xs text-faint">
            {withLesson}/{planned} ready
          </span>
        )}
      </button>

      <ul id={conceptsId} hidden={!open} className="space-y-1 border-t border-subtle px-4 py-3 pl-9">
        {chapter.concepts.map((concept, index) => (
          <li key={concept.slug}>
            <ConceptRow topicSlug={topicSlug} concept={concept} step={index + 1} progress={progress} />
          </li>
        ))}
      </ul>
    </div>
  );
}

function ChapterList({ topicSlug, chapters, progress }: { topicSlug: string; chapters: Chapter[] } & ProgressProp) {
  return (
    <div className="space-y-2">
      {chapters.map((chapter) => (
        <ChapterRow key={chapter.slug} topicSlug={topicSlug} chapter={chapter} progress={progress} />
      ))}
    </div>
  );
}

function TopicSection({ topic, progress }: { topic: Topic } & ProgressProp) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-semibold text-heading">{topic.title}</h3>
      <ChapterList topicSlug={topic.slug} chapters={topic.chapters} progress={progress} />
    </div>
  );
}

// A single-topic track's topic title duplicates the page's own heading (the
// curriculum currently models most tracks as one topic per book/subject), so
// that layer is skipped and its chapters render directly.
export function TrackTopics({ topics, progress }: { topics: Topic[] } & ProgressProp) {
  if (topics.length === 1) {
    return <ChapterList topicSlug={topics[0].slug} chapters={topics[0].chapters} progress={progress} />;
  }
  return (
    <div className="space-y-6">
      {topics.map((topic) => (
        <TopicSection key={topic.slug} topic={topic} progress={progress} />
      ))}
    </div>
  );
}
