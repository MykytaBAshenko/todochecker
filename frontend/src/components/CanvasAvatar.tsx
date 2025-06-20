import React, { useEffect, useRef } from "react";

type CanvasAvatarProps = {
  avatar: number|string;
};

const CanvasAvatar: React.FC<CanvasAvatarProps> = ({ avatar }) => {
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!canvasRef.current || avatar == null) return;

    const canvas = canvasRef.current;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    canvas.width = 32;
    canvas.height = 32;
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    const code = avatar.toString().padStart(8, "0");
    const seed = parseInt(code);
    const prng = mulberry32(seed);

    const blockSize = 2;
    const cols = 16;
    const rows = 16;

    const baseHue = prng() * 360;
    const saturation = 80 + prng() * 20;
    const lightness = 40 + prng() * 20;

    const palette = Array.from({ length: 5 }, (_, i) =>
      `hsl(${(baseHue + i * 72) % 360}, ${saturation}%, ${lightness}%)`
    );

    for (let y = 0; y < rows; y++) {
      for (let x = 0; x < cols; x++) {
        if (prng() > 0.2) {
          const color = palette[Math.floor(prng() * palette.length)];
          ctx.fillStyle = color;
          ctx.fillRect(x * blockSize, y * blockSize, blockSize, blockSize);
        }
      }
    }

    function mulberry32(a: number) {
      return function () {
        let t = (a += 0x6d2b79f5);
        t = Math.imul(t ^ (t >>> 15), t | 1);
        t ^= t + Math.imul(t ^ (t >>> 7), t | 61);
        return ((t ^ (t >>> 14)) >>> 0) / 4294967296;
      };
    }
  }, [avatar]);

  return <canvas ref={canvasRef} />;
};

export default CanvasAvatar;
