import React, { useRef, useEffect } from 'react';
import * as THREE from 'three';

const CustomBackground = () => {
  const mountRef = useRef(null);

  useEffect(() => {
    const scene = new THREE.Scene();
    scene.background = new THREE.Color(0xf0f0f0); // Light gray background
    const camera = new THREE.PerspectiveCamera(75, window.innerWidth / window.innerHeight, 0.1, 1000);
    const renderer = new THREE.WebGLRenderer({ antialias: true });

    renderer.setSize(window.innerWidth, window.innerHeight);
    mountRef.current.appendChild(renderer.domElement);

    const createMountainStroke = () => {
      const points = [];
      const segments = Math.floor(Math.random() * 20) + 15;
      const startX = Math.random() * 60 - 30; // Increased range to ensure coverage of left side
      const startY = Math.random() * 40 - 20; // Increased vertical range

      for (let i = 0; i <= segments; i++) {
        const x = startX + (i / segments) * 60; // Increased width to 60
        const y = startY + Math.sin(i / segments * Math.PI * 2) * (Math.random() * 6 + 3); // More pronounced curves
        points.push(new THREE.Vector3(x, y, 0));
      }

      const curve = new THREE.CatmullRomCurve3(points);
      const geometry = new THREE.BufferGeometry().setFromPoints(curve.getPoints(100));
      
      const material = new THREE.LineBasicMaterial({
        color: new THREE.Color(Math.random() * 0.1 + 0.1, Math.random() * 0.1 + 0.1, Math.random() * 0.1 + 0.2),
        opacity: Math.random() * 0.5 + 0.5,
        transparent: true,
        linewidth: Math.random() * 2 + 1,
      });

      const line = new THREE.Line(geometry, material);
      line.userData.totalLength = line.geometry.attributes.position.count;
      line.userData.currentLength = 0;
      return line;
    };

    const mountainStrokes = [];
    const totalStrokes = 60; // Increased number of strokes for better coverage
    for (let i = 0; i < totalStrokes; i++) {
      const stroke = createMountainStroke();
      stroke.visible = false;
      mountainStrokes.push(stroke);
      scene.add(stroke);
    }

    camera.position.z = 40; // Moved camera further back to show more of the scene
    camera.position.x = 0; // Centering the camera horizontally

    let currentStrokeIndex = 0;
    const drawingSpeed = 1;

    const animate = () => {
      requestAnimationFrame(animate);

      if (currentStrokeIndex < mountainStrokes.length) {
        const currentStroke = mountainStrokes[currentStrokeIndex];
        if (!currentStroke.visible) {
          currentStroke.visible = true;
        }

        currentStroke.userData.currentLength += drawingSpeed;
        const count = Math.min(Math.floor(currentStroke.userData.currentLength), currentStroke.userData.totalLength);

        currentStroke.geometry.setDrawRange(0, count);

        if (count >= currentStroke.userData.totalLength) {
          currentStrokeIndex++;
        }
      } else {
        // Continue drawing new strokes
        const newStroke = createMountainStroke();
        mountainStrokes.push(newStroke);
        scene.add(newStroke);

        // Remove oldest stroke if we have too many
        if (mountainStrokes.length > totalStrokes * 1.5) {
          const oldestStroke = mountainStrokes.shift();
          scene.remove(oldestStroke);
        }

        currentStrokeIndex = mountainStrokes.length - 1;
      }

      // Slight movement of all visible strokes
      mountainStrokes.forEach(stroke => {
        if (stroke.visible) {
          stroke.position.y += (Math.random() - 0.5) * 0.02;
        }
      });

      renderer.render(scene, camera);
    };

    animate();

    const handleResize = () => {
      camera.aspect = window.innerWidth / window.innerHeight;
      camera.updateProjectionMatrix();
      renderer.setSize(window.innerWidth, window.innerHeight);
    };

    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener('resize', handleResize);
      mountRef.current.removeChild(renderer.domElement);
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