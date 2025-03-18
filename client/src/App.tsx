import { MantineProvider, Box } from "@mantine/core";
import "./App.css";
import ShowTodo from "../components/ShowTodo";
import { Link, Route, Routes } from "react-router-dom";
import UserList from "../components/User";

export const ENDPOINT = "http://localhost:4000";

function App() {
  return (
    <MantineProvider>
      <nav className="flex justify-between items-center p-4">
        <Link to="/">
          <span className="text-2xl font-bold text-purple-600">My Todo List</span>
        </Link>
        <Link to="/users">
          <span className="text-2xl font-bold text-purple-600">Users</span>
        </Link>
      </nav>
      <Box
        style={{
          padding: "2rem",
          width: "100%",
          maxWidth: "40rem",
          margin: "0 auto",
        }}
      >
        <Routes>
          <Route
            path="/"
            element={
              <div>
                <ShowTodo />
              </div>
            }
          />
          <Route path="/users" element={<UserList />} />
        </Routes>
      </Box>
    </MantineProvider>
  );
}

export default App;
