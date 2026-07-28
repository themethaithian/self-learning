"use client";

import { Suspense, useCallback, useEffect, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import {
  getCurriculum,
  getLesson,
  setProgress,
  NotFoundError,
  UnauthorizedError,
  type Lesson,
  type Track,
} from "@/lib/api";
import { findNextLesson, locateLessonBreadcrumb, type NextLesson } from "@/lib/curriculum";
import { trackLabel } from "@/lib/trackMeta";
import { LessonBody } from "@/components/LessonBody";
import { RecallCheckCard, type RecallRating } from "@/components/RecallCheckCard";
import { EmptyState } from "@/components/EmptyState";
import { Button, LinkButton } from "@/components/Button";
import { Card } from "@/components/Card";
import { Breadcrumb, type Crumb } from "@/components/Breadcrumb";
import { LessonSkeleton } from "@/components/Skeleton";
import { BookIcon, CheckIcon, WarningIcon } from "@/components/icons";

type ViewState =
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "not-found" }
  | { status: "success"; lesson: Lesson };

type FinishState =
  | { status: "idle" }
  | { status: "saving" }
  | { status: "saved" }
  | { status: "error"; message: string };

const LEARN_CRUMB: Crumb = { label: "Learn", href: "/learn" };

function NextPanel({ next, trackSlug }: { next: NextLesson | null; trackSlug: string }) {
  return (
    <Card className="space-y-3">
      <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Next</h2>
      {next ? (
        <LinkButton
          variant="primary"
          href={`/lesson?topic=${encodeURIComponent(next.topic)}&concept=${encodeURIComponent(next.concept)}`}
        >
          {next.title} →
        </LinkButton>
      ) : (
        <>
          <p className="text-sm text-muted">You&apos;ve reached the end of the available lessons in this track.</p>
          <LinkButton variant="ghost" href={`/learn?track=${encodeURIComponent(trackSlug)}`}>
            Back to track
          </LinkButton>
        </>
      )}
    </Card>
  );
}

