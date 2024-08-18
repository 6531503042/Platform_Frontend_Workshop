import '../styles/globals.css';
import type { AppProps } from 'next/app';
import AppNavbar from '../components/Navbar'

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <AppNavbar />
      <main className="container mx-auto p-4">
        <Component {...pageProps} />
      </main>
    </>
  );
}

export default MyApp;
