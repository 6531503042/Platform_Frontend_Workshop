import { Bar } from 'react-chartjs-2';
import { Chart as ChartJS, CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend } from 'chart.js';
import { useEffect, useState } from 'react';
import api from '../axios';

ChartJS.register(CategoryScale, LinearScale, BarElement, Title, Tooltip, Legend);

interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string;
  }[];
}

const UserStatisticsChart: React.FC = () => {
  const [chartData, setChartData] = useState<ChartData>({
    labels: [],
    datasets: [],
  });

  useEffect(() => {
    fetchChartData();
  }, []);

  const fetchChartData = async () => {
    try {
      const response = await api.get('/user-statistics');
      const data = response.data;

      setChartData({
        labels: data.labels, // assuming backend sends labels for the chart
        datasets: [
          {
            label: 'User Registrations',
            data: data.values, // assuming backend sends values for the chart
            backgroundColor: 'rgba(54, 162, 235, 0.6)',
          },
        ],
      });
    } catch (error) {
      console.error('Failed to fetch user statistics:', error);
    }
  };

  return (
    <div>
      <h2 className="text-xl font-semibold mb-4">User Statistics</h2>
      <Bar
        data={chartData}
        options={{
          responsive: true,
          plugins: {
            legend: {
              position: 'top',
            },
            title: {
              display: true,
              text: 'User Registrations Over Time',
            },
          },
        }}
      />
    </div>
  );
};

export default UserStatisticsChart;
