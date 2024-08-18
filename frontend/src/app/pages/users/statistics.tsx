import { useEffect, useState } from 'react';
import api from '../../utils/api';
import { Chart, ArcElement, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import { Bar } from 'react-chartjs-2';

Chart.register(ArcElement, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

interface UserStatistics {
  _id: string;
  count: number;
}

export default function UserStatistics() {
  const [stats, setStats] = useState<UserStatistics[]>([]);

  useEffect(() => {
    fetchStatistics();
  }, []);

  const fetchStatistics = async () => {
    try {
      const response = await api.get('/user-statistics');
      setStats(response.data);
    } catch (error) {
      console.error('Error fetching user statistics:', error);
    }
  };

  const chartData = {
    labels: stats.map(stat => stat._id),
    datasets: [{
      label: 'User Count',
      data: stats.map(stat => stat.count),
      backgroundColor: 'rgba(75, 192, 192, 0.2)',
      borderColor: 'rgba(75, 192, 192, 1)',
      borderWidth: 1,
    }]
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-semibold mb-4">User Statistics</h1>
      <Bar data={chartData} />
    </div>
  );
}
