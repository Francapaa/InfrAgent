"use client";

import { useState } from "react";

const faqs = [
  {
    q: "How fast does InfrAgent react to incidents?",
    a: "The SDK checks your services every 30 seconds. Once an incident is reported, the AI agent processes it on its next tick cycle, meaning the maximum detection-to-action delay is about 30 seconds.",
  },
  {
    q: "What happens if the agent is temporarily offline?",
    a: "Events are stored in the database with a pending status. When the agent comes back online, it processes all queued events in batch. No incident is ever lost.",
  },
  {
    q: "Is my infrastructure data secure?",
    a: "Absolutely. API keys are hashed with bcrypt. All webhook communications are signed with HMAC-SHA256. Your credentials never travel in plain text and every action is cryptographically verified.",
  },
  {
    q: "Can I control what the AI is allowed to do?",
    a: "Yes. You define which actions the SDK can execute on your server. InfrAgent can only trigger actions you have explicitly authorized in your webhook handler.",
  },
  {
    q: "What services can InfrAgent monitor?",
    a: "Anything with a health endpoint or observable metric. APIs, databases, workers, cron jobs, custom services. You configure what to monitor in the SDK.",
  },
  {
    q: "Do I need to change my existing infrastructure?",
    a: "No. InfrAgent runs alongside your existing stack as a lightweight sidecar. Just install the SDK, point it at your services, and you're ready to go.",
  },
];

export default function FAQ() {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  return (
    <section id="faq" className="px-6 py-24 md:py-32">
      <div className="mx-auto max-w-3xl">
        <div className="text-center">
          <p className="text-sm font-medium tracking-wide text-primary font-mono uppercase">
            FAQ
          </p>
          <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
            Frequently asked questions
          </h2>
        </div>

        <div className="mt-12 flex flex-col divide-y divide-border">
          {faqs.map((faq, i) => (
            <div key={i}>
              <button
                onClick={() => setOpenIndex(openIndex === i ? null : i)}
                className="group flex w-full items-center justify-between py-5 text-left"
              >
                <span className="relative overflow-hidden pr-4 text-base font-medium text-foreground">
                  <span className="block transition-transform duration-300 ease-out group-hover:-translate-y-full">
                    {faq.q}
                  </span>
                  <span className="absolute inset-0 flex items-center translate-y-full transition-transform duration-300 ease-out group-hover:translate-y-0">
                    {faq.q}
                  </span>
                </span>
                <svg
                  width="20"
                  height="20"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  className={`shrink-0 text-muted-foreground transition-transform duration-200 ${
                    openIndex === i ? "rotate-180" : ""
                  }`}
                >
                  <polyline points="6 9 12 15 18 9" />
                </svg>
              </button>
              <div
                className={`overflow-hidden transition-all duration-200 ${
                  openIndex === i ? "max-h-40 pb-5" : "max-h-0"
                }`}
              >
                <p className="text-sm leading-relaxed text-muted-foreground">
                  {faq.a}
                </p>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
