# diagram.py

# diagrams docs - https://diagrams.mingrammer.com/

from diagrams import Cluster, Diagram
from diagrams.k8s.compute import Pod
from diagrams.k8s.storage import PersistentVolume, PersistentVolumeClaim

with Diagram("wizter", show=False):

    
    with Cluster("cluster"):
        Pod("client") >> Pod("server") >> Pod("mongo") >> PersistentVolumeClaim("mongo-pvc")