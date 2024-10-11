const express = require('express');
const dotenv = require('dotenv');
const { initializeFirebase, initializeRedis } = require('./config');
const authRoutes = require('./routes/auth');

dotenv.config();

const app = express();
app.use(express.json());

// Add debug logging middleware
app.use((req, res, next) => {
  console.log(`[${new Date().toISOString()}] ${req.method} ${req.url}`);
  console.log('Headers:', JSON.stringify(req.headers, null, 2));
  console.log('Body:', JSON.stringify(req.body, null, 2));
  next();
});

async function startServer() {
  try {
    // Initialize Firebase
    initializeFirebase();

    // Initialize Redis
    await initializeRedis();

    // Use auth routes
    app.use('/auth', authRoutes);

    // Log all routes
    console.log('All registered routes:');
    app._router.stack.forEach(function(r){
      if (r.route && r.route.path){
        console.log(`${Object.keys(r.route.methods).join(', ').toUpperCase()} /auth${r.route.path}`)
      }
    })

    const PORT = process.env.PORT || 3000;
    app.listen(PORT, () => {
      console.log(`Auth service running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start the server:', error);
    process.exit(1);
  }
}

startServer().catch((error) => {
  console.error('Unhandled error during server startup:', error);
  process.exit(1);
});