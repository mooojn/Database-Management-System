import { useState, useEffect } from "react";

function App() {
  const [message, setMessage] = useState("");
  const [recMessage, setRecMessage] = useState("Loading...");

  useEffect(() => {
    fetch("http://localhost:8080")
      .then((response) => response.json()) // Convert response to JSON
      .then((data) => setRecMessage(data.message)) // Extract "message" field
      .catch((error) => console.error("Error fetching data:", error));
  }, []);

  const sendDataToServer = async () => {
    const data = { message }; // JSON object to send

    const response = await fetch("http://localhost:8080/save", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(data),
    });

    const result = await response.json();
    alert(result.message); // Show success or error message
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

export default App;
