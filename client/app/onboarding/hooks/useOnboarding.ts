"use client";

import { useReducer, useCallback } from "react";
import { useRouter } from "next/navigation";
import { completeProfile } from "@/lib/api/auth";
import { CompleteProfileResponse } from "@/lib/api/types";

interface OnboardingFormData {
  companyName: string;
  webhookUrl: string;
}

interface OnboardingState {
  formData: OnboardingFormData;
  isSubmitting: boolean;
  error: string | null;
  success: boolean;
  profileData: CompleteProfileResponse | null;
}

type OnboardingAction =
  | { type: "SET_FIELD"; field: keyof OnboardingFormData; value: string }
  | { type: "SUBMIT_START" }
  | { type: "SUBMIT_SUCCESS"; data: CompleteProfileResponse }
  | { type: "SUBMIT_ERROR"; error: string }
  | { type: "RESET_ERROR" };

const INITIAL_STATE: OnboardingState = {
  formData: {
    companyName: "",
    webhookUrl: "",
  },
  isSubmitting: false,
  error: null,
  success: false,
  profileData: null,
};

function onboardingReducer(
  state: OnboardingState,
  action: OnboardingAction
): OnboardingState {
  switch (action.type) {
    case "SET_FIELD":
      return {
        ...state,
        formData: {
          ...state.formData,
          [action.field]: action.value,
        },
        error: null,
      };
    case "SUBMIT_START":
      return {
        ...state,
        isSubmitting: true,
        error: null,
      };
    case "SUBMIT_SUCCESS":
      return {
        ...state,
        isSubmitting: false,
        success: true,
        profileData: action.data,
      };
    case "SUBMIT_ERROR":
      return {
        ...state,
        isSubmitting: false,
        error: action.error,
      };
    case "RESET_ERROR":
      return {
        ...state,
        error: null,
      };
    default:
      return state;
  }
}

interface UseOnboardingReturn {
  state: OnboardingState;
  setField: (field: keyof OnboardingFormData, value: string) => void;
  submitForm: () => Promise<void>;
  goToDashboard: () => void;
}

function validateForm(data: OnboardingFormData): string | null {
  if (!data.companyName.trim()) {
    return "El nombre de la empresa/proyecto es requerido";
  }

  if (!data.webhookUrl.trim()) {
    return "El Webhook URL es requerido";
  }

  if (!data.webhookUrl.startsWith("https://")) {
    return "El Webhook URL debe usar HTTPS";
  }

  return null;
}

export function useOnboarding(): UseOnboardingReturn {
  const router = useRouter();
  const [state, dispatch] = useReducer(onboardingReducer, INITIAL_STATE);

  const setField = useCallback(
    (field: keyof OnboardingFormData, value: string) => {
      dispatch({ type: "SET_FIELD", field, value });
    },
    []
  );

  const submitForm = useCallback(async () => {
    const validationError = validateForm(state.formData);

    if (validationError) {
      dispatch({ type: "SUBMIT_ERROR", error: validationError });
      return;
    }

    dispatch({ type: "SUBMIT_START" });

    const response = await completeProfile({
      company_name: state.formData.companyName,
      webhook_url: state.formData.webhookUrl,
    });

    if (response.error || !response.data) {
      dispatch({
        type: "SUBMIT_ERROR",
        error: response.error || "Error al completar el perfil",
      });
      return;
    }

    dispatch({ type: "SUBMIT_SUCCESS", data: response.data });
  }, [state.formData]);

  const goToDashboard = useCallback(() => {
    router.push("/dashboard");
  }, [router]);

  return {
    state,
    setField,
    submitForm,
    goToDashboard,
  };
}
