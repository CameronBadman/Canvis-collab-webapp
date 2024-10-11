const admin = require('firebase-admin');
const redis = require('redis');
const jwt = require('jsonwebtoken');

let redisClient;

exports.initializeFirebase = () => {
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
  redisClient = redis.createClient({
    url: `redis://${process.env.REDIS_HOST}:${process.env.REDIS_PORT}`,
    password: process.env.REDIS_PASSWORD
  });

  redisClient.on('error', (err) => console.log('Redis Client Error', err));

  try {
    await redisClient.connect();
    console.log('Redis connected successfully');
  } catch (error) {
    console.error('Failed to connect to Redis:', error);
    process.exit(1);
  }

  return redisClient;
};

exports.getRedisClient = () => redisClient;

exports.generateToken = (uid) => {
  return jwt.sign({ uid }, process.env.JWT_SECRET, { expiresIn: '1h' });
};