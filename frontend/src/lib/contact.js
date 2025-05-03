export async function handleContactForm(email, subject, message, turnstileToken) {    
    if (!turnstileToken) {
        throw new Error("Must accept the Captcha");
    }

    const response = await fetch("/api/contact", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Titles-CF-Turnstile-Token": turnstileToken,
        },
        body: JSON.stringify({
            email: email,
            subject: subject,
            message: message,
        }),
    });

    if (!response.ok) {
        throw new Error("Error from server");
    }

    return;
}