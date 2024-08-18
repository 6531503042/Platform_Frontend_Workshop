import { Flowbite, Footer as FlowbiteFooter } from 'flowbite-react';

export const Footer = () => {
    return (
        <Flowbite>
            <FlowbiteFooter container={true} className="bg-gray-800 text-white">
                <div className="w-full text-center">
                    <div className="w-full sm:flex sm:items-center sm:justify-between">
                        <FlowbiteFooter.Brand
                            href="https://yourwebsite.com"
                            src="https://flowbite.com/docs/images/logo.svg"
                            alt="Your Logo"
                            name="Your Company"
                        />
                        <FlowbiteFooter.LinkGroup className="mt-4 sm:mt-0">
                            <FlowbiteFooter.Link href="#">About</FlowbiteFooter.Link>
                            <FlowbiteFooter.Link href="#">Privacy Policy</FlowbiteFooter.Link>
                            <FlowbiteFooter.Link href="#">Licensing</FlowbiteFooter.Link>
                            <FlowbiteFooter.Link href="#">Contact</FlowbiteFooter.Link>
                        </FlowbiteFooter.LinkGroup>
                    </div>
                    <FlowbiteFooter.Divider />
                    <FlowbiteFooter.Copyright
                        href="https://yourwebsite.com"
                        by="Your Company™"
                        year={2024}
                    />
                </div>
            </FlowbiteFooter>
        </Flowbite>
    );
};
