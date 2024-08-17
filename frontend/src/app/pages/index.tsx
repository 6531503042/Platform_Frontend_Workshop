// import { useState } from 'react';
// import UserForm from './UserForm';
// import UserList from './UserList';
// import UserStatisticsChart from './UserStatisticsChart';
// import { Navbar } from '../components/Navbar'

// const Home: React.FC = () => {
//   const [reloadCounter, setReloadCounter] = useState(0);

//   const handleUserAdded = () => {
//     setReloadCounter(reloadCounter + 1);
//   };

//   return (
//     <>
//       <Navbar />
//       <div className="container mx-auto mt-8">
//         <h1 className="text-2xl font-semibold mb-4">Add a New User</h1>
//         <UserForm onUserAdded={handleUserAdded} />
//         <h2 className="text-xl font-semibold mt-8 mb-4">User List</h2>
//         <UserList key={reloadCounter} />
//         <h2 className="text-xl font-semibold mt-8 mb-4">User Statistics</h2>
//         <UserStatisticsChart />
//       </div>
//     </>
//   );
// };

// export default Home;
