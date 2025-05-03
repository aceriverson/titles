import { useUser } from "../contexts/UserContext";

import ctxwstrava from '../assets/ctxwstrava.svg';
import ctxwstravawhite from '../assets/ctxwstravawhite.svg';

function Login({ white }) {
    const { user, loading } = useUser();

    return (
        <>
            {loading() ? null : user() && user()?.pic ? (
                <a href="/profile">
                    <img src={user().pic} class="h-8 w-8 rounded-full shadow-md" />
                </a>
            ) : (
                <a
                    href={`https://www.strava.com/oauth/authorize?client_id=110809&response_type=code&redirect_uri=https://${import.meta.env.VITE_HOST_URL}/auth/callback&approval_prompt=auto&scope=activity:write,activity:read_all`}
                >
                    <img src={white ? ctxwstravawhite : ctxwstrava} class="h-8 shadow-md" />
                </a>
            )}
        </>
    );
}


export default Login;