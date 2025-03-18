import { useState } from "react";
import AddTodo from "./AddTodo";
import useSWR from "swr";
import { getTodos, markTodoDone, markTodoUndone, deleteTodo } from "../services/api";
import { CheckCircle, Circle, Trash } from "lucide-react";

export interface Todo {
  id: number;
  title: string;
  body: string;
  done: boolean;
}

// Custom fetcher using our API service
const fetcher = () => getTodos();

function ShowTodo() {
  const { data: todos = [], mutate } = useSWR<Todo[]>("todos", fetcher);
  const [isDeleting, setIsDeleting] = useState<number | null>(null);

  async function handleMarkTodoDone(id: number) {
    try {
      await markTodoDone(id);
      
      // Optimistically update the UI
      mutate(
        todos.map((todo) => 
          todo.id === id ? { ...todo, done: true } : todo
        ),
        false
      );
    } catch (error) {
      console.error("Failed to mark todo as done:", error);
    }
  }

  async function handleMarkTodoUndone(id: number) {
    try {
      await markTodoUndone(id);
      
      // Optimistically update the UI
      mutate(
        todos.map((todo) => 
          todo.id === id ? { ...todo, done: false } : todo
        ),
        false
      );
    } catch (error) {
      console.error("Failed to mark todo as undone:", error);
    }
  }

  async function handleDeleteTodo(id: number) {
    try {
      setIsDeleting(id);
      await deleteTodo(id);
      
      // Optimistically update the UI
      mutate(
        todos.filter((todo) => todo.id !== id),
        false
      );
    } catch (error) {
      console.error("Failed to delete todo:", error);
    } finally {
      setIsDeleting(null);
    }
  }

  return (
    <div className="max-w-md mx-auto p-4">
      <ul className="space-y-3 mb-6">
        {todos.length === 0 && (
          <li className="text-center text-gray-500 py-4">No todos yet. Add one below!</li>
        )}
        
        {todos.map((todo) => (
          <li 
            key={`todo__${todo.id}`}
            className={`bg-white rounded-lg shadow-md p-4 flex items-start ${todo.done ? 'border-l-4 border-green-500' : 'border-l-4 border-purple-500'}`}
          >
            <button
              onClick={() => todo.done ? handleMarkTodoUndone(todo.id) : handleMarkTodoDone(todo.id)}
              className="mt-1 mr-3 flex-shrink-0"
            >
              {todo.done ? (
                <CheckCircle className="h-6 w-6 text-green-500" />
              ) : (
                <Circle className="h-6 w-6 text-purple-500" />
              )}
            </button>
            
            <div className="flex-grow">
              <h3 className={`font-medium ${todo.done ? 'text-gray-500 line-through' : 'text-gray-800'}`}>
                {todo.title}
              </h3>
              {todo.body && (
                <p className={`mt-1 text-sm ${todo.done ? 'text-gray-400 line-through' : 'text-gray-600'}`}>
                  {todo.body}
                </p>
              )}
            </div>
            
            <button
              onClick={() => handleDeleteTodo(todo.id)}
              className="ml-2 flex-shrink-0 text-red-500 hover:text-red-700 transition-colors"
              disabled={isDeleting === todo.id}
            >
              <Trash className="h-5 w-5" />
            </button>
          </li>
        ))}
      </ul>
      
      <AddTodo mutate={mutate} />
    </div>
  );
}

export default ShowTodo;
