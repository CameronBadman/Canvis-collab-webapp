const express = require('express');
const dotenv = require('dotenv');
const { initializeFirebase, initializeRedis } = require('./config');
const authRoutes = require('./routes/auth');
const cors = require('cors');

dotenv.config();

const app = express();
app.use(express.json());


async function startServer() {
  try {
    console.log('Loading environment variables...');
    
    // Log specific variables for debugging (avoid sensitive information)
    console.log('Environment variables loaded:', {
      REDIS_HOST: process.env.REDIS_HOST,
      FIREBASE_SERVICE_ACCOUNT: !!process.env.FIREBASE_SERVICE_ACCOUNT,
    });

    console.log('Initializing Firebase...');
    await initializeFirebase();
    
    console.log('Initializing Redis...');
    await initializeRedis();

    // Use auth routes
    app.use('/', authRoutes);
    console.log('Auth routes registered.');

    const PORT = process.env.PORT || 3000; // Use environment variable for port
    app.listen(PORT, () => {
      console.log(`Auth service running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start the auth service:', error);
    process.exit(1);
  }
}

// Catch unhandled rejections during startup
startServer().catch((error) => {
  console.error('Unhandled error during auth service startup:', error);
  process.exit(1);
});
