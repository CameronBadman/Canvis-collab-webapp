import React, { useRef, useEffect } from 'react';
import * as THREE from 'three';

const CustomBackground = () => {
  const mountRef = useRef(null);
  const sceneRef = useRef(null);
  const rendererRef = useRef(null);
  const linesRef = useRef([]);
  const animationFrameRef = useRef(null);

  useEffect(() => {
    if (!mountRef.current) return;

    const scene = new THREE.Scene();
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 2000);
    const renderer = new THREE.WebGLRenderer({ antialias: true });

    sceneRef.current = scene;
    rendererRef.current = renderer;

    renderer.setSize(window.innerWidth, window.innerHeight);
    mountRef.current.appendChild(renderer.domElement);

    scene.background = new THREE.Color(0xffffff); // White background
    camera.position.z = 1000;

    const createLine = () => {
      const points = [];
      const length = Math.random() * 600 + 400; // Longer lines for more horizontal coverage
      const segments = 6; // More segments for deeper curves
      const startPoint = new THREE.Vector3(
        Math.random() * window.innerWidth - window.innerWidth / 2,
        Math.random() * window.innerHeight - window.innerHeight / 2,
        0
      );
      
      points.push(startPoint);
      let prevPoint = startPoint;
      for (let i = 1; i < segments; i++) {
        const t = i / (segments - 1);
        const newPoint = new THREE.Vector3(
          prevPoint.x + (length / (segments - 1)) + (Math.random() - 0.5) * 50,
          prevPoint.y + (Math.random() - 0.5) * 100,
          0
        );
        points.push(newPoint);
        prevPoint = newPoint;
      }

      const curve = new THREE.CatmullRomCurve3(points);
      const finalPoints = curve.getPoints(200);

      const geometry = new THREE.BufferGeometry().setFromPoints(finalPoints);
      const material = new THREE.LineBasicMaterial({
        color: 0x000000,
        linewidth: 1.5,
        opacity: 0.8,
        transparent: true
      });

      const line = new THREE.Line(geometry, material);
      line.userData.points = finalPoints;
      line.userData.currentPoint = 0;
      line.userData.drawSpeed = Math.random() * 1 + 2;
      line.userData.velocity = new THREE.Vector3(
        (Math.random() - 0.5) * 0.2,
        (Math.random() - 0.5) * 0.05,
        0
      );
      return line;
    };

    const totalLines = 30;
    let currentLineIndex = 0;
    let lastLineAddTime = 0;

    const animate = (time) => {
      animationFrameRef.current = requestAnimationFrame(animate);

      if (currentLineIndex < totalLines) {
        const timeSinceLastLine = time - lastLineAddTime;
        const threshold = Math.max(500, Math.min(2000, currentLineIndex * 50));
        
        if (timeSinceLastLine > threshold) {
          const line = createLine();
          linesRef.current.push(line);
          scene.add(line);
          currentLineIndex++;
          lastLineAddTime = time;
        }
      }

      linesRef.current.forEach((line, index) => {
        if (line.userData.currentPoint < line.userData.points.length) {
          line.geometry.setDrawRange(0, line.userData.currentPoint);
          line.userData.currentPoint += line.userData.drawSpeed;
        }

        line.position.add(line.userData.velocity);

        const boundingBox = new THREE.Box3().setFromObject(line);
        if (boundingBox.max.x < -window.innerWidth / 2 || boundingBox.min.x > window.innerWidth / 2 ||
            boundingBox.max.y < -window.innerHeight / 2 || boundingBox.min.y > window.innerHeight / 2) {
          scene.remove(line);
          linesRef.current.splice(index, 1);
          currentLineIndex--;
        }
      });

      renderer.render(scene, camera);
    };

    animate(0);

    const handleResize = () => {
      if (rendererRef.current) {
        camera.aspect = window.innerWidth / window.innerHeight;
        camera.updateProjectionMatrix();
        rendererRef.current.setSize(window.innerWidth, window.innerHeight);
      }
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      
      if (animationFrameRef.current) {
        cancelAnimationFrame(animationFrameRef.current);
      }

      if (rendererRef.current) {
        rendererRef.current.dispose();
      }

      if (sceneRef.current) {
        sceneRef.current.clear();
      }

      linesRef.current.forEach(line => {
        line.geometry.dispose();
        line.material.dispose();
      });

      if (mountRef.current && rendererRef.current) {
        mountRef.current.removeChild(rendererRef.current.domElement);
      }
    };
  }, []);

  return (
    <div 
      ref={mountRef} 
      style={{ 
        position: 'fixed', 
        top: 0, 
        left: 0, 
        width: '100%', 
        height: '100%', 
        zIndex: -1,
        pointerEvents: 'none'
      }} 
    />
  );
};

export default CustomBackground;