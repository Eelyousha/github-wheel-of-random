import { useState, useEffect } from "react";
import { getRandomProjects, getFilters } from "../api";
import type { Project } from "../api";
import Wheel from "../components/Wheel";
import FilterPanel from "../components/FilterPanel";
import ProjectCard from "../components/ProjectCard";

export default function Home() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [winner, setWinner] = useState<Project | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [limit, setLimit] = useState(20);
  const [lang, setLang] = useState("");
  const [topic, setTopic] = useState("");
  const [minStars, setMinStars] = useState(0);
  const [availableLanguages, setAvailableLanguages] = useState<string[]>([]);
  const [availableTopics, setAvailableTopics] = useState<string[]>([]);

  useEffect(() => {
    getFilters().then((f) => {
      setAvailableLanguages(f.languages);
      setAvailableTopics(f.topics);
    }).catch(() => {});
  }, []);

  const handleSpin = async () => {
    setLoading(true);
    setError("");
    setWinner(null);
    try {
      const data = await getRandomProjects(limit, { lang: lang.trim(), topic: topic.trim(), min_stars: minStars });
      if (data.length === 0) {
        setError("No projects match your filters. Try fewer or broader criteria.");
      }
      setProjects(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to fetch projects");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ maxWidth: 800, margin: "0 auto", padding: 24 }}>
      <h1 style={{ textAlign: "center" }}>GitHub Wheel of Random</h1>
      <FilterPanel
        limit={limit}
        onLimitChange={setLimit}
        lang={lang}
        onLangChange={setLang}
        topic={topic}
        onTopicChange={setTopic}
        minStars={minStars}
        onMinStarsChange={setMinStars}
        availableLanguages={availableLanguages}
        availableTopics={availableTopics}
      />
      <div style={{ textAlign: "center", margin: "16px 0" }}>
        <button onClick={handleSpin} disabled={loading} style={{ padding: "12px 32px", fontSize: 18 }}>
          {loading ? "Loading..." : "Spin!"}
        </button>
      </div>
      {error && (
        <p style={{ textAlign: "center", color: "#e74c3c", margin: "12px 0" }}>{error}</p>
      )}
      {projects.length > 0 && (
        <Wheel projects={projects} onFinish={setWinner} />
      )}
      {winner && (
        <div style={{ marginTop: 32 }}>
          <h2 style={{ textAlign: "center" }}>Winner!</h2>
          <ProjectCard project={winner} />
        </div>
      )}
    </div>
  );
}
