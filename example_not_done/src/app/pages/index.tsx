import Users from './users/users';
import UserStatistics from './users/statistics';

export default function Home() {
  return (
    <div>
      <Users />
      <UserStatistics />
    </div>
  );
}
