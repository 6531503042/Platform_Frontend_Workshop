import { useState, useEffect } from 'react';
import { getUser, updateUser } from '../../serviecs/userService'
import { useRouter } from 'next/router';

const EditUser = () => {
  const [name, setName] = useState('');
  const router = useRouter();
  const { id } = router.query;

  useEffect(() => {
    const fetchUser = async () => {
      if (id) {
        const response = await getUser(id as string);
        setName(response.data.name);
      }
    };
    fetchUser();
  }, [id]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    await updateUser(id as string, { name });
    router.push('/users');
  };

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-4">Edit User</h1>
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
        <button type="submit" className="bg-blue-500 text-white px-4 py-2">Save</button>
      </form>
    </div>
  );
};

export default EditUser;
