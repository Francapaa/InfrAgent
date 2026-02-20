const stats = [
  { value: "< 30s", label: "Detection Time" },
  { value: "24/7", label: "Autonomous Monitoring" },
  { value: "Zero", label: "Human Intervention" },
  { value: "99.9%", label: "Uptime Target" },
];

export default function Stats() {
  return (
    <section className="border-y border-border bg-card">
      <div className="mx-auto grid max-w-7xl grid-cols-2 md:grid-cols-4">
        {stats.map((stat, i) => (
          <div
            key={stat.label}
            className={`flex flex-col items-center justify-center px-6 py-10 ${
              i < stats.length - 1 ? "md:border-r md:border-border" : ""
            } ${i < 2 ? "border-b border-border md:border-b-0" : ""} ${
              i === 0 ? "border-r border-border md:border-r" : ""
            } ${i === 2 ? "border-r border-border md:border-r" : ""}`}
          >
            <span className="text-3xl font-bold tracking-tight text-primary md:text-4xl">
              {stat.value}
            </span>
            <span className="mt-2 text-sm text-muted-foreground">
              {stat.label}
            </span>
          </div>
        ))}
      </div>
    </section>
  );
}
