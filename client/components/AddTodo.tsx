import { useState } from "react";
import { KeyedMutator } from "swr";
import { Todo } from "./ShowTodo";
import { createTodo } from "../services/api";
import { Plus } from "lucide-react";

function AddTodo({ mutate }: { mutate: KeyedMutator<Todo[]> }) {
  const [open, setOpen] = useState(false);
  const [title, setTitle] = useState("");
  const [body, setBody] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  async function handleCreateTodo(e: React.FormEvent) {
    e.preventDefault();
    
    if (title.length < 3) {
      setError("Title must have at least 3 characters");
      return;
    }
    
    try {
      setIsSubmitting(true);
      setError(null);
      
      // Create todo via API service
      await createTodo({ title, body });
      
      // Re-fetch to get the updated list with the new todo
      mutate();

      setTitle("");
      setBody("");
      setOpen(false);
    } catch (error) {
      console.error("Failed to create todo:", error);
      setError("Failed to create todo. Please try again later.");
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <div>
      {open ? (
        <div className="bg-white rounded-lg shadow-md p-4">
          <h2 className="text-lg font-semibold mb-4 text-purple-600">Create New Todo</h2>
          
          <form onSubmit={handleCreateTodo}>
            <div className="mb-4">
              <label htmlFor="title" className="block text-sm font-medium text-gray-700 mb-1">
                Title
              </label>
              <input
                type="text"
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="What do you want to do?"
                required
              />
            </div>
            
            <div className="mb-4">
              <label htmlFor="body" className="block text-sm font-medium text-gray-700 mb-1">
                Detail (optional)
              </label>
              <textarea
                id="body"
                value={body}
                onChange={(e) => setBody(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-purple-500"
                placeholder="Add more details..."
                rows={3}
              />
            </div>
            
            {error && (
              <div className="mb-4 p-2 bg-red-100 border border-red-400 text-red-700 rounded-md text-sm">
                {error}
              </div>
            )}
            
            <div className="flex justify-end space-x-2">
              <button
                type="button"
                onClick={() => setOpen(false)}
                className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-100 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-4 py-2 bg-purple-600 rounded-md text-white hover:bg-purple-700 transition-colors flex items-center"
              >
                {isSubmitting ? "Creating..." : "Create Todo"}
              </button>
            </div>
          </form>
        </div>
      ) : (
        <button
          onClick={() => setOpen(true)}
          className="w-full py-3 bg-purple-600 text-white rounded-md hover:bg-purple-700 transition-colors flex items-center justify-center"
        >
          <Plus className="h-5 w-5 mr-2" />
          Add New Todo
        </button>
      )}
    </div>
  );
}

export default AddTodo;