"use strict";(self.webpackChunkdocs_website=self.webpackChunkdocs_website||[]).push([[443],{5459:(e,t,i)=>{i.r(t),i.d(t,{assets:()=>l,contentTitle:()=>o,default:()=>h,frontMatter:()=>s,metadata:()=>r,toc:()=>c});var n=i(4848),a=i(8453);const s={sidebar_label:"Architecture",sidebar_position:1e3,title:"Architecture"},o=void 0,r={id:"architecture",title:"Architecture",description:"High Level Design",source:"@site/docs/architecture.md",sourceDirName:".",slug:"/architecture",permalink:"/specpipe/architecture",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:1e3,frontMatter:{sidebar_label:"Architecture",sidebar_position:1e3,title:"Architecture"},sidebar:"tutorialSidebar",previous:{title:"Example Applications",permalink:"/specpipe/examples"}},l={},c=[{value:"High Level Design",id:"high-level-design",level:2},{value:"Architecture",id:"architecture",level:2},{value:"NATS",id:"nats",level:3},{value:"Edge Nodes",id:"edge-nodes",level:3},{value:"Applications",id:"applications",level:3},{value:"Health Check Server",id:"health-check-server",level:3},{value:"Controller API",id:"controller-api",level:3},{value:"Information Flows",id:"information-flows",level:3},{value:"1. Data Flow (1)",id:"1-data-flow-1",level:4},{value:"2. Peer to Peer Flow (2)",id:"2-peer-to-peer-flow-2",level:4},{value:"2. Control Flow (3,4)",id:"2-control-flow-34",level:4},{value:"4. Monitoring Flow (5,6)",id:"4-monitoring-flow-56",level:4}];function d(e){const t={a:"a",h2:"h2",h3:"h3",h4:"h4",img:"img",li:"li",p:"p",ul:"ul",...(0,a.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.h2,{id:"high-level-design",children:"High Level Design"}),"\n",(0,n.jsx)(t.p,{children:(0,n.jsx)(t.img,{alt:"SpecPipe overview",src:i(6317).A+"",width:"1030",height:"1019"})}),"\n",(0,n.jsx)(t.h2,{id:"architecture",children:"Architecture"}),"\n",(0,n.jsx)(t.p,{children:(0,n.jsx)(t.img,{alt:"System Architecture",src:i(1549).A+"",width:"1001",height:"550"})}),"\n",(0,n.jsx)(t.h3,{id:"nats",children:"NATS"}),"\n",(0,n.jsx)(t.p,{children:"NATS (Neural Autonomic Transport System) is a lightweight and high-performance messaging system designed for distributed systems, offering simplicity, reliability, and scalability for cloud-native applications."}),"\n",(0,n.jsx)(t.p,{children:"NATS is the backbone of our project. We couple NATS with Jetsream Jetstream to include additional functionality such as message persistence, replication, and delivery guarantees. NATS Jetstream is used in SpecPipe for:"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsx)(t.li,{children:"Sending radio signal data (either raw IQ data, or demodulated data) from an Edge Node to an Application (see definitions below)."}),"\n",(0,n.jsx)(t.li,{children:"Sending configuration and health commands from Applications to Edge nodes to change the configuration of the edge nodes as well as receive health check from the edge nodes."}),"\n"]}),"\n",(0,n.jsxs)(t.p,{children:["To accomplish this goal, we use the following NATS JetStream subjects ",(0,n.jsx)(t.a,{href:"./setup/sever-edge-setup#nats-jetstream-subjects",children:"SpecPipe Subjects"})]}),"\n",(0,n.jsx)(t.h3,{id:"edge-nodes",children:"Edge Nodes"}),"\n",(0,n.jsx)(t.p,{children:"An edge node is a device (such as a laptop or a Raspberyy Pi) that has a software-defined radio attached to it via USB and is connected to the system. The radio is listening at a particular frequency that is initially set when the edge device registers with the system. The frequency can later be changed via commands issues by other Applications via the Controller API. See Control Flow line (3) in the diagram above."}),"\n",(0,n.jsx)(t.h3,{id:"applications",children:"Applications"}),"\n",(0,n.jsx)(t.p,{children:"Applications are software that run on devices connected to the system. An application can access the following data:"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsx)(t.li,{children:"Receiving data from edge node via the NATS (see 'Data Flow' section below) or from an edge node directly via a socket connection (see 'Peer to Peer' section below)."}),"\n",(0,n.jsx)(t.li,{children:"Geting the metadata of edge nodes (such as sampling rate, geolocation) via the Controller API."}),"\n"]}),"\n",(0,n.jsx)(t.p,{children:"Additionally, an application can also update the configuration (such as frequency and sampling rate) of an edge node (see Control Flow line 3) via the Controller API."}),"\n",(0,n.jsxs)(t.p,{children:["The aforementioned features can be accessed by using a Python Software Developemnt Toolkit provided by the SpecPipe project ",(0,n.jsx)(t.a,{href:"/specpipe/setup/application",children:"specpipe-sdk-py"}),"."]}),"\n",(0,n.jsxs)(t.p,{children:[(0,n.jsx)(t.a,{href:"./examples",children:"Here"})," are some example applications build by the SpecPipe development team."]}),"\n",(0,n.jsx)(t.p,{children:"In our diagram above, 'Edge Nodes' and 'Applications' are separate boxes. This isn't meant to imply that an Edge Node cannot be an application. Since an application is just software running on a machine connected to a system, an Edge node can also be running an application. The diagram makes this distinction solely for ease of understanding the various flows."}),"\n",(0,n.jsx)(t.h3,{id:"health-check-server",children:"Health Check Server"}),"\n",(0,n.jsx)(t.p,{children:"The Health Check Server is an example of an Application that can run on devices connected to the system. The purpose of the health check system is to check the health of the nodes."}),"\n",(0,n.jsx)(t.p,{children:"It utilizes NATS to distribute heartbeat signals across all edge devices, ensuring real-time monitoring and robustness within the network infrastructure, bolstering reliability and responsiveness in critical health monitoring scenarios."}),"\n",(0,n.jsx)(t.p,{children:"The server performs this check by having a scheduled task that sends heart beat requests using the NATS subject specpipe-cluster.<sdr_type>.<device_name> to each device in the system (see Monitoring Flow lines 5 and 6)"}),"\n",(0,n.jsx)(t.h3,{id:"controller-api",children:"Controller API"}),"\n",(0,n.jsx)(t.p,{children:"A controller API serves as the interface for managing and orchestrating resources within SpecPipe, enabling Applications to programmatically interact with and manipulate the configuration and behavior of the Edge Nodes."}),"\n",(0,n.jsx)(t.p,{children:"Applications can fetch Edge Node metadata as well as update edge node configuration by API requests to the controller API (see Control Flow lines 3 and 4)"}),"\n",(0,n.jsx)(t.h3,{id:"information-flows",children:"Information Flows"}),"\n",(0,n.jsx)(t.h4,{id:"1-data-flow-1",children:"1. Data Flow (1)"}),"\n",(0,n.jsx)(t.p,{children:"Raw IQ and Demodulated radio data is sent from Edge Nodes via NATS to Applications."}),"\n",(0,n.jsx)(t.h4,{id:"2-peer-to-peer-flow-2",children:"2. Peer to Peer Flow (2)"}),"\n",(0,n.jsx)(t.p,{children:"Raw IQ radio data is sent from Edge Nodes to Applications via a socket. This is useful for sending high bit rate Raw IQ Data since if this data would be sent via NATS, the system would get bogged down."}),"\n",(0,n.jsx)(t.h4,{id:"2-control-flow-34",children:"2. Control Flow (3,4)"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsx)(t.li,{children:"Line 3 : Applications update the configuration of an edge node (such as changing its sampling rate or frequency) via the Controller API. The Controller API then publishes a message on the appropriate NATS subject to change that setting for an Edge Node."}),"\n",(0,n.jsx)(t.li,{children:"Line 4: Applications can get edge node metadata (such as their location and sampling rate) via the Controller API."}),"\n"]}),"\n",(0,n.jsxs)(t.p,{children:["See ",(0,n.jsx)(t.a,{href:"./setup/sever-edge-setup#7-setting-up-controller-plane",children:"Controller API"})," for example calls."]}),"\n",(0,n.jsx)(t.h4,{id:"4-monitoring-flow-56",children:"4. Monitoring Flow (5,6)"}),"\n",(0,n.jsx)(t.p,{children:"Applications can monitor the health of the Edge Nodes by running the health command for a particular edge node. When this command is issued, NATS sending heart beat requests to all the Edge Nodes (line 5). Then, the edge nodes alive respond with a heartbeat (line 6) to the server via NATS"})]})}function h(e={}){const{wrapper:t}={...(0,a.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},1549:(e,t,i)=>{i.d(t,{A:()=>n});const n=i.p+"assets/images/specpipe-architecture-fece4c8e6a209cc80ca980ef344d0ccd.png"},6317:(e,t,i)=>{i.d(t,{A:()=>n});const n=i.p+"assets/images/specpipe-overview-b548bda70326d88c0dc782ad41b3ad15.png"},8453:(e,t,i)=>{i.d(t,{R:()=>o,x:()=>r});var n=i(6540);const a={},s=n.createContext(a);function o(e){const t=n.useContext(s);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function r(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:o(e.components),n.createElement(s.Provider,{value:t},e.children)}}}]);