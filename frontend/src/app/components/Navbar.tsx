// components/AppNavbar.js
import Link from 'next/link';

export default function AppNavbar() {
  return (
    <nav className="bg-gray-800 p-4">
      <div className="container mx-auto flex justify-between items-center">
        <Link href="/" className="text-white text-xl font-semibold">
          MyApp
        </Link>
        <div className="space-x-4">
          <Link href="/" className="text-white hover:text-gray-300">
            Home
          </Link>
          <Link href="/users" className="text-white hover:text-gray-300">
            Users
          </Link>
          <Link href="/statistics" className="text-white hover:text-gray-300">
            Statistics
          </Link>
        </div>
        <button className="bg-blue-500 text-white px-4 py-2 rounded">
          Login
        </button>
      </div>
    </nav>
  );
}
