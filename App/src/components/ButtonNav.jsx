import { useState } from "react";
import './styles/ButtonNav.css';

function ButtonNav() {
    const [transactionActive, setTransactionActive] = useState(false);

    const buttons = [
        'Open Shell',
        'Create Database',
        'Start Transaction',
    ];

    const cancelTransaction = async () => {
        if (!transactionActive) {
            alert("No active transaction to cancel.");
            return;
        }
        try {
            const response = await fetch("http://localhost:8080/cancel", { method: "GET" });

            if (response.ok) {
                setTransactionActive(false);
                alert("Transaction cancelled.");
            } else {
                const errorText = await response.text();
                alert("Failed to cancel transaction: " + errorText);
            }
        } catch (error) {
            console.error("Error cancelling transaction:", error);
            alert("Error cancelling transaction");
        }
    };

    const executeAllTransactions = async () => {
        if (!transactionActive) {
            alert("Start a transaction first!");
            return;
        }
        
        try {
            // Perform multiple operations
            await performOperation("add", 5, 3);
            await performOperation("sub", 10, 4);
            await performOperation("mul", 6, 2);
            await performOperation("div", 8, 2);
    
            // Commit transaction
            await commitTransaction();
            alert("All transactions executed and committed.");
        } catch (error) {
            console.error("Error executing all transactions:", error);
            alert("Error executing all transactions.");
        }
    };
    
    const handleClick = async (val) => {
        if (val === "Start Transaction") {
            try {
                const response = await fetch("http://localhost:8080/start", {
                    method: "GET",
                });

                if (response.ok) {
                    setTransactionActive(true);
                    alert("Transaction started");
                } else {
                    const errorText = await response.text();
                    alert("Failed to start transaction: " + errorText);
                }
            } catch (error) {
                console.error("Error connecting to backend:", error);
                alert("Error connecting to backend");
            }
        }
    };

    const performOperation = async (op, a, b) => {
        if (!transactionActive) {
            alert("Start a transaction first!");
            return;
        }
        try {
            const response = await fetch(`http://localhost:8080/operate?operation=${op}&a=${a}&b=${b}`);
            const result = await response.text();
            alert(result);
        } catch (error) {
            console.error("Error performing operation:", error);
            alert("Error performing operation");
        }
    };

    const commitTransaction = async () => {
        if (!transactionActive) {
            alert("No active transaction to commit.");
            return;
        }
        try {
            const response = await fetch("http://localhost:8080/commit", { method: "GET" });

            if (response.ok) {
                setTransactionActive(false);
                alert("Transaction committed and saved to CSV");
            } else {
                const errorText = await response.text();
                alert("Failed to commit transaction: " + errorText);
            }
        } catch (error) {
            console.error("Error committing transaction:", error);
            alert("Error committing transaction");
        }
    };

    return (
        <div className="buttonNav">
            {buttons.map((val, idx) => (
                <button key={idx} onClick={() => handleClick(val)}>{val}</button>
            ))}
            {transactionActive && (
                <div>
                    <button onClick={() => performOperation("add", 5, 3)}>Add 5 + 3</button>
                    <button onClick={() => performOperation("sub", 10, 4)}>Subtract 10 - 4</button>
                    <button onClick={() => performOperation("mul", 2, 15)}>Subtract 2 * 15</button>
                    <button onClick={() => performOperation("div", 10, 2)}>Subtract 10 / 2</button>

                    <button onClick={commitTransaction} style={{ backgroundColor: "green", color: "white" }}>
                        Commit Transaction
                    </button>
                    <button onClick={cancelTransaction} style={{ backgroundColor: "red", color: "white" }}>
                        Cancel Transaction
                    </button>
                    <button onClick={executeAllTransactions} style={{ backgroundColor: "green", color: "white" }}>
                        Execute All & Save
                    </button>
                </div>
            )}

        </div>
    );
}

export default ButtonNav;
