import { useEffect, useState } from 'react';
import { Bar } from 'react-chartjs-2';
import 'chart.js/auto';

export default function Statistics() {
  const [chartData, setChartData] = useState({});

  useEffect(() => {
    fetch('/api/user-statistics')
      .then((response) => response.json())
      .then((data) => {
        const labels = data.map((item: any) => item._id);
        const counts = data.map((item: any) => item.count);
        
        setChartData({
          labels: labels,
          datasets: [
            {
              label: 'User Registrations by Month',
              data: counts,
              backgroundColor: 'rgba(75,192,192,1)',
              borderColor: 'rgba(75,192,192,1)',
              borderWidth: 1,
            },
          ],
        });
      });
  }, []);

  return (
    <div className="container mx-auto">
      <h1 className="text-2xl font-semibold my-4">User Statistics</h1>
      <div className="bg-white p-4 shadow rounded">
        <Bar data={chartData} />
      </div>
    </div>
  );
}
