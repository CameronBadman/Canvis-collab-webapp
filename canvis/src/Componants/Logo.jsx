import React from 'react';
import { Link } from 'react-router-dom';

const Logo = () => {
  return (
    <Link to="/" className="text-4xl font-bold text-gray-800 tracking-wider no-underline">
      <h2 style={{ fontFamily: 'Arial, sans-serif', textShadow: '2px 2px 4px rgba(0,0,0,0.1)' }}>
        CANVIS
      </h2>
    </Link>
  );
};

export default Logo;