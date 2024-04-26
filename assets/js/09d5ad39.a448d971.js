"use strict";(self.webpackChunkdocs_website=self.webpackChunkdocs_website||[]).push([[864],{7223:(e,t,i)=>{i.r(t),i.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>d,frontMatter:()=>n,metadata:()=>p,toc:()=>r});var s=i(4848),a=i(8453);const n={title:"Example Applications",sidebar_position:3},o="Example SpecPipe Applications",p={id:"examples",title:"Example Applications",description:"1. FM Audio",source:"@site/docs/examples.md",sourceDirName:".",slug:"/examples",permalink:"/specpipe/examples",draft:!1,unlisted:!1,tags:[],version:"current",sidebarPosition:3,frontMatter:{title:"Example Applications",sidebar_position:3},sidebar:"tutorialSidebar",previous:{title:"Setup SpecPipe Application",permalink:"/specpipe/setup/application"},next:{title:"Architecture",permalink:"/specpipe/architecture"}},c={},r=[{value:"1. FM Audio",id:"1-fm-audio",level:3},{value:"2. IQ Visualization &amp; IQEngine Integration",id:"2-iq-visualization--iqengine-integration",level:3},{value:"3. Speech to Text",id:"3-speech-to-text",level:3},{value:"4. Audio Data Mocking &amp; Prometheus Exporter",id:"4-audio-data-mocking--prometheus-exporter",level:3},{value:"5. IQ Data Mocking and ADSB",id:"5-iq-data-mocking-and-adsb",level:3},{value:"6. Speech To Text Web Socket",id:"6-speech-to-text-web-socket",level:3}];function l(e){const t={a:"a",h1:"h1",h3:"h3",p:"p",...(0,a.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.h1,{id:"example-specpipe-applications",children:"Example SpecPipe Applications"}),"\n",(0,s.jsx)(t.h3,{id:"1-fm-audio",children:"1. FM Audio"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/audio_play",children:"This example"})," demonstrates how to subscribe to live audio data streams captured and processed by SpecPipe and play the audio chunks locally."]}),"\n",(0,s.jsx)(t.h3,{id:"2-iq-visualization--iqengine-integration",children:"2. IQ Visualization & IQEngine Integration"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/plot_iq",children:"This example"})," demonstrates how to visualize the IQ spectrum data captured by SpecPipe."]}),"\n",(0,s.jsx)(t.h3,{id:"3-speech-to-text",children:"3. Speech to Text"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/speech2text",children:"This example"})," demonstrates converting speech received via FM to text."]}),"\n",(0,s.jsx)(t.h3,{id:"4-audio-data-mocking--prometheus-exporter",children:"4. Audio Data Mocking & Prometheus Exporter"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/mock_fm",children:"This example"})," streams the content of mock_audio.wav file (sampled at 32 KHz) to NATS JetStream circularly to simulate continuous audio collection. To monitor this data, this example further runs the Prometheus exporter."]}),"\n",(0,s.jsx)(t.h3,{id:"5-iq-data-mocking-and-adsb",children:"5. IQ Data Mocking and ADSB"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/iq_adsb",children:"This example"})," first streams the content of a file (exampledata.bin) containg IQ Data collected at the frequency of 1090 MHz, sampling rate of 2 MHz, and gain of 50 collected for 10 seconds to NATS Jetstream on a loop. Then it creates a client to fetch this IQ Data from Jetstream, processes it via Dump1090, decodes the output, and prints valid ADSB signal content to the console."]}),"\n",(0,s.jsx)(t.h3,{id:"6-speech-to-text-web-socket",children:"6. Speech To Text Web Socket"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/ml4wireless/specpipe/tree/main/_examples/speech-to-text-ws",children:"This example"})," is a simple FastAPI web server fetches FM data from NATS, decodes it into .wav chunks suitable for the SpeechRecognition model, and streams the text to any listeners via websocket"]})]})}function d(e={}){const{wrapper:t}={...(0,a.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(l,{...e})}):l(e)}},8453:(e,t,i)=>{i.d(t,{R:()=>o,x:()=>p});var s=i(6540);const a={},n=s.createContext(a);function o(e){const t=s.useContext(n);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function p(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:o(e.components),s.createElement(n.Provider,{value:t},e.children)}}}]);