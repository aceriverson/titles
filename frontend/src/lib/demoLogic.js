export async function handleDemoForm(url, turnstileToken) {
    if (import.meta.env.MODE === 'development') {
        return "Demo activity title";
        // throw new Error("Invalid URL");
    }
    
    if (!url.includes("strava.com/activities/")) {
        throw new Error("Invalid URL");
    }

    if (!turnstileToken) {
        throw new Error("Must accept the Captcha");
    }

    const raw_gpx = await fetch("https://corsproxy.io/?url=" + url + "/export_gpx", {
        method: "GET",
    });

    if (!raw_gpx.ok) {
        throw new Error("Invalid URL");
    }

    const gpx = await raw_gpx.text();

    if (gpx.startsWith("<!DOCTYPE html>")) {
        throw new Error("Private activity");
    }

    const parser = new DOMParser();
    const gpxDoc = parser.parseFromString(gpx, "application/xml");
    const trackpoints = Array.from(gpxDoc.getElementsByTagName("trkpt")).map((trackpoint) => {
        return {
            lat: parseFloat(trackpoint.getAttribute("lat")),
            lng: parseFloat(trackpoint.getAttribute("lon")),
        };
    });

    const response = await fetch("/api/webhook/demo", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
            "Titles-CF-Turnstile-Token": turnstileToken,
        },
        body: JSON.stringify({
            points: trackpoints
        }),
    });

    if (!response.ok) {
        throw new Error("Error from server");
    }

    const title = await response.text();

    console.log(title);

    return title;
}

export async function applyTitle(title, url) {
    if (import.meta.env.MODE === 'development') {
        return;
    }

    const response = await fetch("/api/webhook/title", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            title: title,
            activity: parseInt(url.split("/").pop()),
        }),
        credentials: "include",
    });

    if (!response.ok) {
        throw new Error("Error from server");
    }

    window.location.href = url;
}