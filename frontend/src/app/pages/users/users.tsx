import { useEffect, useState } from 'react';
import api from '../../utils/api';
import { User } from '../../types/types'; // Ensure correct path
import { Button, Modal, Table } from 'flowbite-react';

export default function Users() {
  const [users, setUsers] = useState<User[]>([]);
  const [newUser, setNewUser] = useState({ name: '', email: '' });
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [modalOpen, setModalOpen] = useState(false);

  useEffect(() => {
    console.log('Fetching users...');
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const response = await api.get('/users');
      console.log('Users fetched:', response.data);
      setUsers(response.data);
    } catch (error) {
      console.error('Error fetching users:', error);
    }
  };

  const handleAddUser = async () => {
    try {
      await api.post('/users', newUser);
      fetchUsers();
      setNewUser({ name: '', email: '' });
    } catch (error) {
      console.error('Error adding user:', error);
    }
  };

  const handleDeleteUser = async (id: string) => {
    try {
      await api.delete(`/users/${id}`);
      fetchUsers();
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  const handleEditUser = async () => {
    if (editingUser) {
      try {
        await api.put(`/users/${editingUser.id}`, editingUser);
        fetchUsers();
        setEditingUser(null);
        setModalOpen(false);
      } catch (error) {
        console.error('Error updating user:', error);
      }
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-semibold mb-4">User Management</h1>

      <div className="mb-4">
        <input
          type="text"
          placeholder="Name"
          value={newUser.name}
          onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
          className="p-2 border rounded mr-2"
        />
        <input
          type="email"
          placeholder="Email"
          value={newUser.email}
          onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
          className="p-2 border rounded mr-2"
        />
        <Button onClick={handleAddUser} color="blue">
          Add User
        </Button>
      </div>

      <Table>
        <Table.Head>
          <Table.HeadCell>ID</Table.HeadCell>
          <Table.HeadCell>Name</Table.HeadCell>
          <Table.HeadCell>Email</Table.HeadCell>
          <Table.HeadCell>Actions</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {users.map((user) => (
            <Table.Row key={user.id}>
              <Table.Cell>{user.id}</Table.Cell>
              <Table.Cell>{user.name}</Table.Cell>
              <Table.Cell>{user.email}</Table.Cell>
              <Table.Cell>
                <Button onClick={() => { setEditingUser(user); setModalOpen(true); }} color="yellow" className="mr-2">
                  Edit
                </Button>
                <Button onClick={() => handleDeleteUser(user.id)} color="red">
                  Delete
                </Button>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table>

      {editingUser && (
        <Modal show={modalOpen} onClose={() => setModalOpen(false)}>
          <Modal.Header>Edit User</Modal.Header>
          <Modal.Body>
            <input
              type="text"
              placeholder="Name"
              value={editingUser.name}
              onChange={(e) => setEditingUser({ ...editingUser, name: e.target.value })}
              className="p-2 border rounded mb-2"
            />
            <input
              type="email"
              placeholder="Email"
              value={editingUser.email}
              onChange={(e) => setEditingUser({ ...editingUser, email: e.target.value })}
              className="p-2 border rounded"
            />
          </Modal.Body>
          <Modal.Footer>
            <Button onClick={handleEditUser} color="yellow">
              Save Changes
            </Button>
            <Button onClick={() => setModalOpen(false)} color="gray">
              Cancel
            </Button>
          </Modal.Footer>
        </Modal>
      )}
    </div>
  );
}
