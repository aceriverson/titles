import { createSignal } from "solid-js";

import Turnstile from "../components/Turnstile";
import { handleContactForm } from "../lib/contact";

function Contact() {
    const [email, setEmail] = createSignal("");
    const [subject, setSubject] = createSignal("");
    const [message, setMessage] = createSignal("");
    const [token, setToken] = createSignal("");
    const [error, setError] = createSignal(null);

    const submitContactForm = async (e) => {
        e.preventDefault();

        if (!email() || !subject() || !message() || !token()) {
            setError("Please fill in all fields and complete the captcha.");
            return;
        }

        setError(null);

        try {
            await handleContactForm(email(), subject(), message(), token())
        } catch (err) {
            setError(err.message);
            console.error("Error submitting form:", err);
            return;
        }

        window.location.href = "/";
    };

    const form = (
        <form onSubmit={submitContactForm} class="max-w-md mx-auto flex flex-col gap-4 items-center">
            <div class="flex flex-col gap-4">
                <div>
                    <input
                        type="email"
                        placeholder="Your email"
                        value={email()}
                        onInput={(e) => setEmail(e.currentTarget.value)}
                        class={`px-4 py-2 border border-gray-300 rounded-sm focus:outline-none focus:ring-2 focus:ring-primary ${error() ? "border-red-500" : "border-gray-300"}`}
                    />
                </div>
                <div>
                    <input
                        type="text"
                        placeholder="Subject"
                        value={subject()}
                        onInput={(e) => setSubject(e.currentTarget.value)}
                        class={`px-4 py-2 border border-gray-300 rounded-sm focus:outline-none focus:ring-2 focus:ring-primary ${error() ? "border-red-500" : "border-gray-300"}`}
                    />
                </div>
                <div>
                    <textarea
                        placeholder="Your message"
                        value={message()}
                        onInput={(e) => setMessage(e.currentTarget.value)}
                        class={`px-4 py-2 border border-gray-300 rounded-sm focus:outline-none focus:ring-2 focus:ring-primary ${error() ? "border-red-500" : "border-gray-300"}`}
                    />
                </div>
                {error() && <p class="text-red-500 text-sm">{error()}</p>}
            </div>
            <button
                type="submit"
                disabled={!email() || !subject() || !message() || !token()}
                class="bg-primary text-white font-medium rounded-sm w-[158px] h-[32px] mx-auto flex items-center justify-center gap-0.5 font-inter hover:bg-[#941127] transition shadow-md"
            >
                Submit
            </button>

            <Turnstile
                siteKey={import.meta.env.VITE_CF_TURNSTILE_SITEKEY}
                onVerify={(t) => setToken(t)}
            />
        </form>
    );


    return (
        <>
            <h3 class="text-2xl font-semibold mb-4 font-inter mx-auto pt-4">Contact</h3>
            {form}
        </>
    )
}

export default Contact;