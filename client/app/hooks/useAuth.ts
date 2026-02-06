"use client";

import { useEffect, useState } from "react";
import { useRouter, usePathname } from "next/navigation";
import { getCurrentUser } from "@/lib/api/auth";
import { UserProfile } from "@/lib/api/types";

interface AuthState {
  isLoading: boolean;
  isAuthenticated: boolean;
  needsProfileCompletion: boolean;
  user: UserProfile | null;
  error: string | null;
}

interface UseAuthReturn extends AuthState {
  refetch: () => Promise<void>;
}

const INITIAL_STATE: AuthState = {
  isLoading: true,
  isAuthenticated: false,
  needsProfileCompletion: false,
  user: null,
  error: null,
};

function checkHasToken(): boolean {
  if (typeof document === "undefined") return false;
  return document.cookie.includes("auth_token=");
}

export function useAuth(): UseAuthReturn {
  const router = useRouter();
  const pathname = usePathname();
  const [state, setState] = useState<AuthState>(INITIAL_STATE);

  const fetchUserData = async (): Promise<void> => {
    const hasToken = checkHasToken();

    if (!hasToken) {
      setState({
        ...INITIAL_STATE,
        isLoading: false,
      });
      return;
    }

    const response = await getCurrentUser();

    if (response.error || !response.data) {
      setState({
        ...INITIAL_STATE,
        isLoading: false,
        error: response.error || "Failed to fetch user data",
      });
      return;
    }

    const userData = response.data;
    const needsCompletion = !userData.company_name || !userData.webhook_url;

    setState({
      isLoading: false,
      isAuthenticated: true,
      needsProfileCompletion: needsCompletion,
      user: userData,
      error: null,
    });

    // Redirigir según el estado del perfil
    if (needsCompletion && pathname === "/dashboard") {
      router.push("/onboarding");
    } else if (!needsCompletion && pathname === "/onboarding") {
      router.push("/dashboard");
    }
  };

  useEffect(() => {
    fetchUserData();
  }, [pathname]);

  return {
    ...state,
    refetch: fetchUserData,
  };
}
