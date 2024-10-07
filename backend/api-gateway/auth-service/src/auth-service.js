const express = require('express');
const dotenv = require('dotenv');
const { initializeFirebase, initializeRedis } = require('./config');
const authRoutes = require('./routes/auth');

dotenv.config();

const app = express();
app.use(express.json());

async function startServer() {
  try {
    // Initialize Firebase
    initializeFirebase();

    // Initialize Redis
    await initializeRedis();

    // Use auth routes
    app.use('/auth', authRoutes());

    const PORT = process.env.PORT || 3000;
    app.listen(PORT, () => {
      console.log(`Server running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start the server:', error);
    process.exit(1);
  }
}

startServer();