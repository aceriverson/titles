function Hero() {
    return (
        <section class="px-6 py-16 text-center gap-8 flex flex-col items-center justify-center">
            <div>
                <h2 class="text-3xl sm:text-5xl font-bold leading-tight font-inter">
                    Automatic Strava Titles
                </h2>
                <p class="mt-4 text-gray-600 max-w-xl mx-auto">
                    Use AI to describe your activity's route and title it automatically. Try it instantly.
                </p>
            </div>
            <a href="#demo" class="mt-6 bg-primary text-white font-medium rounded-sm w-[158px] h-[32px] mx-auto flex items-center justify-center gap-0.5 font-inter hover:bg-[#941127] transition shadow-md">
                Try the Demo
            </a>
        </section>
    )
}

export default Hero;