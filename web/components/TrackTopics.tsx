"use client";

import Link from "next/link";
import { useState } from "react";
import type { Chapter, Concept, ProgressState, Topic } from "@/lib/api";
import { chapterReadStats, conceptReadMarker, countAvailableLessons, progressKey, type ProgressByKey } from "@/lib/curriculum";
import { ChevronIcon, CheckIcon, InProgressIcon } from "@/components/icons";

interface ProgressProps {
  progressByKey: ProgressByKey;
  progressAvailable: boolean;
}

function ConceptMarker({ marker }: { marker: "passed" | "in_progress" }) {
  if (marker === "passed") {
    return (
      <span className="inline-flex shrink-0 items-center text-success-strong" aria-label="Read">
        <CheckIcon />
      </span>
    );
  }
  return (
    <span className="inline-flex shrink-0 items-center text-warning-strong" aria-label="In progress">
      <InProgressIcon />
    </span>
  );
}

function ConceptRow({
  topicSlug,
  concept,
  step,
  progressByKey,
  progressAvailable,
}: { topicSlug: string; concept: Concept; step: number } & ProgressProps) {
  const trailing = concept.has_lesson
    ? typeof concept.est_minutes === "number" && concept.est_minutes > 0
      ? `~${concept.est_minutes} min`
      : null
    : "No lesson yet";

  const state: ProgressState = progressByKey[progressKey(topicSlug, concept.slug)] ?? "not_started";
  const marker = progressAvailable ? conceptReadMarker(concept.has_lesson, state) : "none";

  const rowClassName =
    "-mx-2 flex flex-wrap items-baseline justify-between gap-x-3 gap-y-0.5 rounded-lg px-2 py-1 text-sm";
  const titleClassName = concept.has_lesson ? "text-body" : "text-muted";
  const content = (
    <>
      <span className={`flex min-w-0 items-baseline gap-2 break-words ${titleClassName}`}>
        <span className="text-xs text-faint">{step}.</span>
        {marker !== "none" && <ConceptMarker marker={marker} />}
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

function ChapterRow({ topicSlug, chapter, progressByKey, progressAvailable }: { topicSlug: string; chapter: Chapter } & ProgressProps) {
  const [open, setOpen] = useState(false);
  const conceptsId = `chapter-${chapter.slug}-concepts`;
  const { available, total } = countAvailableLessons(chapter.concepts);
  const readStats = chapterReadStats(topicSlug, chapter.concepts, progressByKey, progressAvailable);

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
        {readStats && readStats.available > 0 ? (
          <span className="flex shrink-0 flex-col items-end text-xs leading-tight text-faint">
            <span>
              {readStats.read}/{readStats.available} read
            </span>
            {readStats.available < readStats.total && (
              <span>
                {readStats.available}/{readStats.total} ready
              </span>
            )}
          </span>
        ) : (
          <span className="shrink-0 text-xs text-faint">
            {available}/{total} ready
          </span>
        )}
      </button>

      <ul
        id={conceptsId}
        className={`space-y-1 border-t border-subtle px-4 py-3 pl-9 ${open ? "" : "hidden"}`}
      >
        {chapter.concepts.map((concept, index) => (
          <li key={concept.slug}>
            <ConceptRow
              topicSlug={topicSlug}
              concept={concept}
              step={index + 1}
              progressByKey={progressByKey}
              progressAvailable={progressAvailable}
            />
          </li>
        ))}
      </ul>
    </div>
  );
}

function ChapterList({ topicSlug, chapters, progressByKey, progressAvailable }: { topicSlug: string; chapters: Chapter[] } & ProgressProps) {
  return (
    <div className="space-y-2">
      {chapters.map((chapter) => (
        <ChapterRow
          key={chapter.slug}
          topicSlug={topicSlug}
          chapter={chapter}
          progressByKey={progressByKey}
          progressAvailable={progressAvailable}
        />
      ))}
    </div>
  );
}

function TopicSection({ topic, progressByKey, progressAvailable }: { topic: Topic } & ProgressProps) {
  return (
    <div className="space-y-2">
      <h3 className="text-sm font-semibold text-heading">{topic.title}</h3>
      <ChapterList topicSlug={topic.slug} chapters={topic.chapters} progressByKey={progressByKey} progressAvailable={progressAvailable} />
    </div>
  );
}

// A single-topic track's topic title duplicates the page's own heading (the
// curriculum currently models most tracks as one topic per book/subject), so
// that layer is skipped and its chapters render directly.
export function TrackTopics({ topics, progressByKey, progressAvailable }: { topics: Topic[] } & ProgressProps) {
  if (topics.length === 1) {
    return (
      <ChapterList
        topicSlug={topics[0].slug}
        chapters={topics[0].chapters}
        progressByKey={progressByKey}
        progressAvailable={progressAvailable}
      />
    );
  }
  return (
    <div className="space-y-6">
      {topics.map((topic) => (
        <TopicSection key={topic.slug} topic={topic} progressByKey={progressByKey} progressAvailable={progressAvailable} />
      ))}
    </div>
  );
}
