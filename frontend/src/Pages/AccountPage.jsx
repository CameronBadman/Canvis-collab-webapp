import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Plus } from 'lucide-react';
import CustomBackground from '../Componants/CustomBackground';

const AccountPage = () => {
  const [canvases, setCanvases] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);

  useEffect(() => {
    // Fetch user's canvases
    // This is a placeholder. Replace with actual API call.
    setCanvases([
      { id: 1, title: 'My First Canvas', lastEdited: '2023-05-15' },
      { id: 2, title: 'Project Ideas', lastEdited: '2023-05-20' },
      { id: 3, title: 'Weekly Planner', lastEdited: '2023-05-22' },
      { id: 4, title: 'Business Model', lastEdited: '2023-05-25' },
      { id: 5, title: 'Vacation Plans', lastEdited: '2023-05-28' },
    ]);
  }, []);

  const nextCanvas = () => {
    setCurrentIndex((prevIndex) => 
      prevIndex + 3 >= canvases.length ? 0 : prevIndex + 3
    );
  };

  const prevCanvas = () => {
    setCurrentIndex((prevIndex) => 
      prevIndex - 3 < 0 ? Math.max(canvases.length - 3, 0) : prevIndex - 3
    );
  };

  const CanvasCard = ({ canvas }) => (
    <div
      key={canvas.id}
      className="bg-white bg-opacity-80 rounded-lg shadow-lg p-6 w-full h-full flex flex-col justify-center items-center transition-all duration-300 ease-in-out hover:bg-opacity-100 hover:shadow-xl cursor-pointer"
    >
      <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">{canvas.title}</h2>
      <p className="text-sm text-gray-600">Last edited: {canvas.lastEdited}</p>
    </div>
  );

  return (
    <div className="relative w-full h-screen overflow-hidden">
      <CustomBackground />
      
      <div className="absolute inset-0 flex flex-col items-center justify-center p-8">
        <h1 className="text-4xl font-bold mb-8 text-center text-gray-800">My Canvases</h1>
        
        <div className="relative flex items-center justify-center w-full h-2/3 mb-8">
          <button 
            onClick={prevCanvas} 
            className="absolute left-4 z-10 bg-gray-800 text-white p-2 rounded-full hover:bg-gray-700 transition duration-300"
            disabled={currentIndex === 0}
          >
            <ChevronLeft size={24} />
          </button>
          
          <div className="overflow-hidden w-full h-full">
            <div 
              className="flex w-full h-full transition-transform duration-300 ease-in-out" 
              style={{ transform: `translateX(-${(currentIndex / 3) * 100}%)` }}
            >
              {canvases.map((canvas, index) => (
                <div key={canvas.id} className="w-1/3 h-full flex-shrink-0 p-2">
                  <CanvasCard canvas={canvas} />
                </div>
              ))}
            </div>
          </div>
          
          <button 
            onClick={nextCanvas} 
            className="absolute right-4 z-10 bg-gray-800 text-white p-2 rounded-full hover:bg-gray-700 transition duration-300"
            disabled={currentIndex + 3 >= canvases.length}
          >
            <ChevronRight size={24} />
          </button>
        </div>
        
        <button className="mt-4 bg-gray-800 text-white py-2 px-4 rounded-full hover:bg-gray-700 transition duration-300 flex items-center text-lg">
          <Plus size={20} className="mr-2" />
          Create New Canvas
        </button>
      </div>
    </div>
  );
};

export default AccountPage;