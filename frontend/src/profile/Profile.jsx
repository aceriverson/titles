import { createEffect, createSignal } from "solid-js";

import "@iconify-icon/solid";

import { useUser } from "../contexts/UserContext";

function Profile() {
    const { user } = useUser();

    const [automaticTitles, setAutomaticTitles] = createSignal(true);
    const [creativeSlider, setCreativeSlider] = createSignal(50);
    const [showAttribution, setShowAttribution] = createSignal(true);
    const [aiDescription, setAiDescription] = createSignal(false);

    createEffect(() => {
        const u = user();
        if (!u) return;

        setAutomaticTitles(u.settings?.automatic_title ?? true);
        setCreativeSlider(u.settings?.tone ?? 50);
        setShowAttribution(u.settings?.attribution ?? true);
        setAiDescription(u.settings?.description ?? false);
    });

    const handleLogoutForm = async (e) => {
        e.preventDefault();
        const response = await fetch("/auth/logout", {
            method: "POST",
            credentials: "include",
        });
        if (!response.ok) throw new Error("Error from server");
        window.location.href = "/";
    };

    const isPro = () => user()?.plan === "pro";

    const saveButton = () => {
        return (
            <button
                class="mt-6 bg-primary text-white font-medium px-6 py-2 rounded shadow hover:bg-[#941127] transition font-inter w-full sm:w-auto"
                onClick={async () => {
                    const payload = {
                        automatic_title: automaticTitles(),
                        tone: creativeSlider(),
                        attribution: showAttribution(),
                        description: aiDescription(),
                    };

                    const res = await fetch("/api/settings", {
                        method: "POST",
                        headers: {
                            "Content-Type": "application/json",
                        },
                        credentials: "include",
                        body: JSON.stringify(payload),
                    });

                    if (!res.ok) {
                        console.error("Failed to save settings");
                        alert("There was a problem saving your settings.");
                    } else {
                        alert("Settings saved successfully.");
                    }
                }}
            >
                Save Settings
            </button>
        )
    }

    return (
        <div class="max-w-xl mx-auto px-4 py-8 space-y-6">
            <h3 class="text-3xl font-bold font-inter text-center">Profile</h3>
            <p class="text-center text-sm text-gray-600 font-inter">
                Plan: <span class="font-semibold uppercase">{user()?.plan.charAt(0).toUpperCase() + user()?.plan.slice(1)}</span>
            </p>

            <div class="flex flex-col sm:flex-row justify-center items-center gap-4">
                <button
                    onClick={handleLogoutForm}
                    type="submit"
                    class="bg-primary text-white font-medium px-6 py-2 rounded shadow hover:bg-[#941127] transition font-inter flex items-center gap-2"
                >
                    <span>Logout</span>
                    <iconify-icon icon="ic:round-logout" width="24" height="24"></iconify-icon>
                </button>
                <a
                    href="https://www.strava.com/settings/apps"
                    target="_blank"
                    rel="noopener noreferrer"
                    class="border border-red-600 text-red-600 px-6 py-2 rounded hover:bg-red-50 transition font-medium font-inter"
                >
                    Disconnect titles.run
                </a>
            </div>

            <div class="border-t pt-6">
                <h4 class="text-xl font-semibold font-inter mb-4">Settings</h4>
                <div class="space-y-4">
                    {/* Automatic Titles */}
                    <label class="flex items-center gap-3">
                        <input
                            type="checkbox"
                            checked={automaticTitles()}
                            onChange={(e) => setAutomaticTitles(e.currentTarget.checked)}
                            class="accent-primary"
                        />
                        <span class="text-sm font-inter">Automatic titles</span>
                    </label>

                    {/* Creative/Terse Slider (PRO only) */}
                    <div class={`space-y-1 ${!isPro() && "opacity-50"}`}>
                        <label for="creative-slider" class="text-sm font-inter">
                            Creative vs. Terse
                        </label>
                        <input
                            id="creative-slider"
                            type="range"
                            min="0"
                            max="100"
                            value={creativeSlider()}
                            onInput={(e) => setCreativeSlider(+e.currentTarget.value)}
                            disabled={!isPro()}
                            class="w-full accent-primary"
                        />
                    </div>

                    {/* Attribution (PRO only) */}
                    <label
                        class={`flex items-center gap-3 ${!isPro() && "opacity-50"}`}
                    >
                        <input
                            type="checkbox"
                            checked={showAttribution()}
                            onChange={(e) => setShowAttribution(e.currentTarget.checked)}
                            disabled={!isPro()}
                            class="accent-primary"
                        />
                        <span class="text-sm font-inter">“Titled via titles.run” attribution</span>
                    </label>

                    {/* AI Description (PRO only) */}
                    <label
                        class={`flex items-center gap-3 ${!isPro() && "opacity-50"}`}
                    >
                        <input
                            type="checkbox"
                            checked={aiDescription()}
                            onChange={(e) => setAiDescription(e.currentTarget.checked)}
                            disabled={!isPro()}
                            class="accent-primary"
                        />
                        <span class="text-sm font-inter">AI activity description (beta)</span>
                    </label>
                </div>
                {saveButton()}
            </div>
        </div>
    );
}

export default Profile;
