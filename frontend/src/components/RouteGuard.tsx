// components/RouteGuard.tsx
import { useUser } from "../contexts/UserContext";
import { useLocation, useNavigate } from "@solidjs/router";
import { JSX, createEffect } from "solid-js";

const PROTECTED_ROUTES = ["/contact"];

export default function RouteGuard(props: { children: JSX.Element }) {
  const { user, loading } = useUser();
  const location = useLocation();
  const navigate = useNavigate();

  createEffect(() => {
    if (loading()) return;

    const path = location.pathname;

    if (!user() && PROTECTED_ROUTES.includes(path)) {
      navigate("/", { replace: true });
      return;
    }

    if (user() && !user()?.terms_accepted && path !== "/terms") {
      navigate("/terms", { replace: true });
      return;
    }
  });

  return props.children;
}
