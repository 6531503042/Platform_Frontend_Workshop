import { useEffect, useState } from 'react';
import { listUsers, deleteUser } from '../../serviecs/userService';

const UserList = () => {
  const [users, setUsers] = useState([]);

  useEffect(() => {
    const fetchUsers = async () => {
      const response = await listUsers();
      setUsers(response.data);
    };
    fetchUsers();
  }, []);

  const handleDelete = async (id: string) => {
    await deleteUser(id);
    setUsers(users.filter((user: any) => user._id !== id));
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">User List</h1>
      <ul className="list-disc list-inside">
        {users.map((user: any) => (
          <li key={user._id} className="flex justify-between">
            <span>{user.name}</span>
            <button className="text-red-500" onClick={() => handleDelete(user._id)}>Delete</button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default UserList;
