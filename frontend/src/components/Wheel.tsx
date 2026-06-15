import { useRef, useEffect, useState, useCallback } from "react";
import type { Project } from "../api";

interface Props {
  projects: Project[];
  onFinish: (project: Project) => void;
}

const COLORS = [
  "#FF6B6B", "#4ECDC4", "#45B7D1", "#96CEB4",
  "#FFEAA7", "#DDA0DD", "#98D8C8", "#F7DC6F",
  "#BB8FCE", "#85C1E9", "#F0B27A", "#82E0AA",
  "#F1948A", "#85929E", "#73C6B6", "#E59866",
];

const CANVAS_SIZE = 400;
const SPIN_DURATION = 4000; // ms

export default function Wheel({ projects, onFinish }: Props) {
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [spinning, setSpinning] = useState(false);
  const rotationRef = useRef(0);
  const animRef = useRef<number | null>(null);

  const draw = useCallback((rotation: number) => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    const n = projects.length;
    if (n === 0) return;

    const arcSize = (2 * Math.PI) / n;
    const center = CANVAS_SIZE / 2;
    const radius = center - 10;

    ctx.clearRect(0, 0, CANVAS_SIZE, CANVAS_SIZE);

    // Draw slices
    for (let i = 0; i < n; i++) {
      const startAngle = rotation + i * arcSize;
      const endAngle = startAngle + arcSize;

      ctx.beginPath();
      ctx.moveTo(center, center);
      ctx.arc(center, center, radius, startAngle, endAngle);
      ctx.closePath();

      ctx.fillStyle = COLORS[i % COLORS.length];
      ctx.fill();
      ctx.strokeStyle = "#fff";
      ctx.lineWidth = 2;
      ctx.stroke();

      // Draw text
      const midAngle = startAngle + arcSize / 2;
      const textRadius = radius * 0.65;
      const x = center + Math.cos(midAngle) * textRadius;
      const y = center + Math.sin(midAngle) * textRadius;

      ctx.save();
      ctx.translate(x, y);
      ctx.rotate(midAngle);
      ctx.fillStyle = "#fff";
      ctx.font = "bold 11px monospace";
      ctx.textAlign = "center";
      ctx.textBaseline = "middle";

      const name = projects[i].name;
      const maxChars = 12;
      const label = name.length > maxChars ? name.slice(0, maxChars) + ".." : name;
      ctx.fillText(label, 0, 0);
      ctx.restore();
    }

    // Draw center circle
    ctx.beginPath();
    ctx.arc(center, center, 20, 0, 2 * Math.PI);
    ctx.fillStyle = "#fff";
    ctx.fill();
    ctx.strokeStyle = "#333";
    ctx.lineWidth = 3;
    ctx.stroke();

    // Draw pointer (top)
    ctx.beginPath();
    ctx.moveTo(center, center - radius - 10);
    ctx.lineTo(center - 12, center - radius + 10);
    ctx.lineTo(center + 12, center - radius + 10);
    ctx.closePath();
    ctx.fillStyle = "#e74c3c";
    ctx.fill();
    ctx.strokeStyle = "#c0392b";
    ctx.lineWidth = 2;
    ctx.stroke();
  }, [projects]);

  useEffect(() => {
    if (!spinning) {
      draw(rotationRef.current);
    }
  }, [draw, spinning, projects]);

  const spin = () => {
    if (spinning || projects.length === 0) return;
    setSpinning(true);

    const n = projects.length;
    const arcSize = (2 * Math.PI) / n;

    // Random final rotation: at least 5 full spins + random offset
    const extraSpins = 5 + Math.random() * 5;
    const randomSlice = Math.floor(Math.random() * n);
    const targetSliceAngle = randomSlice * arcSize;
    // We want the pointer (at top, -PI/2) to land on the target slice
    const totalRotation = extraSpins * 2 * Math.PI + (2 * Math.PI - targetSliceAngle);

    const startRotation = rotationRef.current;
    const targetRotation = startRotation + totalRotation;
    const startTime = performance.now();

    const animate = (currentTime: number) => {
      const elapsed = currentTime - startTime;
      const progress = Math.min(elapsed / SPIN_DURATION, 1);

      // Ease out cubic
      const eased = 1 - Math.pow(1 - progress, 3);
      const currentRotation = startRotation + (targetRotation - startRotation) * eased;

      rotationRef.current = currentRotation;
      draw(currentRotation);

      if (progress < 1) {
        animRef.current = requestAnimationFrame(animate);
      } else {
        // Determine winner
        const normalizedRotation = ((currentRotation % (2 * Math.PI)) + 2 * Math.PI) % (2 * Math.PI);
        // Pointer is at top = -PI/2, so we need to map
        const pointerAngle = (2 * Math.PI - normalizedRotation + (3 * Math.PI / 2)) % (2 * Math.PI);
        const winnerIndex = Math.floor(pointerAngle / arcSize) % n;

        setSpinning(false);
        onFinish(projects[winnerIndex]);
      }
    };

    animRef.current = requestAnimationFrame(animate);
  };

  useEffect(() => {
    return () => {
      if (animRef.current) cancelAnimationFrame(animRef.current);
    };
  }, []);

  return (
    <div style={{ textAlign: "center", margin: "16px 0" }}>
      <canvas
        ref={canvasRef}
        width={CANVAS_SIZE}
        height={CANVAS_SIZE}
        style={{ maxWidth: "100%", cursor: spinning ? "not-allowed" : "pointer" }}
        onClick={spin}
      />
      <p style={{ fontSize: 12, color: "#888" }}>Click the wheel to spin, or press the Spin button</p>
    </div>
  );
}
