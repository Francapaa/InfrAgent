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
            className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              Log In
            </span>
            <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Log In
            </span>
          </Link>
          <Link
            href="/onboarding"
            className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              Get Started
            </span>
            <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Get Started
            </span>
          </Link>
          <a
            href="#features"
            className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              Features
            </span>
            <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              Features
            </span>
          </a>
          <a
            href="#faq"
            className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              FAQ
            </span>
            <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              FAQ
            </span>
          </a>
        </div>

        <p className="text-xs text-muted-foreground">
          {"InfrAgent. All rights reserved."}
        </p>
      </div>
    </footer>
  );
}
