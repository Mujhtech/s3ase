import * as React from "react";
import { Bar,BarChart,CartesianGrid,XAxis } from "recharts";

import {
Card,
CardContent,
CardDescription,
CardHeader,
CardTitle,
} from "~/components/ui/card";
import {
ChartConfig,
ChartContainer,
ChartTooltip,
ChartTooltipContent,
} from "~/components/ui/chart";
import { Usage } from "~/models/usage";

const chartConfig = {
  uploads: { label: "Uploads", color: "hsl(var(--chart-1))" },
  storage_bytes: { label: "Storage added", color: "hsl(var(--chart-2))" },
} satisfies ChartConfig;

function formatBytes(bytes: number) {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const index = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  return `${(bytes / 1024 ** index).toFixed(index === 0 ? 0 : 1)} ${units[index]}`;
}

export function UsageChart({ usage }: { usage: Usage }) {
  const [activeChart, setActiveChart] = React.useState<keyof typeof chartConfig>("uploads");

  return (
    <Card>
      <CardHeader className="flex flex-col items-stretch space-y-0 border-b p-0 sm:flex-row">
        <div className="flex flex-1 flex-col justify-center gap-1 px-6 py-5 sm:py-6">
          <CardTitle>Daily usage</CardTitle>
          <CardDescription>Uploads and storage added over the last 30 days</CardDescription>
        </div>
        <div className="flex">
          {(["uploads", "storage_bytes"] as const).map((chart) => (
            <button
              key={chart}
              type="button"
              data-active={activeChart === chart}
              className="relative z-30 flex flex-1 flex-col justify-center gap-1 border-t px-6 py-4 text-left even:border-l data-[active=true]:bg-muted/50 sm:border-l sm:border-t-0 sm:px-8 sm:py-6"
              onClick={() => setActiveChart(chart)}
            >
              <span className="text-xs text-muted-foreground">{chartConfig[chart].label}</span>
              <span className="text-lg font-bold leading-none sm:text-3xl">
                {chart === "uploads" ? usage.file_count.toLocaleString() : formatBytes(usage.storage_bytes)}
              </span>
            </button>
          ))}
        </div>
      </CardHeader>
      <CardContent className="px-2 sm:p-6">
        <ChartContainer config={chartConfig} className="aspect-auto h-[250px] w-full">
          <BarChart accessibilityLayer data={usage.daily} margin={{ left: 12, right: 12 }}>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={32}
              tickFormatter={(value) => new Date(`${value}T00:00:00Z`).toLocaleDateString("en-US", { month: "short", day: "numeric" })}
            />
            <ChartTooltip
              content={
                <ChartTooltipContent
                  className="w-[170px]"
                  labelFormatter={(value) => new Date(`${value}T00:00:00Z`).toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" })}
                  formatter={(value) => activeChart === "storage_bytes" ? formatBytes(Number(value)) : Number(value).toLocaleString()}
                />
              }
            />
            <Bar dataKey={activeChart} fill={`var(--color-${activeChart})`} />
          </BarChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
