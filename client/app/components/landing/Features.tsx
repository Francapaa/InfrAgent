const features = [
  {
    title: "Auto-Detection",
    description:
      "A lightweight SDK monitors your services every 30 seconds. CPU, memory, health endpoints — if something fails, InfrAgent knows instantly.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="12" cy="12" r="10" />
        <path d="M12 6v6l4 2" />
      </svg>
    ),
  },
  {
    title: "AI Decision Engine",
    description:
      "Powered by advanced LLMs, InfrAgent analyzes the context of each incident — service history, severity, dependencies — and decides the best remediation action.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M12 2a4 4 0 0 0-4 4c0 2 2 3 2 5h4c0-2 2-3 2-5a4 4 0 0 0-4-4z" />
        <line x1="10" y1="17" x2="14" y2="17" />
        <line x1="10" y1="20" x2="14" y2="20" />
        <line x1="11" y1="23" x2="13" y2="23" />
      </svg>
    ),
  },
  {
    title: "Instant Remediation",
    description:
      "Via secure webhooks with HMAC signatures, InfrAgent sends remediation commands directly to your server. Restart services, scale resources, or trigger custom scripts.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polygon points="13 2 3 14 12 14 11 22 21 10 12 10 13 2" />
      </svg>
    ),
  },
  {
    title: "Event-Driven Architecture",
    description:
      "Asynchronous event processing means no events are ever lost. If the agent is temporarily offline, events queue up in the database and get processed when it returns.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="22 12 18 12 15 21 9 3 6 12 2 12" />
      </svg>
    ),
  },
  {
    title: "Secure by Design",
    description:
      "API keys are hashed with bcrypt. Webhook payloads are signed with HMAC-SHA256. Your infrastructure actions are cryptographically verified end-to-end.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
      </svg>
    ),
  },
  {
    title: "Full Visibility",
    description:
      "A real-time dashboard shows every event, every AI decision, and every action taken. Complete audit trail for your ops team with zero blind spots.",
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="3" width="18" height="18" rx="2" ry="2" />
        <line x1="3" y1="9" x2="21" y2="9" />
        <line x1="9" y1="21" x2="9" y2="9" />
      </svg>
    ),
  },
];

export default function Features() {
  return (
    <section id="features" className="px-6 py-24 md:py-32">
      <div className="mx-auto max-w-7xl">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-sm font-medium tracking-wide text-primary font-mono uppercase">
            Capabilities
          </p>
          <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
            Everything your infra needs to stay alive
          </h2>
          <p className="mt-4 text-pretty text-lg leading-relaxed text-muted-foreground">
            From detection to remediation, InfrAgent handles the full incident
            lifecycle autonomously.
          </p>
        </div>

        <div className="mt-16 grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="group rounded-lg border border-border bg-card p-6 transition-colors hover:border-primary/40"
            >
              <div className="flex h-10 w-10 items-center justify-center rounded-md bg-primary/10 text-primary">
                {feature.icon}
              </div>
              <h3 className="mt-4 text-lg font-semibold text-foreground">
                {feature.title}
              </h3>
              <p className="mt-2 text-sm leading-relaxed text-muted-foreground">
                {feature.description}
              </p>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
