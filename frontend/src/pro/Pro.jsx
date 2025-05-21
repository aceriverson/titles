import { useUser } from "../contexts/UserContext";

import Login from '../components/Login';
import MonthlyButton from "./MonthlyButton";
import YearlyButton from "./YearlyButton";

function Pro() {
  const { user, loading } = useUser();

  const subscription = () => {
    if (user() && user().plan === "pro") {
      return (
        <p class="text-left">
          <span class="font-semibold">Status:</span> Active
        </p>
      );
    } else if (user() && user().plan === "free") {
      return (
        <div class="flex flex-col gap-4 sm:flex-row sm:gap-8">
          <YearlyButton user={user().id} />
          <MonthlyButton user={user().id} />
        </div>
      );
    }
    return <Login />;
  }

  return (
    <section class="px-6 py-16 text-center gap-8 flex flex-col items-center justify-center">
      <div class="bg-white rounded-2xl shadow-lg flex flex-col p-6 min-h-64">
        <h2 class="text-3xl sm:text-5xl font-bold leading-tight font-inter text-primary">
          PRO
        </h2>
        <div class="m-auto text-left max-w-xs">
          <ul class="space-y-2 text-gray-600">
            <li>🧠 Advanced AI models</li>
            <li>🎨 Personalized tone</li>
            <li>📈 100 titles per month</li>
          </ul>
        </div>
      </div>
      {subscription()}
    </section>
  );
}

export default Pro;
