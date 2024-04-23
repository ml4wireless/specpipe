import React, { useState, useRef, useEffect } from "react";

const App = () => {
  const [messages, setMessages] = useState([]);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    // Initialize WebSocket connection
    const ws = new WebSocket('ws://0.0.0.0:8000/ws_fm_speech');
    
    // Event listener for receiving messages
    ws.onmessage = (event) => {
      console.log(event.data);
      setMessages((prevMessages) => [...prevMessages, event.data]);
    };
    
    // Clean up function when the component unmounts
    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    // Scroll to the bottom of the message container on new message
    messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const style = {
    container: {
      maxHeight: '200px',
      overflowY: 'auto',
      border: '1px solid #e0e0e0',
      padding: '10px',
      marginTop: '20px',
      backgroundColor: '#f9f9f9',
      borderRadius: '8px'
    },
    header: {
      color: '#333',
      fontFamily: 'Arial, sans-serif'
    },
    message: {
      color: '#555',
      fontFamily: 'Arial, sans-serif',
      fontSize: '14px',
      padding: '5px',
      borderBottom: '1px solid #eee',
      marginBottom: '5px'
    }
  };

  return (
    <>
      <h1 style={style.header}>SpecPipe Speech To Text Live Feed</h1>
      <div style={style.container}>
        {messages.map((message, index) => (
          <p key={index} style={style.message}>{message}</p>
        ))}
        <div ref={messagesEndRef} />
      </div>
    </>
  );
};

export default App;
