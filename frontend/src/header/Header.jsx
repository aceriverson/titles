import { createSignal } from 'solid-js';

import Login from '../components/Login';

import '@iconify-icon/solid';

function Header() {
    const [isOpen, setIsOpen] = createSignal(false);

    const login = <Login />;

    return (
        <header class="w-full px-4 py-3 flex justify-between items-center shadow-sm sticky top-0 bg-white z-50">
            <a href="/">
                <h1 class="text-xl text-primary font-inter tracking-wider">TITLES.run</h1>
            </a>

            {/* Desktop Nav */}
            <nav class="hidden sm:flex gap-5 items-center">
                <div class="sm:flex gap-4">
                    <a href="/pro" class="font-bold text-primary hover:underline">PRO</a>
                    <a href="/#features" class="text-gray-600 hover:text-primary">Features</a>
                    <a href="/contact" class="text-gray-600 hover:text-primary">Contact</a>
                </div>
                {login}
            </nav>

            {/* Hamburger Button */}
            <button
                class="sm:hidden text-primary flex"
                aria-label="Toggle menu"
                onClick={() => setIsOpen(!isOpen())}
            >
                <iconify-icon icon={isOpen() ? "mdi:close" : "mdi:menu"} width="24" height="24" />
            </button>

            {/* Mobile Nav Dropdown */}
            {isOpen() && (
                <div class="absolute top-full left-0 w-full bg-white shadow-md flex flex-col gap-3 p-4 sm:hidden z-40">
                    <a href="/pro" class="font-bold text-primary hover:underline" onClick={() => setIsOpen(false)}>PRO</a>
                    <a href="/#features" class="text-gray-600 hover:text-primary" onClick={() => setIsOpen(false)}>Features</a>
                    <a href="/contact" class="text-gray-600 hover:text-primary" onClick={() => setIsOpen(false)}>Contact</a>
                    {login}
                </div>
            )}
        </header>
    );
}

export default Header;