function LessonView() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const topicSlug = searchParams.get("topic") ?? "";
  const conceptSlug = searchParams.get("concept") ?? "";

  const [state, setState] = useState<ViewState>({ status: "loading" });
  const [tracks, setTracks] = useState<Track[] | null>(null);
  const [ratings, setRatings] = useState<Record<number, RecallRating>>({});
  const [finishState, setFinishState] = useState<FinishState>({ status: "idle" });
  const [alreadyPassed, setAlreadyPassed] = useState(false);
  const cardRefs = useRef<Record<number, HTMLDivElement | null>>({});

  const markInProgress = useCallback(
    async (topic: string, concept: string) => {
      try {
        const entry = await setProgress(topic, concept, "in_progress");
        if (entry.state === "passed") setAlreadyPassed(true);
      } catch (err) {
        if (err instanceof UnauthorizedError) router.replace("/token");
      }
    },
    [router],
  );

  const load = useCallback(async () => {
    if (!topicSlug || !conceptSlug) {
      setState({ status: "error", message: "Missing topic or concept in the URL." });
      return;
    }
    setState({ status: "loading" });
    setRatings({});
    setFinishState({ status: "idle" });
    setAlreadyPassed(false);
    try {
      const lesson = await getLesson(topicSlug, conceptSlug);
      setState({ status: "success", lesson });
      void markInProgress(topicSlug, conceptSlug);
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      if (err instanceof NotFoundError) {
        setState({ status: "not-found" });
        return;
      }
      setState({
        status: "error",
        message: err instanceof Error ? err.message : "Something went wrong",
      });
    }
  }, [topicSlug, conceptSlug, router, markInProgress]);

  useEffect(() => {
    load();
  }, [load]);

  // The curriculum tree only powers the breadcrumb and Next — its failure
  // must never turn a successfully loaded lesson into a page-level error.
  useEffect(() => {
    let cancelled = false;
    getCurriculum()
      .then((res) => {
        if (!cancelled) setTracks(res.tracks);
      })
      .catch((err) => {
        if (err instanceof UnauthorizedError) router.replace("/token");
      });
    return () => {
      cancelled = true;
    };
  }, [router]);

  const handleRate = useCallback((position: number, value: RecallRating) => {
    setRatings((prev) => ({ ...prev, [position]: value }));
  }, []);

  if (state.status === "loading") {
    return (
      <div className="space-y-6">
        <Breadcrumb items={[LEARN_CRUMB]} />
        <LessonSkeleton />
      </div>
    );
  }

  if (state.status === "error") {
    return (
      <div className="space-y-6">
        <Breadcrumb items={[LEARN_CRUMB]} />
        <EmptyState
          icon={<WarningIcon />}
          message={`Could not load this lesson: ${state.message}`}
          action={<Button onClick={load}>Retry</Button>}
        />
      </div>
    );
  }

  if (state.status === "not-found") {
    return (
      <div className="space-y-6">
        <Breadcrumb items={[LEARN_CRUMB]} />
        <EmptyState
          icon={<BookIcon />}
          message="No lesson yet for this concept. Most of the curriculum doesn't have a lesson written yet — check back later."
          action={
            <Button variant="ghost" onClick={() => router.push("/learn")}>
              Back to Learn
            </Button>
          }
        />
      </div>
    );
  }

  const { lesson } = state;
  const location = tracks ? locateLessonBreadcrumb(tracks, topicSlug, conceptSlug) : null;
  const next = tracks ? findNextLesson(tracks, topicSlug, conceptSlug) : null;
  const crumbs: Crumb[] = location
    ? [
        { label: trackLabel(location.track), href: `/learn?track=${encodeURIComponent(location.track)}` },
        { label: location.chapterTitle },
        { label: lesson.title_en },
      ]
    : [LEARN_CRUMB, { label: lesson.title_en }];

  const totalChecks = lesson.recall_checks.length;
  const ratedCount = lesson.recall_checks.filter((check) => ratings[check.position] != null).length;
  const allRated = ratedCount === totalChecks;
  const finished = alreadyPassed || finishState.status === "saved";

  async function finish() {
    setFinishState({ status: "saving" });
    try {
      await setProgress(topicSlug, conceptSlug, "passed");
      setFinishState({ status: "saved" });
    } catch (err) {
      if (err instanceof UnauthorizedError) {
        router.replace("/token");
        return;
      }
      setFinishState({
        status: "error",
        message: err instanceof Error ? err.message : "Something went wrong",
      });
    }
  }

  function handleFinishClick() {
    if (finishState.status === "saving") return;
    if (!allRated) {
      const firstUnrated = lesson.recall_checks.find((check) => ratings[check.position] == null);
      const el = firstUnrated ? cardRefs.current[firstUnrated.position] : null;
      el?.scrollIntoView({ behavior: "smooth", block: "center" });
      el?.focus();
      return;
    }
    void finish();
  }

  return (
    <div className="space-y-8">
      <Breadcrumb items={crumbs} />

      <div className="space-y-2">
        <h1 className="text-xl font-semibold text-heading">{lesson.title_en}</h1>
        <span className="inline-flex items-center rounded-full bg-accent-tint px-3 py-1 text-xs font-medium text-accent-strong">
          {lesson.est_minutes} min read
        </span>
      </div>

      <LessonBody markdown={lesson.body_md} />

      <section className="max-w-[68ch] space-y-3">
        <h2 className="text-xs font-medium uppercase tracking-wide text-faint">References</h2>
        <ul className="space-y-3">
          {lesson.references.map((ref) => (
            <li key={ref.title} className="rounded-xl border border-subtle bg-surface p-4">
              <p className="text-sm font-medium text-heading">{ref.title}</p>
              <p className="text-sm text-muted">{ref.source}</p>
              <p className="mt-1 font-thai text-sm leading-[1.8] text-body">{ref.why}</p>
            </li>
          ))}
        </ul>
      </section>

      {totalChecks > 0 && (
        <section className="max-w-[68ch] space-y-3">
          <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Recall check</h2>
          <div className="space-y-3">
            {lesson.recall_checks.map((check, index) => (
              <RecallCheckCard
                key={check.position}
                ref={(el) => {
                  cardRefs.current[check.position] = el;
                }}
                check={check}
                index={index}
                rating={ratings[check.position] ?? null}
                onRate={(value) => handleRate(check.position, value)}
              />
            ))}
          </div>
        </section>
      )}

      <section className="max-w-[68ch] space-y-4">
        {finished ? (
          <div className="flex items-center gap-2 rounded-xl bg-success/10 px-4 py-3 text-sm font-medium text-success">
            <CheckIcon />
            Lesson finished
          </div>
        ) : (
          <div className="space-y-2">
            <Button aria-disabled={!allRated || finishState.status === "saving"} onClick={handleFinishClick}>
              {finishState.status === "saving" ? "Saving…" : "Finish lesson"}
            </Button>
            {!allRated && (
              <p className="text-xs text-muted">
                Rate all {totalChecks} check{totalChecks === 1 ? "" : "s"} to finish ({ratedCount}/{totalChecks} rated)
              </p>
            )}
            {finishState.status === "error" && (
              <div
                role="alert"
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong"
              >
                <span>Could not save: {finishState.message}</span>
                <Button variant="ghost" onClick={handleFinishClick}>
                  Retry
                </Button>
              </div>
            )}
          </div>
        )}

        {finished && location && <NextPanel next={next} trackSlug={location.track} />}
      </section>
    </div>
  );
}

export default function LessonPage() {
  return (
    <Suspense
      fallback={
        <div className="space-y-6">
          <Breadcrumb items={[LEARN_CRUMB]} />
          <LessonSkeleton />
        </div>
      }
    >
      <LessonView />
    </Suspense>
  );
}
