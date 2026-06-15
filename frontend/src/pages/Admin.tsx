import { useState, useEffect } from "react";
import { getStatus, triggerCrawl, getFilters } from "../api";

export default function Admin() {
  const [status, setStatus] = useState<{ last_crawl: string | null; project_count: number; is_crawling: boolean } | null>(null);
  const [crawling, setCrawling] = useState(false);
  const [filters, setFilters] = useState<{ languages: string[]; topics: string[] } | null>(null);
  const [adminKey, setAdminKey] = useState("");

  const fetchStatus = async () => {
    try {
      const s = await getStatus();
      setStatus(s);
    } catch (err) {
      console.error(err);
    }
  };

  const fetchFilters = async () => {
    try {
      const f = await getFilters();
      setFilters(f);
    } catch (err) {
      console.error(err);
    }
  };

  useEffect(() => {
    fetchStatus();
    fetchFilters();
  }, []);

  const handleCrawl = async () => {
    setCrawling(true);
    try {
      await triggerCrawl(adminKey);
      setTimeout(fetchStatus, 1000);
    } catch (err) {
      console.error(err);
    } finally {
      setCrawling(false);
    }
  };

  return (
    <div style={{ maxWidth: 800, margin: "0 auto", padding: 24 }}>
      <h1>Admin</h1>

      <section style={{ marginBottom: 32 }}>
        <h2>Crawler Control</h2>
        <div style={{ display: "flex", gap: 8, alignItems: "center", marginBottom: 12 }}>
          <input
            type="password"
            placeholder="Admin key"
            value={adminKey}
            onChange={(e) => setAdminKey(e.target.value)}
            style={{ width: 200 }}
          />
          <button onClick={handleCrawl} disabled={crawling || (status?.is_crawling ?? false) || !adminKey}>
            {crawling || status?.is_crawling ? "Crawling..." : "Crawl Now"}
          </button>
        </div>
        {status && (
          <div style={{ marginTop: 12 }}>
            <p>Projects in DB: <strong>{status.project_count}</strong></p>
            <p>Last crawl: <strong>{status.last_crawl ? new Date(status.last_crawl).toLocaleString() : "never"}</strong></p>
            <p>Currently crawling: <strong>{status.is_crawling ? "yes" : "no"}</strong></p>
          </div>
        )}
      </section>

      <section>
        <h2>Available Filters</h2>
        {filters ? (
          <>
            <div>
              <h3>Languages ({filters.languages.length})</h3>
              <div style={{ display: "flex", flexWrap: "wrap", gap: 6 }}>
                {filters.languages.map((l) => (
                  <span key={l} style={{ background: "#e0e0e0", padding: "2px 8px", borderRadius: 4, fontSize: 13 }}>{l}</span>
                ))}
              </div>
            </div>
            <div style={{ marginTop: 16 }}>
              <h3>Topics ({filters.topics.length})</h3>
              <div style={{ display: "flex", flexWrap: "wrap", gap: 6 }}>
                {filters.topics.map((t) => (
                  <span key={t} style={{ background: "#e0e0e0", padding: "2px 8px", borderRadius: 4, fontSize: 13 }}>{t}</span>
                ))}
              </div>
            </div>
          </>
        ) : (
          <p>Loading filters...</p>
        )}
      </section>
    </div>
  );
}
