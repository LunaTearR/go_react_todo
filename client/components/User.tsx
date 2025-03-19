// src/components/UserList.jsx
import { useState, useEffect } from "react";
import { deleteUser, getUsers } from "../services/api";
import React from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faUserPlus } from "@fortawesome/free-solid-svg-icons";
import ModalAddUser from "./ModalUser";
import { mutate } from "swr";
import Swal from "sweetalert2";
import { Trash } from "lucide-react";

export interface User {
  id: number;
  username: string;
  email: string;
}

const UserList = () => {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isDeleting, setIsDeleting] = useState<number | null>(null);
  const [swalProps, setSwalProps] = useState({});
  const fetchUsers = async () => {
    try {
      setLoading(true);
      const data = await getUsers();
      setUsers(data);
      setError(null);
    } catch (err) {
      setError("Failed to fetch users");
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  useEffect(() => {
    const handleUserAdded = () => {
      fetchUsers();
    };
    window.addEventListener("user-added", handleUserAdded);
    return () => {
      window.removeEventListener("user-added", handleUserAdded);
    };
  }, []);

  if (loading) {
    return <div className="text-center py-4">Loading...</div>;
  }

  if (error) {
    return <div className="text-center py-4 text-red-500">{error}</div>;
  }

  if (users.length === 0) {
    return (
      <div className="text-center py-4 text-gray-500">
        No users yet. Add one above!
      </div>
    );
  }
  const mutateUsers = async () => {
    await fetchUsers();
    return users;
  };

  async function handleDeleteUser(id: number) {
    const result = await Swal.fire({
      title: "Are you sure?",
      text: "You won't be able to recover this user!",
      icon: "warning",
      showCancelButton: true,
      confirmButtonColor: "#3085d6",
      cancelButtonColor: "#d33",
      confirmButtonText: "Yes, delete it!",
    });
    if (result.isConfirmed) {
      try {
        setIsDeleting(id);

        Swal.fire({
          title: "Deleting...",
          text: "Please wait",
          allowOutsideClick: false,
          didOpen: () => {
            Swal.showLoading();
          },
        });
        await deleteUser(id);

        mutate(
          users.map((user) => user.id !== id),
          false
        );

        Swal.fire("Deleted!", "User has been deleted successfully.", "success");
      } catch (error) {
        console.error("Failed to delete user:", error);

        Swal.fire("Error!", "Failed to delete the user.", "error");
      } finally {
        fetchUsers();
        setIsDeleting(null);
      }
    }
  }
  return (
    <>
      <div className="flex justify-end py-2">
        <button
          type="button"
          onClick={() => setIsModalOpen(true)}
          className="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-100 transition-colors justify-items-end"
        >
          <FontAwesomeIcon icon={faUserPlus} />
          <span className="p-2">Add user</span>
        </button>
      </div>
      <ModalAddUser
        mutate={mutateUsers}
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
      />
      <div className="bg-white shadow overflow-hidden sm:rounded-lg">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                ID
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Username
              </th>
              <th
                scope="col"
                className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
              >
                Email
              </th>
            </tr>
          </thead>

          <tbody className="bg-white divide-y divide-gray-200">
            {users?.map((user) => (
              <tr key={user.id}>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {user.id}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  {user.username}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {user.email}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm">
                  <button
                    onClick={() => handleDeleteUser(user.id)}
                    className="ml-2 flex-shrink-0 text-red-500 hover:text-red-700 transition-colors"
                    disabled={isDeleting === user.id}
                  >
                    <Trash className="h-5 w-5" />
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </>
  );
};

export default UserList;
