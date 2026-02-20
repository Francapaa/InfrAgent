const steps = [
  {
    number: "01",
    title: "Install the SDK",
    description:
      "Drop our lightweight SDK into your server. It starts monitoring your services immediately with a single configuration call.",
    code: `sdk := infragent.New(
  "https://api.infragent.dev",
  "agent_key_xyz789",
  "whsec_secret456",
)
sdk.MonitorAndReport()`,
  },
  {
    number: "02",
    title: "SDK detects incidents",
    description:
      "Every 30 seconds, the SDK checks your services. Health endpoints, response times, resource usage. When something is wrong, it reports to InfrAgent.",
    code: `[CRITICAL] auth-service — down (503)
> Reporting event to InfrAgent...
> Event received: event-a1b2c3`,
  },
  {
    number: "03",
    title: "AI analyzes and decides",
    description:
      "InfrAgent's AI engine evaluates the incident context, past patterns, and severity to determine the optimal remediation action.",
    code: `> Analyzing incident context...
> Service: auth-service
> History: 3 restarts in 30 days
> Decision: restart auth-service`,
  },
  {
    number: "04",
    title: "Automatic remediation",
    description:
      "Via a cryptographically signed webhook, the action is sent to your server and executed. The service recovers — no pager, no human.",
    code: `> Executing: restart auth-service
> Webhook signed (HMAC-SHA256)
> Response: { "ok": true }
[RESOLVED] auth-service — healthy`,
  },
];

export default function HowItWorks() {
  return (
    <section
      id="how-it-works"
      className="border-t border-border bg-card px-6 py-24 md:py-32"
    >
      <div className="mx-auto max-w-7xl">
        <div className="mx-auto max-w-2xl text-center">
          <p className="text-sm font-medium tracking-wide text-primary font-mono uppercase">
            How It Works
          </p>
          <h2 className="mt-3 text-balance text-3xl font-bold tracking-tight text-foreground md:text-4xl">
            From incident to resolution in seconds
          </h2>
          <p className="mt-4 text-pretty text-lg leading-relaxed text-muted-foreground">
            A fully autonomous loop that keeps your services running without
            human intervention.
          </p>
        </div>

        <div className="mt-16 flex flex-col gap-12 lg:gap-16">
          {steps.map((step, i) => (
            <div
              key={step.number}
              className={`flex flex-col items-start gap-8 lg:flex-row lg:items-center ${
                i % 2 !== 0 ? "lg:flex-row-reverse" : ""
              }`}
            >
              {/* Text */}
              <div className="flex-1">
                <span className="font-mono text-sm font-bold text-primary">
                  {step.number}
                </span>
                <h3 className="mt-2 text-2xl font-bold text-foreground">
                  {step.title}
                </h3>
                <p className="mt-3 max-w-md text-base leading-relaxed text-muted-foreground">
                  {step.description}
                </p>
              </div>

              {/* Code block */}
              <div className="w-full flex-1 overflow-hidden rounded-lg border border-border bg-background">
                <div className="flex items-center gap-2 border-b border-border px-4 py-2.5">
                  <span className="h-2.5 w-2.5 rounded-full bg-red-500/70" />
                  <span className="h-2.5 w-2.5 rounded-full bg-yellow-500/70" />
                  <span className="h-2.5 w-2.5 rounded-full bg-green-500/70" />
                </div>
                <pre className="overflow-x-auto p-5 font-mono text-sm leading-relaxed text-muted-foreground">
                  <code>{step.code}</code>
                </pre>
              </div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
