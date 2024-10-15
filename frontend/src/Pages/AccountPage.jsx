import React, { useState, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Plus, X } from 'lucide-react';
import CustomBackground from '../Componants/CustomBackground';
import { canvasService } from '../services/CanvasService';

const AccountPage = () => {
  const [canvases, setCanvases] = useState([]);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isCreatingCanvas, setIsCreatingCanvas] = useState(false);
  const [newCanvasName, setNewCanvasName] = useState('');

  useEffect(() => {
    fetchCanvases();
  }, []);

  const fetchCanvases = async () => {
    try {
      setIsLoading(true);
      const fetchedCanvases = await canvasService.getUserCanvases();
      setCanvases(fetchedCanvases);
      setIsLoading(false);
    } catch (err) {
      console.error('Error fetching canvases:', err);
      setError('Failed to load canvases. Please try again later.');
      setIsLoading(false);
    }
  };

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

  const handleCreateCanvas = () => {
    setIsCreatingCanvas(true);
    setNewCanvasName('');
  };

  const handleCancelCreate = () => {
    setIsCreatingCanvas(false);
    setNewCanvasName('');
  };

  const handleNewCanvasNameChange = (e) => {
    setNewCanvasName(e.target.value);
  };

  const handleSubmitNewCanvas = async (e) => {
    e.preventDefault();
    if (newCanvasName.trim()) {
      try {
        const newCanvas = await canvasService.createCanvas(newCanvasName);
        setCanvases([...canvases, newCanvas]);
        setIsCreatingCanvas(false);
        setNewCanvasName('');
      } catch (err) {
        console.error('Error creating canvas:', err);
        alert('Failed to create a new canvas. Please try again.');
      }
    }
  };

  const CanvasCard = ({ canvas }) => (
    <div
      key={canvas.id}
      className="bg-white bg-opacity-80 rounded-lg shadow-lg p-6 w-full h-full flex flex-col justify-center items-center transition-all duration-300 ease-in-out hover:bg-opacity-100 hover:shadow-xl cursor-pointer"
    >
      <h2 className="text-2xl font-bold text-gray-800 mb-2 text-center">{canvas.name}</h2>
      <p className="text-sm text-gray-600">Last edited: {new Date(canvas.updated_at).toLocaleDateString()}</p>
    </div>
  );

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

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
              {canvases.map((canvas) => (
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
        
        {!isCreatingCanvas ? (
          <button 
            onClick={handleCreateCanvas}
            className="mt-4 bg-gray-800 text-white py-2 px-4 rounded-full hover:bg-gray-700 transition duration-300 flex items-center text-lg"
          >
            <Plus size={20} className="mr-2" />
            Create New Canvas
          </button>
        ) : (
          <form onSubmit={handleSubmitNewCanvas} className="mt-4 flex items-center">
            <input
              type="text"
              value={newCanvasName}
              onChange={handleNewCanvasNameChange}
              placeholder="Enter canvas name"
              className="px-3 py-2 border border-gray-300 rounded-l-md focus:outline-none focus:ring-2 focus:ring-gray-500"
              autoFocus
            />
            <button
              type="submit"
              className="bg-gray-800 text-white py-2 px-4 rounded-r-md hover:bg-gray-700 transition duration-300"
            >
              Create
            </button>
            <button
              type="button"
              onClick={handleCancelCreate}
              className="ml-2 bg-gray-300 text-gray-800 py-2 px-4 rounded-md hover:bg-gray-400 transition duration-300"
            >
              Cancel
            </button>
          </form>
        )}
      </div>
    </div>
  );
};

export default AccountPage;