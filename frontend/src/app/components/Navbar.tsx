import Link from 'next/link';
import { Navbar } from 'flowbite-react';

const MyNavbar: React.FC = () => {
  return (
    <Navbar fluid={true} rounded={true}>
      <Navbar.Brand href="/">
        <span className="self-center text-xl font-semibold whitespace-nowrap dark:text-white">
          MyApp
        </span>
      </Navbar.Brand>
      <Navbar.Toggle />
      <Navbar.Collapse>
        <Navbar.Link href="/" active={true}>
          Home
        </Navbar.Link>
        <Navbar.Link href="/users">Users</Navbar.Link>
        <Navbar.Link href="/about">About</Navbar.Link>
      </Navbar.Collapse>
    </Navbar>
  );
};

export default MyNavbar;
