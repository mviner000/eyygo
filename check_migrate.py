import sqlite3

def main():
    # Replace with your database file path
    db_path = "db_migrate.sqlite3"

    # Connect to the SQLite database
    try:
        conn = sqlite3.connect(db_path)
    except sqlite3.Error as e:
        print(f"Failed to connect to database: {e}")
        return

    # Check if the database is reachable
    try:
        conn.execute("SELECT 1")
    except sqlite3.Error as e:
        print(f"Failed to ping database: {e}")
        conn.close()
        return

    # Query the schema for the auth_users table
    query = "PRAGMA table_info(eyygo_session)"
    try:
        cursor = conn.execute(query)
    except sqlite3.Error as e:
        print(f"Failed to query table info: {e}")
        conn.close()
        return

    # Print the field names and their data types
    print("Field Name\tData Type")
    print("------------------------")

    rows = cursor.fetchall()
    for row in rows:
        # row[1] is the field name and row[2] is the data type
        print(f"{row[1]}\t\t{row[2]}")

    # Close the database connection
    conn.close()

if __name__ == "__main__":
    main()
