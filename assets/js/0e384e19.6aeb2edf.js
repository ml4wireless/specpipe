"use strict";(self.webpackChunkdocs_website=self.webpackChunkdocs_website||[]).push([[976],{1512:(e,i,t)=>{t.r(i),t.d(i,{assets:()=>c,contentTitle:()=>o,default:()=>l,frontMatter:()=>s,metadata:()=>r,toc:()=>p});var n=t(4848),a=t(8453);const s={sidebar_position:1,slug:"/",title:"Introduction"},o="SpecPipe: A scalable AI/ML-facilitating data pipeline for spectrum",r={id:"intro",title:"Introduction",description:"By: Omair Alam, Will Almy, Alice Lee and Huo-Ming Shu",source:"@site/docs/intro.md",sourceDirName:".",slug:"/",permalink:"/specpipe/",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:1,frontMatter:{sidebar_position:1,slug:"/",title:"Introduction"},sidebar:"tutorialSidebar",next:{title:"Architecture",permalink:"/specpipe/architecture"}},c={},p=[{value:"Setup Guide",id:"setup-guide",level:2},{value:"Setting up SpecPipe Server and Edge Device",id:"setting-up-specpipe-server-and-edge-device",level:3},{value:"1. Setting up NATS on Amazon Web Services",id:"1-setting-up-nats-on-amazon-web-services",level:4},{value:"2. Connecting a Software Defined Radio (SDR) to a machine",id:"2-connecting-a-software-defined-radio-sdr-to-a-machine",level:4},{value:"3. Setting up SpecPipe Edge Device",id:"3-setting-up-specpipe-edge-device",level:4},{value:"4. Creating Monitoring Dashboard",id:"4-creating-monitoring-dashboard",level:4},{value:"Building Specpipe Applicataion",id:"building-specpipe-applicataion",level:3}];function d(e){const i={a:"a",h1:"h1",h2:"h2",h3:"h3",h4:"h4",hr:"hr",p:"p",...(0,a.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(i.h1,{id:"specpipe-a-scalable-aiml-facilitating-data-pipeline-for-spectrum",children:"SpecPipe: A scalable AI/ML-facilitating data pipeline for spectrum"}),"\n",(0,n.jsx)(i.p,{children:"By: Omair Alam, Will Almy, Alice Lee and Huo-Ming Shu"}),"\n",(0,n.jsx)(i.hr,{}),"\n",(0,n.jsx)(i.p,{children:"In today's interconnected world, radio spectrum signals surround us, yet there exist noticeable limitations in the data systems created to access, monitor, perform AI experiments, and contribute to this analog data."}),"\n",(0,n.jsxs)(i.p,{children:["To democratize the access and usage of spectrum data, we have built ",(0,n.jsx)(i.a,{href:"https://github.com/ml4wireless/specpipe",children:"SpecPipe"}),", a distributed AI/ML data pipeline. This platform\u2019s core values of accessibility, extensibility and scalability ensure that individual users can start to work with radio data with inexpensive hardware, minimal configuration, and a smooth onboarding process."]}),"\n",(0,n.jsx)(i.p,{children:"We have accomplished this goal of improving access to spectrum data by building SpecPipe as an open-source project free for people to access and use, with easy to follow documentation, and a plethora of startup examples that allow users to understand our framework interactively."}),"\n",(0,n.jsxs)(i.p,{children:["For more details on the architecture of SpecPipe ",(0,n.jsx)(i.a,{href:"./architecture",children:"click here"}),"."]}),"\n",(0,n.jsx)(i.hr,{}),"\n",(0,n.jsx)(i.h2,{id:"setup-guide",children:"Setup Guide"}),"\n",(0,n.jsx)(i.h3,{id:"setting-up-specpipe-server-and-edge-device",children:"Setting up SpecPipe Server and Edge Device"}),"\n",(0,n.jsxs)(i.p,{children:["Note: This only applies if you want to setup to send Radio Data as an ",(0,n.jsx)(i.a,{href:"./architecture#edge-nodes",children:"Edge Node"}),". If you only want to build an application that accesses data, you don't need to have an SDR."]}),"\n",(0,n.jsx)(i.h4,{id:"1-setting-up-nats-on-amazon-web-services",children:"1. Setting up NATS on Amazon Web Services"}),"\n",(0,n.jsxs)(i.p,{children:["Note: This is only required if you want to deploy your server on AWS. If you want a local NATS installation, ",(0,n.jsx)(i.a,{href:"./#3-setting-up-specpipe-edge-device",children:"step 3"})," will guide you through that."]}),"\n",(0,n.jsxs)(i.p,{children:["The ",(0,n.jsx)(i.a,{href:"/specpipe/NATS",children:"NATS AWS tutorial"})," tutorial will help you setup NATS on AWS that you will be able to then use when setting up your SpecPipe Edge devices in ",(0,n.jsx)(i.a,{href:"./#3-setting-up-specpipe-edge-device",children:"step 3"}),"."]}),"\n",(0,n.jsx)(i.h4,{id:"2-connecting-a-software-defined-radio-sdr-to-a-machine",children:"2. Connecting a Software Defined Radio (SDR) to a machine"}),"\n",(0,n.jsx)(i.p,{children:"The SDR will be collecting data in real-time."}),"\n",(0,n.jsxs)(i.p,{children:["For the hardware setup, see: ",(0,n.jsx)(i.a,{href:"https://www.youtube.com/watch?v=uM8NkB2nIis",children:"https://www.youtube.com/watch?v=uM8NkB2nIis"})]}),"\n",(0,n.jsxs)(i.p,{children:["You can follow this ",(0,n.jsx)(i.a,{href:"https://www.youtube.com/watch?v=bT2WZhKBkRk",children:"software setup"})," to confirm that your SDR is functioning correctly. SpecPipe doesn't make use of this software so you won't need it for anything besides this."]}),"\n",(0,n.jsx)(i.h4,{id:"3-setting-up-specpipe-edge-device",children:"3. Setting up SpecPipe Edge Device"}),"\n",(0,n.jsxs)(i.p,{children:["Please follow ",(0,n.jsx)(i.a,{href:"https://github.com/ml4wireless/specpipe?tab=readme-ov-file#getting-started",children:"these instructions"})," to setup an Edge device for sending Radio Data using the SDR added in Step 1."]}),"\n",(0,n.jsx)(i.h4,{id:"4-creating-monitoring-dashboard",children:"4. Creating Monitoring Dashboard"}),"\n",(0,n.jsxs)(i.p,{children:["To create a monitoring dashboard for your system, follow ",(0,n.jsx)(i.a,{href:"/specpipe/dashboard",children:"this"})," tutorial."]}),"\n",(0,n.jsx)(i.h3,{id:"building-specpipe-applicataion",children:"Building Specpipe Applicataion"}),"\n",(0,n.jsxs)(i.p,{children:["This step only applies if you would like to build a ",(0,n.jsx)(i.a,{href:"./architecture#applications",children:"Specpipe Application"})]}),"\n",(0,n.jsxs)(i.p,{children:["Please follow ",(0,n.jsx)(i.a,{href:"https://github.com/ml4wireless/specpipe-sdk-py?tab=readme-ov-file#installation--usage",children:"these instructions"})," to create a SpecPipe Application."]})]})}function l(e={}){const{wrapper:i}={...(0,a.R)(),...e.components};return i?(0,n.jsx)(i,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},8453:(e,i,t)=>{t.d(i,{R:()=>o,x:()=>r});var n=t(6540);const a={},s=n.createContext(a);function o(e){const i=n.useContext(s);return n.useMemo((function(){return"function"==typeof e?e(i):{...i,...e}}),[i,e])}function r(e){let i;return i=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:o(e.components),n.createElement(s.Provider,{value:i},e.children)}}}]);