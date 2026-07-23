"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { getToken } from "@/lib/token";

export default function RootPage() {
  const router = useRouter();

  useEffect(() => {
    router.replace(getToken() ? "/read" : "/token");
  }, [router]);

  return null;
}
