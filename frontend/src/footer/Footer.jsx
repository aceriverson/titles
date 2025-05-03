import poweredbystrava from '../assets/poweredbystrava.svg';

function Footer() {
    return (
        <>
            <div class="p-7"></div>
            <footer id="contact" class="mt-auto px-6 py-6 flex justify-center text-sm text-gray-500 border-t fixed bottom-0 w-full bg-white z-50">
                <img src={poweredbystrava} class="h-3" />
            </footer>
        </>
    )
}

export default Footer;