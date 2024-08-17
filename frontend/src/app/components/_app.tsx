// pages/_app.tsx
import { AppProps } from 'next/app';
import Navbar from '../components/Navbar';
import '../styles/globals.css';

const MyApp = ({ Component, pageProps }: AppProps) => {
  return (
    <>
      <Navbar />
      <main className="container mx-auto p-4">
        <Component {...pageProps} />
      </main>
    </>
  );
};

export default MyApp;
