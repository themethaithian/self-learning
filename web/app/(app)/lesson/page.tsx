"use client";

import { Suspense, useCallback, useEffect, useMemo, useRef, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import {
  getCurriculum,
  getLesson,
  postAttempt,
  setProgress,
  NotFoundError,
  UnauthorizedError,
  type Lesson,
  type Track,
} from "@/lib/api";
import { findNextLesson, locateLessonBreadcrumb, type NextLessonResult } from "@/lib/curriculum";
import { trackLabel } from "@/lib/trackMeta";
import { LessonBody } from "@/components/LessonBody";
import { RecallCheckCard, type AttemptSaveStatus, type CompletedAttempt } from "@/components/RecallCheckCard";
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

interface AttemptPayload extends CompletedAttempt {
  question: string;
}

// short_answer's Pass/Not yet can flip back and forth before the user
// settles on one — submitting every intermediate value would append one
// permanently-ordered row per flip-flop into an append-only table with no
// upsert, and Q-2c's "latest rating" read can't break ties within the same
// second (recall_attempts.created_at is TIMESTAMP, second precision).
// Debouncing collapses a burst of changes into the single settled value.
// mcq never needs this: its inputs freeze the instant it reaches reveal, so
// it always has exactly one value to submit and does so immediately.
//
// Not exported: Next.js's App Router restricts page.tsx to its recognised
// exports (default component, metadata, ...) — an arbitrary named export
// here fails the generated route type-check. Tests duplicate this literal.
const SHORT_ANSWER_SUBMIT_DEBOUNCE_MS = 1000;

const LEARN_CRUMB: Crumb = { label: "Learn", href: "/learn" };

function NextPanelContent({
  next,
  trackSlug,
  emphasis,
}: {
  next: NextLessonResult;
  trackSlug: string;
  emphasis: "primary" | "secondary";
}) {
  switch (next.kind) {
    case "next":
      return (
        <LinkButton
          variant={emphasis === "primary" ? "primary" : "ghost"}
          href={`/lesson?topic=${encodeURIComponent(next.topic)}&concept=${encodeURIComponent(next.concept)}`}
          className="max-w-full"
        >
          {/* Concept titles are plain strings of unknown length/script. */}
          <span className="min-w-0 break-words">{next.title} →</span>
        </LinkButton>
      );
    case "end-of-track":
    case "not-found":
      return (
        <>
          <p className="text-sm text-muted">You&apos;ve reached the end of the available lessons in this track.</p>
          <LinkButton variant="ghost" href={`/learn?track=${encodeURIComponent(trackSlug)}`}>
            Back to track
          </LinkButton>
        </>
      );
    default: {
      const exhaustive: never = next;
      return exhaustive;
    }
  }
}

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
      <NextPanelContent next={next} trackSlug={trackSlug} emphasis={emphasis} />
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
  const [finishedChecks, setFinishedChecks] = useState<Record<number, boolean>>({});
  const [finishState, setFinishState] = useState<FinishState>({ status: "idle" });
  const [alreadyPassed, setAlreadyPassed] = useState(false);
  const [retryToken, setRetryToken] = useState(0);
  const [attemptStatus, setAttemptStatus] = useState<Record<number, AttemptSaveStatus>>({});
  const cardRefs = useRef<Record<number, HTMLDivElement | null>>({});
  // Retry (and a debounce flush) re-sends the latest payload already
  // computed for this position — RecallCheckCard only reports a completed
  // attempt once (or once per short_answer correction), so neither must
  // require the user to redo the check to reconstruct it.
  const pendingAttemptsRef = useRef<Record<number, AttemptPayload>>({});
  // Pending short_answer submissions waiting out the debounce window (see
  // SHORT_ANSWER_SUBMIT_DEBOUNCE_MS) — flushed early on Finish and on
  // navigating away so a settled answer is never lost to a page transition.
  const debounceTimersRef = useRef<Record<number, ReturnType<typeof setTimeout>>>({});
  const finishedBannerRef = useRef<HTMLDivElement>(null);
  // The single source of truth for "which lesson is actually on screen right
  // now" — finish() is owned by a click handler, not the effect below, so it
  // can't rely on that effect's own `cancelled` flag (that guard only covers
  // async work the effect itself started). Cleared the instant navigation
  // starts, (re)set once the new lesson is confirmed current.
  const currentIdentityRef = useRef<{ topic: string; concept: string } | null>(null);
  // Identifies a *visit*, not a lesson: revisiting the same topic/concept
  // bumps this too, unlike currentIdentityRef which is only slug-keyed and
  // so cannot tell "this response belongs to the visit that started it"
  // apart from "this response's slugs happen to match what's on screen now"
  // — a real gap when the user leaves and comes back to the same lesson
  // while an earlier visit's request is still in flight.
  const loadGenerationRef = useRef(0);

  const handleCheckFinishedChange = useCallback((position: number, finished: boolean) => {
    setFinishedChecks((prev) => (prev[position] === finished ? prev : { ...prev, [position]: finished }));
  }, []);

  const currentLesson = state.status === "success" ? state.lesson : null;

  // A stale response is one whose *visit* has ended, not one whose slugs no
  // longer match the screen — see loadGenerationRef above.
  const submitAttempt = useCallback(
    async (position: number, payload: AttemptPayload, opts?: { keepalive?: boolean }) => {
      pendingAttemptsRef.current[position] = payload;
      const identity = currentIdentityRef.current;
      if (!identity) {
        setAttemptStatus((prev) => ({ ...prev, [position]: "error" }));
        return;
      }
      const generation = loadGenerationRef.current;
      const stillCurrent = () => loadGenerationRef.current === generation;

      setAttemptStatus((prev) => ({ ...prev, [position]: "saving" }));

      try {
        const body = {
          question: payload.question,
          confidence: payload.confidence,
          outcome: payload.outcome,
          selected_option: payload.selectedOption,
        };
        // keepalive is only threaded through as an actual 4th argument when
        // requested (pagehide/visibilitychange) — every other call keeps the
        // exact 3-argument shape callers already assert on.
        const result = opts?.keepalive
          ? await postAttempt(identity.topic, identity.concept, body, { keepalive: true })
          : await postAttempt(identity.topic, identity.concept, body);
        if (!stillCurrent()) return;
        setAttemptStatus((prev) => ({ ...prev, [position]: result.kind === "ok" ? "saved" : "error" }));
      } catch (err) {
        if (err instanceof UnauthorizedError) router.replace("/token");
      }
    },
    [router],
  );

  const clearPendingTimer = useCallback((position: number) => {
    const timer = debounceTimersRef.current[position];
    if (timer === undefined) return;
    clearTimeout(timer);
    delete debounceTimersRef.current[position];
  }, []);

  // Unconditional: called both for an actual pending timer (quiet period,
  // Finish, navigate-away) and for Retry, which has no timer of its own but
  // must still cancel one — a short_answer correction made after a failure
  // schedules a fresh timer while the old "error" status (and its Retry
  // button) is still on screen, and Retry firing alongside that timer later
  // is two submissions of two different values under one check_key.
  const flushPendingAttempt = useCallback(
    (position: number, opts?: { keepalive?: boolean }) => {
      clearPendingTimer(position);
      const payload = pendingAttemptsRef.current[position];
      if (payload) void submitAttempt(position, payload, opts);
    },
    [clearPendingTimer, submitAttempt],
  );

  const flushAllPendingAttempts = useCallback(
    (opts?: { keepalive?: boolean }) => {
      Object.keys(debounceTimersRef.current).forEach((key) => flushPendingAttempt(Number(key), opts));
    },
    [flushPendingAttempt],
  );

  const handleAttemptReady = useCallback(
    (position: number, attempt: CompletedAttempt) => {
      const check = currentLesson?.recall_checks.find((c) => c.position === position);
      if (!check) {
        // The position this claims to belong to isn't in the lesson
        // currently on screen — nothing valid to retry, but silently
        // dropping it would leave no trace that anything went wrong.
        setAttemptStatus((prev) => ({ ...prev, [position]: "error" }));
        return;
      }
      const payload: AttemptPayload = { question: check.question, ...attempt };
      pendingAttemptsRef.current[position] = payload;

      if (check.type !== "short_answer") {
        void submitAttempt(position, payload);
        return;
      }

      clearPendingTimer(position);
      debounceTimersRef.current[position] = setTimeout(() => {
        delete debounceTimersRef.current[position];
        void submitAttempt(position, payload);
      }, SHORT_ANSWER_SUBMIT_DEBOUNCE_MS);
    },
    [currentLesson, submitAttempt, clearPendingTimer],
  );

  // Routed through flushPendingAttempt, not a bare submitAttempt call: a
  // short_answer correction made after a failure schedules a fresh debounce
  // timer while the old "error" status (and Retry button) is still on
  // screen — clicking Retry must cancel that timer, not race it, or the
  // timer's later fire and Retry's own submission both land as two rows.
  const handleRetrySave = useCallback(
    (position: number) => {
      if (attemptStatus[position] === "saving") return;
      if (!pendingAttemptsRef.current[position]) return;
      flushPendingAttempt(position);
    },
    [attemptStatus, flushPendingAttempt],
  );

  // One effect owns both the lesson fetch and the mark-in-progress write so
  // they share a single `cancelled` flag. Navigating lesson -> Next -> Back
  // -> Forward re-runs this effect without remounting the page (same route,
  // new searchParams) — without this guard a late response for the OLD
  // lesson could set state for whatever lesson is on screen now.
  useEffect(() => {
    loadGenerationRef.current += 1;

    if (!topicSlug || !conceptSlug) {
      setState({ status: "error", message: "Missing topic or concept in the URL." });
      return;
    }

    let cancelled = false;
    currentIdentityRef.current = null;
    setState({ status: "loading" });
    setFinishedChecks({});
    setFinishState({ status: "idle" });
    setAlreadyPassed(false);
    setAttemptStatus({});
    pendingAttemptsRef.current = {};
    debounceTimersRef.current = {};

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
      currentIdentityRef.current = { topic: lesson.topic, concept: lesson.concept };

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
      // A settled short_answer rating waiting out its debounce window must
      // not be lost just because the user moved on before it fired.
      flushAllPendingAttempts();
    };
  }, [topicSlug, conceptSlug, router, retryToken, flushAllPendingAttempts]);

  // A hard navigation (reload, close tab, switch app on mobile) doesn't run
  // this component's own cleanup the way an in-SPA lesson change does —
  // pagehide + visibilitychange:hidden is the reliable pair for catching
  // that (beforeunload alone is unreliable on mobile Safari, the reviewing
  // device). keepalive keeps the request alive past unload; sendBeacon can't
  // carry the Authorization header, so fetch(...,{keepalive:true}) is used
  // instead (see api.ts's postAttempt).
  useEffect(() => {
    function flushForUnload() {
      flushAllPendingAttempts({ keepalive: true });
    }
    function handleVisibilityChange() {
      if (document.visibilityState === "hidden") flushForUnload();
    }
    document.addEventListener("visibilitychange", handleVisibilityChange);
    window.addEventListener("pagehide", flushForUnload);
    return () => {
      document.removeEventListener("visibilitychange", handleVisibilityChange);
      window.removeEventListener("pagehide", flushForUnload);
    };
  }, [flushAllPendingAttempts]);

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
  // Focus goes to the banner itself, not the Next panel below it — a screen
  // reader announces whatever receives focus, and that's more reliable than
  // counting on a freshly-mounted role="status" region being read at all.
  useEffect(() => {
    if (finishState.status === "saved") finishedBannerRef.current?.focus({ preventScroll: true });
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
  const finishedCount = lesson.recall_checks.filter((check) => finishedChecks[check.position]).length;
  const allFinished = finishedCount === totalChecks;
  const attemptStatusValues = Object.values(attemptStatus);
  const savingAttemptCount = attemptStatusValues.filter((s) => s === "saving").length;
  const hasFailedAttempt = attemptStatusValues.includes("error");

  async function finish() {
    const identity = { topic: lesson.topic, concept: lesson.concept };
    const stillCurrent = () =>
      currentIdentityRef.current?.topic === identity.topic && currentIdentityRef.current?.concept === identity.concept;

    setFinishState({ status: "saving" });
    try {
      await setProgress(identity.topic, identity.concept, "passed");
      // The user may have navigated to a different lesson while this PUT was
      // in flight (Next is visible before Finish resolves, by design) — a
      // late response for the lesson we started this from must never mutate
      // whatever's on screen now.
      if (!stillCurrent()) return;
      setFinishState({ status: "saved" });
    } catch (err) {
      if (!stillCurrent()) return;
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
    if (!allFinished) {
      const firstUnfinished = lesson.recall_checks.find((check) => !finishedChecks[check.position]);
      const el = firstUnfinished ? cardRefs.current[firstUnfinished.position] : null;
      el?.scrollIntoView({ behavior: "smooth", block: "center" });
      el?.focus({ preventScroll: true });
      return;
    }
    // A short_answer rating chosen just before clicking Finish may still be
    // sitting in its debounce window — Finish must not race it.
    flushAllPendingAttempts();
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
                onFinishedChange={handleCheckFinishedChange}
                onAttemptReady={handleAttemptReady}
                saveStatus={attemptStatus[check.position] ?? "idle"}
                onRetrySave={() => handleRetrySave(check.position)}
              />
            ))}
          </div>
        </section>
      )}

      <section className="max-w-[68ch] space-y-4">
        <div className="space-y-2">
          {savingAttemptCount > 0 && (
            // Finish only writes lesson_progress — it says nothing about
            // whether recall_attempts caught up, so "Lesson finished" alone
            // would read as fully recorded even while a save is still in
            // flight (its own resolution is silently dropped if the user
            // has since navigated away, per submitAttempt's stillCurrent).
            <div role="status" className="rounded-xl border border-subtle bg-page px-4 py-3 text-sm text-muted">
              Saving {savingAttemptCount} recall attempt{savingAttemptCount === 1 ? "" : "s"}…
            </div>
          )}
          {hasFailedAttempt && (
            <div
              role="alert"
              className="rounded-xl border border-danger/30 bg-danger/10 px-4 py-3 text-sm text-danger-strong"
            >
              Some recall attempts didn&apos;t save — retry them above, or your recall history for this lesson will be
              incomplete.
            </div>
          )}
          {finished ? (
            <div
              ref={finishedBannerRef}
              tabIndex={-1}
              role="status"
              className="flex items-center gap-2 rounded-xl bg-success/10 px-4 py-3 text-sm font-medium text-success-strong focus:outline-none focus:ring-2 focus:ring-focus focus:ring-offset-2"
            >
              <CheckIcon />
              Lesson finished
            </div>
          ) : (
            <>
              <Button
                aria-disabled={!allFinished || finishState.status === "saving"}
                aria-describedby={!allFinished ? "finish-hint" : undefined}
                onClick={handleFinishClick}
              >
                {finishState.status === "saving" ? "Saving…" : "Finish lesson"}
              </Button>
              {!allFinished && (
                <p id="finish-hint" className="text-xs text-muted">
                  Complete all {totalChecks} check{totalChecks === 1 ? "" : "s"} to finish ({finishedCount}/{totalChecks} done)
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
