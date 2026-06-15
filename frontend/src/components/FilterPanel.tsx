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

export default function FilterPanel({ limit, onLimitChange, lang, onLangChange, topic, onTopicChange, minStars, onMinStarsChange }: Props) {
  return (
    <div style={{ display: "flex", flexWrap: "wrap", gap: 16, justifyContent: "center", margin: "16px 0" }}>
      <label>
        Limit (X):
        <input type="range" min={2} max={100} value={limit} onChange={(e) => onLimitChange(Number(e.target.value))} />
        <span style={{ marginLeft: 8 }}>{limit}</span>
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
        <input type="number" min={0} value={minStars} onChange={(e) => onMinStarsChange(Number(e.target.value))} style={{ marginLeft: 8, width: 100 }} />
      </label>
    </div>
  );
}
