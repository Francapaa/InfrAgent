import Link from "next/link";

export default function Hero() {
  return (
    <section className="relative flex min-h-screen flex-col items-center justify-center overflow-hidden px-6 pt-20">
      <div className="pointer-events-none absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[600px] w-[600px] rounded-full bg-primary/5 blur-3xl" />

      <div className="relative z-10 mx-auto flex max-w-4xl flex-col items-center text-center">
        <div className="mb-8 inline-flex items-center gap-2 rounded-full border border-primary/30 bg-primary/10 px-4 py-1.5">
          <span className="h-2 w-2 rounded-full bg-primary animate-pulse" />
          <span className="text-xs font-medium tracking-wide text-primary font-mono uppercase">
            Autonomous Infrastructure Agent
          </span>
        </div>

        <h1 className="text-balance text-4xl font-bold leading-tight tracking-tight text-foreground sm:text-5xl md:text-6xl lg:text-7xl">
          Your infrastructure,{" "}
          <span className="text-primary">self-healing.</span>
        </h1>

        <p className="mt-6 max-w-2xl text-pretty text-lg leading-relaxed text-muted-foreground md:text-xl">
          InfrAgent monitors your servers 24/7, detects failures in real-time,
          and uses AI to decide and execute the fix — automatically.
          No human intervention required.
        </p>

        <div className="mt-10 flex flex-col items-center gap-4 sm:flex-row">
          <Link
            href="/onboarding"
            className="group relative h-11 overflow-hidden rounded-md bg-primary px-8 py-3 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              Start Monitoring
            </span>
            <span className="absolute inset-0 flex items-center justify-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Start Monitoring
            </span>
          </Link>
          <a
            href="#how-it-works"
            className="group relative flex h-11 items-center justify-center gap-2 overflow-hidden rounded-md border border-border px-8 py-3 text-sm font-medium text-foreground transition-colors hover:border-muted-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              See How It Works
            </span>
            <span className="absolute inset-0 flex items-center justify-center gap-2 translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              See How It Works
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <path d="M12 5v14" />
                <path d="m19 12-7 7-7-7" />
              </svg>
            </span>
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
              className="transition-transform group-hover:translate-y-0.5"
            >
              <path d="M12 5v14" />
              <path d="m19 12-7 7-7-7" />
            </svg>
          </a>
        </div>

        <div className="mt-16 w-full max-w-2xl overflow-hidden rounded-lg border border-border bg-card">
          <div className="flex items-center gap-2 border-b border-border px-4 py-3">
            <span className="h-3 w-3 rounded-full bg-red-500/70" />
            <span className="h-3 w-3 rounded-full bg-yellow-500/70" />
            <span className="h-3 w-3 rounded-full bg-green-500/70" />
            <span className="ml-4 text-xs text-muted-foreground font-mono">
              infragent-monitor
            </span>
          </div>
          <div className="p-5 font-mono text-sm leading-relaxed text-left">
            <p className="text-muted-foreground">
              <span className="text-primary">$</span> infragent status
            </p>
            <p className="mt-2 text-green-400">
              {"[OK]"} payments-api &mdash; healthy (200)
            </p>
            <p className="text-green-400">
              {"[OK]"} worker-service &mdash; healthy (200)
            </p>
            <p className="text-red-400">
              {"[CRITICAL]"} auth-service &mdash; down (503)
            </p>
            <p className="mt-2 text-primary">
              {">"} AI analyzing incident...
            </p>
            <p className="text-primary">
              {">"} Decision: restart auth-service
            </p>
            <p className="text-primary">
              {">"} Executing action via webhook...
            </p>
            <p className="mt-2 text-green-400">
              {"[RESOLVED]"} auth-service &mdash; healthy (200)
            </p>
          </div>
        </div>
      </div>
    </section>
  );
}
