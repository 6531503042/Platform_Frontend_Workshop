import { useState } from 'react';
import { createUser } from '../../services/userService';

const CreateUser = () => {
  const [name, setName] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await createUser({ name });
    setName('');
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Create User</h1>
      <form onSubmit={handleSubmit}>
        <div className="mb-4">
          <label className="block text-sm font-medium">Name</label>
          <input
            type="text"
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mt-1 block w-full"
          />
        </div>
        <button type="submit" className="bg-blue-500 text-white px-4 py-2">Create</button>
      </form>
    </div>
  );
};

export default CreateUser;
