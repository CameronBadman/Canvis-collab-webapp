function Sidebar({ randomId, currentTool, setCurrentTool, setStrokeWidth, setStrokeColor }) {
    return (
        <div
            className="sidebar"
            style={{
                width: '16rem', // Equivalent to w-64
                height: '100%', // Full height of the parent container
                backgroundColor: '#f3f4f6', // Tailwind's gray-100
                padding: '1.25rem', // Tailwind's p-5
                boxShadow: '0px 4px 6px rgba(0, 0, 0, 0.1)', // Tailwind's shadow-md
                display: 'flex', // Flex layout ensures child alignment
                flexDirection: 'column',
                boxSizing: 'border-box', // Ensures padding doesn't overflow
            }}
        >
            <h1 style={{ fontSize: '1.25rem', fontWeight: 'bold', marginBottom: '1.25rem' }}>
                Canvas Tools
            </h1>
            <p style={{ marginBottom: '1.25rem', fontSize: '0.875rem', color: '#4b5563' }}>
                Page ID: <span style={{ fontWeight: '600' }}>{randomId}</span>
            </p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: '1rem' }}>
                <button
                    style={{
                        width: '100%',
                        padding: '0.5rem 0',
                        borderRadius: '0.375rem',
                        color: '#fff',
                        backgroundColor: currentTool === 'pen' ? '#3b82f6' : '#4b5563',
                    }}
                    onClick={() => setCurrentTool('pen')}
                >
                    Pen Tool
                </button>
                <button
                    style={{
                        width: '100%',
                        padding: '0.5rem 0',
                        borderRadius: '0.375rem',
                        color: '#fff',
                        backgroundColor: currentTool === 'eraser' ? '#10b981' : '#4b5563',
                    }}
                    onClick={() => setCurrentTool('eraser')}
                >
                    Eraser Tool
                </button>
                <button
                    style={{
                        width: '100%',
                        padding: '0.5rem 0',
                        borderRadius: '0.375rem',
                        color: '#fff',
                        backgroundColor: currentTool === 'rectangle' ? '#ef4444' : '#4b5563',
                    }}
                    onClick={() => setCurrentTool('rectangle')}
                >
                    Rectangle Tool
                </button>
            </div>
            <div style={{ marginTop: '1.5rem' }}>
                <label htmlFor="strokeWidth" style={{ fontSize: '0.875rem', fontWeight: '500' }}>
                    Stroke Width
                </label>
                <input
                    type="range"
                    id="strokeWidth"
                    min="1"
                    max="10"
                    defaultValue="2"
                    style={{ width: '100%', marginTop: '0.25rem' }}
                    onChange={(e) => setStrokeWidth(Number(e.target.value))}
                />
            </div>
            <div style={{ marginTop: '1rem' }}>
                <label htmlFor="strokeColor" style={{ fontSize: '0.875rem', fontWeight: '500' }}>
                    Stroke Color
                </label>
                <input
                    type="color"
                    id="strokeColor"
                    defaultValue="#000000"
                    style={{
                        width: '100%',
                        marginTop: '0.25rem',
                        height: '2rem',
                        cursor: 'pointer',
                    }}
                    onChange={(e) => setStrokeColor(e.target.value)}
                />
            </div>
        </div>
    );
}

export default Sidebar;
