// src/contexts/UserContext.tsx
import { createContext, createSignal, useContext, JSX, onMount } from "solid-js";

type User = {
    id: number;
    pic: string;
    name: string;
    plan: string;
    terms_accepted: boolean;
};

const UserContext = createContext<{
    user: () => User | null;
    loading: () => boolean;
    refetch: () => Promise<void>;
}>();

export function UserProvider(props: { children: JSX.Element }) {
    const [user, setUser] = createSignal<User | null>(null);
    const [loading, setLoading] = createSignal(true);

    async function fetchUser() {
        try {
            setLoading(true);
            const res = await fetch("/api/user", { credentials: "include" });
            if (res.ok) {
                const data = await res.json();
                setUser(data);
            } else {
                setUser(null);
            }
        } catch (err) {
            console.error("Failed to fetch user", err);
            setUser(null);
        } finally {
            setLoading(false);
        }
    }

    onMount(fetchUser);

    return (
        <UserContext.Provider value={{ user, loading, refetch: fetchUser }}>
            {props.children}
        </UserContext.Provider>
    );
}

export function useUser() {
    const context = useContext(UserContext);
    if (!context) throw new Error("useUser must be used within a UserProvider");
    return context;
}
