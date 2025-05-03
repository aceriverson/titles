import { useUser } from "../contexts/UserContext";

function Terms() {
    const { user } = useUser();

    const lastUpdated = "May 3, 2025";

    const submitAcceptForm = async (e) => {
        e.preventDefault();

        const response = await fetch("/api/accept_terms", {
            method: "POST",
            credentials: "include",
        });
    
        if (!response.ok) {
            throw new Error("Error from server");
        }

        window.location.href = "/";
    };

    const mustAccept = (
        <div class="container mx-auto px-4 py-8">
            You must accept the terms of service to use this website. Please read the terms carefully and click "Accept" to continue.
        </div>
    );

    const acceptForm = (
        <div class="container mx-auto px-4 py-8">
            <form onSubmit={submitAcceptForm} class="flex flex-col gap-4">
                <p class="text-sm text-gray-500">
                    By clicking "Accept", you agree to our Terms of Service and Privacy Policy.
                </p>
                <button type="submit" class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600">
                    Accept
                </button>
            </form>
        </div>
    )

    const terms = (
        <div class="container mx-auto px-4 py-8">
            <h1 class="text-3xl font-bold mb-4">Terms of Service</h1>
            <p class="mb-4">Last updated: {lastUpdated}</p>
            <p class="mb-4">
                By using our website and services, you agree to comply with and be bound by these Terms of Service. If you do not agree with any part of these terms, please do not use our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">1. Acceptance of Terms</h2>
            <p class="mb-4">
                By accessing or using our website, you confirm that you have read, understood, and agree to be bound by these Terms of Service and our Privacy Policy.
            </p>
            <h2 class="text-2xl font-semibold mb-2">2. Changes to Terms</h2>
            <p class="mb-4">
                We reserve the right to modify these Terms of Service at any time. Any changes will be effective immediately upon posting on our website. Your continued use of our services after any changes indicates your acceptance of the new terms.
            </p>
            <h2 class="text-2xl font-semibold mb-2">3. User Responsibilities</h2>
            <p class="mb-4">
                You agree to use our services only for lawful purposes and in a manner that does not infringe the rights of others or restrict or inhibit anyone else's use of our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">4. Intellectual Property</h2>
            <p class="mb-4">
                All content, trademarks, and other intellectual property on our website are the property of our company or our licensors. You may not use, reproduce, or distribute any content without our prior written consent.
            </p>
            <h2 class="text-2xl font-semibold mb-2">5. Limitation of Liability</h2>
            <p class="mb-4">
                To the fullest extent permitted by law, we shall not be liable for any direct, indirect, incidental, special, consequential, or punitive damages arising from your use of our services or any content obtained through our website.
            </p>
            <h2 class="text-2xl font-semibold mb-2">6. Governing Law</h2>
            <p class="mb-4">
                These Terms of Service shall be governed by and construed in accordance with the laws of the jurisdiction in which our company is located, without regard to its conflict of law principles.
            </p>
            <h2 class="text-2xl font-semibold mb-2">7. Contact Us</h2>
            <p class="mb-4">
                If you have any questions or concerns about these Terms of Service, please contact us at <a href="/contact" class="text-blue-500 hover:underline">titles.run/contact</a>.
            </p>
            <p class="mb-4">
                By using our services, you acknowledge that you have read and understood these Terms of Service and agree to be bound by them.
            </p>
            <h1 class="text-3xl font-bold mb-4">Privacy Policy</h1>
            <p class="mb-4">
                We value your privacy and are committed to protecting your personal information. This Privacy Policy outlines how we collect, use, and safeguard your information when you use our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">1. Information We Collect</h2>
            <p class="mb-4">
                We may collect personal information, such as your name, email address, and activity data, when you use our services. This information is used to provide and improve our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">2. Use of Information</h2>
            <p class="mb-4">
                We may use your information to communicate with you, provide customer support, and improve our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">3. Sharing of Information</h2>
            <p class="mb-4">
                We do not sell or rent your personal information to third parties. We may share your information with trusted third-party service providers who assist us in operating our website and providing our services.
            </p>
            <h2 class="text-2xl font-semibold mb-2">4. Data Security</h2>
            <p class="mb-4">
                We take reasonable measures to protect your personal information from unauthorized access, use, or disclosure. However, no method of transmission over the internet or electronic storage is 100% secure.
            </p>
            <h2 class="text-2xl font-semibold mb-2">5. Your Rights</h2>
            <p class="mb-4">
                You have the right to access, correct, or delete your personal information. If you wish to exercise these rights, please contact us at <a href="/contact" class="text-blue-500 hover:underline">titles.run/contact</a>.
            </p>
            <h2 class="text-2xl font-semibold mb-2">6. Changes to This Privacy Policy</h2>
            <p class="mb-4">
                We may update this Privacy Policy from time to time. We will notify you of any changes by posting the new Privacy Policy on our website. Your continued use of our services after any changes indicates your acceptance of the new policy.
            </p>
            <p class="mb-4">
                Thank you for using our services!
            </p>
        </div>
    );

    return (
        <>
            {user() && !user()?.terms_accepted && mustAccept}
            {terms}
            {user() && !user()?.terms_accepted && acceptForm}
        </>
    )
}

export default Terms;