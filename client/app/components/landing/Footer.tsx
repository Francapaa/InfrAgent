import Link from "next/link";

export default function Footer() {
  return (
    <footer className="border-t border-border px-6 py-10">
      <div className="mx-auto flex max-w-7xl flex-col items-center justify-between gap-6 md:flex-row">
        <div className="flex items-center gap-2">
          <div className="flex h-6 w-6 items-center justify-center rounded bg-primary">
            <svg
              width="12"
              height="12"
              viewBox="0 0 24 24"
              fill="none"
              stroke="hsl(var(--primary-foreground))"
              strokeWidth="2.5"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M12 2L2 7l10 5 10-5-10-5z" />
              <path d="M2 17l10 5 10-5" />
              <path d="M2 12l10 5 10-5" />
            </svg>
          </div>
          <span className="text-sm font-semibold text-foreground">
            InfrAgent
          </span>
        </div>

        <div className="flex items-center gap-6">
          <Link
            href="/login"
            className="text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            Log In
          </Link>
          <Link
            href="/onboarding"
            className="text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            Get Started
          </Link>
          <a
            href="#features"
            className="text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            Features
          </a>
          <a
            href="#faq"
            className="text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            FAQ
          </a>
        </div>

        <p className="text-xs text-muted-foreground">
          {"InfrAgent. All rights reserved."}
        </p>
      </div>
    </footer>
  );
}
