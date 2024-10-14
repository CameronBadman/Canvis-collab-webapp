const admin = require('firebase-admin');
const Redis = require('ioredis');
const jwt = require('jsonwebtoken');

let redisClient;

exports.initializeFirebase = () => {
  console.log('Initializing Firebase...');
  try {
    const firebaseConfig = JSON.parse(process.env.FIREBASE_SERVICE_ACCOUNT);
    
    if (firebaseConfig.private_key) {
      firebaseConfig.private_key = firebaseConfig.private_key.replace(/\\n/g, '\n');
    }

    admin.initializeApp({
      credential: admin.credential.cert(firebaseConfig)
    });
    console.log('Firebase initialized successfully');
  } catch (error) {
    console.error('Error initializing Firebase:', error);
    process.exit(1);
  }
};

exports.initializeRedis = async () => {
  console.log('Connecting to Redis...');
  redisClient = new Redis({
    host: process.env.REDIS_HOST,
    port: process.env.REDIS_PORT,
    password: process.env.REDIS_PASSWORD
  });

  redisClient.on('error', (err) => console.error('Redis Client Error:', err));

  try {
    await redisClient.ping();
    console.log('Redis connected successfully');
  } catch (error) {
    console.error('Failed to connect to Redis:', error);
    process.exit(1);
  }

  return redisClient;
};

exports.getRedisClient = () => redisClient;

exports.generateToken = (uid) => {
  console.log(`Generating token for UID: ${uid}`);
  return jwt.sign({ uid }, process.env.JWT_SECRET, { expiresIn: '1h' });
}; 
