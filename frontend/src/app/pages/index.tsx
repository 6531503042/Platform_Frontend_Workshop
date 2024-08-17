import { useState } from 'react';
import MyNavbar from '../components/Navbar';
import UserForm from '../components/UserForm';
import UserList from '../components/UserList';
import UserStatisticsChart from '../components/UserStatisticsChart';

const Home: React.FC = () => {
  const [reloadCounter, setReloadCounter] = useState(0);

  const handleUserAdded = () => {
    setReloadCounter(reloadCounter + 1);
  };

  return (
    <>
      <MyNavbar />
      <div className="container mx-auto mt-8">
        <h1 className="text-2xl font-semibold mb-4">Add a New User</h1>
        <UserForm onUserAdded={handleUserAdded} />
        <h2 className="text-xl font-semibold mt-8 mb-4">User List</h2>
        <UserList key={reloadCounter} />
        <h2 className="text-xl font-semibold mt-8 mb-4">User Statistics</h2>
        <UserStatisticsChart />
      </div>
    </>
  );
};

export default Home;
