
import React, { useState, useCallback, useEffect } from 'react';
import { Excalidraw } from '@excalidraw/excalidraw';
import { debounce } from 'lodash';

const CanvasComponent = () => {
  const [elements, setElements] = useState([]);
  const [appState, setAppState] = useState({});

  // Memoize the debounced function
  const handleChange = useCallback(
    debounce((newElements, newAppState, _files) => {
      setElements(newElements);
      setAppState(newAppState);
    }, 300),
    []
  );

  // Use useEffect to log the updates when the state changes
  useEffect(() => {
    console.log('Updated elements:', elements);
    console.log('Updated appState:', appState);
  }, [elements, appState]); // This effect runs whenever elements or appState change

  return (
    <div style={{ width: '100%', height: '500px' }}>
      <Excalidraw
        initialData={{ elements, appState }}
        onChange={handleChange}
      />
    </div>
  );
};

export default CanvasComponent;
