# **Canvas WebApp**

**Canvis** is a collaborative drawing platform designed to support large numbers of users on the same canvas with ease. Users can access the platform through a simple 10-letter URL.

---

## **Tech Stack**

### **Frontend**
- **React (Vite)**: A fast and modern JavaScript framework for building user interfaces.
- **Three.js**: for rendering the pretty 3D graphics background 
- **Axios**: HTTP client for making API requests.

### **Backend**
- **Databases**:
  - **Cassandra DB**: holds Users and Canvas long-term storage 
  - **Redis TTL (Time-to-Live)**: Used for Canvases

- **Server**:
  - **GoLang (mux)**: A web framework for building high-performance REST APIs.
  - **JWT HS256**: JSON Web Tokens for secure authentication and authorization.

- **Message Broker**:
  - **Apache Kafka**: used to register vector data and take canvases from Redis TTL to cassandra storage after TTL is doen 

### **Cloud Infrastructure**
- **AWS Cognito**: A user authentication and management service.
- **Terraform**: Infrastructure-as-code tool to provision AWS resources.
- **AWS EKS**: for deployment 
