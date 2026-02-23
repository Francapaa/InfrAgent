import Link from "next/link";

export default function CTA() {
  return (
    <section className="relative border-t border-border bg-card px-6 py-24 md:py-32 overflow-hidden">
      <div className="pointer-events-none absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 h-[600px] w-[600px] rounded-full bg-yellow-500/10 blur-3xl" />
      <div className="mx-auto max-w-3xl text-center relative z-10">
        <h2 className="text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
          Stop babysitting your servers.
        </h2>
        <p className="mt-4 text-pretty text-lg leading-relaxed text-muted-foreground">
          Let InfrAgent handle the 3 AM outages. Set it up once, and your
          infrastructure takes care of itself.
        </p>
        <div className="mt-10 flex flex-col items-center justify-center gap-4 sm:flex-row">
          <Link
            href="/onboarding"
            className="group relative h-11 overflow-hidden rounded-md bg-primary px-8 py-3 text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-[150%]">
              Get Started for Free
            </span>
            <span className="absolute inset-0 flex items-center justify-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Get Started for Free
            </span>
          </Link>
          <Link
            href="/login"
            className="group relative flex h-11 items-center justify-center gap-2 overflow-hidden rounded-md border border-border px-8 py-3 text-sm font-medium text-foreground transition-colors hover:border-muted-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-[150%]">
              Log In to Dashboard
            </span>
            <span className="absolute inset-0 flex items-center justify-center gap-2 translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Log In to Dashboard
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
                <line x1="5" y1="12" x2="19" y2="12" />
                <polyline points="12 5 19 12 12 19" />
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
              className="transition-transform group-hover:translate-x-0.5"
            >
              <line x1="5" y1="12" x2="19" y2="12" />
              <polyline points="12 5 19 12 12 19" />
            </svg>
          </Link>
        </div>
      </div>
    </section>
  );
}
