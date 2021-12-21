# Kubernetes Cluster Deployment and Monitoring

This project demonstrates a comprehensive Kubernetes deployment using Amazon EKS (v1.31) on AWS with two `t2.large` worker nodes. It showcases essential Kubernetes concepts, including deploying a MariaDB database with persistent storage, a highly available web server, and custom configurations. Additionally, it includes a custom-built Golang monitoring application to track pod lifecycle events.

The deployment is managed via Helm charts, ensuring repeatability and efficiency.

## Features

1. **Kubernetes Cluster Deployment**:
   - Cluster created using Amazon EKS with essential add-ons.
   - Worker nodes with sufficient resources to handle workloads.

2. **Database Deployment**:
   - MariaDB cluster deployed on Kubernetes with persistent storage (provided by Amazon EBS CSI Driver).

3. **Web Server Deployment**:
   - Multi-replica web server setup using Nginx.
   - Custom web server configuration mounted to pods.
   - A variable updated in an init container which show a custom value.

4. **Monitoring**:
   - Custom Golang application monitoring pod lifecycle events.
   - Logs events such as pod creation, deletion, and updates in real-time.

5. **Helm Chart Management**:
   - All components deployed using Helm for consistency and scalability.

## Add-Ons Used

To support the deployment on Amazon EKS, the following add-ons are enabled:
- **Amazon VPC CNI**: Handles networking for Kubernetes pods in the EKS cluster.
- **kube-proxy**: Manages network rules on nodes.
- **Amazon EKS Pod Identity Agent**: Ensures IAM role-based permissions for pods.
- **Amazon EBS CSI Driver**: Provides persistent storage for the MariaDB cluster.
- **CoreDNS**: Offers internal DNS resolution for the cluster.

---

The following sections provide step-by-step guidance for deploying and managing each component.

