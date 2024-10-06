import React, { useRef, useMemo, useEffect } from 'react';
import { useFrame, useThree } from '@react-three/fiber';
import { Points, PointMaterial } from '@react-three/drei';
import * as THREE from 'three';

const ParticleField = () => {
  const ref = useRef();
  const { size, viewport } = useThree();
  const aspect = size.width / viewport.width;

  const [positions, colors] = useMemo(() => {
    const positions = new Float32Array(5000 * 3);
    const colors = new Float32Array(5000 * 3);
    for (let i = 0; i < 5000; i++) {
      const x = (Math.random() - 0.5) * 10 * aspect;
      const y = (Math.random() - 0.5) * 10;
      const z = (Math.random() - 0.5) * 10;
      positions.set([x, y, z], i * 3);
      
      const r = Math.random() * 0.3 + 0.7; // More white
      const g = Math.random() * 0.3 + 0.7; // More white
      const b = Math.random(); // Full range for blue
      colors.set([r, g, b], i * 3);
    }
    return [positions, colors];
  }, [aspect]);

  const geometry = useRef(new THREE.BufferGeometry());

  useEffect(() => {
    geometry.current.setAttribute('position', new THREE.BufferAttribute(positions, 3));
    geometry.current.setAttribute('color', new THREE.BufferAttribute(colors, 3));
  }, [positions, colors]);

  useFrame((state) => {
    if (ref.current) {
      const time = state.clock.getElapsedTime();
      const positionAttribute = geometry.current.attributes.position;
      for (let i = 0; i < positionAttribute.count; i++) {
        const i3 = i * 3;
        positionAttribute.array[i3] += Math.sin(time + i * 0.1) * 0.001;
        positionAttribute.array[i3 + 1] += Math.cos(time + i * 0.1) * 0.001;
      }
      positionAttribute.needsUpdate = true;
    }
  });

  return (
    <Points ref={ref} geometry={geometry.current}>
      <PointMaterial vertexColors size={0.05} sizeAttenuation={true} transparent={true} />
    </Points>
  );
};

export default ParticleField;