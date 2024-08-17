import { Button, Label, TextInput } from 'flowbite-react';
import { useState } from 'react';
import api from '../axios';

interface UserFormProps {
  onUserAdded: () => void;
}

const UserForm: React.FC<UserFormProps> = ({ onUserAdded }) => {
  const [email, setEmail] = useState<string>('');
  const [password, setPassword] = useState<string>('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await api.post('/users', { email, password });
      setEmail('');
      setPassword('');
      onUserAdded(); // Refresh the user list
    } catch (error) {
      console.error('Failed to add user:', error);
    }
  };

  return (
    <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
      <div>
        <Label htmlFor="email" value="Email" />
        <TextInput
          id="email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required={true}
        />
      </div>
      <div>
        <Label htmlFor="password" value="Password" />
        <TextInput
          id="password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required={true}
        />
      </div>
      <Button type="submit">Add User</Button>
    </form>
  );
};

export default UserForm;
