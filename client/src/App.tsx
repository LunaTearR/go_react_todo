import { MantineProvider, Box } from "@mantine/core";
import "./App.css";
import ShowTodo from "../components/ShowTodo";
import { Link, Route, Routes } from "react-router-dom";
import UserList from "../components/User";

export const ENDPOINT = "http://localhost:4000";

function App() {
  return (
    <MantineProvider>
      <Box
        style={{
          padding: "2rem",
          width: "100%",
          maxWidth: "40rem",
          margin: "0 auto",
        }}
      >
        <nav>
          <Link to="/">Home</Link> | <Link to="/users">Users</Link>
        </nav>

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
