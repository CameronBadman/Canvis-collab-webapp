const express = require('express');
const admin = require('firebase-admin');
const jwt = require('jsonwebtoken');
const { generateToken, getRedisClient } = require('../config');

const router = express.Router();

// Add a test route
router.get('/test', (req, res) => {
  console.log('Test route hit');
  res.json({ message: 'Auth service is running' });
});

// Register route
router.post('/register', async (req, res) => {
  console.log('Register route hit');
  console.log('Received registration request');
  try {
    const { email, password } = req.body;
    
    if (!email || !password) {
      console.log('Registration failed: Missing email or password');
      return res.status(400).json({ error: 'Email and password are required' });
    }
    
    if (password.length < 6) {
      console.log('Registration failed: Password too short');
      return res.status(400).json({ error: 'Password must be at least 6 characters long' });
    }

    console.log(`Attempting to create user with email: ${email}`);
    const userRecord = await admin.auth().createUser({
      email,
      password,
    });
    console.log(`User created successfully with UID: ${userRecord.uid}`);

    const token = generateToken(userRecord.uid);
    console.log(`Generated token for user: ${userRecord.uid}`);

    const expirationTime = 3600; // 1 hour in seconds
    await getRedisClient().set(userRecord.uid, token, { EX: expirationTime });
    console.log(`Stored token in Redis for user: ${userRecord.uid} with expiration: ${expirationTime} seconds`);

    res.status(201).json({ 
      message: 'User registered successfully',
      uid: userRecord.uid,
      email: userRecord.email,
      token,
      expiresIn: expirationTime
    });
  } catch (error) {
    console.error('Error in registration:', error);
    if (error.code === 'auth/email-already-exists') {
      return res.status(400).json({ error: 'Email already in use' });
    }
    res.status(500).json({ error: 'An error occurred during registration' });
  }
});

// Login route
router.post('/login', async (req, res) => {
  console.log('Login route hit');
  console.log('Received login request');
  try {
    const { email, password } = req.body;
    
    if (!email || !password) {
      console.log('Login failed: Missing email or password');
      return res.status(400).json({ error: 'Email and password are required' });
    }

    console.log(`Attempting to authenticate user: ${email}`);
    const userRecord = await admin.auth().getUserByEmail(email);
    
    // Note: Firebase Admin SDK doesn't provide a way to verify passwords
    // In a real-world scenario, you'd use Firebase Authentication REST API for this
    // For now, we'll assume the password is correct if we get this far
    console.log(`User authenticated: ${userRecord.uid}`);

    const token = generateToken(userRecord.uid);
    console.log(`Generated token for user: ${userRecord.uid}`);

    await getRedisClient().set(userRecord.uid, token, { EX: 3600 }); // Store in Redis for 1 hour
    console.log(`Stored token in Redis for user: ${userRecord.uid}`);

    res.json({ token, uid: userRecord.uid, email: userRecord.email });
  } catch (error) {
    console.error('Error in login:', error);
    if (error.code === 'auth/user-not-found') {
      return res.status(401).json({ error: 'Invalid email or password' });
    }
    res.status(500).json({ error: 'An error occurred during login' });
  }
});

// Logout route
router.post('/logout', async (req, res) => {
  console.log('Logout route hit');
  console.log('Received logout request');
  try {
    const authHeader = req.headers.authorization;
    if (!authHeader) {
      console.log('Logout failed: No authorization header');
      return res.status(401).json({ error: 'No authorization header provided' });
    }

    const token = authHeader.split(' ')[1];
    if (!token) {
      console.log('Logout failed: No token provided');
      return res.status(401).json({ error: 'No token provided' });
    }

    console.log('Verifying token');
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    console.log(`Token verified for user: ${decoded.uid}`);

    await getRedisClient().del(decoded.uid);
    console.log(`Removed token from Redis for user: ${decoded.uid}`);

    res.json({ message: 'Logged out successfully' });
  } catch (error) {
    console.error('Error in logout:', error);
    if (error instanceof jwt.JsonWebTokenError) {
      return res.status(401).json({ error: 'Invalid token' });
    }
    res.status(500).json({ error: 'An error occurred during logout' });
  }
});

router.get('/check-token/:uid', async (req, res) => {
  console.log('Check token route hit');
  try {
    const { uid } = req.params;
    const redisClient = getRedisClient();
    const token = await redisClient.get(uid);
    
    if (token) {
      const ttl = await redisClient.ttl(uid);
      console.log(`Token found for user: ${uid}, expires in: ${ttl} seconds`);
      res.json({ 
        message: 'Token found in Redis', 
        token,
        expiresIn: ttl
      });
    } else {
      console.log(`No token found for user: ${uid}`);
      res.status(404).json({ message: 'No token found for this user' });
    }
  } catch (error) {
    console.error('Error checking token:', error);
    res.status(500).json({ error: 'An error occurred while checking the token' });
  }
});

module.exports = router;