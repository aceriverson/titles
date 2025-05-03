function Profile() {
    const handleLogoutForm = async (e) => {
        e.preventDefault();

        const response = await fetch("/auth/logout", {
            method: "POST",
            credentials: "include",
        });

        if (!response.ok) {
            throw new Error("Error from server");
        }

        window.location.href = "/";
    }

    const logout = (
        <div class="flex justify-center">
            <button
                onClick={handleLogoutForm}
                type="submit"
                class="bg-primary text-white font-medium rounded-sm w-[158px] h-[32px] mx-auto flex items-center justify-center gap-0.5 font-inter hover:bg-[#941127] transition shadow-md"
            >
                <span>Logout</span>
                <iconify-icon icon="ic:round-logout" width="24" height="32"></iconify-icon>
            </button>
        </div>
    )

    return (
        <>
            <h3 class="text-2xl font-semibold mb-4 font-inter mx-auto pt-4">Profile</h3>
            {logout}
        </>
    );
}

export default Profile;