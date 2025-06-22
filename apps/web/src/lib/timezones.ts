const fallbackTimezones = [
  "UTC",
  "Europe/London",
  "Europe/Paris",
  "America/New_York",
  "Asia/Tokyo",
  "Australia/Sydney",
];

function getSupportedTimezones(): string[] {
  if (
    typeof Intl !== "undefined" &&
    "supportedValuesOf" in Intl &&
    typeof (Intl as { supportedValuesOf?: unknown }).supportedValuesOf ===
      "function"
  ) {
    // @ts-expect-error: supportedValuesOf is not yet in TS types
    return Intl.supportedValuesOf("timeZone");
  }
  return fallbackTimezones;
}

const timezones = getSupportedTimezones();

function getTimezoneOffsetMinutes(timezone: string): number {
  try {
    const now = new Date();
    const tzDate = new Date(
      now.toLocaleString("en-US", { timeZone: timezone })
    );
    const utcDate = new Date(now.toLocaleString("en-US", { timeZone: "UTC" }));
    return (tzDate.getTime() - utcDate.getTime()) / 60000; // in minutes
  } catch {
    return 0;
  }
}

export function getTimezoneOffsetLabel(timezone: string): string {
  try {
    const now = new Date();
    const tzDate = new Date(now.toLocaleString('en-US', { timeZone: timezone }));
    const utcDate = new Date(now.toLocaleString('en-US', { timeZone: 'UTC' }));
    const diff = (tzDate.getTime() - utcDate.getTime()) / 60000; // in minutes
    const sign = diff >= 0 ? '+' : '-';
    const absDiff = Math.abs(diff);
    const hours = Math.floor(absDiff / 60);
    const minutes = absDiff % 60;
    if (diff === 0) return 'UTC±0';
    return `UTC${sign}${hours}${minutes ? ':' + String(minutes).padStart(2, '0') : ''}`;
  } catch {
    return 'UTC';
  }
}

export const sortedTimezones = [...timezones].sort((a, b) => {
  const offsetA = getTimezoneOffsetMinutes(a);
  const offsetB = getTimezoneOffsetMinutes(b);
  // Group by sign: negatives first, then zero, then positives
  if (offsetA < 0 && offsetB >= 0) return -1;
  if (offsetA >= 0 && offsetB < 0) return 1;
  if (offsetA === 0 && offsetB !== 0) return -1;
  if (offsetA !== 0 && offsetB === 0) return 1;
  return offsetA - offsetB;
});
