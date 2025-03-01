import { useState, useEffect } from "react";

function Server() {
  const [message, setMessage] = useState("");
  const [recMessage, setRecMessage] = useState("Loading...");

  useEffect(() => {
    fetch("http://localhost:8080")
      .then((response) => response.json()) 
      .then((data) => setRecMessage(data.message)) 
      .catch((error) => console.error("Error fetching data:", error));
  }, []);

  const sendDataToServer = async () => {
    const data = { message }; 

    const response = await fetch("http://localhost:8080/save", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();
    alert(result.message); 
  };

  return (
    <div>
     <div style={{ textAlign: "center", marginTop: "50px" }}>
      <h1>Message</h1>
      <p>{recMessage}</p>
    </div>
      <h1>Send Data to C Backend</h1>
      <input
        type="text"
        placeholder="Enter message"
        value={message}
        onChange={(e) => setMessage(e.target.value)}
      />
      <button onClick={sendDataToServer}>Send</button>
    </div>
  );
}

export default Server;
