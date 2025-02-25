import { MantineProvider, Box, List, ThemeIcon } from "@mantine/core";
import useSWR from "swr";
import "./App.css";
import AddToDo from "../components/AddTodo";
import { CheckCircleFillIcon } from "@primer/octicons-react";

export interface Todo {
  id: number;
  title: string;
  body: string;
  done: boolean;
}

export const ENDPOINT = "http://localhost:4000";

const fetcher = async (url: string) => {
  const res = await fetch(`${ENDPOINT}/${url}`);
  if (!res.ok) throw new Error("Failed to fetch data");
  return res.json();
};

function App() {
  const { data, mutate } = useSWR<Todo[]>("api/todos", fetcher);

  const todos = Array.isArray(data) ? data : [];

  async function markTodoDone(id: number) {
    const updated = await fetch(`${ENDPOINT}/api/todos/${id}/done`, {
      method: "PATCH",
    }).then((r) => r.json());

    mutate(updated);
  }

  async function markTodoUndone(id: number) {
    const updated = await fetch(`${ENDPOINT}/api/todos/${id}/undone`, {
      method: "PATCH",
    }).then((r) => r.json());

    mutate(updated);
  }

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
        <List>
          {todos.map((todo) => {
            return (
              <List.Item
                className="flex items-center text-purple-600 hover:text-blue-400 py-2"
                onClick={() => todo.done? markTodoUndone(todo.id) : markTodoDone(todo.id)}
                key={`todo__${todo.id}`}
              >
                <ThemeIcon
                  className={
                    todo.done ? "text-teal-400 mr-2 inline-flex" : "text-gray-500 mr-2 inline-flex"
                  }
                  size={24}
                >
                  <CheckCircleFillIcon size={24} />
                </ThemeIcon>
                <span className="text-center flex-1">{todo.title}</span>
              </List.Item>
            );
          })}
        </List>

        <AddToDo mutate={mutate}></AddToDo>
      </Box>
    </MantineProvider>
  );
}

export default App;
