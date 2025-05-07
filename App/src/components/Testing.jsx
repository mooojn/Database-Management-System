import { useState } from "react";
import './styles/ButtonNav.css';

function Testing() {
  const [queryText, setQueryText] = useState("");
  const [tables, setTables] = useState({});
  const [visible, setVisible] = useState(false);
  const [currentDatabase, setCurrentDatabase] = useState(null); // ✅ store selected DB

  const parseQuery = async (query) => {
    const tokens = query.trim().replace(/;$/, "").split(/\s+/);
    const command = tokens[0]?.toUpperCase();

    switch (command) {
      case "CREATE":
        if (tokens[1]?.toUpperCase() === "DATABASE") {
          await handleDatabaseAction("create_db", tokens[2]);
        } else if (tokens[1]?.toUpperCase() === "TABLE") {
          if (!currentDatabase) return alert("Use `USE dbName;` to select a database first.");
          const tableName = tokens[2];
          const columnsMatch = query.match(/\(([^)]+)\)/);
          const columns = columnsMatch ? columnsMatch[1].split(",").map(col => col.trim()) : [];
          await handleTableAction("create_table", tableName, columns);
        }
        break;

      case "DELETE":
        if (tokens[1]?.toUpperCase() === "DATABASE") {
          await handleDatabaseAction("delete_db", tokens[2]);
        } else if (tokens[1]?.toUpperCase() === "TABLE") {
          if (!currentDatabase) return alert("Use `USE dbName;` to select a database first.");
          await handleTableAction("delete_table", tokens[2]);
        }
        break;

      case "MODIFY":
        if (tokens[1]?.toUpperCase() === "DATABASE") {
          const oldName = tokens[2];
          const newName = tokens[3];
          await handleDatabaseAction("modify_db", oldName, newName);
        }
        break;

      case "USE": // ✅ Handle `USE database_name;`
        const dbToUse = tokens[1];
        if (!dbToUse) return alert("Specify a database to use.");
        setCurrentDatabase(dbToUse);
        alert(`Now using database: ${dbToUse}`);
        break;

      case "SELECT": // ✅ Basic SELECT query
        if (!currentDatabase) return alert("Select a database first with `USE dbName;`");
        const selectTable = tokens[tokens.length - 1];
        await handleTableAction("select_table", selectTable);
        break;

      case "UPDATE": // ✅ Basic UPDATE query
        if (!currentDatabase) return alert("Select a database first with `USE dbName;`");
        const updateTable = tokens[1];
        const setMatch = query.match(/SET\s+(.*?)\s+WHERE/i);
        const whereMatch = query.match(/WHERE\s+(.*)/i);
        const setClause = setMatch?.[1];
        const whereClause = whereMatch?.[1];
        if (!setClause || !whereClause) return alert("Malformed UPDATE query.");

        await handleTableAction("update_table", updateTable, [], { setClause, whereClause });
        break;

      default:
        alert("Unknown query or syntax error.");
        break;
    }
  };

  const handleDatabaseAction = async (action, dbName, newName = null) => {
    if (!dbName) return alert("Database name missing.");
    let url = `http://localhost:8080/${action}`;
    let options = { method: "POST" };

    if (action === "modify_db") {
      options.headers = { "Content-Type": "application/json" };
      options.body = JSON.stringify({ oldName: dbName, newName });
    } else {
      url += `?dbName=${encodeURIComponent(dbName)}`;
    }

    const response = await fetch(url, options);
    if (response.ok) {
      alert(`Database ${action}ed successfully.`);
      loadDatabase(); // refresh
    } else {
      alert(`Failed to ${action} database.`);
    }
  };

  const handleTableAction = async (action, tableName, columns = [], extras = {}) => {
    if (!currentDatabase) return alert("No database selected.");

    const baseUrl = `http://localhost:8080/${action}?dbName=${encodeURIComponent(currentDatabase)}&tableName=${tableName}`;
    let url = baseUrl;
    const options = { method: "POST" };

    if (action === "create_table") {
      const colParam = columns.filter(Boolean).map(col => `columns=${col}`).join("&");
      url = `${baseUrl}&${colParam}`;
    } else if (action === "update_table") {
      options.headers = { "Content-Type": "application/json" };
      options.body = JSON.stringify({
        set: extras.setClause,
        where: extras.whereClause
      });
    }

    const response = await fetch(url, options);
    if (response.ok) {
      alert(`Table ${action}ed successfully.`);
      loadDatabase(); // refresh
    } else {
      alert(`Failed to ${action} table.`);
    }
  };

  const loadDatabase = async () => {
    try {
      const res = await fetch('http://localhost:8080/load-database');
      const data = await res.json();
      setTables(data);
      setVisible(true);
    } catch (err) {
      console.error('Failed to fetch database:', err);
    }
  };

  const styles = {
    container: {
      padding: "20px",
      fontFamily: "monospace",
      background: "#1e1e1e",
      color: "#c3f3c3",
      minHeight: "100vh",
    },
    consoleBox: {
      background: "#2b2b2b",
      padding: "20px",
      borderRadius: "10px",
      marginBottom: "20px",
    },
    textarea: {
      width: "100%",
      height: "100px",
      background: "#1e1e1e",
      color: "#c3f3c3",
      border: "1px solid #333",
      borderRadius: "4px",
      padding: "10px",
      fontSize: "16px",
      resize: "none",
    },
    button: {
      marginTop: "10px",
      padding: "10px 20px",
      background: "#00bc8c",
      border: "none",
      borderRadius: "5px",
      color: "#fff",
      cursor: "pointer",
    },
    tableViewer: {
      background: "#2b2b2b",
      padding: "20px",
      borderRadius: "10px",
      marginTop: "20px",
    },
    table: {
      width: "100%",
      borderCollapse: "collapse",
      color: "#fff",
    },
    th: {
      backgroundColor: "#00bc8c",
      padding: "10px",
    },
    td: {
      padding: "10px",
      borderBottom: "1px solid #444",
    },
  };

  return (
    <div style={styles.container}>
      <div style={styles.consoleBox}>
        <h2>Query Console {currentDatabase && `(Using: ${currentDatabase})`}</h2>
        <textarea
          style={styles.textarea}
          placeholder='E.g., CREATE DATABASE testdb;'
          value={queryText}
          onChange={(e) => setQueryText(e.target.value)}
        />
        <button style={styles.button} onClick={() => parseQuery(queryText)}>Run Query</button>
      </div>

      <div style={styles.tableViewer}>
        <h2>Database Viewer</h2>
        <button style={styles.button} onClick={loadDatabase}>Load Database</button>
        {visible && (
          <table style={styles.table}>
            <thead>
              <tr>
                <th style={styles.th}>Table Name</th>
                <th style={styles.th}>Columns</th>
              </tr>
            </thead>
            <tbody>
              {Object.entries(tables || {}).map(([tableName, table]) => (
                <tr key={tableName}>
                  <td style={styles.td}>{tableName}</td>
                  <td style={styles.td}>{(table.columns || []).join(", ")}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

export default Testing;
