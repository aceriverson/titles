import { onMount, createSignal } from "solid-js";

export default function Turnstile({ siteKey, onVerify, ref }) {
    let turnstileRef;
    const [widgetId, setWidgetId] = createSignal(null);

    async function loadTurnstileScript() {
        if (window.turnstile) return;

        const existing = document.querySelector("script[src*='challenges.cloudflare.com']");
        if (existing) return new Promise((resolve) => existing.addEventListener("load", () => resolve()));

        return new Promise((resolve, reject) => {
            const script = document.createElement("script");
            script.src = "https://challenges.cloudflare.com/turnstile/v0/api.js";
            script.async = true;
            script.defer = true;
            script.onload = () => resolve();
            script.onerror = () => reject(new Error("Failed to load Turnstile"));
            document.body.appendChild(script);
        });
    }

    function renderTurnstile() {
        if (window.turnstile && turnstileRef) {
            const id = window.turnstile.render(turnstileRef, {
                sitekey: siteKey,
                callback: onVerify,
            });
            setWidgetId(id);
        }
    }

    function reset() {
        const id = widgetId();
        if (window.turnstile && id) {
            window.turnstile.reset(id);
        }
    }

    onMount(async () => {
        await loadTurnstileScript();
        renderTurnstile();
        ref?.({ reset });
    });

    return <div ref={turnstileRef}></div>;
}
