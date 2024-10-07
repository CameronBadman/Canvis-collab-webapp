// src/Canvas/Canvas.jsx
import React, { useState, useRef, useCallback, useEffect } from 'react';

const CustomSVGDrawing = ({ strokeWidth = 2, strokeColor = "black", currentTool }) => {
  const [isDrawing, setIsDrawing] = useState(false);
  const [lines, setLines] = useState([]);
  const [dimensions, setDimensions] = useState({ width: '100%', height: '100%' });
  const svgRef = useRef(null);
  const containerRef = useRef(null);

  useEffect(() => {
    const updateDimensions = () => {
      if (containerRef.current) {
        setDimensions({
          width: containerRef.current.offsetWidth,
          height: containerRef.current.offsetHeight
        });
      }
    };

    window.addEventListener('resize', updateDimensions);
    updateDimensions();

    return () => window.removeEventListener('resize', updateDimensions);
  }, []);

  const getCoordinates = useCallback((event) => {
    if (!svgRef.current) return { x: 0, y: 0 };
    const svgRect = svgRef.current.getBoundingClientRect();
    return {
      x: event.clientX - svgRect.left,
      y: event.clientY - svgRect.top
    };
  }, []);

  const startDrawing = useCallback((event) => {
    setIsDrawing(true);
    const { x, y } = getCoordinates(event);
    if (currentTool === 'pen') {
      setLines(prevLines => [...prevLines, [{ x, y }]]);
    } else if (currentTool === 'eraser') {
      erase(x, y);
    }
  }, [getCoordinates, currentTool]);

  const draw = useCallback((event) => {
    if (!isDrawing) return;
    const { x, y } = getCoordinates(event);
    if (currentTool === 'pen') {
      setLines(prevLines => {
        const newLines = [...prevLines];
        newLines[newLines.length - 1].push({ x, y });
        return newLines;
      });
    } else if (currentTool === 'eraser') {
      erase(x, y);
    }
  }, [isDrawing, getCoordinates, currentTool]);

  const endDrawing = useCallback(() => {
    setIsDrawing(false);
  }, []);

  const erase = useCallback((x, y) => {
    const eraserRadius = 10; // Adjust this value to change eraser size

    setLines(prevLines => {
      return prevLines.flatMap(line => {
        const newLine = [];
        let segmentStart = line[0];

        for (let i = 1; i < line.length; i++) {
          const segmentEnd = line[i];
          const distToSegment = distanceToLineSegment({ x, y }, segmentStart, segmentEnd);

          if (distToSegment > eraserRadius) {
            newLine.push(segmentStart);
            if (i === line.length - 1) {
              newLine.push(segmentEnd);
            }
          } else {
            if (newLine.length > 0) {
              return [newLine, [segmentEnd]];
            }
            return [[segmentEnd]];
          }

          segmentStart = segmentEnd;
        }

        return newLine.length > 1 ? [newLine] : [];
      });
    });
  }, []);

  const distanceToLineSegment = (point, lineStart, lineEnd) => {
    const A = point.x - lineStart.x;
    const B = point.y - lineStart.y;
    const C = lineEnd.x - lineStart.x;
    const D = lineEnd.y - lineStart.y;

    const dot = A * C + B * D;
    const lenSq = C * C + D * D;
    let param = -1;
    if (lenSq !== 0) param = dot / lenSq;

    let xx, yy;

    if (param < 0) {
      xx = lineStart.x;
      yy = lineStart.y;
    } else if (param > 1) {
      xx = lineEnd.x;
      yy = lineEnd.y;
    } else {
      xx = lineStart.x + param * C;
      yy = lineStart.y + param * D;
    }

    const dx = point.x - xx;
    const dy = point.y - yy;
    return Math.sqrt(dx * dx + dy * dy);
  };

  return (
    <div ref={containerRef} style={{ width: '100%', height: '100%', position: 'absolute', top: 0, left: 0 }}>
      <svg
        ref={svgRef}
        width={dimensions.width}
        height={dimensions.height}
        onMouseDown={startDrawing}
        onMouseMove={draw}
        onMouseUp={endDrawing}
        onMouseLeave={endDrawing}
        style={{ display: 'block', cursor: currentTool === 'eraser' ? 'crosshair' : 'default' }}
      >
        <rect width="100%" height="100%" fill="white" />
        {lines.map((line, lineIndex) => (
          <polyline
            key={lineIndex}
            points={line.map(point => `${point.x},${point.y}`).join(' ')}
            fill="none"
            stroke={strokeColor}
            strokeWidth={strokeWidth}
          />
        ))}
      </svg>
    </div>
  );  
};

export default CustomSVGDrawing;