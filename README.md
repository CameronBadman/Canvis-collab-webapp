# **Canvas WebApp**

**Canvis** is a collaborative drawing platform designed to support large numbers of users on the same canvas with ease. Users can access the platform through a simple 10-letter URL.

---

## **Tech Stack**

### **Frontend**
- **React (Vite)**: A fast and modern JavaScript framework for building user interfaces.
- **Three.js**: A JavaScript library to render 3D graphics for immersive drawing experiences.
- **Axios**: A promise-based HTTP client for making API requests.

### **Backend**
- **Databases**:
  - **Cassandra DB**: A distributed NoSQL database for handling large volumes of data.
  - **Redis**: An in-memory data structure store for caching and real-time data.
  - **Redis TTL (Time-to-Live)**: Used for expiring data after a certain period of inactivity.

- **Server**:
  - **GoLang (mux)**: A web framework for building high-performance REST APIs.
  - **JWT HS256**: JSON Web Tokens for secure authentication and authorization.

- **Message Broker**:
  - **Apache Kafka**: A distributed event streaming platform to manage real-time data feeds and event-driven communication.

### **Cloud Infrastructure**
- **AWS Cognito**: A user authentication and management service.
- **Terraform**: Infrastructure-as-code tool to provision AWS resources.
- **AWS EKS**: Managed Kubernetes service for deploying and scaling the application on AWS.
