import { Table } from 'flowbite-react';
import api from '../axios';
import { useEffect, useState } from 'react';

interface User {
  _id: string;
  email: string;
  password: string;
}

const UserList: React.FC = () => {
  const [users, setUsers] = useState<User[]>([]);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await api.get('/users');
      setUsers(response.data);
    } catch (error) {
      console.error('Failed to fetch users:', error);
    }
  };

  const handleDeleteUser = async (id: string) => {
    try {
      await api.delete(`/users/${id}`);
      fetchUsers(); // Refresh the list after deletion
    } catch (error) {
      console.error('Failed to delete user:', error);
    }
  };

  return (
    <Table>
      <Table.Head>
        <Table.HeadCell>Email</Table.HeadCell>
        <Table.HeadCell>Actions</Table.HeadCell>
      </Table.Head>
      <Table.Body>
        {users.map((user) => (
          <Table.Row key={user._id}>
            <Table.Cell>{user.email}</Table.Cell>
            <Table.Cell>
              <button
                className="font-medium text-red-600 dark:text-red-500 hover:underline ml-4"
                onClick={() => handleDeleteUser(user._id)}
              >
                Delete
              </button>
            </Table.Cell>
          </Table.Row>
        ))}
      </Table.Body>
    </Table>
  );
};

export default UserList;
