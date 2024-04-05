---
sidebar_label: 'Architecture'
sidebar_position: 1000
title: Architecture
---

The following is the detailed data and control flow of SpecPipe.

![System Architecture](/img/specpipe-architecture.png)

### NATS

NATS (Neural Autonomic Transport System) is a lightweight and high-performance messaging system designed for distributed systems, offering simplicity, reliability, and scalability for cloud-native applications.

NATS is the backbone of our project. We couple NATS with Jetsream Jetstream to include additional functionality such as message persistence, replication, and delivery guarantees. NATS Jetstream is used in SpecPipe for:
* Sending radio signal data (either IQ data, or demodulated data) from an Edge Node to an Application (see definitions below).
* Sending configuration and health commands from Applications to Edge nodes to change the configuration of the edge nodes as well as receive health check from the edge nodes.

To accomplish this goal, we use the following NATS JetStream subjects [SpecPipe Subjects](./setup/sever-edge-setup#nats-jetstream-subjects)

### Edge Nodes

An edge node is a device (such as a laptop or a Raspberyy Pi) that has a software-defined radio attached to it via USB and is connected to the system. The radio is listening at a particular frequency that is initially set when the edge device registers with the system. The frequency can later be changed via commands issues by other Applications via the Controller API. See Control Flow line (3) in the diagram above.

### Applications

Applications are software that run on devices connected to the system. An application can access the following data:
* Receiving data from edge node via the NATS (see 'Data Flow' section below) or from an edge node directly via a socket connection (see 'Peer to Peer' section below).
* Geting the metadata of edge nodes (such as sampling rate, geolocation) via the Controller API.

Additionally, an application can also update the configuration (such as frequency and sampling rate) of an edge node (see Control Flow line 3) via the Controller API.

The aforementioned features can be accessed by using a Python Software Developemnt Toolkit provided by the SpecPipe project [specpipe-sdk-py](./setup/application.md). 

[Here](./examples) are some example applications build by the SpecPipe development team.

In our diagram above, 'Edge Nodes' and 'Applications' are separate boxes. This isn't meant to imply that an Edge Node cannot be an application. Since an application is just software running on a machine connected to a system, an Edge node can also be running an application. The diagram makes this distinction solely for ease of understanding the various flows.


### Health Check Server

The Health Check Server is an example of an Application that can run on devices connected to the system. The purpose of the health check system is to check the health of the nodes.

It utilizes NATS to distribute heartbeat signals across all edge devices, ensuring real-time monitoring and robustness within the network infrastructure, bolstering reliability and responsiveness in critical health monitoring scenarios.

The server performs this check by having a scheduled task that sends heart beat requests using the NATS subject specpipe-cluster.&lt;sdr_type&gt;.&lt;device_name&gt; to each device in the system (see Monitoring Flow lines 5 and 6)

### Controller API

A controller API serves as the interface for managing and orchestrating resources within SpecPipe, enabling Applications to programmatically interact with and manipulate the configuration and behavior of the Edge Nodes.

Applications can fetch Edge Node metadata as well as update edge node configuration by API requests to the controller API (see Control Flow lines 3 and 4)


### Information Flows

#### 1. Data Flow (1)

Raw IQ and Demodulated radio data is sent from Edge Nodes via NATS to Applications.

#### 2. Peer to Peer Flow (2)

Raw IQ radio data is sent from Edge Nodes to Applications via a socket. This is useful for sending high bit rate Raw IQ Data since if this data would be sent via NATS, the system would get bogged down.

#### 2. Control Flow (3,4)

* Line 3 : Applications update the configuration of an edge node (such as changing its sampling rate or frequency) via the Controller API. The Controller API then publishes a message on the appropriate NATS subject to change that setting for an Edge Node.
* Line 4: Applications can get edge node metadata (such as their location and sampling rate) via the Controller API.

See [Controller API](./setup/sever-edge-setup#7-setting-up-controller-plane) for example calls.

#### 4. Monitoring Flow (5,6)
Applications can monitor the health of the Edge Nodes by running the health command for a particular edge node. When this command is issued, NATS sending heart beat requests to all the Edge Nodes (line 5). Then, the edge nodes alive respond with a heartbeat (line 6) to the server via NATS
