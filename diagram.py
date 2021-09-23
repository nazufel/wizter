
# diagram.py

# diagrams docs - https://diagrams.mingrammer.com/

from diagrams import Cluster, Diagram
from diagrams.k8s.compute import Pod
from diagrams.k8s.podconfig import CM
from diagrams.k8s.storage import PersistentVolume, PersistentVolumeClaim

with Diagram("wizter", show=False):

    
    with Cluster("cluster"):
        pv = PersistentVolume("mongo-pv")
        with Cluster("wizards namespace"):
            client = Pod("client") 
            client_cm = CM("client")
            server = Pod("server") 
            server_cm = CM("server")
            mongo = Pod("mongo")
            pvc = PersistentVolumeClaim("mongo-pvc")

    client >> server >> mongo >> pvc
    pvc >> pv
    client >> client_cm
    server >> server_cm