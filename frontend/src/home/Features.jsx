import { createEffect, createSignal } from "solid-js";
import { useUser } from "../contexts/UserContext";

import Login from '../components/Login';

function Features() {
    const { user, loading } = useUser();

    const [freePlanText, setFreePlanText] = createSignal(null);

    createEffect(() => {
        if (user() && user().plan === "none") {
            setFreePlanText(
                <>
                    <p>Accept the <a href="/terms">Terms & Conditions</a> to use titles.run</p>
                </>
            );
        }
        else if (user() && user().plan === "free") {
            setFreePlanText(
                <>
                    <p>You are currently on the free plan.</p>
                    <p>Upgrade to PRO for more features.</p>
                </>
            );
        } else if (user() && user().plan === "pro") {
            setFreePlanText(
                <p class="text-gray-600">You are currently on the Pro plan.</p>
            );
        } else if (user() === null){
            setFreePlanText(
                <Login white={true} />
            );
        } else {
            setFreePlanText(
                null
            );
        }
    })

    return (
        <section id="features" class="px-6 py-12 bg-gray-100 text-center">
            <h3 class="text-2xl font-semibold mb-4 font-inter">Features</h3>
            <p class="max-w-2xl mx-auto text-gray-600 mb-12">
                Combine your Strava activity data with state of the art AI models to create unique titles. Automatically title your activities as soon as they upload.
            </p>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 max-w-4xl w-full mx-auto">

                <div class="bg-white rounded-2xl shadow-lg flex flex-col p-6 min-h-64">
                    <h2 class="text-xl font-semibold text-gray-800 mb-4">Free</h2>
                    <div class="flex-1 mb-4 mx-auto text-left max-w-xs">
                        <ul class="space-y-2 text-gray-600">
                            <li>✅ AI generated</li>
                            <li>✅ Automatic titles</li>
                            <li>✅ 20 titles per month</li>
                        </ul>
                    </div>

                    <div class="mt-auto mx-auto text-sm text-gray-600">
                        {loading() ? null : freePlanText()}
                    </div>
                </div>

                <div class="bg-white rounded-2xl shadow-lg flex flex-col p-6 min-h-64">
                    <h2 class="text-xl font-bold text-primary mb-4">Pro</h2>
                    <div class="flex-1 mb-4 mx-auto text-left max-w-xs">
                        <ul class="space-y-2 text-gray-600">
                            <li>🧠 Advanced AI models</li>
                            <li>🎨 Personalized tone</li>
                            <li>📈 100 titles per month</li>
                        </ul>
                    </div>

                    <a href="/pro" class="mt-auto bg-primary text-white font-medium rounded-sm w-[158px] h-[32px] mx-auto flex items-center justify-center gap-0.5 font-inter hover:bg-[#941127] transition shadow-md">
                        <span>titles.run PRO</span>
                        <iconify-icon icon="basil:arrow-right-solid" width="24" height="32"></iconify-icon>
                    </a>
                </div>

            </div>
        </section>
    )
}

export default Features;