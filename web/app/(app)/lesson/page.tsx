"use client";

import { Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
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
import { findNextLesson, locateLessonBreadcrumb, type NextLessonResult } from "@/lib/curriculum";
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

interface NavInfo {
  location: { track: string; chapterTitle: string };
  next: NextLessonResult;
}

const LEARN_CRUMB: Crumb = { label: "Learn", href: "/learn" };

function NextPanel({
  next,
  trackSlug,
  emphasis,
}: {
  next: NextLessonResult;
  trackSlug: string;
  emphasis: "primary" | "secondary";
}) {
  return (
    <Card className="space-y-3">
      <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Next</h2>
      {next.kind === "next" ? (
        <LinkButton
          variant={emphasis === "primary" ? "primary" : "ghost"}
          href={`/lesson?topic=${encodeURIComponent(next.topic)}&concept=${encodeURIComponent(next.concept)}`}
          className="max-w-full"
        >
          {/* Concept titles are plain strings of unknown length/script — break-words
              lets the browser wrap inside the button instead of overflowing it. */}
          <span className="min-w-0 break-words">{next.title} →</span>
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
  const [retryToken, setRetryToken] = useState(0);
  const cardRefs = useRef<Record<number, HTMLDivElement | null>>({});
  const nextPanelRef = useRef<HTMLDivElement>(null);

  // One effect owns both the lesson fetch and the mark-in-progress write so
  // they share a single `cancelled` flag. Navigating lesson -> Next -> Back
  // -> Forward re-runs this effect without remounting the page (same route,
  // new searchParams) — without this guard a late response for the OLD
  // lesson could set state for whatever lesson is on screen now.
  useEffect(() => {
    if (!topicSlug || !conceptSlug) {
      setState({ status: "error", message: "Missing topic or concept in the URL." });
      return;
    }

    let cancelled = false;
    setState({ status: "loading" });
    setRatings({});
    setFinishState({ status: "idle" });
    setAlreadyPassed(false);

    async function run() {
      let lesson: Lesson;
      try {
        lesson = await getLesson(topicSlug, conceptSlug);
      } catch (err) {
        if (cancelled) return;
        if (err instanceof UnauthorizedError) {
          router.replace("/token");
          return;
        }
        if (err instanceof NotFoundError) {
          setState({ status: "not-found" });
          return;
        }
        setState({ status: "error", message: err instanceof Error ? err.message : "Something went wrong" });
        return;
      }
      if (cancelled) return;
      setState({ status: "success", lesson });

      try {
        const entry = await setProgress(lesson.topic, lesson.concept, "in_progress");
        if (cancelled) return;
        if (entry.state === "passed") setAlreadyPassed(true);
      } catch (err) {
        if (cancelled) return;
        if (err instanceof UnauthorizedError) router.replace("/token");
        // Any other failure here is safe to drop silently: the upsert's
        // COALESCE keeps first_passed_at correct regardless of whether this
        // visit gets recorded, and Finish still writes progress explicitly.
      }
    }

    run();
    return () => {
      cancelled = true;
    };
  }, [topicSlug, conceptSlug, router, retryToken]);

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

  const currentLesson = state.status === "success" ? state.lesson : null;

  const navInfo = useMemo<NavInfo | null>(() => {
    if (!tracks || !currentLesson) return null;
    const location = locateLessonBreadcrumb(tracks, currentLesson.topic, currentLesson.concept);
    if (!location) return null;
    return { location, next: findNextLesson(tracks, currentLesson.topic, currentLesson.concept) };
  }, [tracks, currentLesson]);

  const finished = alreadyPassed || finishState.status === "saved";

  // Finish unmounts on success, which would otherwise drop focus to <body>;
  // only the user's own click should steal focus, not a page load that
  // starts out already-passed, so this keys off finishState, not `finished`.
  useEffect(() => {
    if (finishState.status === "saved") nextPanelRef.current?.focus({ preventScroll: true });
  }, [finishState.status]);

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
          action={<Button onClick={() => setRetryToken((t) => t + 1)}>Retry</Button>}
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
  const crumbs: Crumb[] = navInfo
    ? [
        { label: trackLabel(navInfo.location.track), href: `/learn?track=${encodeURIComponent(navInfo.location.track)}` },
        { label: navInfo.location.chapterTitle },
        { label: lesson.title_en },
      ]
    : [LEARN_CRUMB, { label: lesson.title_en }];

  const totalChecks = lesson.recall_checks.length;
  const ratedCount = lesson.recall_checks.filter((check) => ratings[check.position] != null).length;
  const allRated = ratedCount === totalChecks;

  async function finish() {
    setFinishState({ status: "saving" });
    try {
      await setProgress(lesson.topic, lesson.concept, "passed");
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
      el?.focus({ preventScroll: true });
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
        <div className="space-y-2">
          {finished ? (
            <div
              role="status"
              className="flex items-center gap-2 rounded-xl bg-success/10 px-4 py-3 text-sm font-medium text-success-strong"
            >
              <CheckIcon />
              Lesson finished
            </div>
          ) : (
            <>
              <Button
                aria-disabled={!allRated || finishState.status === "saving"}
                aria-describedby={!allRated ? "finish-hint" : undefined}
                onClick={handleFinishClick}
              >
                {finishState.status === "saving" ? "Saving…" : "Finish lesson"}
              </Button>
              {!allRated && (
                <p id="finish-hint" className="text-xs text-muted">
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
            </>
          )}
        </div>

        <div ref={nextPanelRef} tabIndex={-1} className="rounded-2xl focus:outline-none focus:ring-2 focus:ring-focus focus:ring-offset-2">
          {navInfo ? (
            <NextPanel next={navInfo.next} trackSlug={navInfo.location.track} emphasis={finished ? "primary" : "secondary"} />
          ) : (
            finished && (
              <Card className="space-y-3">
                <h2 className="text-xs font-medium uppercase tracking-wide text-faint">Next</h2>
                <p className="text-sm text-muted">Couldn&apos;t figure out what&apos;s next right now.</p>
                <LinkButton variant="ghost" href="/learn">
                  Back to Learn
                </LinkButton>
              </Card>
            )
          )}
        </div>
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
