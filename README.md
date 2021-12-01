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
   - A variable updated in an init container which shows a custom value.

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



# Kubernetes Cluster Deployment

## EKS Cluster Setup

1. **Configure AWS CLI**  
   Set up AWS CLI using the provided Access Key ID and Secret Access Key with sufficient privileges:
   ```bash
   aws configure
   ```

2. **Clone the Repository**  
   Clone the repository containing the Terraform configuration for the EKS cluster:
   ```bash
   git clone https://github.com/emiakia/aws_eks_provider_simple
   cd aws_eks_provider_simple
   ```

3. **Initialize Terraform**  
   Initialize the Terraform configuration:
   ```bash
   terraform init
   ```

4. **Plan the Deployment**  
   Review the resources Terraform will create:
   ```bash
   terraform plan
   ```

5. **Apply the Terraform Configuration**  
   Deploy the EKS cluster by applying the configuration:
   ```bash
   terraform apply
   ```

   Wait for several minutes until the cluster creation completes.

6. **Verify Cluster Creation**  
   After deployment, verify the cluster with the following command:
   ```bash
   aws eks list-clusters
   ```

7. **Update kubeconfig**  
   Set the kubeconfig to interact with the new EKS cluster:
   ```bash
   aws eks update-kubeconfig --region <region> --name <cluster-name>
   ```

8. **Verify Connectivity**  
   Ensure the Kubernetes cluster is up and running by checking the nodes:
   ```bash
   kubectl get nodes
   ```

   The cluster is now ready for further steps.

## Clean Up (Post-Demo)

After the demo or testing, run the following command to avoid incurring extra costs:
```bash
terraform destroy
```

## Database Deployment

In this section, we will deploy a MariaDB instance using KubeDB, which simplifies managing databases on Kubernetes.

### 1. Install Helm
To begin, install Helm by following the official installation instructions for your environment:  
[Helm Installation Guide](https://helm.sh/docs/intro/install/)

### 2. Install KubeDB

1. Add the KubeDB Helm chart repository:
   ```bash
   helm repo add appscode https://charts.appscode.com/stable/
   helm repo update
   ```

2. Install the KubeDB Operator:
   ```bash
   helm install kubedb-operator appscode/kubedb --version v2024.8.21 --namespace kubedb --create-namespace --set-file global.license=/root/license.txt
   ```

   The `license.txt` file must be obtained from the AppsCode License Server:
   - Visit the [KubeDB Setup Page](https://kubedb.com/docs/v2024.11.18/setup/install/kubedb/)
   - Click on "Download a FREE license from AppsCode License Server"
   - Fill in the required fields.
   - Retrieve the lincense from your email and make the license.txt in suitable path.

   The **Kubernetes Cluster ID** can be found using the following command:
   ```bash
   kubectl get ns kube-system -o=jsonpath='{.metadata.uid}'
   ```

3. To verify that the KubeDB Operator is running:
   ```bash
   kubectl --namespace kubedb get pods
   ```

### 3. Install the MariaDB Cluster

1. Navigate to the `mydbcluster-helm` directory:
   ```bash
   cd mydbcluster-helm
   ```

2. Install the MariaDB cluster using Helm:
   ```bash
   helm install -f values.yaml -n demo mydbcluster .
   ```

   Alternatively, you can use this command if you prefer to install in the `demo` namespace:
   ```bash
   helm install -f values.yaml --namespace demo --create-namespace mydbcluster .
   ```

### 4. Verify the Database Deployment

1. Verify the MariaDB pods are running:
   ```bash
   kubectl get pod -n demo
   ```

   You should see three pods named `data-mydbcluster-mariadb-[0..2]`.

### Verifying MariaDB Cluster and Data Consistency

To verify your MariaDB cluster setup and ensure that data is consistent across nodes, follow these steps:

1. **Connect to the First Node (Pod) in the Cluster:**
   Use the following `kubectl exec` command to access the first MariaDB node:

   ```bash
   kubectl exec -it -n demo data-mydbcluster-mariadb-0 -- bash

2. **Login to MariaDB: Once inside the pod, log in to the MariaDB instance with the root credentials:**

   ```bash
    mysql -u${MYSQL_ROOT_USERNAME} -p${MYSQL_ROOT_PASSWORD}

3. **Verify Databases: After logging in, check the existing databases with the following command:**

   ```bash
   show databases;

4. **Create a Test Database: Create a new database to test the cluster's functionality:**

   ```bash
   create database test;
5. **reate a Table and Insert Data: Switch to the test database, create a table, and insert a sample record:**

   ```bash
   use test;
   create table t1 (id integer);
   insert into t1 values(1);
6. **Exit the Pod: To exit the MariaDB shell, use Ctrl+D. To exit the pod shell, type:**

   ```bash
   exit

7. **Verify Data on Other Pods: To ensure data is replicated across the cluster, repeat the above steps on another MariaDB pod in the cluster. First, connect to the second pod:**

   ```bash
   kubectl exec -it -n demo data-mydbcluster-mariadb-1 -- bash
Then, follow the same steps to log in to MariaDB and verify that the test database and its data exist.

   ```bash
   mysql -u${MYSQL_ROOT_USERNAME} -p${MYSQL_ROOT_PASSWORD}
   show databases;
   use test;
   select * from t1;

   You should see the data inserted in the first pod selected in the second pod.

