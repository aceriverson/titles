import { createEffect, createSignal } from "solid-js";
import '@iconify-icon/solid';

import { useUser } from "../contexts/UserContext";
import Turnstile from "../components/Turnstile";
import { handleDemoForm, applyTitle } from "../lib/demoLogic";

function Demo() {
    const { user } = useUser();

    const [url, setUrl] = createSignal("");
    const [token, setToken] = createSignal("");
    const [status, setStatus] = createSignal("form");
    const [activityTitle, setActivityTitle] = createSignal("");
    const [demoError, setDemoError] = createSignal("");

    let turnstileApi;

    createEffect(() => {
        if (status() === "form") {
            turnstileApi?.reset();
        }
    })

    const reset = () => {
        turnstileApi?.reset();
        setUrl("");
        setStatus("form");
        setActivityTitle("");
        setDemoError("");
    }

    const submitDemo = async (event) => {
        event.preventDefault();

        const inputUrl = url().trim();
        if (!inputUrl) return;

        const turnstileToken = token();

        try {
            setStatus("waiting");
            const res = await handleDemoForm(inputUrl, turnstileToken);
            setActivityTitle(res);
            setStatus("done");
        } catch (err) {
            console.error("Error submitting form:", err);
            setDemoError(err.message);
            setStatus("form");
        }
    }

    const error = (
        <p class="text-left"><span class="font-semibold text-red-700 ">Error:</span> {demoError()}</p>
    )

    const form = (
        <form
            onSubmit={submitDemo}
            class="max-w-md mx-auto flex flex-col gap-4"
        >
            <div class="max-w-md mx-auto flex flex-col items-center sm:flex-row gap-4">
                <div>
                    <input
                        type="text"
                        placeholder="https://strava.com/activities/..."
                        value={url()}
                        onInput={(e) => setUrl(e.currentTarget.value)}
                        class={`flex-1 px-4 py-2 border border-gray-300 rounded-sm focus:outline-none focus:ring-2 focus:ring-primary ${demoError() ? "border-red-500" : "border-gray-300"}`}
                    />
                    {demoError() && error}
                </div>
                <button
                    type="submit"
                    disabled={!url()}
                    class="bg-primary text-white font-medium rounded-sm w-[158px] h-[32px] mx-auto flex items-center justify-center gap-0.5 font-inter hover:bg-[#941127] transition shadow-md"
                >
                    Submit
                </button>
            </div>
            <Turnstile
                siteKey={import.meta.env.VITE_CF_TURNSTILE_SITEKEY}
                onVerify={(t) => setToken(t)}
                ref={(api) => (turnstileApi = api)}
            />
        </form>
    )

    const waiting = (
        <>
            <p class="mt-4 text-gray-500">Generating title...</p>
            <iconify-icon icon="eos-icons:three-dots-loading" width="24" height="24"></iconify-icon>
        </>
    )

    const result = (
        <div class="mt-12 max-w-md mx-auto">
            <div class="mt-6 max-w-md mx-auto bg-white border border-gray-300 text-gray-800 rounded-lg p-4 shadow-sm">
                <p class="text-lg font-mono break-words font-semibold">
                    "{activityTitle()}"
                </p>
            </div>
            <div class="relative w-full flex justify-center">
                <button onClick={() => reset()} class="absolute left-0"><iconify-icon icon="system-uicons:reset" width="18" height="18" class="px-1 py-1 hover:border-b-2"></iconify-icon></button>
                { user() && <button onClick={() => applyTitle(activityTitle(), url())} class="text-md text-semibold hover:underline pr-2 text-primary">Apply Title</button> }
                <a href={url()} target="_blank" class="absolute right-0 text-sm hover:underline pr-2">View on Strava</a>
            </div>
        </div>
    )

    const demoMap = {
        "form": form,
        "waiting": waiting,
        "done": result,
        "error": error,
    }


    return (
        <section id="demo" class="px-6 py-12 text-center h-96">
            <h3 class="text-2xl font-semibold mb-4 font-inter">Live Demo</h3>
            <p class="text-gray-600 mb-6">Enter a Strava activity URL and see titles.run in action.</p>
            {demoMap[status()]}
        </section>
    )
}

export default Demo;