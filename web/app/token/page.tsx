"use client";

import { useEffect, useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import { getToken, setToken } from "@/lib/token";
import { Card } from "@/components/Card";
import { Button } from "@/components/Button";

export default function TokenPage() {
  const router = useRouter();
  const [value, setValue] = useState("");

  useEffect(() => {
    if (getToken()) {
      router.replace("/today");
    }
  }, [router]);

  function handleSubmit(event: FormEvent) {
    event.preventDefault();
    const trimmed = value.trim();
    if (!trimmed) return;
    setToken(trimmed);
    router.push("/today");
  }

  return (
    <main className="flex min-h-screen items-center justify-center bg-page px-6">
      <Card className="w-full max-w-sm">
        <h1 className="text-xl font-semibold text-heading">Self-Improve</h1>
        <p className="mt-2 text-sm text-muted">Paste your API bearer token to continue.</p>

        <form onSubmit={handleSubmit} className="mt-6 space-y-4">
          <div>
            <label htmlFor="token" className="text-xs font-medium uppercase tracking-wide text-faint">
              Bearer token
            </label>
            <input
              id="token"
              type="password"
              autoComplete="off"
              spellCheck={false}
              value={value}
              onChange={(event) => setValue(event.target.value)}
              placeholder="changeme-bearer-token"
              className="mt-1 w-full rounded-xl border border-subtle bg-surface px-3 py-2 text-sm text-body focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-focus focus-visible:ring-offset-2"
            />
          </div>
          <Button type="submit" className="w-full" disabled={value.trim().length === 0}>
            Continue
          </Button>
        </form>
      </Card>
    </main>
  );
}
