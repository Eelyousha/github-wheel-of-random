interface Props {
  limit: number;
  onLimitChange: (v: number) => void;
  lang: string;
  onLangChange: (v: string) => void;
  topic: string;
  onTopicChange: (v: string) => void;
  minStars: number;
  onMinStarsChange: (v: number) => void;
}

function parsePositiveInt(value: string, fallback: number): number {
  const trimmed = value.trim();
  if (trimmed === "") return fallback;
  const n = Number(trimmed);
  if (!Number.isFinite(n) || n < 0 || !Number.isInteger(n)) return fallback;
  return n;
}

export default function FilterPanel({ limit, onLimitChange, lang, onLangChange, topic, onTopicChange, minStars, onMinStarsChange }: Props) {
  return (
    <div style={{ display: "flex", flexWrap: "wrap", gap: 16, justifyContent: "center", margin: "16px 0" }}>
      <label>
        Limit (X):
        <input
          type="text"
          inputMode="numeric"
          placeholder="2–100"
          value={limit}
          onChange={(e) => {
            const v = parsePositiveInt(e.target.value, 20);
            onLimitChange(Math.min(Math.max(v, 2), 100));
          }}
          style={{ marginLeft: 8, width: 70 }}
        />
      </label>
      <label>
        Language:
        <input type="text" placeholder="e.g. Go" value={lang} onChange={(e) => onLangChange(e.target.value)} style={{ marginLeft: 8, width: 120 }} />
      </label>
      <label>
        Topic:
        <input type="text" placeholder="e.g. ai" value={topic} onChange={(e) => onTopicChange(e.target.value)} style={{ marginLeft: 8, width: 120 }} />
      </label>
      <label>
        Min Stars:
        <input
          type="text"
          inputMode="numeric"
          placeholder="0"
          value={minStars > 0 ? minStars : ""}
          onChange={(e) => {
            const v = parsePositiveInt(e.target.value, 0);
            onMinStarsChange(Math.min(v, 1_000_000_000));
          }}
          style={{ marginLeft: 8, width: 100 }}
        />
      </label>
    </div>
  );
}
