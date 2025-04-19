import { useState } from "react";
import './styles/ButtonNav.css';

function Testing() {
    const [dbName, setDbName] = useState("");
    const [tableName, setTableName] = useState("");
    const [columns, setColumns] = useState([""]); // Array to store column names
    const [tables, setTables] = useState({});
  const [visible, setVisible] = useState(false);

  const loadDatabase = async () => {
    try {
      const res = await fetch('http://localhost:8080/load-database');
      const data = await res.json();
      console.log('Received Tables:', data);
      setTables(data);
      setVisible(true);
    } catch (err) {
      console.error('Failed to fetch database:', err);
    }
  };
  
    const handleDatabaseAction = async (action) => {
        if (!dbName) {
            alert("Please provide a database name.");
            return;
        }
        try {
            let url = `http://localhost:8080/${action}?dbName=${dbName}`;
            const response = await fetch(url, { method: "POST" });
            if (response.ok) {
                alert(`Database ${action}ed successfully.`);
            } else {
                const errorText = await response.text();
                alert(`Failed to ${action} database: ${errorText}`);
            }
        } catch (error) {
            console.error("Error performing database action:", error);
            alert("Error performing database action");
        }
    };

    const handleTableAction = async (action) => {
        if (!tableName) {
            alert("Please provide a table name.");
            return;
        }

        // Join columns as query parameters to pass them to the backend
        const columnsQuery = columns.filter(col => col).join("&columns="); // Filtering out empty columns

        try {
            let url = `http://localhost:8080/${action}?dbName=${dbName}&tableName=${tableName}&columns=${columnsQuery}`;
            const response = await fetch(url, { method: "POST" });
            if (response.ok) {
                alert(`Table ${action}ed successfully.`);
            } else {
                const errorText = await response.text();
                alert(`Failed to ${action} table: ${errorText}`);
            }
        } catch (error) {
            console.error("Error performing table action:", error);
            alert("Error performing table action");
        }
    };

    const handleAddColumn = () => {
        setColumns([...columns, ""]); // Add an empty column input field
    };

    const handleColumnChange = (index, value) => {
        const updatedColumns = [...columns];
        updatedColumns[index] = value;
        setColumns(updatedColumns);
    };

    const handleTransaction = async (action) => {
        try {
            const response = await fetch(`http://localhost:8080/${action}`, { method: "GET" });
            if (response.ok) {
                alert(`Transaction ${action}ed`);
            } else {
                const errorText = await response.text();
                alert(`Failed to ${action} transaction: ${errorText}`);
            }
        } catch (error) {
            console.error("Error performing transaction action:", error);
            alert(`Error performing transaction action: ${error}`);
        }
    };

    return (
        <div className="buttonNav">
            <h2>Database Operations</h2>
            <input
                type="text"
                placeholder="Enter Database Name"
                value={dbName}
                onChange={(e) => setDbName(e.target.value)}
            />
            <div>
                <button onClick={() => handleDatabaseAction("create_db")}>Create Database</button>
                <button onClick={() => handleDatabaseAction("delete_db")}>Delete Database</button>
                <button onClick={() => handleDatabaseAction("modify_db")}>Modify Database</button>
            </div>

            <h2>Table Operations</h2>
            <input
                type="text"
                placeholder="Enter Table Name"
                value={tableName}
                onChange={(e) => setTableName(e.target.value)}
            />

            {/* Column input fields */}
            <div>
                <h3>Columns</h3>
                {columns.map((column, index) => (
                    <div key={index}>
                        <input
                            type="text"
                            placeholder={`Column ${index + 1}`}
                            value={column}
                            onChange={(e) => handleColumnChange(index, e.target.value)}
                        />
                    </div>
                ))}
                <button onClick={handleAddColumn}>Add Column</button>
            </div>

            <div>
                <button onClick={() => handleTableAction("create_table")}>Create Table</button>
                <button onClick={() => handleTableAction("delete_table")}>Delete Table</button>
                <button onClick={() => handleTableAction("edit_table")}>Edit Table</button>
            </div>

            <h2>Transaction Operations</h2>
            <div>
                <button onClick={() => handleTransaction("start")}>Start Transaction</button>
                <button onClick={() => handleTransaction("commit")}>Commit Transaction</button>
                <button onClick={() => handleTransaction("rollback")}>Rollback Transaction</button>
            </div>
            <div>
      <button onClick={loadDatabase}>Show Database</button>

      {visible && (
        <div style={{ marginTop: '20px' }}>
          <h2>Loaded Tables</h2>
          <table style={{ marginTop: "20px" }}>
        <tbody>
          {Object.entries(tables || {}).map(([tableName, table]) => (
            <tr key={tableName}>
              <td>{tableName}</td>
              <td>{(table.columns || []).join(", ")}</td>
            </tr>
          ))}
        </tbody>
      </table>

        </div>
      )}
    </div>
        </div>
    );
}

export default Testing;
