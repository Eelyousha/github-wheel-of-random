import type { Project } from "../api";

interface Props {
  project: Project;
}

const COLORS = ["#f0db4f", "#e34c26", "#563d7c", "#3178c6", "#00add8", "#f29111", "#cc6699", "#6cc24a"];

export default function ProjectCard({ project }: Props) {
  const color = COLORS[project.topics.length % COLORS.length] || "#888";

  return (
    <div style={{
      border: "1px solid #ddd",
      borderRadius: 8,
      padding: 16,
      display: "flex",
      gap: 16,
      alignItems: "center",
      boxShadow: "0 2px 8px rgba(0,0,0,0.1)",
      maxWidth: 500,
      margin: "0 auto",
    }}>
      {project.avatar_url && (
        <img src={project.avatar_url} alt={project.owner} style={{ width: 48, height: 48, borderRadius: "50%" }} />
      )}
      <div style={{ flex: 1 }}>
        <a href={project.html_url} target="_blank" rel="noopener noreferrer" style={{ fontSize: 18, fontWeight: "bold", color: "#0366d6" }}>
          {project.full_name}
        </a>
        <p style={{ margin: "4px 0", color: "#555", fontSize: 14 }}>{project.description || "No description"}</p>
        <div style={{ display: "flex", gap: 12, fontSize: 13, color: "#777", alignItems: "center" }}>
          <span>⭐ {project.stars.toLocaleString()}</span>
          {project.language && <span style={{ background: color, width: 12, height: 12, borderRadius: "50%", display: "inline-block" }} />}
          {project.language && <span>{project.language}</span>}
        </div>
        {project.topics.length > 0 && (
          <div style={{ display: "flex", flexWrap: "wrap", gap: 4, marginTop: 6 }}>
            {project.topics.slice(0, 5).map((t) => (
              <span key={t} style={{ background: "#e8f0fe", color: "#1a73e8", padding: "1px 6px", borderRadius: 10, fontSize: 11 }}>
                {t}
              </span>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
