"use client";

import { useState } from "react";
import Link from "next/link";

export default function Navbar() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <nav className="fixed top-0 left-0 right-0 z-50 border-b border-border bg-background/80 backdrop-blur-md">
      <div className="mx-auto flex max-w-7xl items-center justify-between px-6 py-4">
        <Link href="/" className="flex items-center gap-2">
          <div className="flex h-8 w-8 items-center justify-center rounded-md bg-primary">
            <svg
              width="18"
              height="18"
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
          <span className="text-lg font-bold tracking-tight text-foreground">
            InfrAgent
          </span>
        </Link>

        <div className="hidden items-center gap-8 md:flex">
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
            href="#how-it-works"
            className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
              How It Works
            </span>
            <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
              How It Works
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

        <div className="hidden items-center gap-3 md:flex">
          <Link
            href="/login"
            className="group relative overflow-hidden rounded-md px-4 py-2 text-sm font-medium text-foreground transition-colors hover:bg-secondary"
          >
            <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-[140%]">
              Log In
            </span>
            <span className="absolute inset-0 flex items-center justify-center translate-y-full transition-transform duration-300 ease-out group-hover:-translate-y-0">
              Log In
            </span>
          </Link>
        </div>

        <button
          onClick={() => setMobileOpen(!mobileOpen)}
          className="flex h-10 w-10 items-center justify-center rounded-md text-foreground md:hidden"
          aria-label="Toggle menu"
        >
          {mobileOpen ? (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="18" y1="6" x2="6" y2="18" />
              <line x1="6" y1="6" x2="18" y2="18" />
            </svg>
          ) : (
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <line x1="3" y1="12" x2="21" y2="12" />
              <line x1="3" y1="6" x2="21" y2="6" />
              <line x1="3" y1="18" x2="21" y2="18" />
            </svg>
          )}
        </button>
      </div>

      {mobileOpen && (
        <div className="border-t border-border bg-background px-6 pb-6 pt-4 md:hidden">
          <div className="flex flex-col gap-4">
            <a
              href="#features"
              onClick={() => setMobileOpen(false)}
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
              href="#how-it-works"
              onClick={() => setMobileOpen(false)}
              className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
                How It Works
              </span>
              <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
                How It Works
              </span>
            </a>
            <a
              href="#faq"
              onClick={() => setMobileOpen(false)}
              className="group relative h-5 overflow-hidden text-sm text-muted-foreground transition-colors hover:text-foreground"
            >
              <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
                FAQ
              </span>
              <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
                FAQ
              </span>
            </a>
            <div className="flex flex-col gap-2 pt-4 border-t border-border">
              <Link
                href="/login"
                className="group relative overflow-hidden rounded-md px-4 py-2 text-center text-sm font-medium text-foreground transition-colors hover:bg-secondary"
              >
                <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
                  Log In
                </span>
                <span className="absolute inset-0 flex items-center justify-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
                  Log In
                </span>
              </Link>
              <Link
                href="/onboarding"
                className="group relative overflow-hidden rounded-md bg-primary px-4 py-2 text-center text-sm font-semibold text-primary-foreground transition-colors hover:bg-primary/90"
              >
                <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
                  Get Started
                </span>
                <span className="absolute inset-0 flex items-center justify-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
                  Get Started
                </span>
              </Link>
            </div>
          </div>
        </div>
      )}
    </nav>
  );
}
