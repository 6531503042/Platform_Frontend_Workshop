// components/Navbar.tsx
import Link from 'next/link';
import { Navbar as FlowbiteNavbar } from 'flowbite-react';

const Navbar = () => {
  return (
    <FlowbiteNavbar fluid={true} rounded={true}>
      <FlowbiteNavbar.Brand href="/">
        <span className="self-center whitespace-nowrap text-xl font-semibold dark:text-white">
          My App
        </span>
      </FlowbiteNavbar.Brand>
      <FlowbiteNavbar.Toggle />
      <FlowbiteNavbar.Collapse>
        <FlowbiteNavbar.Link href="/" active={true}>
          Home
        </FlowbiteNavbar.Link>
        <FlowbiteNavbar.Link href="/user">
          User
        </FlowbiteNavbar.Link>
        <FlowbiteNavbar.Link href="/chart">
          Chart
        </FlowbiteNavbar.Link>
      </FlowbiteNavbar.Collapse>
    </FlowbiteNavbar>
  );
};

export default Navbar;
